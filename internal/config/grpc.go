package config

import (
	"errors"
	"net"
	"os"
)

const (
	hostEnv = "GRPC_HOST"
	portEnv = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig достает из config-файла данные о GRPC-сервере: название хоста и номер порта
func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(hostEnv)
	if len(host) == 0 {
		return nil, errors.New("grpc host env not found")
	}
	port := os.Getenv(portEnv)
	if len(port) == 0 {
		return nil, errors.New("grpc port env not found")
	}

	cfg := &grpcConfig{
		host: host,
		port: port,
	}

	return cfg, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
