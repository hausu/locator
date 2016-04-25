package services
import (
    "github.com/gin-gonic/gin"
    "gopkg.in/olivere/elastic.v3"
    "log"
    "os"
)

func ElasticMiddleWare() gin.HandlerFunc {
    return func (c *gin.Context) {
        client, err := elastic.NewClient(
            elastic.SetURL(os.Getenv("ELASTIC_HOST")),
            elastic.SetSniff(false),
            elastic.SetTraceLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
        )

        if err != nil {
            panic(err)
        }

        c.Set("elastic", client)
        c.Next()
    }
}