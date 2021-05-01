package controllers

import (
	"book/dto"
	"book/entity"
	"book/helper"
	"book/services"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

type BookController interface {
	All(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Insert(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bookController struct {
	bookService services.BookServices
	jwtService  services.JWTService
}

func NewBookController(bookServ services.BookServices, jwtService services.JWTService) BookController {
	return &bookController{
		bookService: bookServ,
		jwtService:  jwtService,
	}
}

func (b bookController) All(ctx *gin.Context) {
	var books []entity.Book = b.bookService.All()
	res := helper.BuildResponse(true, "OK!", books)
	ctx.JSON(http.StatusOK, res)
}

func (b bookController) FindById(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		res := helper.BuildErrorResponse("param id not found", err.Error(), helper.EmptyObj{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entity.Book = b.bookService.FindByID(id)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Data not found", "no data with given id", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
	} else {
		res := helper.BuildResponse(true, "OK!", book)
		ctx.JSON(http.StatusOK, res)
	}
}

func (b bookController) Insert(ctx *gin.Context) {
	var dtoBookCreate dto.BookCreateDTO
	errDTO := ctx.ShouldBind(&dtoBookCreate)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := ctx.GetHeader("Authorization")
		userID := b.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			dtoBookCreate.UserID = convertedUserID
		}
		result := b.bookService.Insert(dtoBookCreate)
		response := helper.BuildResponse(true, "OK!", result)
		ctx.JSON(http.StatusOK, response)
	}
}

func (b bookController) Update(ctx *gin.Context) {

	var bookUpdate dto.BookUpdateDTO
	bookUpdate.ID, _ = strconv.ParseUint(ctx.Param("id"), 10, 64)

	//check exist or not
	var book entity.Book = b.bookService.FindByID(bookUpdate.ID)
	if (book == entity.Book{}) {
		res := helper.BuildErrorResponse("Book not found", "Book not found", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	errDTO := ctx.ShouldBind(&bookUpdate)

	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	userID := b.getUserIDByToken(authHeader)
	convertedUserID, err := strconv.ParseUint(userID, 10, 64)

	if b.bookService.IsAllowedToEdit(userID, bookUpdate.ID) {
		if err == nil {
			bookUpdate.UserID = convertedUserID
		}
		result := b.bookService.Update(bookUpdate)
		response := helper.BuildResponse(true, "book updated", result)
		ctx.JSON(http.StatusOK, response)
	} else {
		response := helper.BuildErrorResponse("you dont have permisson", "you are not the owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, response)
	}
}

func (b bookController) Delete(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("failed to get id", "no param id were found", helper.EmptyObj{})
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var oldBook entity.Book = b.bookService.FindByID(id)

	if (oldBook == entity.Book{}) {
		res := helper.BuildErrorResponse("book not found", "id not found", helper.EmptyObj{})
		ctx.JSON(http.StatusNotFound, res)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	userID := b.getUserIDByToken(authHeader)

	if b.bookService.IsAllowedToEdit(userID, oldBook.ID){
		b.bookService.Delete(oldBook)
		res := helper.BuildResponse(true, "Book deleted", helper.EmptyObj{})
		ctx.JSON(http.StatusOK, res)
		return
	}else{
		res := helper.BuildErrorResponse("you dont have permission", "you are not owner", helper.EmptyObj{})
		ctx.JSON(http.StatusForbidden, res)
		return
	}

}

/////
func (b *bookController) getUserIDByToken(header string) string {

	splitAuthHeader := strings.Fields(header)
	token := strings.Join(splitAuthHeader[1:], "")
	aToken, err := b.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
