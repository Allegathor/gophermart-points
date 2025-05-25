package srv

import (
	"gophermart-points/internal/repo/pgsql"
	"gophermart-points/internal/srv/external"
	"gophermart-points/internal/srv/handlers"
	"gophermart-points/internal/srv/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Srv struct {
	db     *pgsql.PgSQL
	router *gin.Engine
	logger *zap.SugaredLogger
	http.Server
}

func New(addr string, db *pgsql.PgSQL, l *zap.SugaredLogger) *Srv {

	return &Srv{
		db:     db,
		router: gin.Default(),
		logger: l,
		Server: http.Server{
			Addr: addr,
		},
	}
}

func (s *Srv) MountHandlers(authKey string, q *external.OrderProcessing) {
	api := handlers.NewAPI(s.db, s.logger)
	orderAPI := handlers.NewOrderAPI(s.db, q)
	userAPI := handlers.NewUserAPI(authKey, s.db)

	apiGr := s.router.Group("/api")
	{
		userGr := apiGr.Group("/user")
		{
			userGr.POST("/register", middlewares.RestrictJSON(), userAPI.UserReg)
			userGr.POST("/login", middlewares.RestrictJSON(), userAPI.UserLogin)

			authGr := userGr.Group("")
			authGr.Use(middlewares.CheckAuth(authKey))
			authGr.POST("/orders", middlewares.RestrictText(), orderAPI.RegOrder)
			authGr.GET("/orders", orderAPI.Orders)
			authGr.GET("/balance", api.GetBalance)
			authGr.POST("/balance/withdraw", middlewares.RestrictJSON(), api.Withdraw)
			authGr.GET("/withdrawals", api.Withdrawals)
		}
	}

	s.Handler = s.router.Handler()
}
