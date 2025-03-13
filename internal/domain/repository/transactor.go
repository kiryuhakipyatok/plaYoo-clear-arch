package repository

import (
	"context"
	"gorm.io/gorm"
	"log"
)

type Transactor interface {
	WithinTransaction(context.Context, func(c context.Context) (any, error)) (any, error)
}

type transactor struct {
	DB *gorm.DB
}

func NewTransactor(db *gorm.DB) Transactor {
	return &transactor{
		DB: db,
	}
}

type txKey struct{}

func injectTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return nil
}

func (t *transactor) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) (any, error)) (any, error) {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	tx := t.DB.WithContext(txCtx).Begin()

	res, err := tFunc(injectTx(ctx, tx))
	if err != nil {
		if err := tx.Rollback().Error; err != nil {
			log.Printf("cannot rollback transaction: %v", err)
		}

		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("cannot commit transaction: %v", err)
	}

	return res, nil
}

func (t *transactor) tx(ctx context.Context) *gorm.DB {
	if tx := extractTx(ctx); tx != nil {
		return tx.WithContext(ctx)
	}

	return t.DB.WithContext(ctx)
}
