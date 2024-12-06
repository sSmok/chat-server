package model

// Message модель сообщения в слое репозитория
type Message struct {
	ID   int64
	Info MessageInfo
}

// MessageInfo информация о сообщении в слое репозитория
type MessageInfo struct {
	ChatID int64
	UserID int64
	Text   string
}
