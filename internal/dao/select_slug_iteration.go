package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"regexp"

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

func (repository *SelectSlugIterationRepository) SelectSlugIteration(
	ctx context.Context, data SelectSlugIterationData,
) (models.Slug, int, error) {
	tx, err := lib.PostgresContext(ctx)
	if err != nil {
		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("get postgres client: %w", err))
	}

	output := new(struct {
		Slug models.Slug `bun:"slug"`
	})

	reg, err := regexp.CompilePOSIX(`^` + regexp.QuoteMeta(string(data.Slug)) + `-([0-9]+)$`)
	if err != nil {
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

	if err = query.Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return data.Slug + "-1", 1, nil
		}

		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("select slug iteration: %w", err))
	}

	// Capture the index of the last iteration.
	index := 1

	_, err = fmt.Sscanf(string(output.Slug), string(data.Slug)+"-%d", &index)
	if err != nil {
		return "", 0, NewErrSelectSlugIteration(fmt.Errorf("parse slug iteration: %w", err))
	}

	return models.Slug(fmt.Sprintf("%s-%d", data.Slug, index+1)), index + 1, nil
}

func NewSelectSlugIterationRepository() *SelectSlugIterationRepository {
	return &SelectSlugIterationRepository{}
}
