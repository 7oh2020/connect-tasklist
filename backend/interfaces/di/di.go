package di

import (
	"time"

	"github.com/7oh2020/connect-tasklist/backend/app/handler"
	"github.com/7oh2020/connect-tasklist/backend/app/usecase"
	"github.com/7oh2020/connect-tasklist/backend/app/util/auth"
	"github.com/7oh2020/connect-tasklist/backend/app/util/clock"
	"github.com/7oh2020/connect-tasklist/backend/app/util/contextkey"
	"github.com/7oh2020/connect-tasklist/backend/app/util/identification"
	"github.com/7oh2020/connect-tasklist/backend/domain/service"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/model/db"
	"github.com/7oh2020/connect-tasklist/backend/infrastructure/persistence/sqlc"
)

func InitUser(issuer string, keyPath string, qry db.Querier, duration time.Duration) (*handler.UserHandler, error) {
	tm, err := auth.NewTokenManager(issuer, keyPath)
	if err != nil {
		return nil, err
	}
	rpu := sqlc.NewSQLCUserRepository(qry)
	svu := service.NewUserService(rpu)
	uca := usecase.NewAuthUsecase(rpu, tm, duration)
	ucu := usecase.NewUserUsecase(svu)
	hdr := handler.NewUserHandler(uca, ucu)
	return hdr, nil
}

func InitTask(qry db.Querier) (*handler.TaskHandler, error) {
	im := identification.NewUUIDManager()
	cm := clock.NewClockManager()
	cr := contextkey.NewContextReader()
	rpt := sqlc.NewSQLCTaskRepository(qry)
	svt := service.NewTaskService(rpt)
	uct := usecase.NewTaskUsecase(svt)
	hdr := handler.NewTaskHandler(im, cm, cr, uct)
	return hdr, nil
}
