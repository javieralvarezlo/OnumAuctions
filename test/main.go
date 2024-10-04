package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.PUT("/", func(c *gin.Context) {

		body, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		fmt.Println("Body:", string(body))

		c.JSON(200, gin.H{"status": "success", "body": string(body)})
	})

	router.Run(":8081")
}
