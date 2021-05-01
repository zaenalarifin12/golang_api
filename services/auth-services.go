package services

import (
	"book/dto"
	"book/entity"
	"book/repository"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthServices interface {
	VerifyCredential(email string, password string) interface{}
	CreateUser(user dto.AuthRegisterDTO) entity.User
	//FindByEmail
	IsDuplicateEmail(email string) bool
}

type authServices struct {
	userRepository repository.UserRepository
}

func NewAuthServices(userRep repository.UserRepository) AuthServices {
	return &authServices{
		userRepository: userRep,
	}
}

func  (service *authServices) CreateUser(user dto.AuthRegisterDTO) entity.User  {
	userToCreate := entity.User{}
	err := smapping.FillStruct(&userToCreate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed To Map %v", err)
	}
	res := service.userRepository.InsertUser(userToCreate)
	return res
}

func (service *authServices) VerifyCredential(email string, password string) interface{} {
	res := service.userRepository.VerifyCredential(email, password)

	if u, ok := res.(entity.User); ok {
		comparedPassword := comparedPassword(u.Password, []byte(password))
		if u.Email == email && comparedPassword {
			return res
		}
		return false
	}
	return false
}

func (service *authServices) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}

//

func comparedPassword(hashPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)

	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
