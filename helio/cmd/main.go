// Helio API
//
// API for retrieving Butterfly observation data from Florida through the years of 2012-2022
// that is has been observed using the iNaturalist website/app (www.iNaturalist.com)
//
// Schemes: http
// Host: localhost:3000
// BasePath: /entities
// Version: 1.0.0
//
// Contact:
// <mbcarruthers@crimson.ua.edu> https://github.com/mbcarruthers
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"mbcarruthers/helio/dataservice/db"
	"mbcarruthers/helio/routes"
)

const (
	addr = ":8000" // change to 8000 to use with react client
)

var (
	btrflydb *db.DataStore
)

func init() {
	btrflydb = db.NewDataStore(db.Defaultdb)
}

func main() {
	defer func(btrflydb *db.DataStore, ctx context.Context) {
		err := btrflydb.Close(ctx)
		if err != nil {
			log.Printf("database didnt close properly:%s \n", err.Error())
		}
	}(btrflydb, context.Background())

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://*", "http://", "*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length", "Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	entities := r.Group("/entities")
	{
		btrflyHandler := routes.NewEntityRouteHandler(btrflydb)
		entities.POST("/", btrflyHandler.NewEntityHandler) // Note: All mutable operations will move to authorized
		entities.GET("/:id", btrflyHandler.GetEntityById)
		entities.GET("/", btrflyHandler.ListEntityHandler)
		entities.PUT("/:id", btrflyHandler.UpdateEntityHandler)              // Note: All mutable operations will move to authorized
		entities.DELETE("/:id", btrflyHandler.DeleteEntityHandler)           // Note: All mutable operations will move to authorized
		entities.GET("/search", btrflyHandler.SearchEntitiesWithinDateRange) // Todo: Be able to get dates not within a string values.
	}

	if err := r.Run(addr); err != nil {
		log.Fatalf("Error running at port %s\n%s",
			addr,
			err.Error())
	}

}
