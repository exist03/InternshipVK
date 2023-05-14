package handlers

import (
	"database/sql"
	"fmt"
	"fsm/internal/keyboards"
	"fsm/pkg/models/mysql"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
)

var (
	InputGetServiceState = InputSG.New("InputGetServiceState")
)

func initGetHandlers(db *sql.DB, manager *fsm.Manager) {
	manager.Bind("/get", fsm.DefaultState, list(db))
	manager.Bind(&keyboards.GetBtn, fsm.DefaultState, list(db))
	manager.Bind(tele.OnText, InputGetServiceState, getRecord(keyboards.ConfirmBtn, db))
}

func list(db *sql.DB) fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		username := c.Sender().Username
		formModel := mysql.FormModel{DB: db}
		servList, err := formModel.GetList(username)
		if err != nil {
			log.Println(err)
			state.Set(fsm.DefaultState)
			return c.Send("Что-то случилось. Повторите попытку позднее", keyboards.OnStartKB())
		}
		if len(servList) == 0 {
			state.Set(fsm.DefaultState)
			return c.Send("У вас еще нет ни одной записи", keyboards.OnStartKB())
		}
		state.Set(InputGetServiceState)
		return c.Send("Выберите сервис", keyboards.ServersKB(servList))
	}
}
func getRecord(confirmBtn tele.Btn, db *sql.DB) fsm.Handler {
	m := &tele.ReplyMarkup{}
	m.Inline(
		m.Row(confirmBtn),
	)
	return func(c tele.Context, state fsm.FSMContext) error {
		username := c.Sender().Username
		service := c.Text()
		formModel := mysql.FormModel{DB: db}
		login, password := formModel.Get(username, service)
		if login == "" {
			return c.Send("Что-то случилось. Повторите попытку позднее", keyboards.OnStartKB())
		}
		return c.Send(fmt.Sprintf("Логин: %s\nПароль: %s", login, password), m)
	}
}
