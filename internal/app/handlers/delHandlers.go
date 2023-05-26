package handlers

import (
	"database/sql"
	"fsm/internal/app/service"
	"fsm/internal/keyboards"
	"fsm/pkg/models/mysql"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
)

<<<<<<< HEAD:internal/app/handlers/delHandlers.go
func (h *Handler) initDelHandlers(db *sql.DB, manager *fsm.Manager) {
=======
var InputDeleteServiceState = InputSG.New("InputDeleteServiceState")

func initDelHandlers(db *sql.DB, manager *fsm.Manager) {
>>>>>>> main:internal/handlers/delHandlers.go
	manager.Bind("/del", fsm.DefaultState, onStartDelete(keyboards.CancelBtn))
	manager.Bind(&keyboards.DelBtn, fsm.DefaultState, onStartDelete(keyboards.CancelBtn))

	manager.Bind(tele.OnText, service.InputDeleteServiceState, delRecord(db))
}

func onStartDelete(cancelBtn tele.Btn) fsm.Handler {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(cancelBtn))
	return func(c tele.Context, state fsm.FSMContext) error {
		state.Set(service.InputDeleteServiceState)
		return c.Send("Введите название сервиса", menu)
	}
}

func delRecord(db *sql.DB) fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Set(fsm.DefaultState)
		username := c.Sender().Username
		service := c.Text()
		formModel := mysql.FormModel{DB: db}
		err := formModel.Delete(username, service)
		if err != nil {
			log.Println(err)
			return c.Send("Что-то случилось. Повторите попытку позднее")
		}
		return c.Send("Запись удалена", keyboards.OnStartKB())
	}
}
