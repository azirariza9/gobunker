package usecase

import (
	"context"
	"fmt"
	"gobunker/utils/service"
)

type authenticationUsecase struct {
	userUseCase UserUsecase
	jwtService  service.JwtService
}

type AuthenticationUsecase interface {
	LoginHandler(ctx context.Context, email string, password string) (string, error)
}

func NewAuthenticationUsecase(uc UserUsecase, jwtService service.JwtService) AuthenticationUsecase {
	return &authenticationUsecase{userUseCase: uc, jwtService: jwtService}
}

func (a *authenticationUsecase) LoginHandler(ctx context.Context, email string, password string) (string, error) {

	user, err := a.userUseCase.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = service.ComparePassword(user.Password, password)
	if err != nil {
		return "", fmt.Errorf("err : usecase : password salah")
	}

	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
