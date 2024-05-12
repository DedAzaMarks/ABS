package postgres

import (
	"database/sql"

	"github.com/google/uuid"
)

type _ struct {
	db *sql.DB
}

func (p *_) GetById(uuid.UUID) {

}

func (p *_) GetAll() {

}

func (p *_) Insert() {
}
func (p *_) Delete() {

}
