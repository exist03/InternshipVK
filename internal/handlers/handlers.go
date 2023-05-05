package handlers

import (
	"database/sql"
	"fmt"
	"fsm/internal/keyboards"
	"fsm/pkg/models/mysql"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
	"unicode/utf8"
)

var (
	InputSG            = fsm.NewStateGroup("reg")
	InputServiceState  = InputSG.New("inputService")
	InputLoginState    = InputSG.New("login")
	InputPasswordState = InputSG.New("password")
	InputConfirmState  = InputSG.New("confirm")
)

func InitHandlers(bot *tele.Group, db *sql.DB, manager *fsm.Manager) {
	initDelHandlers(db, manager)
	initGetHandlers(db, manager)
	bot.Handle("/start", onStart)
	manager.Bind("/set", fsm.DefaultState, onStartRegister(keyboards.CancelBtn))
	manager.Bind("/cancel", fsm.AnyState, onCancelForm(keyboards.SetBtn))

	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.FSMContext) error {
		s := state.State()
		return c.Send(s.String())
	})

	// buttons
	manager.Bind(&keyboards.SetBtn, fsm.DefaultState, onStartRegister(keyboards.CancelBtn))
	manager.Bind(&keyboards.CancelBtn, fsm.AnyState, onCancelForm(keyboards.SetBtn))

	// form
	manager.Bind(tele.OnText, InputServiceState, onInputServiceRegister)
	manager.Bind(tele.OnText, InputLoginState, onInputLogin)
	manager.Bind(tele.OnText, InputPasswordState, onInputPassword(keyboards.ConfirmBtn, keyboards.ResetFormBtn, keyboards.CancelInlineBtn))
	manager.Bind(&keyboards.ConfirmBtn, InputConfirmState, onInputConfirm(db), deleteAfterHandler)
	manager.Bind(&keyboards.ResetFormBtn, InputConfirmState, onInputResetForm, deleteAfterHandler)
	manager.Bind(&keyboards.CancelInlineBtn, InputConfirmState, onCancelForm(keyboards.SetBtn), deleteAfterHandler)
}

func onStart(c tele.Context) error {
	log.Println("new user", c.Sender().ID)
	return c.Send(
		"Добро пожаловать в бот для ваших паролей\n"+
			"Отправьте /set чтобы добавить сервис\n"+
			"Отправьте /cancel чтобы омтенить ввод", keyboards.OnStartKB())

}

func onStartRegister(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancelBtn))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(InputServiceState)
		return c.Send("Введите название сервиса", menu)
	}
}

func onInputServiceRegister(c tele.Context, state fsm.FSMContext) error {
	service := c.Message().Text
	go state.Update("inputService", service)
	go state.Set(InputLoginState)
	return c.Send(fmt.Sprintf("Супер, %s. Теперь введи логин", service))
}

func onInputLogin(c tele.Context, state fsm.FSMContext) error {
	login := c.Message().Text

	go state.Update("login", login)
	go state.Set(InputPasswordState)

	return c.Send("Отлично! Теперь введи пароль")
}

func onInputPassword(confirmBtn, resetBtn, cancelBtn tele.Btn) fsm.Handler {
	m := &tele.ReplyMarkup{}
	m.Inline(
		m.Row(confirmBtn),
		m.Row(resetBtn, cancelBtn),
	)

	return func(c tele.Context, state fsm.FSMContext) error {
		go state.Update("password", c.Message().Text)
		go state.Set(InputConfirmState)
		service := state.MustGet("inputService")
		login := state.MustGet("login")
		c.Delete()
		return c.Send(fmt.Sprintf(
			"Проверьте правильность:\n"+
				"Сервис: %s\n"+
				"Логин: %s\n"+
				"Пароль: %s\n",
			service,
			login,
			c.Message().Text,
		), m)
	}
}

func onInputConfirm(db *sql.DB) fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Finish(true)
		service := state.MustGet("inputService")
		login := state.MustGet("login")
		password := state.MustGet("password")
		formModel := mysql.FormModel{DB: db}
		formModel.Insert(c.Sender().Username, service, login, password)
		return c.Send("Запись сохраненна", tele.RemoveKeyboard)
	}

	//if NoSQL use this
	//data, _ := json.Marshal(map[string]interface{}{
	//	"inputService": service,
	//	"login":        login,
	//	"password":     password,
	//})
	//log.Printf("new form: %s", data)
	//username := "@" + c.Sender().Username + " " // whitespace for formatting

}

func onCancelForm(regBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(regBtn))

	return func(c tele.Context, state fsm.FSMContext) error {
		go state.Finish(true)
		return c.Send("Данные удалены", menu)
	}
}

func onInputResetForm(c tele.Context, state fsm.FSMContext) error {
	go state.Set(InputServiceState)
	c.Send("Хорошо! Начнем сначала.")
	return c.Send("Введите название сервиса")
}

func EditFormMessage(old, new string) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			strOffset := utf8.RuneCountInString(old)
			if nLen := utf8.RuneCountInString(new); nLen > 1 {
				strOffset -= nLen - 1
			}

			entities := make(tele.Entities, len(c.Message().Entities))
			for i, entity := range c.Message().Entities {
				entity.Offset -= strOffset
				entities[i] = entity
			}
			defer func() {
				err := c.EditOrSend(strings.Replace(c.Message().Text, old, new, 1), entities)
				if err != nil {
					c.Bot().OnError(err, c)
				}
			}()
			return next(c)
		}
	}
}

func deleteAfterHandler(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func(c tele.Context) {
			if err := c.Delete(); err != nil {
				c.Bot().OnError(err, c)
			}
		}(c)
		return next(c)
	}
}
