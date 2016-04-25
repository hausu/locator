package objects

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
