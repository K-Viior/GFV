package main

import (
	"GFV/router"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	router.Init(engine)

	engine.Run("0.0.0.0:6162")

}
