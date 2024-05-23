package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/DedAzaMarks/ABS/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var tableNames = []string{"user_devices", "users", "devices"}

const (
	createUsers       = `CREATE TABLE users (user_id SERIAL PRIMARY KEY);`
	createDevices     = `CREATE TABLE devices (device_id SERIAL PRIMARY KEY);`
	createUserDevices = `
CREATE TABLE user_devices (
    user_id INT,
    device_id INT,
    PRIMARY KEY (user_id, device_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES devices(device_id) ON DELETE CASCADE
);`
)

const (
	addUser = `
INSERT INTO users(id) 
VALUES($1) 
ON CONFLICT DO NOTHING;`

	checkUserExists = `
SELECT user_id 
FROM users
WHERE username=$1`
	addClient = `
INSERT INTO devices
DEFAULT VALUES RETURNING device_id;`
	addUserClientPair = `
INSERT INTO user_devices (user_id, device_id) 
VALUES ($1, $2)
ON CONFLICT (user_id, device_id) 
DO NOTHING;`

	getUser = `SELECT * FROM users WHERE id = $1;`
)

type Postgres struct {
	db *sql.DB
}

func (p *Postgres) AddNewUser(ctx context.Context, userID string) error {
	_, err := p.db.ExecContext(ctx, addUser, userID)
	if err != nil {
		log.Println("ERROR: failed to add new user", err)
		return err
	}
	return nil
}

func (p *Postgres) AddNewClient(ctx context.Context, userID, clientID string) error {
	p.db.QueryRowContext(ctx, checkUserExists, userID).Scan(&userID)
	return nil
}

func (p *Postgres) GetUser(ctx context.Context, userID string) (*domain.TGUser, error) {
	user := domain.TGUser{}
	err := p.db.QueryRowContext(ctx, getUser, userID).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("nothing found")
			return nil, sql.ErrNoRows
		}
		log.Println("ERROR: failed to query user", err)
		return nil, err
	}
	rows, err := p.db.QueryContext(ctx, addClient, userID, userID)
	if err != nil {

	}
	for rows.Next() {

	}
	return &user, nil
}

func NewPostgres(ctx context.Context, dsn string) (*Postgres, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Printf("error on connecting to database: %v", err)
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	if err := createIfTableNotExists(ctx, db, createUsers, createDevices, createUserDevices); err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func createIfTableNotExists(ctx context.Context, db *sql.DB, queries ...string) error {
	for _, query := range queries {
		_, err := db.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}

func dropIfExists(ctx context.Context, db *sql.DB, tableName string) error {
	_, err := db.ExecContext(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", tableName))
	if err != nil {
		log.Printf("error on drop %s table: %v", tableName, err)
		return err
	}
	return nil
}
