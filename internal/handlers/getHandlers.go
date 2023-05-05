package handlers

import (
	"database/sql"
	"fmt"
	"fsm/internal/keyboards"
	"fsm/pkg/models/mysql"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

var (
	InputGetServiceState   = InputSG.New("InputGetServiceState")
	ConfirmGetServiceState = InputSG.New("ConfirmGetServiceState")
)

func initGetHandlers(db *sql.DB, manager *fsm.Manager) {
	manager.Bind("/get", fsm.DefaultState, onStartGet(keyboards.CancelBtn))
	manager.Bind(&keyboards.GetBtn, fsm.DefaultState, onStartGet(keyboards.CancelBtn))

	manager.Bind(tele.OnText, InputGetServiceState, getRecord(keyboards.ConfirmBtn, db))
}

func onStartGet(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancelBtn))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(InputGetServiceState)
		return c.Send("Введите название сервиса", menu)
	}
}
func getRecord(confirmBtn tele.Btn, db *sql.DB) fsm.Handler {
	m := &tele.ReplyMarkup{}
	m.Inline(
		m.Row(confirmBtn),
	)
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Set(fsm.DefaultState)
		username := c.Sender().Username
		service := c.Text()
		formModel := mysql.FormModel{DB: db}
		login, password := formModel.Get(username, service)
		if login == "" {
			return c.Send("Что-то случилось. Повторите попытку позднее")
		}
		return c.Send(fmt.Sprintf("Логин: %s\nПароль: %s", login, password), m)
	}
}

//func onInputPassword(confirmBtn, resetBtn, cancelBtn tele.Btn) fsm.Handler {
//	m := &tele.ReplyMarkup{}
//	m.Inline(
//		m.Row(confirmBtn),
//		m.Row(resetBtn, cancelBtn),
//	)
//
//	return func(c tele.Context, state fsm.FSMContext) error {
//		go state.Update("password", c.Message().Text)
//		go state.Set(InputConfirmState)
//		service := state.MustGet("inputService")
//		login := state.MustGet("login")
//		c.Delete()
//		return c.Send(fmt.Sprintf(
//			"Проверьте правильность:\n"+
//				"Сервис: %s\n"+
//				"Логин: %s\n"+
//				"Пароль: %s\n",
//			service,
//			login,
//			c.Message().Text,
//		), m)
//	}
//}
