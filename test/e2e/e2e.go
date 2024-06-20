package e2e

import (
	"fmt"
	"log"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/stori-test/internal/config"
	"github.com/jordanlanch/stori-test/internal/core/usecase"
	"github.com/jordanlanch/stori-test/internal/infrastructure/email"
	"github.com/jordanlanch/stori-test/internal/infrastructure/repository"
	"github.com/jordanlanch/stori-test/internal/interface/api/controller"
	"github.com/jordanlanch/stori-test/internal/interface/api/router"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup(t *testing.T, cassetteName string) (expect *httpexpect.Expect, teardown func()) {
	t.Helper()

	env := config.NewEnv("../.envtest")
	if err := env.Validate(); err != nil {
		t.Fatalf("Environment validation failed: %v", err)
	}

	redisOptions := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", env.RedisHost, env.RedisPort),
		DB:   0,
	}

	if env.RedisPassword != "" {
		redisOptions.Password = env.RedisPassword
	}

	redisClient := redis.NewClient(redisOptions)

	// Create new VCR cassette
	rec, err := recorder.New(cassetteName)
	if err != nil {
		log.Fatal(err)
	}

	// Use the recorder for all requests
	httpClient := rec.GetDefaultClient()

	// Setup Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup application components
	dbRepo := repository.NewDBTransactionRepository(db, env.CSVFilePath)
	cacheRepo := repository.NewCacheTransactionRepository(redisClient, env.CacheDurationSec)
	emailService := &email.EmailService{}
	transactionUseCase := usecase.NewTransactionUseCase(dbRepo, cacheRepo, emailService, redisClient, env.RateLimit, env.RedisTimeoutSec, env.CacheDurationSec)
	transactionController := &controller.TransactionController{
		UseCase: transactionUseCase,
	}

	router := router.SetupRouter(transactionController)

	srv := httptest.NewUnstartedServer(router)
	listener, err := net.Listen("tcp", "127.0.0.1:42783")
	if err != nil {
		if listener, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			panic(fmt.Sprintf("httptest: failed to listen on a port: %v", err))
		}
	}
	srv.Listener = listener
	srv.Start()
	expect = httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  srv.URL,
		Client:   httpClient,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	return expect, func() {
		redisClient.Close()
		rec.Stop()
		srv.Close()
	}
}
