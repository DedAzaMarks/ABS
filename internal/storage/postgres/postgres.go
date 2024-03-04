package postgres

import (
	"database/sql"

	"github.com/google/uuid"
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) GetById(uuid.UUID) {

}

func (p *Postgres) GetAll() {

}

func (p *Postgres) Insert() {

}

func (p *Postgres) Delete() {

}
