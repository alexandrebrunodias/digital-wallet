package uow

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type Repository func(tx *sql.Tx) interface{}

type UnitOfWorkInterface interface {
	Add(name string, repository Repository)
	Remove(name string)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(unitOfWork *UnitOfWork) error) error
	CommitOrRollback() error
	RollBack() error
}

type UnitOfWork struct {
	Tx           *sql.Tx
	Db           *sql.DB
	repositories map[string]Repository
}

func NewUnitOfWork(db *sql.DB) *UnitOfWork {
	return &UnitOfWork{
		Db:           db,
		repositories: make(map[string]Repository),
	}
}

func (u *UnitOfWork) Add(name string, repository Repository) {
	u.repositories[name] = repository
}

func (u *UnitOfWork) Remove(name string) {
	delete(u.repositories, name)
}

func (u *UnitOfWork) GetRepository(ctx context.Context, name string) (interface{}, error) {
	var err error
	if u.Tx == nil {
		u.Tx, err = u.Db.BeginTx(ctx, nil)
		if err != nil {
			return nil, err
		}
	}
	return u.repositories[name](u.Tx), nil
}

func (u *UnitOfWork) Do(ctx context.Context, fn func(unitOfWork *UnitOfWork) error) error {
	if u.Tx != nil {
		return errors.New("transaction is already started")
	}

	var err error
	u.Tx, err = u.Db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(u)
	if err != nil {
		errRollBack := u.RollBack()
		if errRollBack != nil {
			return errors.New(fmt.Sprintf(
				"transaction error: %s | rollback error: %s",
				err.Error(),
				errRollBack.Error()),
			)
		}
		return err
	}
	return u.CommitOrRollback()
}

func (u *UnitOfWork) CommitOrRollback() error {
	err := u.Tx.Commit()
	if err != nil {
		errRollBack := u.Tx.Rollback()
		if errRollBack != nil {
			return errors.New(fmt.Sprintf(
				"transaction error: %s | rollback error: %s",
				err.Error(),
				errRollBack.Error()),
			)
		}
		return err
	}
	u.Tx = nil
	return nil
}

func (u *UnitOfWork) RollBack() error {
	if u.Tx != nil {
		err := u.Tx.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}
