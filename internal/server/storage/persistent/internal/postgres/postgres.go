package postgres

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/DedAzaMarks/ABS/internal/domain"
	myerrors "github.com/DedAzaMarks/ABS/internal/domain/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

var tableNames = []string{"user_devices", "users", "devices"}

type Postgres struct {
	db *pgx.Conn
}

func (p *Postgres) SaveUser(ctx context.Context, userDTO *domain.UserDTO) error {
	userInsert, err := p.db.Query(ctx, `
INSERT INTO USR (UserID,SessionKey,State)
VALUES ($1,$2,$3)
ON CONFLICT (UserID) DO UPDATE
  SET UserID = excluded.UserID,
      SessionKey = excluded.SessionKey,
      State = excluded.State;`,
		userDTO.ID, userDTO.SessionKey, userDTO.State)
	if err != nil {
		return fmt.Errorf("adding user: %w", err)
	}
	userInsert.Close()
	if err := userInsert.Err(); err != nil {
		return fmt.Errorf("adding new user: %w", err)
	}
	return nil
}

func (p *Postgres) LoadUser(ctx context.Context, userID int64) (*domain.UserDTO, error) {
	dto := &domain.UserDTO{}
	if err := p.db.QueryRow(ctx, `
SELECT u.UserID,u.SessionKey,u.State
FROM USR as u
WHERE u.UserID = $1;`, userID).Scan(&dto.ID, &dto.SessionKey, &dto.State); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.ErrorUserNotFound
		}
		return nil, fmt.Errorf("loading user: %w", err)
	}
	rows, err := p.db.Query(ctx, `
SELECT DISTINCT d.DeviceID,d.DeviceName
FROM USR as u
INNER JOIN USER_DEVICE as ud 
    ON $1 = ud.UserID
INNER JOIN DEVICE as d
	ON ud.DeviceID = d.DeviceID;
`, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.ErrorDeviceNotFound
		}
		return nil, fmt.Errorf("loading user: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var devDTO domain.DeviceDTO
		if err := rows.Scan(&devDTO.ID, &devDTO.Name); err != nil {
			return nil, fmt.Errorf("scanning user: %w", err)
		}
		dto.Devices = append(dto.Devices, devDTO)
	}
	return dto, nil
}

func (p *Postgres) GetUsersByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]*domain.UserDTO, error) {
	rows, err := p.db.Query(ctx, `
select distinct U.UserID,U.SessionKey,U.State
from USR as U
inner join USER_DEVICE as ud
    on U.UserID = ud.UserID
inner join Device as d
    on ud.DeviceID = $1;`, deviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.ErrorDeviceNotFound
		}
		return nil, fmt.Errorf("loading users: %w", err)
	}
	var users []*domain.UserDTO
	for rows.Next() {
		var user domain.UserDTO
		if err := rows.Scan(&user.ID, &user.SessionKey, &user.State); err != nil {
			return nil, fmt.Errorf("scanning users: %w", err)
		}
		users = append(users, &user)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scanning users: %w", err)
	}
	for _, u := range users {
		loadUser, err := p.LoadUser(ctx, u.ID)
		if err != nil {
			return nil, fmt.Errorf("loading user: %w", err)
		}
		u.Devices = append(u.Devices, loadUser.Devices...)
	}
	return users, nil
}

func (p *Postgres) AddNewDevice(ctx context.Context, userID int64, deviceID uuid.UUID, deviceName string) error {
	log.Printf("deviceID:%d; deviceID.String():%s, deviceID.String().len():%d; deviceName:%d", len(deviceID[:]), deviceID.String(), len(deviceID.String()), len(deviceName))
	if err := p.db.QueryRow(ctx, `
INSERT INTO USER_DEVICE(UserID,DeviceID) VALUES ($1,$2);`, userID, deviceID.String()).Scan(); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("adding new user_device: %w", err)
		}
	}
	if err := p.db.QueryRow(ctx, `
INSERT INTO DEVICE(DeviceID, DeviceName) VALUES($1,$2)`, deviceID, deviceName).Scan(); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("adding new device: %w", err)
		}
	}
	return nil
}

const (
	host     = "rc1a-9dhe5sg7m2nxcrew.mdb.yandexcloud.net"
	port     = 6432
	user     = "user"
	password = "password"
	dbname   = "db"
	ca       = "/Users/m.bordyugov/.postgresql/root.crt"
)

func NewPostgres(ctx context.Context) (*Postgres, error) {
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("failed to read ca certificate: %w", err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, fmt.Errorf("failed to append root certificate")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=verify-full target_session_attrs=read-write",
		host, port, dbname, user, password)
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	connConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}
	var usersRows, devicesRows, usersDevicesRows pgx.Rows
	if usersRows, err = conn.Query(ctx, `
CREATE TABLE IF NOT EXISTS USR(
    UserID BIGINT PRIMARY KEY,
    SessionKey CHAR(8),
    State TEXT
);`); err != nil {
		return nil, fmt.Errorf("unable to create users table: %w", err)
	}
	usersRows.Close()
	if usersRows.Err() != nil {
		return nil, fmt.Errorf("unable to create users table: %w", usersRows.Err())
	}

	if devicesRows, err = conn.Query(ctx, `
CREATE TABLE IF NOT EXISTS DEVICE(
    DeviceID CHAR(36) PRIMARY KEY,
    DeviceName TEXT
);`); err != nil {
		return nil, fmt.Errorf("unable to create devices table: %w", err)
	}
	devicesRows.Close()
	if devicesRows.Err() != nil {
		return nil, fmt.Errorf("unable to create devices table: %w", devicesRows.Err())
	}

	if usersDevicesRows, err = conn.Query(ctx, `
CREATE TABLE IF NOT EXISTS USER_DEVICE(
    UserID BIGINT,
    DeviceID CHAR(36),
    PRIMARY KEY(UserID, DeviceID)
);`); err != nil {
		return nil, fmt.Errorf("unable to create users table: %w", err)
	}
	usersDevicesRows.Close()
	if usersDevicesRows.Err() != nil {
		return nil, fmt.Errorf("unable to create users devices table: %w", usersDevicesRows.Err())
	}

	return &Postgres{db: conn}, nil
}
