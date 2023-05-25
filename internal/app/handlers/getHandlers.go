package handlers

import (
	"database/sql"
	"fmt"
	"fsm/internal/app/repository"
	"fsm/internal/app/service"
	"fsm/internal/keyboards"
	"fsm/pkg/models/mysql"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
)

func (h *Handler) initGetHandlers(db *sql.DB, manager *fsm.Manager) {
	manager.Bind("/get", fsm.DefaultState, onStartGet(h.repo))
	manager.Bind(&keyboards.GetBtn, fsm.DefaultState, onStartGet(h.repo))
	manager.Bind(tele.OnText, service.InputGetServiceState, getRecord(keyboards.ConfirmBtn, db))
}

func onStartGet(repo *repository.Repository) fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		username := c.Sender().Username
		servList, err := repo.GetList(username)
		if err != nil {
			log.Println(err)
			state.Set(fsm.DefaultState)
			return c.Send("Что-то случилось. Повторите попытку позднее", keyboards.OnStartKB())
		}
		if len(servList) == 0 {
			state.Set(fsm.DefaultState)
			return c.Send("У вас еще нет ни одной записи", keyboards.OnStartKB())
		}
		state.Set(service.InputGetServiceState)
		return c.Send("Выберите сервис", keyboards.ServersKB(servList))
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
