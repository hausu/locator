package tasks

import (
    "github.com/gin-gonic/gin"
    "gopkg.in/olivere/elastic.v3"
    "os"
    "github.com/hausu/locator/objects"
    "log"
    "encoding/json"
    "strconv"
)

type SearchResults struct {
    Hits int64
    Areas *[]objects.AreaType
}

func Search(c *gin.Context) {
    lat := c.Query("lat")
    lon := c.Query("lon")

    client, err := elastic.NewClient(
    elastic.SetURL(os.Getenv("ELASTIC_HOST")),
    elastic.SetSniff(false),
    elastic.SetTraceLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
    )

    if err != nil {
        c.JSON(500, err.Error())
        return
    }

    q := elastic.NewMatchAllQuery()
    geoQ := elastic.NewGeoDistanceQuery("location")
    qLat, _ := strconv.ParseFloat(lat, 64)
    qLon, _ := strconv.ParseFloat(lon, 64)

    geoQ.Lat(qLat)
    geoQ.Lon(qLon)
    geoQ.Distance("10km")

    r, err := client.Search().Index("areas").Type("area").Query(q).PostFilter(geoQ).Pretty(true).Do()

    if err != nil {
        c.JSON(500, err.Error())
        return
    }

    var results SearchResults

    if r.Hits != nil {
        areas := make([]objects.AreaType, 0)
        for _, hit := range r.Hits.Hits {
            var a objects.AreaType
            err = json.Unmarshal(*hit.Source, &a)

            if err == nil {
                areas = append(areas, a)
            }
        }
        results.Areas = &areas
    }

    results.Hits = r.Hits.TotalHits

    c.JSON(200, results)
}