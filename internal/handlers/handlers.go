package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"fsm/internal/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
	"strings"
	"unicode/utf8"
)

var (
	InputSG            = fsm.NewStateGroup("reg")
	InputServiceState  = InputSG.New("inputService")
	InputLoginState    = InputSG.New("age")
	InputPasswordState = InputSG.New("hobby")
	InputConfirmState  = InputSG.New("confirm")
)

func InitHandlers(bot *tele.Group, db *sql.DB, manager *fsm.Manager) {
	bot.Handle("/start", onStart)
	manager.Bind("/set", fsm.DefaultState, onStartRegister(keyboards.CancelBtn))
	manager.Bind("/cancel", fsm.AnyState, OnCancelForm(keyboards.SetBtn))

	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.FSMContext) error {
		s := state.State()
		return c.Send(s.String())
	})

	// buttons
	manager.Bind(&keyboards.SetBtn, fsm.DefaultState, onStartRegister(keyboards.CancelBtn))
	manager.Bind(&keyboards.CancelBtn, fsm.AnyState, OnCancelForm(keyboards.SetBtn))

	// form
	manager.Bind(tele.OnText, InputServiceState, onInputServiceRegister)
	manager.Bind(tele.OnText, InputLoginState, onInputLogin)
	manager.Bind(tele.OnText, InputPasswordState, onInputPassword(keyboards.ConfirmBtn, keyboards.ResetFormBtn, keyboards.CancelInlineBtn))
	manager.Bind(&keyboards.ConfirmBtn, InputConfirmState, OnInputConfirm, EditFormMessage("Now check y", "Y"))
	manager.Bind(&keyboards.ResetFormBtn, InputConfirmState, OnInputResetForm, EditFormMessage("Now check your", "Your old"))
	manager.Bind(&keyboards.CancelInlineBtn, InputConfirmState, OnCancelForm(keyboards.SetBtn), DeleteAfterHandler)
}

func onStart(c tele.Context) error {
	log.Println("new user", c.Sender().ID)
	return c.Send(
		"Добро пожаловать в бот для ваших паролей\n"+
			"Отправьте /set чтобы добавить сервис\n"+
			"Отправьте /cancel чтобы омтенить ввод", keyboards.OnStartKB())

}

func onStartRegister(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(cancelBtn))
	menu.ResizeKeyboard = true
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
		senderName := state.MustGet("inputService")
		senderAge := state.MustGet("age")
		c.Delete()
		return c.Send(fmt.Sprintf(
			"Проверьте правильность:\n"+
				"Сервис: %q\n"+
				"Логин: %d\n"+
				"Пароль: %q\n",
			senderName,
			senderAge,
			c.Message().Text,
		), m)
	}
}

func OnInputConfirm(c tele.Context, state fsm.FSMContext) error {
	defer state.Finish(true)

	senderName := state.MustGet("inputService")
	senderAge := state.MustGet("login")
	senderHobby := state.MustGet("password")

	data, _ := json.Marshal(map[string]interface{}{
		"inputService": senderName,
		"login":        senderAge,
		"password":     senderHobby,
	})
	log.Printf("new form: %s", data)
	//username := "@" + c.Sender().Username + " " // whitespace for formatting
	return c.Send("Form accepted", tele.RemoveKeyboard)
}

func OnCancelForm(regBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{}
	menu.Reply(menu.Row(regBtn))
	menu.ResizeKeyboard = true

	return func(c tele.Context, state fsm.FSMContext) error {
		go state.Finish(true)
		return c.Send("Данные удалены", menu)
	}
}

func OnInputResetForm(c tele.Context, state fsm.FSMContext) error {
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

func DeleteAfterHandler(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func(c tele.Context) {
			if err := c.Delete(); err != nil {
				c.Bot().OnError(err, c)
			}
		}(c)
		return next(c)
	}
}
