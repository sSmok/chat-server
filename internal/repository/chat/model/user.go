package model

// User модель пользователя в слое репозитория
type User struct {
	ID   int64
	Info UserInfo
}

// UserInfo информация о пользователе в слое репозитория
type UserInfo struct {
	Name string
}
