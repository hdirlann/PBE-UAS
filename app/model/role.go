package model

import "time"

type Role struct {
    ID        string    `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Desc      string    `db:"description" json:"description"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
}
