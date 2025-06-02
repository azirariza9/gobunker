package database

import (
	"context"
	"database/sql"
	"fmt"
	"gobunker/repository"
)

type TxManager struct {
	exec repository.SQLExecutor
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{exec: db}
}

func (tm *TxManager) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	db, ok := tm.exec.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("database : cannot start transaction from a transaction")
	}

	tx, err := db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("database : %w", err)
	}
	return tx, nil
}
