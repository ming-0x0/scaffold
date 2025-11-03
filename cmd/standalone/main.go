package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
	"github.com/ming-0x0/scaffold/internal/infra/logger/slog/bunslog"
	"github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/pgdialect"
    "github.com/uptrace/bun/driver/pgdriver"

)

func main() {
	logger := sloglogger.New()
	ctx := context.Background()

	// Open database
	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN("postgres://scaffold:password@localhost:5432/scaffold?sslmode=disable"),
	))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Create Bun instance
	db.AddQueryHook(bunslog.New(bunslog.Logger{
		Logger:            logger,
		SlowThreshold:     time.Second,
		IgnoreNoRowsError: false,
		LogLevel:          sloglogger.Info,
	}))

	// Define model
	type User struct {
		ID   int64  `bun:",pk,autoincrement"`
		Name string `bun:",notnull"`
	}

	// Create table
	db.NewCreateTable().Model((*User)(nil)).Exec(ctx)

	// Insert user
	user := &User{Name: "John Doe"}
	db.NewInsert().Model(user).Exec(ctx)

	// Query user
	err := db.NewSelect().Model(user).Where("id = ?", user.ID).Scan(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User: %+v\n", user)
}
