package main

import (
	"context"
	"learn_bot_admin_panel/config"
	"learn_bot_admin_panel/external/bot_api"
	"learn_bot_admin_panel/external/rest_api"
	"learn_bot_admin_panel/external/rest_api/middleware"
	"learn_bot_admin_panel/internal/chanel_bus"
	"learn_bot_admin_panel/internal/entity/profile"
	"learn_bot_admin_panel/internal/transaction"
	"learn_bot_admin_panel/rimport"
	"learn_bot_admin_panel/tools/logger"
	"learn_bot_admin_panel/uimport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	config := config.NewConfig()

	// подключение к redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})
	defer rdb.Close()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Panic("ошибка при пинге redis: ", err)
	}

	// подключение к postgres
	pgdb, err := sqlx.Connect("pgx", config.PostgresURL)
	if err != nil {
		log.Fatalln("не удалось подключиться к базе postgres: ", err)
	}
	defer pgdb.Close()

	if err := pgdb.Ping(); err != nil {
		log.Fatal("ошибка при пинге postgres : ", err)
	}

	// инициализация логгера
	hook := logger.NewPostgresHook(pgdb)
	logger, err := logger.InitLogger(hook)
	if err != nil {
		log.Fatalln("Не удалось инициализировать логгер:", err)
	}

	// инициализация session manager
	sessionManager := transaction.NewSQLSessionManager(pgdb)

	// инициализация бота
	b, err := bot.New(config.BotToken)
	if err != nil {
		log.Panic("ошибка при чтении токена бота: ", err)
	}
	defer b.Close(ctx)

	// Настройка HTTP-сервера
	ginConfig := cors.Config{
		AllowOrigins:     []string{"http://admin-panel.local:3000", "https://admin-panel.local:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Device-ID"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	router := gin.Default()
	router.Use(cors.New(ginConfig))

	v1Router := router.Group("/api/v1")
	v1Router.Use(cors.New(ginConfig))
	srv := &http.Server{
		Addr:    config.ServerPort,
		Handler: router,
	}

	// инициализация event bus авторизация
	authChan := chanel_bus.NewBusChanel[profile.User]()
	twoStepAuthChan := chanel_bus.NewBusChanel[profile.PasswordLoginResponse]()
	// инициализация репо
	ri := rimport.NewRepositoryImports(config, rdb)

	// инициализация usecase
	ui := uimport.NewUsecaseImport(ri, logger, authChan, twoStepAuthChan, config, b)

	//
	middleware := middleware.NewAuthMiddleware(ui.Jwt)

	// инициализация rest handler
	rest_api.NewProfileHandler(ui, v1Router, config, logger, middleware, sessionManager)
	rest_api.NewBotUsersHandler(ui, v1Router, config, logger, middleware, sessionManager)

	// инициализация tg bot handler
	bot_api.NewBotProfileHandler(ui, b, config, logger, sessionManager)

	wg.Add(2)
	// запуск gin
	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panic("Ошибка при запуске HTTP-сервера: ", err)
		}
	}()

	// запуск api бота
	go func() {
		defer wg.Done()

		b.Start(ctx)
	}()

	// Ожидание сигнала завершения
	<-ctx.Done()
	log.Println("Останавливаем приложение...")

	// Даём время на завершение запросов
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Остановка HTTP-сервера
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Println("Ошибка при остановке HTTP-сервера:", err)
	}

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Println("Приложение остановлено")
}
