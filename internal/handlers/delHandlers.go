package handlers

import (
	"database/sql"
	"fsm/internal/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

var (
	InputDeleteServiceState = InputSG.New("InputDeleteServiceState")
)

func initDelHandlers(db *sql.DB, manager *fsm.Manager) {
	manager.Bind("/del", fsm.DefaultState, onStartDelete(keyboards.CancelBtn))
	manager.Bind(&keyboards.DelBtn, fsm.DefaultState, onStartDelete(keyboards.CancelBtn))

	manager.Bind(tele.OnText, InputDeleteServiceState, onStartGet(keyboards.CancelBtn))
}

func onStartDelete(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancelBtn))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(InputDeleteServiceState)
		return c.Send("Введите название сервиса", menu)
	}
}

func delRecord() fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Set(fsm.DefaultState)
		username := c.Sender().Username
		service := c.Text()
	}
}
