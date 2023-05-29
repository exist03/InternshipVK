package service

import (
	"database/sql"
	"fmt"
	"fsm/internal/app/repository"
	"fsm/internal/pkg/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

var (
	InputSG                 = fsm.NewStateGroup("reg")
	InputServiceState       = InputSG.New("inputService")
	InputLoginState         = InputSG.New("login")
	InputPasswordState      = InputSG.New("password")
	InputConfirmState       = InputSG.New("confirm")
	InputDeleteServiceState = InputSG.New("InputDeleteServiceState")
	InputGetServiceState    = InputSG.New("InputGetServiceState")
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) OnStartRegister(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancelBtn))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(InputServiceState)
		return c.Send("Введите название сервиса", menu)
	}
}

func (s *Service) OnInputServiceRegister(c tele.Context, state fsm.FSMContext) error {
	service := c.Message().Text
	go state.Update("inputService", service)
	go state.Set(InputLoginState)
	return c.Send(fmt.Sprintf("Супер. Теперь введи логин"))
}

func (s *Service) OnInputLogin(c tele.Context, state fsm.FSMContext) error {
	login := c.Message().Text

	go state.Update("login", login)
	go state.Set(InputPasswordState)

	return c.Send("Отлично! Теперь введи пароль")
}

func (s *Service) OnInputPassword(confirmBtn, resetBtn, cancelBtn tele.Btn) fsm.Handler {
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

func (s *Service) OnInputConfirm(db *sql.DB) fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Finish(true)
		service := state.MustGet("inputService")
		login := state.MustGet("login")
		password := state.MustGet("password")
		repo := repository.New(db)
		err := repo.Insert(c.Sender().Username, service, login, password)
		if err != nil {
			return err
		}
		return c.Send("Запись сохраненна", keyboards.OnStartKB())
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

func (s *Service) OnCancelForm() fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		go state.Finish(true)
		return c.Send("Данные удалены", keyboards.OnStartKB())
	}
}

func (s *Service) OnInputResetForm(c tele.Context, state fsm.FSMContext) error {
	go state.Set(InputServiceState)
	c.Send("Хорошо! Начнем сначала.")
	return c.Send("Введите название сервиса")
}

func (s *Service) DeleteAfterHandler(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func(c tele.Context) {
			if err := c.Delete(); err != nil {
				c.Bot().OnError(err, c)
			}
		}(c)
		return next(c)
	}
}
