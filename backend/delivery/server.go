package delivery

import (
	"fmt"
	"gobunker/config"
	"gobunker/database"
	"gobunker/delivery/controller"
	"gobunker/middleware"
	"gobunker/repository"
	"gobunker/usecase"
	"gobunker/utils/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host       string
	engine     *gin.Engine
	jwtService service.JwtService
	authUC     usecase.AuthenticationUsecase
}

func NewServer() *Server {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(fmt.Errorf("failed to calling config: %v", err))
	}

	db, err := config.ConnectDB(cfg)
	if err != nil {
		panic(fmt.Errorf("failed to connect db : %v", err))
	}

	jwtService := service.NewJwtService(cfg.TokenConfig)

	txManager := database.NewTxManager(db)
	userRepo := repository.NewUserRepository(db)

	userUsecase := usecase.NewUserUsecase(userRepo, *txManager)
	authUsecase := usecase.NewAuthenticationUsecase(userUsecase, jwtService)

	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())

	host := fmt.Sprintf(":%s", cfg.ApiPort)
	return &Server{
		host:       host,
		engine:     engine,
		jwtService: jwtService,
		authUC:     authUsecase,
	}
}

func (s *Server) initRoute() {
	rg := s.engine.Group("/api")
	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)
	controller.NewAuthController(s.authUC, rg).Route()
	controller.NewChatController(rg, authMiddleware).Route()
}

func (s *Server) Run() {
	s.initRoute()
	s.engine.SetTrustedProxies([]string{"127.0.0.1"})

	err := s.engine.Run(s.host)
	if err != nil {
		panic(fmt.Errorf("failed to start server on host %s: %v", s.host, err))
	}
}
