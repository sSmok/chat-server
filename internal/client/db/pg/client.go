package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/client/db"
)

type pgClient struct {
	masterDB db.DB
}

// NewPGClient - конструктор клиента для работы с PostgreSQL
func NewPGClient(ctx context.Context, dsn string) (db.ClientI, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &pgClient{masterDB: NewPG(pool)}, nil
}

func (client *pgClient) DB() db.DB {
	return client.masterDB
}

func (client *pgClient) Close() error {
	if client.masterDB != nil {
		client.masterDB.Close()
	}

	return nil
}
