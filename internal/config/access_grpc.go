package config

import (
	"errors"
	"net"
	"os"
)

const (
	hostEnv = "ACCESS_GRPC_HOST"
	portEnv = "ACCESS_GRPC_PORT"
)

type accessGRPCConfig struct {
	host string
	port string
}

// NewAccessGRPCConfig достает из config-файла данные о GRPC-сервере для проверки прав доступа: название хоста и номер порта
func NewAccessGRPCConfig() (AccessGRPCConfigI, error) {
	host := os.Getenv(hostEnv)
	if len(host) == 0 {
		return nil, errors.New("host env not found")
	}
	port := os.Getenv(portEnv)
	if len(port) == 0 {
		return nil, errors.New("port env not found")
	}

	return &accessGRPCConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *accessGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
