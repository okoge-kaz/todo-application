package db

// schema.go provides data models in DB
import (
	"time"
)

// Task corresponds to a row in `tasks` table
type Task struct {
	ID          uint64    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	IsDone      bool      `db:"is_done"`
	Status      string    `db:"status"`
	DueDate     time.Time `db:"due_date"`
	Priority    int       `db:"priority"`
}

type User struct {
	ID       uint64 `db:"id"`
	Name     string `db:"name"`
	Password []byte `db:"password"`
}
