package main

import (
	accountcontroller "bank/internal/account/controller"
	accountrepository "bank/internal/account/repository"
	accountservice "bank/internal/account/service"
	authcontroller "bank/internal/auth/controller"
	authrepository "bank/internal/auth/repository"
	authservice "bank/internal/auth/service"
	cardcontroller "bank/internal/card/controller"
	cardrepository "bank/internal/card/repository"
	cardservice "bank/internal/card/service"
	"bank/internal/config"
	creditcontroller "bank/internal/credit/controller"
	creditservice "bank/internal/credit/service"
	"bank/internal/server"
	"bank/internal/storage"
	transactioncontroller "bank/internal/transaction/controller"
	transactionrepository "bank/internal/transaction/repository"
	transactionservice "bank/internal/transaction/service"
	"bank/pkg/errors"
	"bank/pkg/log"
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(dsn string, log log.Logger) {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		log.Error("migration error", "error", err)
		os.Exit(1)
	}

	err = m.Up()

	switch {
	case errors.Is(err, migrate.ErrNoChange):
		log.Info("migration not needed, schema in actual state")
		return

	case err != nil:
		log.Error("migration failed", "error", err)
		os.Exit(1)
	}

	log.Info("migration success")
}

func main() {
	configPath := flag.String("config", "", "configuration file. Can be YAML or JSON file")
	flag.Parse()

	// Загрузка конфига и его обогащение данными из переменных окружения.
	cfg := start(config.NewConfigFromFile(*configPath))
	cfg = config.EnvEnrichment(cfg)

	logger := start(log.New(cfg.Logger))

	// Менеджер сервисов.
	ctx := context.Background()
	manager := NewManager(&ctx, logger)

	dsn := storage.BuildDSN(&cfg.Storage.Postgres)
	runMigrations(dsn, logger)

	logger.Info("server starting", "build", manager.build)
	defer logger.Info("server stopped")

	storageLog := logger.WithGroup("storage")

	// Инициализация хранилища.
	st := start(storage.NewStorage(
		storage.WithConfig(cfg.Storage.Postgres),
		storage.WithLogger(storageLog),
	))

	// Отдельная группа логгеров для серверов
	serverLog := logger.WithGroup("http")

	// Родительский логгер для подсистем внутри сервиса auth.
	authlog := serverLog.WithGroup("auth")

	// Инициализация репозитория auth.
	authRepo := start(authrepository.NewRepository(
		authrepository.WithStorage(st),
		authrepository.WithLogger(authlog.WithGroup("repository")),
	))

	// Инициализация сервиса auth.
	authService := start(authservice.NewAuthService(
		authservice.WithAuthLogger(authlog.WithGroup("service")),
		authservice.WithAuthRepository(authRepo),
	))

	// Инициализация контроллера auth.
	authController := start(authcontroller.NewHandler(
		authcontroller.WithLogger(authlog.WithGroup("controller")),
		authcontroller.WithService(authService),
	))

	// Родительский логгер для подсистем внутри сервиса account.
	accountLog := serverLog.WithGroup("account")

	// Инициализация репозитория Account.
	accountRepo := start(accountrepository.NewRepository(
		accountrepository.WithStorage(st),
		accountrepository.WithLogger(accountLog.WithGroup("repository")),
	))

	// Инициализация сервиса Account.
	accountService := start(accountservice.NewService(
		accountservice.WithLogger(accountLog.WithGroup("service")),
		accountservice.WithRepository(accountRepo),
		accountservice.WithAuthService(authService),
	))

	// Инициализация контроллера Account.
	accountController := start(accountcontroller.NewHandler(
		accountcontroller.WithLogger(accountLog.WithGroup("controller")),
		accountcontroller.WithService(accountService),
	))

	// Родительский логгер для подсистем внутри сервиса card.
	cardLog := serverLog.WithGroup("card")

	// Инициализация репозитория Card.
	cardRepo := start(cardrepository.NewRepository(
		cardrepository.WithStorage(st),
		cardrepository.WithPGPConfig(cfg.Crypt.PGP),
		cardrepository.WithLogger(cardLog.WithGroup("repository")),
	))

	// Инициализация сервиса Card.
	cardService := start(cardservice.NewService(
		cardservice.WithLogger(cardLog.WithGroup("service")),
		cardservice.WithRepository(cardRepo),
		cardservice.WithAuthService(authService),
		cardservice.WithAccountService(accountService),
	))

	// Инициализация контроллера Card.
	cardController := start(cardcontroller.NewHandler(
		cardcontroller.WithLogger(cardLog.WithGroup("controller")),
		cardcontroller.WithService(cardService),
	))

	// Родительский логгер для подсистем внутри сервиса credit.
	creditLog := serverLog.WithGroup("credit")

	// Инициализация сервиса Credit.
	creditService := start(creditservice.NewService(
		creditservice.WithLogger(creditLog.WithGroup("service")),
	))

	// Инициализация контроллера Credit.
	controllers := start(creditcontroller.NewHandler(
		creditcontroller.WithLogger(creditLog.WithGroup("controller")),
		creditcontroller.WithService(creditService),
	))

	// Родительский логгер для подсистем внутри сервиса transaction.
	transactionLog := serverLog.WithGroup("card")

	// Инициализация репозитория Draw.
	transactionRepo := start(transactionrepository.NewRepository(
		transactionrepository.WithStorage(st),
		transactionrepository.WithLogger(transactionLog.WithGroup("repository")),
	))

	// Инициализация сервиса Draw.
	transactionService := start(transactionservice.NewService(
		transactionservice.WithLogger(transactionLog.WithGroup("service")),
		transactionservice.WithRepository(transactionRepo),
		transactionservice.WithAccountService(accountService),
		transactionservice.WithCardService(cardService),
	))

	// Инициализация контроллера Draw.
	transactionController := start(transactioncontroller.NewHandler(
		transactioncontroller.WithLogger(transactionLog.WithGroup("controller")),
		transactioncontroller.WithService(transactionService),
	))

	// Инициализация HTTP сервера.
	http := start(server.New(cfg.Server.HTTP,
		server.WithLogger(serverLog.WithGroup("server")),
		server.WithController(authController),
		server.WithController(accountController),
		server.WithController(cardController),
		server.WithController(transactionController),
	))

	go manager.run(http.ListenAndServe)

	logger.Info("server started")

	select {
	case <-manager.quit: // Ждем пока все сервисы не остановятся.
	case <-sigint(): // Или сигнал Interrupt.
	}

	if controllers != nil {
		logger.Info("controllers stopping")
	}

	logger.Info("server stopping")
}

// startErr завершает работу программы с ошибкой, если err != nil.
func startErr(err error, name string) {
	if err != nil {
		fmt.Fprintf(flag.CommandLine.Output(), "Failed to init %s:\n  %v\n\n", name, err)
		flag.Usage()
		os.Exit(1)
	}
}

// start проверяет ошибку, и если она не nil, то завершает программу.
// Это позволяет проводить инициализацию без однотипного кода.
func start[T any](svc T, err error) T {
	name := fmt.Sprintf("%T", svc)
	startErr(err, name)

	return svc
}

// sigint создаёт сигнал, который принимает события [os.Interrupt].
//
//go:noinline
func sigint() <-chan os.Signal {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	return sigint
}
