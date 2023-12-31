package main

import (
	"github.com/Godzizizilla/douyin-simple/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// go service.RunMessageServer()

	// init database connection
	database.InitDB()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
