package tasks

import (
    "github.com/gin-gonic/gin"
    "gopkg.in/olivere/elastic.v3"
    "github.com/hausu/locator/objects"
    "encoding/json"
    "strconv"
    "github.com/hashicorp/golang-lru"
    "fmt"
)

type SearchResults struct {
    Hits int64
    Areas *[]objects.AreaType
}

var lruCache 	*lru.Cache

func init() {
    lruCache, _ = lru.New(65000)
}

func Search(c *gin.Context) {
    lat := c.Query("lat")
    lon := c.Query("lon")

    cacheKey := fmt.Sprintf("%s_%s", lat, lon)

    if lruCache.Contains(cacheKey) {
        buffer, _ := lruCache.Get(cacheKey)
        c.JSON(200, buffer)
        return
    }

    client, _ := c.MustGet("elastic").(*elastic.Client)

    q := elastic.NewMatchAllQuery()
    geoQ := elastic.NewGeoDistanceQuery("location")
    qLat, _ := strconv.ParseFloat(lat, 64)
    qLon, _ := strconv.ParseFloat(lon, 64)

    geoQ.Lat(qLat)
    geoQ.Lon(qLon)
    geoQ.Distance("15km")

    sorter := elastic.NewGeoDistanceSort("location").Point(qLat, qLon)

    r, err := client.Search().Index("areas").Type("area").Query(q).PostFilter(geoQ).SortBy(sorter).From(0).Size(1).Pretty(true).Do()

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

    _ = lruCache.Add(cacheKey, results)

    c.JSON(200, results)
}