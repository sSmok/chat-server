package model

// Chat модель чата в сервисном слое
type Chat struct {
	ID   int64
	Info ChatInfo
}

// ChatInfo информация о чате в сервисном слое
type ChatInfo struct {
	Name    string
	UserIDs []int64
}
