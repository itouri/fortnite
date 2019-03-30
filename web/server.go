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

	// LEAR
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	// Add usecase

	e.Start(viper.GetString("server.address"))
}
