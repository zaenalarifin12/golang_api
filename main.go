package main

import (
	"book/config"
	"book/controllers"
	"book/middlewares"
	"book/repository"
	"book/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	db             *gorm.DB                   = config.SetupDatabaseConnection()
	userRepository repository.UserRepository  = repository.NewUserRepository(db)
	bookRepository repository.BookRepository  = repository.NewBookRepository(db)
	authService    services.AuthServices      = services.NewAuthServices(userRepository)
	bookService    services.BookServices      = services.NewBookService(bookRepository)
	jwtService     services.JWTService        = services.NewJWTService()
	authController controllers.AuthController = controllers.NewAuthController(authService, jwtService)
	bookController controllers.BookController = controllers.NewBookController(bookService, jwtService)
)

func main() {

	defer config.CloseDatabaseConnection(db)

	r := gin.Default()

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", authController.Register)
		authGroup.POST("/login", authController.Login)
	}

	bookGroup := r.Group("/api/books", middlewares.AuthorizeJWT(jwtService))
	{
		bookGroup.GET("/", bookController.All)
		bookGroup.POST("/", bookController.Insert)
		bookGroup.GET("/:id", bookController.FindById)
		bookGroup.PUT("/:id", bookController.Update)
		bookGroup.DELETE("/:id", bookController.Delete)

	}

	r.Run(":3000")
}
