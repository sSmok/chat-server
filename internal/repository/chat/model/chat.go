package model

// Chat модель чата в слое репозитория
type Chat struct {
	ID   int64
	Info ChatInfo
}

// ChatInfo информация о чате в слое репозитория
type ChatInfo struct {
	Name    string
	UserIDs []int64
}
