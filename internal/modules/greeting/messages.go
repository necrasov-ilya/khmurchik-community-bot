package greeting

const DefaultMessage = "Я робот, и Паша заставил меня каждое утро желать вам охуенного дня!\n\nОбнял, покружил, на место поставил! 😘"

type MessageTemplate struct {
	ID      int64
	Text    string
	Enabled bool
	Weight  int
}
