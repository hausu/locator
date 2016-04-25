# Locator Service

This service is used to help us to locate a coordinates. We have another repo with some json files that describes different areas. This is used later to locate a coordinates into those areas.

For instance, we have all the boroughs of London. Given any coordinate inside London we want to know to which borouch aproximately that coordinates belongs to.

## Stack used

### Go

We used Go to build this microservice. This is a really small service built on top of Gin and some other packages like:

- Gin
- Godeps
- Golang-lru
- Elastic.v3

We use the lru package to give some cache functionality to this service.

### Elastic

We store those services on Elastic. We use the geocode filter in order to given a coordinates and a radius find the best match. We have to improve this when we have several points to describe a borough.
