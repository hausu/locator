package tasks

import (
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/olivere/elastic.v3"
    "github.com/hausu/locator/objects"
)

func ImportAreas(c *gin.Context) {
    file, err := ioutil.ReadFile("/Users/ssola/Workspace/go/src/github.com/hausu/locator_data/gb/london.json")

    if err != nil {
        c.JSON(500, err.Error())
        return
    }

    client, _ := c.MustGet("elastic").(*elastic.Client)

    _, _ = client.DeleteIndex("areas").Do()

    mapping := `{
        "settings":{
            "number_of_shards":1,
            "number_of_replicas":0
        },
        "mappings":{
            "area":{
                "properties":{
                    "city":{
                        "type":"string"
                    },
                    "slug":{
                        "type":"string"
                    },
                    "name":{
                        "type":"string"
                    },
                    "country":{
                        "type":"string"
                    },
                    "location":{
                        "type":"geo_point"
                    },
                    "suggest_field":{
                        "type":"completion",
                        "payloads":true
                    }
                }
            }
        }
    }`

    _, err = client.CreateIndex("areas").BodyString(mapping).Do()

    if err != nil {
        c.JSON(500, err.Error())
        return
    }

    var areas []objects.AreaType
    json.Unmarshal(file, &areas)

    for _, area := range areas {
        _, err = client.Index().Index("areas").Type("area").BodyJson(area).Do()

        if err != nil {
            c.JSON(500, err.Error())
            return
        }
    }

    c.JSON(200, areas)
}