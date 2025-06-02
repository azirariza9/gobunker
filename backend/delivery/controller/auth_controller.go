package controller

import (
	"gobunker/model/dto"
	"gobunker/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authController struct {
	authUC usecase.AuthenticationUsecase
	rg     *gin.RouterGroup
}

func (a *authController) Route() {
	a.rg.POST("/login", a.loginHandler)
}

func NewAuthController(authUc usecase.AuthenticationUsecase, rg *gin.RouterGroup) *authController {
	return &authController{authUC: authUc, rg: rg}
}

func (a *authController) loginHandler(c *gin.Context) {
	var payload dto.UserDTO

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := a.authUC.LoginHandler(c.Request.Context(), payload.Email, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "Login Success",
		Token:   token,
	})
}
