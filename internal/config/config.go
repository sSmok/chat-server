package config

// AccessGRPCConfigI - интерфейс конфигурации GRPC-клиента для проверки прав доступа пользователя
type AccessGRPCConfigI interface {
	Address() string
}
