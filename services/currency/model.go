package currency

import (
	"database/sql"
	"time"
)

type Currency struct {
	Type      string
	ISO       string
	Chain     sql.NullString
	CreatedAt time.Time
}
