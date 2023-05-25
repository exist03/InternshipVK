package handlers

import (
	"database/sql"
	"fsm/internal/app/repository"
	"fsm/internal/app/service"
	"fsm/internal/keyboards"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	tele "gopkg.in/telebot.v3"
	"log"
)

type Handler struct {
	serv *service.Service
	repo *repository.Repository
}

func New(serv *service.Service, repo *repository.Repository) *Handler {
	return &Handler{serv: serv,
		repo: repo}
}
func (h *Handler) Init(bot *tele.Group, db *sql.DB, manager *fsm.Manager) {
	h.initDelHandlers(db, manager)
	h.initGetHandlers(db, manager)
	bot.Handle("/start", onStart)
	manager.Bind("/set", fsm.DefaultState, h.serv.OnStartRegister(keyboards.CancelBtn))
	manager.Bind("/cancel", fsm.AnyState, h.serv.OnCancelForm())

	manager.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.FSMContext) error {
		s := state.State()
		return c.Send(s.String())
	})

	// buttons
	manager.Bind(&keyboards.SetBtn, fsm.DefaultState, h.serv.OnStartRegister(keyboards.CancelBtn))
	manager.Bind(&keyboards.CancelBtn, fsm.AnyState, h.serv.OnCancelForm())

	// form
	manager.Bind(tele.OnText, service.InputServiceState, h.serv.OnInputServiceRegister)
	manager.Bind(tele.OnText, service.InputLoginState, h.serv.OnInputLogin)
	manager.Bind(tele.OnText, service.InputPasswordState, h.serv.OnInputPassword(keyboards.ConfirmBtn, keyboards.ResetFormBtn, keyboards.CancelInlineBtn))
	manager.Bind(&keyboards.ConfirmBtn, service.InputConfirmState, h.serv.OnInputConfirm(db), h.serv.DeleteAfterHandler)
	manager.Bind(&keyboards.ResetFormBtn, service.InputConfirmState, h.serv.OnInputResetForm, h.serv.DeleteAfterHandler)
	manager.Bind(&keyboards.CancelInlineBtn, service.InputConfirmState, h.serv.OnCancelForm(), h.serv.DeleteAfterHandler)
}

func onStart(c tele.Context) error {
	log.Println("new user", c.Sender().ID)
	return c.Send(
		"Добро пожаловать в бот для ваших паролей\n"+
			"Отправьте /set чтобы добавить сервис\n"+
			"Отправьте /cancel чтобы омтенить действие\n"+
			"Отправьте /get чтобы получить запись", keyboards.OnStartKB())

}
