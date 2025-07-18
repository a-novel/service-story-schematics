package main

import (
	"context"
	"log"

	"github.com/a-novel/golib/postgres"

	"github.com/a-novel/service-story-schematics/migrations"
	"github.com/a-novel/service-story-schematics/models/config"
)

func main() {
	ctx, err := postgres.NewContext(context.Background(), config.PostgresPresetDefault)
	if err != nil {
		log.Fatalf("failed to create context: %v", err)
	}

	err = postgres.RunMigrationsContext(ctx, migrations.Migrations)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}
