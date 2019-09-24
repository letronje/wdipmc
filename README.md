# [w]here [d]o [i] [p]ark [m]y [c]ar ??

## Importing carpark information

Env var `CARPARK_DSN` is expected to contain the DSN(Data source name) for the MySQL DB to be used.

Env var `CARPARK_CSV` is expected to contain the path of the csv file with carpark info.
( CSV Available @ https://data.gov.sg/dataset/hdb-carpark-information )

Example:

`export CARPARK_DSN="root:@/wdipmc?charset=utf8&parseTime=True&loc=Local"`
`CARPARK_CSV=hdb-carpark-information.csv go run cmd/importcarparks/main.go`

`> Carparks after import: 2113`

## Updating Availability information
`go run cmd/updateavailability/main.go`

`> Success: 1940 , Failures: 76`

## Tests
Tests involving DB operations assume the DSN for test DB to present in `CARPARK_DSN_TEST`

Example:
`export CARPARK_DSN_TEST="root:@/wdipmc_test?charset=utf8&parseTime=True&loc=Local"`

## Running the API
Run `go run server.go` to start the api server. 

## How are `nearest` carparks found ?

Using MySQL's `ST_Distance_Sphere` function which `Returns the mimimum spherical distance between two points on a sphere(earth)`

Available from MySQL 5.7.6 onwards

Doc: https://dev.mysql.com/doc/refman/5.7/en/spatial-convenience-functions.html#function_st-distance-sphere

## Libs/frameworks used

`svy21` https://github.com/cgcai/SVY21/tree/master/Go is used for `SVY21` -> `(lat,lng)` conversion

`gorm` https://gorm.io/ is used as the ORM

`gin` https://github.com/gin-gonic/gin is the http framework used for the endpoint.

`govendor` is used for vendoring.
