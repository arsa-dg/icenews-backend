package handler

import "context"

func (h handler) DBMigrate() {
	h.DB.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS users(
			id varchar(255) primary key, 
			username varchar(255), 
			password varchar(255), 
			name varchar(255), 
			bio varchar(255), 
			web varchar(255), 
			picture varchar(255)
		)
	`)
}
