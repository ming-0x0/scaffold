package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

	router := gin.Default()
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	})

	router.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, user)
	})

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}

	logger.Info("Server started on :8080")
}
