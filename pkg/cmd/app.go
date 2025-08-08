package cmdpkg

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/a-novel/golib/otel"
	"github.com/a-novel/golib/postgres"
	authmodels "github.com/a-novel/service-authentication/models"
	jkconfig "github.com/a-novel/service-json-keys/models/config"
	jkpkg "github.com/a-novel/service-json-keys/pkg"

	"github.com/a-novel/service-story-schematics/internal/api"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/services"
	apimodels "github.com/a-novel/service-story-schematics/models/api"
	"github.com/a-novel/service-story-schematics/models/config"
)

func App[Otel otel.Config, Pg postgres.Config](
	ctx context.Context, config config.App[Otel, Pg],
) error {
	// =================================================================================================================
	// DEPENDENCIES
	// =================================================================================================================
	otel.SetAppName(config.App.Name)

	err := otel.InitOtel(config.Otel)
	if err != nil {
		return fmt.Errorf("init otel: %w", err)
	}
	defer config.Otel.Flush()

	// Don't override the context if it already has a bun.IDB
	_, err = postgres.GetContext(ctx)
	if err != nil {
		ctx, err = postgres.NewContext(ctx, config.Postgres)
		if err != nil {
			return fmt.Errorf("init postgres: %w", err)
		}
	}

	jkClient, err := jkpkg.NewAPIClient(ctx, config.DependenciesConfig.JSONKeysURL)
	if err != nil {
		return fmt.Errorf("create JSON keys client: %w", err)
	}

	accessTokenVerifier, err := jkpkg.NewClaimsVerifier[authmodels.AccessTokenClaims](
		jkClient,
		jkconfig.JWKSPresetDefault,
	)
	if err != nil {
		return fmt.Errorf("create access token verifier: %w", err)
	}

	// =================================================================================================================
	// DAO
	// =================================================================================================================

	selectSlugIterationDAO := dao.NewSelectSlugIterationRepository()

	insertBeatsSheetDAO := dao.NewInsertBeatsSheetRepository()
	insertLoglineDAO := dao.NewInsertLoglineRepository()
	listBeatsSheetsDAO := dao.NewListBeatsSheetsRepository()
	listLoglinesDAO := dao.NewListLoglinesRepository()
	selectBeatsSheetDAO := dao.NewSelectBeatsSheetRepository()
	selectLoglineDAO := dao.NewSelectLoglineRepository()
	selectLoglineBySlugDAO := dao.NewSelectLoglineBySlugRepository()

	expandBeatDAO := daoai.NewExpandBeatRepository(&config.OpenAI)
	expandLoglineDAO := daoai.NewExpandLoglineRepository(&config.OpenAI)
	generateBeatsSheetDAO := daoai.NewGenerateBeatsSheetRepository(&config.OpenAI)
	generateLoglinesDAO := daoai.NewGenerateLoglinesRepository(&config.OpenAI)
	regenerateBeatsDAO := daoai.NewRegenerateBeatsRepository(&config.OpenAI)

	// =================================================================================================================
	// SERVICES
	// =================================================================================================================

	selectStoryPlanService := services.NewSelectStoryPlanService()

	createBeatsSheetService := services.NewCreateBeatsSheetService(
		services.NewCreateBeatsSheetServiceSource(
			insertBeatsSheetDAO,
			selectStoryPlanService,
			selectLoglineDAO,
		),
	)
	createLoglineService := services.NewCreateLoglineService(
		services.NewCreateLoglineServiceSource(
			insertLoglineDAO,
			selectSlugIterationDAO,
		),
	)
	expandBeatService := services.NewExpandBeatService(
		services.NewExpandBeatServiceSource(
			expandBeatDAO,
			selectBeatsSheetDAO,
			selectLoglineDAO,
			selectStoryPlanService,
		),
	)
	expandLoglineService := services.NewExpandLoglineService(expandLoglineDAO)
	generateBeatsSheetService := services.NewGenerateBeatsSheetService(
		services.NewGenerateBeatsSheetServiceSource(
			generateBeatsSheetDAO,
			selectLoglineDAO,
			selectStoryPlanService,
		),
	)
	generateLoglinesService := services.NewGenerateLoglinesService(generateLoglinesDAO)
	listBeatsSheetsService := services.NewListBeatsSheetsService(
		services.NewListBeatsSheetsServiceSource(
			listBeatsSheetsDAO,
			selectLoglineDAO,
		),
	)
	listLoglinesService := services.NewListLoglinesService(listLoglinesDAO)
	regenerateBeatsService := services.NewRegenerateBeatsService(
		services.NewRegenerateBeatsServiceSource(
			regenerateBeatsDAO,
			selectBeatsSheetDAO,
			selectLoglineDAO,
			selectStoryPlanService,
		),
	)
	selectBeatsSheetService := services.NewSelectBeatsSheetService(
		services.NewSelectBeatsSheetServiceSource(
			selectBeatsSheetDAO,
			selectLoglineDAO,
		),
	)
	selectLoglineService := services.NewSelectLoglineService(
		services.NewSelectLoglineServiceSource(
			selectLoglineDAO,
			selectLoglineBySlugDAO,
		),
	)

	// =================================================================================================================
	// SETUP ROUTER
	// =================================================================================================================

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(config.API.Timeouts.Request))
	router.Use(middleware.RequestSize(config.API.MaxRequestSize))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.API.Cors.AllowedOrigins,
		AllowedHeaders:   config.API.Cors.AllowedHeaders,
		AllowCredentials: config.API.Cors.AllowCredentials,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		MaxAge: config.API.Cors.MaxAge,
	}))
	router.Use(config.Otel.HTTPHandler())

	handler := &api.API{
		CreateBeatsSheetService: createBeatsSheetService,
		CreateLoglineService:    createLoglineService,

		ExpandBeatService:    expandBeatService,
		ExpandLoglineService: expandLoglineService,

		GenerateBeatsSheetService: generateBeatsSheetService,
		GenerateLoglinesService:   generateLoglinesService,

		ListBeatsSheetsService: listBeatsSheetsService,
		ListLoglinesService:    listLoglinesService,

		RegenerateBeatsService: regenerateBeatsService,

		SelectBeatsSheetService: selectBeatsSheetService,
		SelectLoglineService:    selectLoglineService,

		JKClient:     jkClient,
		OpenAIClient: &config.OpenAI,
	}

	securityHandler, err := api.NewSecurity(accessTokenVerifier, config.PermissionsConfig)
	if err != nil {
		return fmt.Errorf("create security handler: %w", err)
	}

	apiServer, err := apimodels.NewServer(handler, securityHandler)
	if err != nil {
		return fmt.Errorf("new api server: %w", err)
	}

	router.Mount("/v1/", http.StripPrefix("/v1", apiServer))

	httpServer := &http.Server{
		Addr:              ":" + strconv.Itoa(config.API.Port),
		Handler:           router,
		ReadTimeout:       config.API.Timeouts.Read,
		ReadHeaderTimeout: config.API.Timeouts.ReadHeader,
		WriteTimeout:      config.API.Timeouts.Write,
		IdleTimeout:       config.API.Timeouts.Idle,
		BaseContext:       func(_ net.Listener) context.Context { return ctx },
	}

	err = httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}
