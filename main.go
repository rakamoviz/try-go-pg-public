package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	fmt.Println("before")
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println("after")
	fmt.Println(q)
	return nil
}

type Execution struct {
	tableName struct{} `pg:"execution, discard_unknown_columns"`

	ID           int64         `pg:"id, pk"`
	PartnerName  string        `pg:"partner, fk:name"`
	Async        bool          `pg:"is_async`
	Timestamp    time.Time     `pg:"timestamp, default:now()"`
	LoadStatuses []*LoadStatus `pg:"rel:has-many, join_fk:fk_execution"`
}

type LoadStatus struct {
	tableName struct{} `pg:"load_status, discard_unknown_columns"`

	ID          int64     `pg:"id,pk"`
	Event       string    `pg:"event"`
	Status      string    `pg:"status"`
	Description string    `pg:"status"`
	Timestamp   time.Time `pg:"timestamp, default:now()"`
	ExecutionID int64     `pg:"fk_execution"`
	//Execution   *Execution `pg:"rel:has-one, fk:fk_execution"`
}

func ExampleDB_Model() {
	db := pg.Connect(&pg.Options{
		Database: "try_gopg",
		User:     "postgres",
		Password: "secret",
	})

	defer db.Close()

	db.AddQueryHook(dbLogger{})

	/*
		err := createSchema(db)
		if err != nil {
			panic(err)
		}

		execution1 := &Execution{
			PartnerName: "partnerName 1",
			Async:       true,
			Timestamp:   time.Now(),
		}
		_, err = db.Model(execution1).Insert()
		if err != nil {
			panic(err)
		}

		fmt.Println(execution1)

		execution2 := &Execution{
			PartnerName: "partnerName 2",
			Async:       true,
			Timestamp:   time.Now(),
		}
		_, err = db.Model(execution2).Insert()
		if err != nil {
			panic(err)
		}

		fmt.Println(execution2)

		loadStatus1 := &LoadStatus{
			Event:       "event 1",
			Status:      "status 1",
			Description: "description 1",
			Timestamp:   time.Now(),
			ExecutionID: execution1.ID,
		}
		_, err = db.Model(loadStatus1).Insert()
		if err != nil {
			panic(err)
		}

		fmt.Println(loadStatus1)

		loadStatus2 := &LoadStatus{
			Event:       "event 2",
			Status:      "status 2",
			Description: "description 2",
			Timestamp:   time.Now(),
			ExecutionID: execution1.ID,
		}
		_, err = db.Model(loadStatus2).Insert()
		if err != nil {
			panic(err)
		}

		fmt.Println(loadStatus2)

		loadStatus3 := &LoadStatus{
			Event:       "event 3",
			Status:      "status 3",
			Description: "description 3",
			Timestamp:   time.Now(),
			ExecutionID: execution2.ID,
		}
		_, err = db.Model(loadStatus3).Insert()
		if err != nil {
			panic(err)
		}

		fmt.Println(loadStatus3)
	*/

	/*
		execution1 := new(Execution)
		err := db.Model(execution1).Where("id = ?", 1).Relation("LoadStatuses").Relation("LoadStatuses.Execution").Select()

		if err != nil {
			panic(err)
		}

		fmt.Println(execution1)
		fmt.Println(execution1.LoadStatuses[1].Execution)

		loadStatus1 := new(LoadStatus)
		err = db.Model(loadStatus1).Where("load_status.id = ?", 1).Relation("Execution").Select()

		if err != nil {
			panic(err)
		}

		fmt.Println(loadStatus1)
	*/

	execution1 := new(Execution)
	err := db.Model(execution1).
		Where("partner = ?", "MMS_TECHNOLOGY_BR").
		Order("timestamp DESC").
		Limit(1).
		Relation("LoadStatuses", func(q *orm.Query) (*orm.Query, error) {
			return q.Order("timestamp DESC"), nil
		}).
		Select()

	if err != nil {
		panic(err)
	}

	fmt.Println(execution1)
	fmt.Println(execution1.LoadStatuses)
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	var u *Execution
	var s *LoadStatus
	models := []interface{}{
		u,
		s,
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ExampleDB_Model()
}
