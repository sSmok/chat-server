package model

// User модель пользователя в сервисном слое
type User struct {
	ID   int64
	Info UserInfo
}

// UserInfo информация о пользователе в сервисном слое
type UserInfo struct {
	Name string
}
