package models

type Session struct {
	OneTimeId     string `json:"one_time_id"`
	CreatedAtUnix int64  `json:"created_at_unix"`
}
