package main

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
	bunDB "github.com/ming-0x0/scaffold/internal/infra/db"
	sloglogger "github.com/ming-0x0/scaffold/internal/infra/logger/slog"
)

func main() {
	logger := sloglogger.New()
	ctx := context.Background()

	// Open database
	db, err := bunDB.NewPostgreSQLDB(ctx, logger)
	if err != nil {
		panic(err)
	}
	defer bunDB.ClosePostgreSQLDB(ctx, db, logger)

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
	if err != nil {
		panic(err)
	}

	fmt.Printf("User: %+v\n", user)
}
