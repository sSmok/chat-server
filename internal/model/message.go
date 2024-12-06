package model

// Message модель сообщения в сервисном слое
type Message struct {
	ID   int64
	Info MessageInfo
}

// MessageInfo информация о сообщении в сервисном слое
type MessageInfo struct {
	ChatID int64
	UserID int64
	Text   string
}
