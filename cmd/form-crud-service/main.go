package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Koyo-os/form-crud-service/internal/app"
	"github.com/Koyo-os/form-crud-service/internal/repo"
	"github.com/Koyo-os/form-crud-service/internal/server"
	"github.com/Koyo-os/form-crud-service/internal/service"
	"github.com/Koyo-os/form-crud-service/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getRepo(logger *logger.Logger) (*repo.Repository, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var (
		err error
		db  *gorm.DB
	)

	for range 7 {
		db, err = gorm.Open(mysql.Open(dsn))
		if err == nil {
			break
		} else {
			logger.Error("failed connect to db",
				zap.Error(err),
				zap.String("db_user", os.Getenv("DB_USER")),
				zap.String("db_password", os.Getenv("DB_PASSWORD")),
				zap.String("db_host", os.Getenv("DB_HOST")),
				zap.String("db_port", os.Getenv("DB_PORT")),
				zap.String("db_name", os.Getenv("DB_NAME")))
		}

		time.Sleep(5 * time.Second)
	}

	if err != nil{
		logger.Error("final fail connect to db", zap.Error(err))

		return nil, err
	}

	return repo.NewRepository(db, logger), nil
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	repo, err := getRepo(logger.Get())
	if err != nil{
		return
	}

	service := service.NewService(repo)

	server := server.NewServer(service)

	app := app.NewApp(server)

	if err = app.Start(ctx);err != nil{
		return
	}
}
