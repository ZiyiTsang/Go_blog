package main

import (
	"Go_blog/app/http/middlewares"
	"Go_blog/bootstrap"
	"Go_blog/pkg/logTool"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"runtime"
)

var router *mux.Router
var db *sql.DB

func main() {

	defer func() {
		err := recover()
		if err != nil {
			switch err.(type) {
			//if is runtime error....
			case runtime.Error:
				fmt.Println("runtime:error", err)
			//other error....
			default:
				fmt.Println("other error", err)
			}
			os.Exit(-1)
		}
		fmt.Println("Thank you for using!")
		os.Exit(0)
	}()
	fmt.Println("initiating...")
	//DBTool.Initialize()
	//db = DBTool.DB
	bootstrap.SetupDB()
	router = bootstrap.SetupRoute()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println("DB can not close")
			os.Exit(-1)
		}
	}(db)
	router.Use(middlewares.ForceHTML)

	go func() {
		err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
		panic(err)
	}()
	var choose string
	for {
		out := 0
		fmt.Println("Do you want to exit?")
		_, err := fmt.Scan(&choose)
		if err != nil {
			logTool.CheckError(err)
		}
		switch choose {
		case "0":
			out = 1
		default:
			fmt.Println("input again")
		}
		if out == 1 {
			break
		}
	}
}
