package handlers

import (
	"fmt"
	"fsm/internal/app/repository"
	"fsm/internal/app/service"
	"fsm/internal/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) initGetHandlers(manager *fsm.Manager) {
	manager.Bind("/get", fsm.DefaultState, servList(h.repo, "get"))
	manager.Bind(&keyboards.GetBtn, fsm.DefaultState, servList(h.repo, "get"))
	manager.Bind(tele.OnText, service.InputGetServiceState, getRecord(keyboards.ConfirmBtn, h.repo))
}

func getRecord(confirmBtn tele.Btn, repo *repository.Repository) fsm.Handler {
	m := &tele.ReplyMarkup{}
	m.Inline(
		m.Row(confirmBtn),
	)
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Set(fsm.DefaultState)
		username := c.Sender().Username
		serv := c.Text()
		login, password := repo.Get(username, serv)
		if login == "" {
			return c.Send("Что-то случилось. Повторите попытку позднее")
		}
		return c.Send(fmt.Sprintf("Логин: %s\nПароль: %s", login, password), m)
	}
}
