package api

import (
	"fmt"
	"net/http"

	db "github.com/Mersock/golang-sample-bank/db/sqlc"
	"github.com/Mersock/golang-sample-bank/token"
	"github.com/Mersock/golang-sample-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	//regis validator currency
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRounter()
	return server, nil
}

func (server *Server) setUpRounter() {
	router := gin.Default()
	router.GET("/ping", pingRes)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	server.router = router

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func pingRes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func errRes(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
