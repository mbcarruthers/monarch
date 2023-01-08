package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"mbcarruthers/helio/model"
)

//Todo: Implement a pgxPool connection && set a logger for pgx

// DataConfig is a type created to limit errors created when instantiating a database
type DataConfig string

// DataConfig.String() is created to facilitate passing a DataConfig as a parameter to a function
func (d DataConfig) String() string {
	return string(d)
}

var (
	Defaultdb DataConfig = "postgresql://root@cockroach:26257/defaultdb?sslmode=disable" // Note: default database configuration. Nothing fancy, just for testing.
)

// DataStore represents a Cochroachdb connection and facilitates operations surrounding it.
type DataStore struct {
	Conn *pgx.Conn
}

// NewDataStore creates a new default Database Connection.
func NewDataStore(dataConfig DataConfig) *DataStore {
	config, err := pgx.ParseConfig(dataConfig.String())
	if err != nil {
		log.Fatalf("Error setting database configuration! %+v \n", err)
	}
	// set a default name for the session
	config.RuntimeParams["application_name"] = "$ helio"

	conn, err := pgx.ConnectConfig(context.Background(), config) // Note: Might not be a bad place for a Circuit Breaker Pattern
	if err != nil {
		log.Fatalf("Error connecting to database! %+v \n", err)
	}

	return &DataStore{
		Conn: conn,
	}
}

// DataStore.Close() extends the pgx.Conn Close() function to be defered outside of the structure.
func (d *DataStore) Close(ctx context.Context) error {
	return d.Conn.Close(ctx)
}

// DataStore.CreateAndInsert() function to Create and Insert information into the a temporary database produced by docker-compose. For testing.
// Note: specifically for testing called upon in the creation of a new EntityRouteHandler
func (d *DataStore) CreateAndInsert(observations []model.Entity, ctx context.Context) error {
	// preparedStatements is created to make creation + insertion a bit easier to read.
	preparedStatements := map[string]string{
		"database": "CREATE DATABASE observations",
		"table": "CREATE TABLE observations.fl_lepidoptera(" +
			"id INT8 PRIMARY KEY NOT NULL," +
			"taxon_id INT NOT NULL, " +
			"uuid UUID UNIQUE NOT NULL," +
			"place_guess STRING NOT NULL," +
			"species_guess STRING NOT NULL," +
			"latitude string NOT NULL," +
			"longitude string NOT NULL," +
			"observed_on DATE NOT NULL," +
			"time_zone STRING NOT NULL);",
	} // Todo: Fix the errors below to not expose SQL information
	if _, err := d.Conn.Exec(ctx, preparedStatements["database"]); err != nil {
		return fmt.Errorf("Error creating database\n %+v", err)
	} else if _, err = d.Conn.Exec(ctx, preparedStatements["table"]); err != nil {
		return fmt.Errorf("Error creating table \n %+v\n", err)
	} else {
		log.Println("Database and table created")
		tx, err := d.Conn.Begin(ctx)

		if err != nil {
			return fmt.Errorf("Error beginning transaction\n $+v \n", err)
		}
		// rollback if something Went wrong before commit
		defer func(t pgx.Tx, c context.Context) {
			if err := t.Rollback(c); err != nil && err != pgx.ErrTxClosed {
				log.Printf("Error, Rolled Back. Connection Closed \n %+v\n", err)
			}
		}(tx, ctx)
		// insert items into database
		for _, item := range observations { // Todo: Change the name 'item' to 'entity'
			insertTransaction := fmt.Sprintf("INSERT INTO observations.fl_lepidoptera(id,taxon_id,uuid,place_guess,species_guess,latitude,longitude,observed_on,time_zone)" +
				"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)")

			_, err := tx.Exec(ctx, insertTransaction, item.Id, item.TaxonId, item.Uuid, item.PlaceGuess, item.SpeciesGuess, item.Latitude, item.Longitude, item.ObservedOn, item.TimeZone)
			if err != nil {
				return fmt.Errorf("Error executing \n %v \n", err)
			}
		}
		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("Error committing insert operations\n %+v \n", err)
		}
	}
	return nil
}

// InsertNewEntity function to insert a new model.Entity into observations.fl_lepidoptera
// Note: Used within the EntityRouteHandler.NewEntityHandler
func (d *DataStore) InsertNewEntity(entity model.Entity, ctx context.Context) error {
	// Note:Upon insertion, even though UUID is NOT NULL, it will generate a zero value for uuid(000-000...).
	tx, err := d.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Commit(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Println("Error Commiting change.")
		}
	}(tx, ctx)
	insertTransaction := fmt.Sprintf("INSERT INTO observations.fl_lepidoptera(id,taxon_id,uuid,place_guess,species_guess,latitude,longitude,observed_on,time_zone)" +
		"VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)")
	_, err = tx.Exec(ctx, insertTransaction, entity.Id, entity.TaxonId, entity.Uuid, entity.PlaceGuess, entity.SpeciesGuess, entity.Latitude, entity.Longitude, entity.ObservedOn, entity.TimeZone)
	if err != nil {
		log.Printf("Err inserting element\n %s \n", err.Error())
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		log.Printf("Err commiting insert of %d\n %s \n",
			entity.Id, err.Error())
		return fmt.Errorf("CommitErr")
	}
	return nil
}

// GetEntityById requests an entity by its observation id from the database.
// Note: Used within the EntityRouteHandler.GetEntityById
func (d *DataStore) GetEntityById(id int, ctx context.Context) (model.Entity, error) {
	var entity model.Entity
	selectStatement := fmt.Sprintf("SELECT (id,taxon_id,uuid,place_guess,species_guess,latitude,longitude,observed_on,time_zone)" +
		"FROM observations.fl_lepidoptera WHERE id = $1")
	if err := d.Conn.QueryRow(ctx, selectStatement, id).Scan(&entity); err != nil {
		log.Printf("Error finding %d \n %s \n", id, err.Error())
		return model.Entity{}, fmt.Errorf("err not found")
	} else {
		return entity, nil
	}
}

// ListAllEntities requests all information within the database of observations.fl_lepidoptera
// Note: Made primarily for EntityRouteHandler.ListEntityHandler
func (d *DataStore) ListAllEntities(ctx context.Context) ([]model.Entity, error) {
	selectStatement := "SELECT id,taxon_id,uuid,place_guess,species_guess,latitude,longitude,observed_on,time_zone FROM observations.fl_lepidoptera"
	rows, err := d.Conn.Query(ctx, selectStatement)
	if err != nil {
		log.Println("Error executing query for listing all elements\n %s\n",
			err.Error())
		return nil, fmt.Errorf("err execute")
	}
	defer rows.Close()
	entities := []model.Entity{}
	for rows.Next() {
		var entity model.Entity
		if err := rows.Scan(&entity.Id, &entity.TaxonId, &entity.Uuid, &entity.PlaceGuess, &entity.SpeciesGuess, &entity.Latitude, &entity.Longitude, &entity.ObservedOn, &entity.TimeZone); err != nil {
			log.Printf("Error Scanning through entities\n %s \n %s",
				err.Error(),
				rows.Err())
			return nil, fmt.Errorf("error scanning entities")
		} else {
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

// UpdateEntityById updates the database entry by id
// Note: Made to be used in UpdateEntityHandler
func (d *DataStore) UpdateEntityById(id int, entity model.Entity, ctx context.Context) error {
	tx, err := d.Conn.Begin(ctx)
	if err != nil {
		log.Printf("Error Beginning update Query\n %s \n",
			err.Error())
		return fmt.Errorf("err execute")
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Commit(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("Error Commiting.\n %+v \n",
				err.Error())
		}
	}(tx, ctx)
	tag, err := tx.Exec(ctx, "UPDATE observations.fl_lepidoptera SET "+
		"place_guess = $1,species_guess= $2, latitude = $3, longitude = $4, observed_on = $5,"+
		"time_zone = $6 WHERE id = $7", entity.PlaceGuess, entity.SpeciesGuess, entity.Latitude, entity.Longitude,
		entity.ObservedOn, entity.TimeZone, id)
	if err != nil {
		log.Printf("Err executing Update \n %s \n", err.Error())
		return fmt.Errorf("ErrExecute")
	} else if tag.RowsAffected() == 0 {
		return fmt.Errorf("Err not found")
	} else {
		//return entity, tx.Commit(ctx) // <- what it was, should i keep it that way?
		if err = tx.Commit(ctx); err != nil { // Note: Should this even happen?
			log.Printf("Error commiting update\n %s \n",
				err.Error())
			return fmt.Errorf("could not persist data")
		} else {
			return nil
		}
	}
}

// DeleteEntity deletes an entity within the database by id but cross-references the id with the id in the request body
// Note: Made to be used with the DeleteEntityHandler
func (d *DataStore) DeleteEntityById(id int, ctx context.Context) error {
	tx, err := d.Conn.Begin(ctx) // To conform to the name? or pass with model
	if err != nil {              // and cross-reference the id to the model?
		log.Printf("Error beginning deletion \n %s \n",
			err.Error())
		return fmt.Errorf("err connect")
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		if err := tx.Commit(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Println(err.Error())
		}
	}(tx, ctx)

	if tag, err := tx.Exec(ctx, "DELETE FROM observations.fl_lepidoptera WHERE id = $1", id); err != nil {
		log.Printf("Err deleting %d \n %s\n",
			id, err.Error())
		return err
	} else if tag.RowsAffected() == 0 {
		return fmt.Errorf("err not found")
	} else {
		return tx.Commit(ctx)
	}
}

// GetEntitiesByTaxonId accept a query parameter named taxon_id within a request-url query string and returns the
// results of the entities with that taxon id value and a nil error or , on error, it returns nil and the error
// Route GET /entities/search?
func (d *DataStore) GetEntitiesByTaxonId(taxon int, ctx context.Context) ([]model.Entity, error) { // StoppingPoint- testing
	queryStatement := fmt.Sprintf("SELECT id, taxon_id, uuid, place_guess, species_guess, latitude, longitude, observed_on, time_zone FROM observations.fl_lepidoptera WHERE taxon_id =$1")
	if rows, err := d.Conn.Query(ctx, queryStatement, taxon); err != nil { // Note: Make sure to change error to
		return nil, err
	} else {
		defer rows.Close()
		var entities []model.Entity
		for rows.Next() {
			var entity model.Entity
			if err := rows.Scan(&entity.Id, &entity.TaxonId, &entity.Uuid,
				&entity.PlaceGuess, &entity.SpeciesGuess, &entity.Latitude,
				&entity.Longitude, &entity.ObservedOn, &entity.TimeZone); err != nil {
				return nil, err
			} else {
				entities = append(entities, entity)
			}
		}
		return entities, nil
	}
}

// GetEntitiesByTaxonIdWithinDateRange is a very long , and aptly named function to retrieve entities based on taxon_id(species) within
// a set of date ranges
// Route /entities/search?
func (d *DataStore) GetEntitiesByTaxonIdWithinDateRange(taxon_id int, date_one pgtype.Date, date_two pgtype.Date, ctx context.Context) ([]model.Entity, error) {
	log.Println("GetEntitiesByTaxonIdWithinDateRange called")
	if date_one.Time.After(date_two.Time) {
		date_two, date_one = date_one, date_two // Swap values just in case date_two is greater than date_one. I'm programming for me, so I'm preparing for idiocy
	}
	queryStatement := fmt.Sprintf("SELECT id,taxon_id, uuid, place_guess, species_guess, latitude, longitude, observed_on , time_zone FROM observations.fl_lepidoptera WHERE taxon_id = $1 AND observed_on BETWEEN $2 AND $3")

	rows, err := d.Conn.Query(ctx, queryStatement, taxon_id, date_one, date_two)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []model.Entity
	for rows.Next() {
		var entity model.Entity
		if err := rows.Scan(&entity.Id, &entity.TaxonId, &entity.Uuid,
			&entity.PlaceGuess, &entity.SpeciesGuess, &entity.Latitude,
			&entity.Longitude, &entity.ObservedOn, &entity.TimeZone); err != nil {
			return nil, err
		} else {
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

// GetEntitiesWithinRange queries all entities within a date range provided
// Route /entities/search?
func (d *DataStore) GetEntitiesWithinRange(date_one pgtype.Date, date_two pgtype.Date, ctx context.Context) ([]model.Entity, error) {
	log.Println("GetEntitiesWithinRange invoked")
	if date_one.Time.After(date_two.Time) {
		date_two, date_one = date_one, date_two // Swap values just in case date_two is greater than date_one. I'm programming for me, so I'm preparing for idiocy
	}
	queryStatement := fmt.Sprintf("SELECT id,taxon_id, uuid, place_guess, species_guess, latitude, longitude, observed_on , time_zone FROM observations.fl_lepidoptera WHERE observed_on BETWEEN $1 AND $2")

	rows, err := d.Conn.Query(ctx, queryStatement, date_one, date_two)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entities []model.Entity
	for rows.Next() {
		var entity model.Entity
		if err := rows.Scan(&entity.Id, &entity.TaxonId, &entity.Uuid,
			&entity.PlaceGuess, &entity.SpeciesGuess, &entity.Latitude,
			&entity.Longitude, &entity.ObservedOn, &entity.TimeZone); err != nil {
			return nil, err
		} else {
			entities = append(entities, entity)
		}
	}
	return entities, nil
}

// GetEntitiesWithinYear queries all given Entities from a given year
// Todo: Need to do some testing on this one. And throw it into a route
func (d *DataStore) GetEntitiesWithinYear(year pgtype.Date, ctx context.Context) ([]model.Entity, error) { // Note: Should consider changing datatype of year parameter
	_year := year.Time.Year()
	queryStatement := fmt.Sprintf("SELECT id, taxon_id, uuid, place_guess, species_guess, latitude, longitude, observed_on, time_zone FROM observations.fl_lepidoptera WHERE date_part('year',observed_on) = $1")
	rows, err := d.Conn.Query(ctx, queryStatement, _year) // query entities by year
	if err != nil {
		// return err if there is something wrong with query or database
		return nil, err
	}
	defer rows.Close()
	var entities []model.Entity

	for rows.Next() { // loop through results from query
		var entity model.Entity // scan entities into variable
		if err := rows.Scan(&entity.Id, &entity.TaxonId, &entity.Uuid, &entity.PlaceGuess, &entity.SpeciesGuess, &entity.Latitude, &entity.Longitude, &entity.ObservedOn, &entity.TimeZone); err != nil {
			return nil, err
		} else {
			entities = append(entities, entity) // append entities into slice
		}
	}
	return entities, nil
}
