package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/spf13/viper"

	"github.com/itouri/fortnite/web/middleware"
	_playerRepo "github.com/itouri/fortnite/web/player/repository"
	_playerUcase "github.com/itouri/fortnite/web/player/usecase"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	dbHost := viper.GetString(`database.host`)
	dbPort := viper.GetString(`database.port`)
	dbUser := viper.GetString(`database.user`)
	dbPass := viper.GetString(`database.pass`)
	dbName := viper.GetString(`database.name`)

	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	//LEARN
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")

	// LEARN
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)
	if err != nil && viper.GetBool("debug") {
		fmt.Println(err)
	}

	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer dbConn.Close()

	e := echo.New()
	mid := middleware.InitMiddleware()
	e.Use(mid.CORS)

	// Add repo
	playerRepo = _playerRepo.NewPlayerRepository(dbConn)

	// LEARN
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	// Add usecase
	playerUcase = _playerUcase.NewPlayerUsecase(playerRepo, timeoutContext)

	e.Start(viper.GetString("server.address"))
}
