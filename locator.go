package main

import (
    "github.com/gin-gonic/gin"
    "github.com/hausu/locator/tasks"
)

func main() {
    r := gin.Default()

    v1 := r.Group("api/v1")
    {
        v1.GET("/ping", func(c *gin.Context) {
            c.JSON(200, gin.H{
                "message": "pong",
            })
        })

        v1.GET("/import", tasks.ImportAreas)
        v1.GET("/search", tasks.Search)
    }

    r.Run()
}