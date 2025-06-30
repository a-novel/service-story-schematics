package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"

	"github.com/getsentry/sentry-go"

	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/models"
)

var ErrSelectSlugIteration = errors.New("SelectSlugIterationRepository.SelectSlugIteration")

func NewErrSelectSlugIteration(err error) error {
	return errors.Join(err, ErrSelectSlugIteration)
}

type SelectSlugIterationData struct {
	Slug models.Slug

	Table string

	Filter map[string][]any
	Order  []string
}

type SelectSlugIterationRepository struct{}

func NewSelectSlugIterationRepository() *SelectSlugIterationRepository {
	return &SelectSlugIterationRepository{}
}

func (repository *SelectSlugIterationRepository) SelectSlugIteration(
	ctx context.Context, data SelectSlugIterationData,
) (models.Slug, int, error) {
	span := sentry.StartSpan(ctx, "SelectSlugIterationRepository.SelectSlugIteration")
	defer span.Finish()

	span.SetData("slug", data.Slug)
	span.SetData("table", data.Table)
	span.SetData("filter", data.Filter)
	span.SetData("order", data.Order)

	tx, err := lib.PostgresContext(span.Context())
	if err != nil {
		span.SetData("postgres.context.error", err.Error())

		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("get postgres client: %w", err))
	}

	output := new(struct {
		Slug models.Slug `bun:"slug"`
	})

	reg, err := regexp.CompilePOSIX(`^` + regexp.QuoteMeta(string(data.Slug)) + `-([0-9]+)$`)
	if err != nil {
		span.SetData("regex.compile.error", err.Error())

		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("compile regex: %w", err))
	}

	query := tx.NewSelect().
		Model(output).
		ModelTableExpr(data.Table).
		Where("slug ~ ?", reg.String()).
		Limit(1)

	for key, values := range data.Filter {
		query = query.Where(key, values...)
	}

	for _, order := range data.Order {
		query = query.Order(order)
	}

	err = query.Scan(span.Context())
	if err != nil {
		span.SetData("scan.error", err.Error())

		if errors.Is(err, sql.ErrNoRows) {
			return data.Slug + "-1", 1, nil
		}

		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("select slug iteration: %w", err))
	}

	// Capture the index of the last iteration.
	index := 1

	_, err = fmt.Sscanf(string(output.Slug), string(data.Slug)+"-%d", &index)
	if err != nil {
		span.SetData("parse.slug.iteration.error", err.Error())

		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("parse slug iteration: %w", err))
	}

	return models.Slug(fmt.Sprintf("%s-%d", data.Slug, index+1)), index + 1, nil
}
