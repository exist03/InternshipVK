package keyboards

import tele "gopkg.in/telebot.v3"

var (
	SetBtn          = tele.Btn{Text: "📝 Создать запись"}
	GetBtn          = tele.Btn{Text: "📝 Получить запись"}
	CancelBtn       = tele.Btn{Text: "❌ Закрыть"}
	DelBtn          = tele.Btn{Text: "❌ Удалить запись"}
	ConfirmBtn      = tele.Btn{Text: "✅ Подтвердить", Unique: "confirm"}
	ResetFormBtn    = tele.Btn{Text: "🔄 Обновить запись", Unique: "reset"}
	CancelInlineBtn = tele.Btn{Text: "❌ Закрыть", Unique: "cancel"}
)

func OnStartKB() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(menu.Row(SetBtn, GetBtn),
		menu.Row(DelBtn))
	return menu
}
