package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/sbxb/loyalty/models"
	"github.com/sbxb/loyalty/storage"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// DBStorage defines a database storage implemented as a wrapper
// around any database/sql implementation
type DBStorage struct {
	db         *sql.DB
	userTable  string
	orderTable string
}

// DBStorage implements Storage interface
var _ storage.Storage = (*DBStorage)(nil)

// if it takes more than 2 seconds to ping the database, then database
// is considered unavailable
const pingTimeout = 2 * time.Second

func NewDBStorage(dsn string) (*DBStorage, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DBStorage: empty dsn")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("DBStorage: Open: %v", err)
	}

	// ping the database before returning DBStorage instance
	ctx, cancel := context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("DBStorage: Ping: %v", err)
	}

	// create all the necessary tables in the database
	userTable := "users"
	orderTable := "orders"
	if err := createTables(db, userTable, orderTable); err != nil {
		db.Close()
		return nil, fmt.Errorf("DBStorage: Create Tables: %v", err)
	}

	return &DBStorage{db: db, userTable: userTable, orderTable: orderTable}, nil
}

func createTables(db *sql.DB, userTable, orderTable string) error {
	userTableQuery := `CREATE TABLE IF NOT EXISTS ` + userTable + ` (
		id INT primary key GENERATED ALWAYS AS IDENTITY,
		login VARCHAR(128) NOT NULL UNIQUE,
		hash VARCHAR(256) NOT NULL
	)`
	orderTableQuery := `CREATE TABLE IF NOT EXISTS ` + orderTable + ` (
		id INT primary key GENERATED ALWAYS AS IDENTITY,
		number TEXT NOT NULL UNIQUE,
		status VARCHAR(16) NOT NULL,
		accrual BIGINT NOT NULL DEFAULT 0,
		uploaded_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
		user_id INT NOT NULL REFERENCES ` + userTable + ` (id) ON DELETE CASCADE		
	)`
	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("DBStorage: createTables: %v", err)
	}
	defer tx.Rollback()

	tables := []string{userTableQuery, orderTableQuery}
	for _, tableName := range tables {
		if _, err := tx.Exec(tableName); err != nil {
			return fmt.Errorf("DBStorage: createTables: %v", err)
		}
	}

	return tx.Commit()
}

// tests use Truncate() to reset changes
func (st *DBStorage) TruncateTables() error {
	tx, err := st.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("DBStorage: createTables: %v", err)
	}
	defer tx.Rollback()

	tables := []string{st.userTable, st.orderTable}
	for _, tableName := range tables {
		if _, err := tx.Exec(`TRUNCATE ` + tableName + ` RESTART IDENTITY CASCADE`); err != nil {
			return fmt.Errorf("DBStorage: truncateTables: %v", err)
		}
	}

	return tx.Commit()
}

func (st *DBStorage) AddUser(ctx context.Context, user *models.User) error {
	AddURLQuery := `INSERT INTO ` + st.userTable + `(login, hash) VALUES($1, $2)`
	result, err := st.db.ExecContext(ctx, AddURLQuery, user.Login, user.Hash)
	if err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return storage.ErrLoginAlreadyExists
		}
		return fmt.Errorf("DBStorage: AddUser: %v", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DBStorage: AddUser: %v", err)
	}
	if rows != 1 {
		return fmt.Errorf("DBStorage: AddUser: expected to affect 1 row, affected %d", rows)
	}

	return nil
}

func (st *DBStorage) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	dbUser := &models.User{}

	GetURLQuery := `SELECT id, login, hash FROM ` + st.userTable + ` WHERE login=$1`
	err := st.db.QueryRowContext(ctx, GetURLQuery, user.Login).Scan(
		&dbUser.ID, &dbUser.Login, &dbUser.Hash,
	)
	switch {
	case err == sql.ErrNoRows:
		return nil, storage.ErrLoginMissing
	case err != nil:
		return nil, fmt.Errorf("DBStorage: GetUser: %v", err)
	default:
		return dbUser, nil
	}
}

func (st *DBStorage) AddOrder(ctx context.Context, order *models.Order, userID int) error {
	AddOrderQuery := `INSERT INTO ` + st.orderTable + `(number, status, user_id) 
		VALUES($1, $2, $3)`
	_, err := st.db.ExecContext(ctx, AddOrderQuery, order.Number, order.Status, userID)

	if err == nil {
		return nil
	}
	if !strings.Contains(err.Error(), "SQLSTATE 23505") {
		return fmt.Errorf("DBStorage: AddOrder: %v", err)
	}

	// conflict due to unique constraint: check who owns the order
	var uid int
	CheckOrderQuery := `SELECT user_id FROM ` + st.orderTable + ` WHERE 
		number = $1`
	err = st.db.QueryRowContext(ctx, CheckOrderQuery, order.Number).Scan(&uid)

	if err != nil {
		return fmt.Errorf("DBStorage: AddOrder: %v", err)
	}

	return storage.NewExistingOrderError(uid)
}

func (st *DBStorage) GetOrders(ctx context.Context, userID int) ([]*models.Order, error) {
	res := []*models.Order{}

	GetOrdersQuery := `SELECT number, status, accrual, uploaded_at FROM ` + st.orderTable + `
		WHERE user_id = $1 ORDER BY uploaded_at ASC`

	rows, err := st.db.QueryContext(ctx, GetOrdersQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("DBStorage: GetOrders: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		order := &models.Order{}
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.Uploaded_at)
		if err != nil {
			return nil, fmt.Errorf("DBStorage: GetOrders: %v", err)
		}
		res = append(res, order)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("DBStorage: GetOrders: %v", err)
	}

	return res, nil
}

func (st *DBStorage) Close() error {
	if st.db == nil {
		return nil
	}

	if err := st.db.Close(); err != nil {
		return err
	}

	st.db = nil

	return nil
}
