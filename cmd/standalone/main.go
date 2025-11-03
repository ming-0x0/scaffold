package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"
	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
	"github.com/ming-0x0/scaffold/internal/infra/logger/slog/bunslog"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func main() {
	logger := sloglogger.New()
	ctx := context.Background()

	// Open database
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:")
	if err != nil {
		panic(err)
	}

	// Create Bun instance
	db := bun.NewDB(sqldb, sqlitedialect.New())
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
	err = db.NewSelect().Model(user).Where("id = ?", user.ID).Scan(ctx)
	fmt.Printf("User: %+v\n", user)
}
