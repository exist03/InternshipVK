package handlers

import (
	"fsm/internal/app/repository"
	"fsm/internal/app/service"
	"fsm/internal/pkg/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
)

func (h *Handler) initDelHandlers(manager *fsm.Manager) {
	manager.Bind("/del", fsm.DefaultState, servList(h.repo, "del"))
	manager.Bind(&keyboards.DelBtn, fsm.DefaultState, servList(h.repo, "del"))

	manager.Bind(tele.OnText, service.InputDeleteServiceState, delRecord(h.repo))
}

func delRecord(repo *repository.Repository) fsm.Handler {
	return func(c tele.Context, state fsm.FSMContext) error {
		defer state.Set(fsm.DefaultState)
		username := c.Sender().Username
		serv := c.Text()
		err := repo.Delete(username, serv)
		if err != nil {
			log.Println(err)
			return c.Send("Что-то случилось. Повторите попытку позднее")
		}
		return c.Send("Запись удалена", keyboards.OnStartKB())
	}
}
