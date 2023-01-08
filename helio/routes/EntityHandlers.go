// Package routes provides HTTP Routing for create, read, update , and delete utilities surrounding the database of butterfly observations(Entites).
package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mbcarruthers/helio/dataservice/db"
	"mbcarruthers/helio/model"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//Note: Crud operations are more than likely not necessary except by an admin. Therefore they will be by
// by authorization in the future. They really are not needed, I really wanted to as a means to get familiar
// with the jackc/pgx Postgres library  as well as weed out any issues I will encounter with the database layout and application layer.

//Note: any mutating operations must include a body with the same identification number as within the url
// as a means to make sure the correct information is being modified.
// Todo: Write OpenAPI compliant documentation (https://ogen.dev/docs/intro/) <- try that.. anything but Swagger

// EntityRouteHandler struct manages routes surrounding a particular entity
type EntityRouteHandler struct {
	btrflydb *db.DataStore
}

// NewEntityRouteHandler constructs a new EntityRouteHandler with a lepidoptera database (and until all functions are made to work with the database-a btrfly array)
func NewEntityRouteHandler(bfdb *db.DataStore) *EntityRouteHandler {
	entities := make([]model.Entity, 0)
	file, _ := os.ReadFile("data/monarch.json")
	_ = json.Unmarshal([]byte(file), &entities)

	// Note: Comment Out / Uncomment following 3 lines if the database does/doesn't need to be created.It will just err & continue,however
	if err := bfdb.CreateAndInsert(entities, context.Background()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateAndInsertError \n %+v\n", err)
	}
	return &EntityRouteHandler{
		btrflydb: bfdb,
	}
}

// NewEntityHandler POST /entities
// Returns entity posted to database
// Produces and Consumes - application/json
// Responses:
// 200 - Successful Operation
// 400 - Invalid Input
// 500 - Error inserting Entity into database
func (e *EntityRouteHandler) NewEntityHandler(c *gin.Context) {
	var btrfly model.Entity
	if err := c.ShouldBindJSON(&btrfly); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "malformed request",
		})
		return
	} else {
		btrfly.Id = 555555                                                              // Note: Represents a dummy value. Will not be allowed once this is in authorized & a permanent database.
		btrfly.Uuid, _ = uuid.FromBytes([]byte("00000000-0000-0000-0000-000000000000")) // Note: It will be created anyway by crdb may as well make it concrete.
		if err := e.btrflydb.InsertNewEntity(btrfly, context.Background()); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "error inserting element",
			})
			return
		}
		c.JSON(http.StatusOK, btrfly)
	}
}

// GetEntityById GET /entities/:id
// Returns the entity based upon its id value which is a required integer parameter found in the url path
// Produces and consumes - application/json
// Responses:
// 200 - Successful Operation
// 400 - Invalid input
// 404 - Entity Not Found
func (e *EntityRouteHandler) GetEntityById(c *gin.Context) {
	id, err := strconv.Atoi(strings.ReplaceAll(c.Param("id"), " ", "")) // remove any spaces left by accident
	if err != nil {
		// Send error if identification cannot be parsed
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "err parsing: invalid syntax",
		})
		return
	}
	entity, err := e.btrflydb.GetEntityById(id, context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, entity)
}

// ListEntityHandler GET /entities
// Returns a collection of Entities within the database.
// Produces and Consumes - application/json
// Responses:
// 200 - Successful operation. Returns all entities within the database.
// 500 - Internal Database Error
func (e *EntityRouteHandler) ListEntityHandler(c *gin.Context) {
	if entities, err := e.btrflydb.ListAllEntities(context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, entities)
		return
	}
}

// UpdateEntityHandler Method PUT /entities/:id
// updates the values of a given Entity, based upon its id.
// the id is a required integer parameter found in the url path /entities/:id using the PUT method.
// UpdateEntityHandler Will not create a new Entity if one does not exist.
// Produces and Consumes - application/json
// Returns:
// 200 - Successful operation. Returns update Entity
// 400 - Invalid Input
// 500 - Internal database error
func (e *EntityRouteHandler) UpdateEntityHandler(c *gin.Context) {
	id, err := strconv.Atoi(strings.ReplaceAll(c.Param("id"), " ", "")) // remove any spaces left by accident
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "malformed input",
		})
		return
	}
	var btrfly model.Entity
	if err := c.ShouldBindJSON(&btrfly); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := e.btrflydb.UpdateEntityById(id, btrfly, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "error updating",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "update successful",
		})
	}

}

// DeleteEntityHandler DELETE /entities/:id
// Deletes an Entity in the database based upon its id  value which is an integer value  as a parameter provided within the path URL.
// Produces and Consumes: application/json
// responses:
// 200 - Successful Operation
// 400 - Invalid input / No Request body(todo:Remove that condition and change function signature of crdb function to just an integer)
// 500 - Database error
func (e *EntityRouteHandler) DeleteEntityHandler(c *gin.Context) {
	id, err := strconv.Atoi(strings.ReplaceAll(c.Param("id"), " ", "")) // in case any space is accidentally left in postman
	if err != nil {
		// if there is a problem extracting the observation identification from URL.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err = e.btrflydb.DeleteEntityById(id, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("deleted observation %d from database", id),
		})
		return
	}
}

// hope to be Route /entities/search?taxon_id=XXX&date1=yyyy-mm-dd&date2=yyyy-mm-dd
// TestEntitySearchHandler is a temporary name for a handler that searches for two dates or a taxon_id
// Note: Dear god. It works fine but don't use it anymore.
func (e *EntityRouteHandler) SearchEntityTaxonIdAndDateRange(c *gin.Context) { // StoppingPoint
	var queryEntity model.EntityQuery
	err := c.ShouldBindQuery(&queryEntity)
	if err != nil {
		// return an error if nothing is able to bind
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if queryEntity.Date1.Valid && queryEntity.Date2.Valid {
		// if there is a taxonId
		if queryEntity.TaxonId != 0 {
			if entities, err := e.btrflydb.GetEntitiesByTaxonIdWithinDateRange(queryEntity.TaxonId, queryEntity.Date1, queryEntity.Date2, context.Background()); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // return as error if the query falls through
				return
			} else {
				c.JSON(http.StatusOK, entities) // return entities by taxon_id within date range if not
				return
			}
		} else {
			entities, err := e.btrflydb.GetEntitiesWithinRange(queryEntity.Date1, queryEntity.Date2, context.Background())
			if err != nil {
				// if any entities(regardless of taxon) cannot be found within a range
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "message": "could not get items within range"})
				return
			} else {
				// return entities(regardless of taxon_id) of that particular range
				c.JSON(http.StatusOK, entities)
				return
			}
		}
	}
	if queryEntity.TaxonId != 0 {
		entities, err := e.btrflydb.GetEntitiesByTaxonId(queryEntity.TaxonId, context.Background())
		if err != nil {
			// returns not found(404) if the data cannot be found
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		} else {
			// returns entities by taxon_id on a successful operation
			c.JSON(http.StatusOK, entities)
			return
		}
	}
}

// SearchEntitiesWithinDateRange - search entities within a given date range.
// Route /entities/search?date1=yyyy-mm-dd&date2=yyyy-mm-dd
func (e *EntityRouteHandler) SearchEntitiesWithinDateRange(c *gin.Context) {
	var dateQuery model.DateQuery
	err := c.ShouldBindQuery(&dateQuery)
	if err != nil {
		// if there is an error in formatting
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err.Error(),
		})
		return
	}
	// get Entities Within Date range
	if entities, err := e.btrflydb.GetEntitiesWithinRange(dateQuery.Date1, dateQuery.Date2, context.Background()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
		return
	} else {
		c.JSON(http.StatusOK, entities)
	}
}
