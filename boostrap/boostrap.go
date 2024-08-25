package boostrap

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/billing-engine/config"
	"github.com/billing-engine/internal/repository"
	"github.com/billing-engine/internal/service"
	"github.com/spf13/viper"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	configuration *config.Config
	once          sync.Once
)

func initConfig() *config.Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("../../../.")
		viper.AddConfigPath("../../../../.")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file, %s", err)
		}

		if err := viper.Unmarshal(&configuration); err != nil {
			log.Fatalf("Unable to decode into struct, %v", err)
		}
	})

	return configuration
}

func Boostrap() *config.AppConfig {
	cfg := initConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("error open sql", err)
	}

	gormConfig := &gorm.Config{}

	gormDB, err := gormtrace.Open(mysql.New(mysql.Config{
		Conn: db,
	}), gormConfig)
	if err != nil {
		log.Fatal("error open gorm", err)
	}

	repo := initRepo(gormDB)

	service := service.NewService(
		repo,
	)

	return &config.AppConfig{
		Config:  cfg,
		Service: service,
	}
}

func initRepo(gormDB *gorm.DB) *repository.Repository {
	loanRepo := repository.NewLoanRepository(gormDB)
	userRepo := repository.NewUserRepository(gormDB)
	payLoanRepo := repository.NewPayLoanRepository(gormDB)

	return &repository.Repository{
		Loan:    loanRepo,
		User:    userRepo,
		PayLoan: payLoanRepo,
	}
}
