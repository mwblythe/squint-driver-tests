package driver_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/mwblythe/squint"
	"github.com/stretchr/testify/suite"
)

type PostgresSuite struct {
	DriverSuite
}

func (s *PostgresSuite) SetupSuite() {
	s.DriverSuite.driver = "postgres"
	s.DriverSuite.dsn = "host=localhost user=postgres password=squint dbname=squint sslmode=disable"
	s.builder = squint.NewBuilder(squint.BindDollar())

	s.DriverSuite.SetupSuite()

	if err := s.db.Ping(); err != nil {
		s.T().Fatal(err)
	}

	_, err := s.db.Exec(`drop table if exists people`)
	if err == nil {
		_, err = s.db.Exec(`
			create table people (
				id   integer primary key,
				name text not null
			)
		`)
	}

	if err != nil {
		s.T().Fatal(err)
	}
}

func TestPostgres(t *testing.T) {
	suite.Run(t, &PostgresSuite{})
}
