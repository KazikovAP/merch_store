package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

func StartServer(ctx context.Context, port string, handler http.Handler) {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Сервер запущен на порту :" + port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Останавливаем сервер...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	if err := srv.Shutdown(shutdownCtx); err != nil {
		cancel()
		log.Fatalf("Ошибка при остановке сервера: %v", err)
	}

	cancel()

	log.Println("Сервер успешно остановлен.")
}
