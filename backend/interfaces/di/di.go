package di

import (
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app/handler"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	"github.com/7oh2020/connect-tasklist/backend/domain/service"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/model/db"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/sqlc"
	"github.com/7oh2020/connect-tasklist/backend/util/auth"
	"github.com/7oh2020/connect-tasklist/backend/util/clock"
	"github.com/7oh2020/connect-tasklist/backend/util/contextkey"
	"github.com/7oh2020/connect-tasklist/backend/util/identification"
)

func InitUser(qry db.Querier) *handler.UserHandler {
	repo := sqlc.NewSQLCUserRepository(qry)
	srv := service.NewUserService(repo)
	uc := usecase.NewUserUsecase(srv)
	return handler.NewUserHandler(uc)
}

func InitTask(qry db.Querier) *handler.TaskHandler {
	im := identification.NewUUIDManager()
	cm := clock.NewClockManager()
	cr := contextkey.NewContextReader()
	repo := sqlc.NewSQLCTaskRepository(qry)
	srv := service.NewTaskService(repo, im, cm)
	uc := usecase.NewTaskUsecase(srv)
	return handler.NewTaskHandler(uc, cr)
}

func InitAuth(issuer string, keyPath string, qry db.Querier, timeout time.Duration) (*handler.AuthHandler, error) {
	tm, err := auth.NewTokenManager(issuer, keyPath)
	if err != nil {
		return nil, err
	}
	repo := sqlc.NewSQLCUserRepository(qry)
	uc := usecase.NewAuthUsecase(repo, tm, timeout)
	return handler.NewAuthHandler(uc), nil
}
