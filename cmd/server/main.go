package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/spf13/viper"

	http "app/internal/delivery/http"
	repo "app/internal/repository/mysql"
	service "app/internal/service"
)

func init() {
	viper.SetConfigFile(`config.json`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
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
	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}

	}()

	e := echo.New()

	//limit_type routes
	mysqlrlt := repo.NewMysqlRepositoryLimitType(dbConn)
	slt := service.NewServiceLimitType(mysqlrlt)
	http.NewLimitTypeHandler(e, slt)

	//loan routes
	mysqlrlo := repo.NewMysqlRepositoryLoan(dbConn)
	mysqlrli := repo.NewMysqlRepositoryLimit(dbConn)
	sl := service.NewServiceLoan(mysqlrlo, mysqlrli)
	http.NewLoanHandler(e, sl)

	//user routes
	mysqlru := repo.NewMysqlRepositoryUser(dbConn)
	su := service.NewServiceUser(mysqlru)
	http.NewUserHandler(e, su)

	//user_profile routes
	mysqlrup := repo.NewMysqlRepositoryUserProfile(dbConn)
	sup := service.NewServiceUserProfile(mysqlrup)
	http.NewUserProfileHandler(e, sup)

	log.Fatal(e.Start(viper.GetString("server.address")))

}
