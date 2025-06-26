package dao_test

import (
	"context"
	"os"

	"github.com/a-novel/service-story-schematics/internal/lib"
)

var ctx context.Context

func init() {
	var err error

	//nolint:fatcontext
	ctx, err = lib.NewPostgresContext(context.Background(), os.Getenv("DAO_DSN"))
	if err != nil {
		panic(err)
	}
}
