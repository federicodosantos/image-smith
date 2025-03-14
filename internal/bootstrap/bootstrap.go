package bootstrap

import (
	"log"
	"net/http"
	"os"

	"github.com/federicodosantos/image-smith/internal/delivery"
	"github.com/federicodosantos/image-smith/internal/repository"
	"github.com/federicodosantos/image-smith/internal/usecase"
	"github.com/federicodosantos/image-smith/pkg/jwt"
	"github.com/federicodosantos/image-smith/pkg/util"
	"github.com/jmoiron/sqlx"
)

type Bootstrap struct {
	db     *sqlx.DB
	router *http.ServeMux
}

func NewBootstrap(db *sqlx.DB, router *http.ServeMux) *Bootstrap {
	return &Bootstrap{
		db:     db,
		router: router,
	}
}

func (b *Bootstrap) InitApp() {
	// initialize jwt service
	jwtService, err := jwt.NewJwt(os.Getenv("JWT_SECRET_KEY"), os.Getenv("JWT_EXPIRED"))
	if err != nil {
		log.Printf("cannot initialize jwt service due to %s", err.Error())
	}

	//initialize repositories
	userRepo := repository.NewUserRepository(b.db)

	//initialize usecases
	userUsecase := usecase.NewUserUsecase(userRepo, jwtService)

	//initialize handlers
	userHandler := delivery.NewUserHandler(userUsecase)

	//initialize routes
	delivery.UserRoutes(b.router, userHandler)

	util.HealthCheck(b.router, b.db)
}
