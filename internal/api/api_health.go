package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-faster/jx"
	"github.com/uptrace/bun"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"
	jkApiModels "github.com/a-novel/service-json-keys/models/api"

	"github.com/a-novel/service-story-schematics/models/api"
)

func (api *API) Ping(_ context.Context) (apimodels.PingRes, error) {
	return &apimodels.PingOK{Data: strings.NewReader("pong")}, nil
}

func (api *API) reportPostgres(ctx context.Context) apimodels.Dependency {
	ctx, span := otel.Tracer().Start(ctx, "api.reportPostgres")
	defer span.End()

	pg, err := postgres.GetContext(ctx)
	if err != nil {
		_ = otel.ReportError(span, err)

		return apimodels.Dependency{
			Name:   "postgres",
			Status: apimodels.DependencyStatusDown,
		}
	}

	pgdb, ok := pg.(*bun.DB)
	if !ok {
		_ = otel.ReportError(span, fmt.Errorf("retrieve postgres context: invalid type %T", pg))

		return apimodels.Dependency{
			Name:   "postgres",
			Status: apimodels.DependencyStatusDown,
		}
	}

	err = pgdb.Ping()
	if err != nil {
		_ = otel.ReportError(span, err)

		return apimodels.Dependency{
			Name:   "postgres",
			Status: apimodels.DependencyStatusDown,
		}
	}

	otel.ReportSuccessNoContent(span)

	return apimodels.Dependency{
		Name:   "postgres",
		Status: apimodels.DependencyStatusUp,
	}
}

func (api *API) reportJSONKeys(ctx context.Context) apimodels.Dependency {
	ctx, span := otel.Tracer().Start(ctx, "api.reportJSONKeys")
	defer span.End()

	rawRes, err := api.JKClient.Ping(ctx)
	if err != nil {
		_ = otel.ReportError(span, err)

		return apimodels.Dependency{
			Name:   "json-keys",
			Status: apimodels.DependencyStatusDown,
		}
	}

	_, ok := rawRes.(*jkApiModels.PingOK)
	if !ok {
		_ = otel.ReportError(span, fmt.Errorf("ping JSON keys: unexpected response type %T", rawRes))

		return apimodels.Dependency{
			Name:   "json-keys",
			Status: apimodels.DependencyStatusDown,
		}
	}

	otel.ReportSuccessNoContent(span)

	return apimodels.Dependency{
		Name:   "json-keys",
		Status: apimodels.DependencyStatusUp,
	}
}

func (api *API) reportOpenAI(ctx context.Context) apimodels.Dependency {
	ctx, span := otel.Tracer().Start(ctx, "api.reportOpenAI")
	defer span.End()

	_, err := api.OpenAIClient.Client().Models.Get(ctx, api.OpenAIClient.Model)
	if err != nil {
		_ = otel.ReportError(span, err)

		return apimodels.Dependency{
			Name:   "openai",
			Status: apimodels.DependencyStatusDown,
		}
	}

	otel.ReportSuccessNoContent(span)

	return apimodels.Dependency{
		Name:   "openai",
		Status: apimodels.DependencyStatusUp,
		AdditionalProps: map[string]jx.Raw{
			"base_url": []byte(strconv.Quote(api.OpenAIClient.BaseURL)),
			"model":    []byte(strconv.Quote(api.OpenAIClient.Model)),
		},
	}
}

func (api *API) Healthcheck(ctx context.Context) (apimodels.HealthcheckRes, error) {
	return &apimodels.Health{
		Postgres: api.reportPostgres(ctx),
		JsonKeys: api.reportJSONKeys(ctx),
		Openai:   api.reportOpenAI(ctx),
	}, nil
}
