package driver_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/mwblythe/squint"
	"github.com/mwblythe/squint/driver"
	"github.com/stretchr/testify/suite"
)

type H map[string]interface{}
type Bits []interface{}

var ctx = context.TODO()

func (b Bits) Split() (string, Bits) {
	return b[0].(string), b[1:]
}

type DriverSuite struct {
	suite.Suite
	driver  string          // driver name to wrap
	dsn     string          // dsn to open
	builder *squint.Builder // builder to test with
	db      *sql.DB         // wrapped db handle
	count   int64           // insert count
}

func (s *DriverSuite) SetupSuite() {
	var err error

	if s.builder == nil {
		s.builder = squint.NewBuilder()
	}

	driver.Register(s.driver, driver.Builder(s.builder))

	s.db, err = sql.Open("squint-"+s.driver, s.dsn)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DriverSuite) TearDownSuite() {
	s.db.Close()
}

func (s *DriverSuite) TestDriver() {
	if !s.Run("Ping", func() {
		s.Nil(s.db.PingContext(ctx))
	}) {
		s.T().Skip()
	}

	s.Run("Exec", func() {
		for n := 0; n < 5; n++ {
			s.InsertPerson()
		}
	})

	s.Run("QueryRow", func() {
		s.GetPerson()
	})

	s.Run("Query", func() {
		s.GetPeople()
	})

	s.Run("Prepared", func() {
		s.Prepared()
	})

	s.Run("Transaction", func() {
		s.Transaction()
	})
}

func (s *DriverSuite) InsertPerson() {
	s.count++

	query, args := Bits{
		"insert into people",
		H{
			"id":   s.count,
			"name": fmt.Sprintf("user-%d", s.count),
		},
	}.Split()

	_, err := s.db.ExecContext(ctx, query, args...)
	s.Nil(err)
}

func (s *DriverSuite) GetPerson() {
	var name string
	query, args := Bits{"select name from people where id =", s.count}.Split()
	row := s.db.QueryRowContext(ctx, query, args...)

	s.Nil(row.Scan(&name))
	s.True(strings.HasPrefix(name, "user-"))
}

func (s *DriverSuite) GetPeople() {
	query, args := Bits{"select id, name from people", "order by id desc"}.Split()
	res, err := s.db.QueryContext(ctx, query, args...)
	if !s.Nil(err) {
		return
	}
	defer res.Close()

	count := s.count

	for res.Next() {
		var id int64
		var name string
		s.Nil(res.Scan(&id, &name))
		s.Equal(count, id)
		s.NotEmpty(name)

		count--
	}

	s.Empty(count)
}

func (s *DriverSuite) Transaction() {
	// start a transaction
	tx, err := s.db.Begin()
	if !s.Nil(err) {
		return
	}

	// delete a row, check that it worked
	res, err := tx.ExecContext(ctx, "delete from people where id =", s.count)
	s.Nil(err)
	affected, err := res.RowsAffected()
	s.Nil(err)
	s.NotEmpty(affected)

	// rollback so it never happened
	s.Nil(tx.Rollback())

	// confirm count is unchanged
	var count int64
	row := s.db.QueryRowContext(ctx, "select count(*) from people")
	s.Nil(row.Scan(&count))
	s.Equal(s.count, count)
}

func (s *DriverSuite) Prepared() {
	s.Run("WithPlaceholders", func() {
		stmt, err := s.db.PrepareContext(ctx, "select name from people where id = "+s.builder.Binder(1))
		if !s.Nil(err) {
			return
		}

		var name string
		row := stmt.QueryRowContext(ctx, s.count)
		s.Nil(row.Scan(&name))
		s.NotEmpty(name)
	})

	s.Run("WithoutPlaceholders", func() {
		stmt, err := s.db.PrepareContext(ctx, "select count(*) from people")
		if !s.Nil(err) {
			return
		}

		var count int64
		row := stmt.QueryRowContext(ctx)
		s.Nil(row.Scan(&count))
		s.Equal(s.count, count)
	})
}
