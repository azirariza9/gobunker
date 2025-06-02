package usecase

import (
	"context"
	"gobunker/database"
	"gobunker/model"
	"gobunker/repository"
)

type userUsecase struct {
	userRepo  repository.UserRepository
	txManager database.TxManager
}

type UserUsecase interface {
	GetUserByName(ctx context.Context, name string) (model.User, error)
	GetUserByEmail(ctx context.Context, email string) (model.User, error)
}

func NewUserUsecase(userRepo repository.UserRepository, txManager database.TxManager) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (u *userUsecase) GetUserByName(ctx context.Context, name string) (model.User, error) {
	user, err := u.userRepo.GetUserByName(ctx, name)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (u *userUsecase) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return user, nil

}
