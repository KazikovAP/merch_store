package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/KazikovAP/merch_store/internal/config"
	"github.com/KazikovAP/merch_store/internal/database"
	"github.com/KazikovAP/merch_store/internal/router"
	"github.com/KazikovAP/merch_store/internal/server"
	"github.com/KazikovAP/merch_store/internal/services"
)

func main() {
	if err := config.Initialize(); err != nil {
		log.Fatal("Не удалось инициализировать конфигурацию:", err)
	}

	serverCfg := config.LoadServerConfig()
	dbCfg := config.LoadDatabaseConfig()
	authCfg := config.LoadAuthConfig()

	db, err := database.SetupDatabase(&dbCfg)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка при закрытии базы данных: %v", err)
		}
	}()

	handler := services.SetupServicesAndHandlers(db, authCfg)

	r := router.SetupRouter(handler, serverCfg, authCfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Println("Получен сигнал для остановки...")
		cancel()
	}()

	log.Println("Приложение успешно запущено!")
	server.StartServer(ctx, serverCfg.Port, r)
}
