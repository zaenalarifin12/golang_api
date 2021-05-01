package controllers

import (
	"book/dto"
	"book/entity"
	"book/helper"
	"book/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService services.AuthServices
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthServices, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dto.AuthLoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	authResult := c.authService.VerifyCredential(loginDTO.Email, loginDTO.Password)

	if authResult == false {
		res := helper.BuildErrorResponse("Failed to process request", "credential wrong", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	if u, ok := authResult.(entity.User); ok {
		generatedToken := c.jwtService.GenerateToken(strconv.FormatUint(u.ID, 10))
		u.Token = generatedToken
		response := helper.BuildResponse(true, "OK!", u)
		ctx.JSON(http.StatusOK, response)
	}

}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dto.AuthRegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helper.BuildErrorResponse("Failed To process request", errDTO.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helper.BuildErrorResponse("Failed to process request", "duplicate email", helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)

	} else {
		createdUser := c.authService.CreateUser(registerDTO)
		token := c.jwtService.GenerateToken(strconv.FormatUint(createdUser.ID, 10))
		createdUser.Token = token
		response := helper.BuildResponse(true, "OK!", createdUser)
		ctx.JSON(http.StatusOK, response)
	}

}
