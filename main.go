package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	articleHttpDeliveryHTTP "github.com/caudaganesh/go-pbt-demo/internal/article/delivery/http"
	articleRepo "github.com/caudaganesh/go-pbt-demo/internal/article/repository"
	articleUC "github.com/caudaganesh/go-pbt-demo/internal/article/usecase"
	authorRepo "github.com/caudaganesh/go-pbt-demo/internal/author/repository"
	"github.com/caudaganesh/go-pbt-demo/internal/middleware"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
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
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
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

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	authorRepo := authorRepo.NewMysqlAuthorRepository(dbConn)
	ar := articleRepo.NewMysqlArticleRepository(dbConn)

	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	au := articleUC.NewArticleUsecase(ar, authorRepo, timeoutContext)
	articleHttpDeliveryHTTP.NewArticleHandler(e, au)

	log.Fatal(e.Start(viper.GetString("server.address")))
}
