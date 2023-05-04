package handlers

import (
	"database/sql"
	"fsm/internal/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

var (
	InputGetServiceState = InputSG.New("InputGetServiceState")
)

func initGetHandlers(db *sql.DB, manager *fsm.Manager) {
	manager.Bind("/get", fsm.DefaultState, onStartGet(keyboards.CancelBtn))
	manager.Bind(&keyboards.GetBtn, fsm.DefaultState, onStartGet(keyboards.CancelBtn))

	manager.Bind(tele.OnText, InputGetServiceState, onStartGet(keyboards.CancelBtn))
}

func onStartGet(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancelBtn))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(InputGetServiceState)
		return c.Send("Введите название сервиса", menu)
	}
}
func getRecord() fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Set(fsm.DefaultState)
		username := c.Sender().Username
		service := c.Text()
	}
}
