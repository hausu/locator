package tasks

import (
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "encoding/json"
    "gopkg.in/olivere/elastic.v3"
    "os"
)

type Location struct {
    Lat float32 `json:"lat"`
    Lon float32 `json:"lon"`
}

type AreaType struct {
    City string `json:"city"`
    Country string `json:"country"`
    Name string `json:"name"`
    Slug string `json:"slug"`
    Area float32 `json:"area"`
    Location *Location `json:"location"`
}

func ImportAreas(c *gin.Context) {
    file, e := ioutil.ReadFile("")

    if e != nil {
        c.JSON(500, e.Error())
        return
    }

    client, err := elastic.NewClient(
        elastic.SetURL(os.Getenv("ELASTIC_HOST")),
        elastic.SetSniff(false),
    )

    if err != nil {
        c.JSON(500, err.Error())
        return
    }

    _, err = client.DeleteIndex("areas").Do()

    if err != nil {
        c.JSON(500, err.Error())
        return
    }

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

    var areas []AreaType
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