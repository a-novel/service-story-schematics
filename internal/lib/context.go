package lib

import (
	"fmt"

	"github.com/a-novel-kit/context"
	pgctx "github.com/a-novel-kit/context/pgbun"

	"github.com/a-novel/story-schematics/migrations"
)

func NewAgoraContext(parentCTX context.Context) (context.Context, error) {
	ctx, err := pgctx.NewContext(parentCTX, &migrations.Migrations)
	if err != nil {
		return nil, fmt.Errorf("create postgres context: %w", err)
	}

	return ctx, nil
}
