package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/getsentry/sentry-go/attribute"
	sentryhttp "github.com/getsentry/sentry-go/http"
	sentryotel "github.com/getsentry/sentry-go/otel"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	authapiclient "github.com/a-novel/service-authentication/api/apiclient"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/internal/services"
)

const SentryFlushTimeout = 2 * time.Second

func main() {
	ctx := context.Background()
	// =================================================================================================================
	// LOAD DEPENDENCIES (EXTERNAL)
	// =================================================================================================================
	err := sentry.Init(config.SentryClient)
	if err != nil {
		log.Fatalf("initialize sentry: %v", err)
	}
	defer sentry.Flush(SentryFlushTimeout)

	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sentryotel.NewSentrySpanProcessor()))
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(sentryotel.NewSentryPropagator())

	logger := sentry.NewLogger(ctx)
	logger.SetAttributes(
		attribute.String("app", "agora"),
		attribute.String("service", "authentication"),
	)

	logger.Info(ctx, "starting application")

	// Initialize all contexts at once.
	ctx, err = lib.NewAgoraContext(context.Background(), config.DSN)
	if err != nil {
		logger.Fatalf(ctx, "initialize agora context: %v", err)
	}

	// =================================================================================================================
	// LOAD REPOSITORIES (INTERNAL)
	// =================================================================================================================

	// REPOSITORIES ----------------------------------------------------------------------------------------------------

	selectSlugIterationDAO := dao.NewSelectSlugIterationRepository()

	insertBeatsSheetDAO := dao.NewInsertBeatsSheetRepository()
	insertLoglineDAO := dao.NewInsertLoglineRepository()
	insertStoryPlanDAO := dao.NewInsertStoryPlanRepository()
	listBeatsSheetsDAO := dao.NewListBeatsSheetsRepository()
	listLoglinesDAO := dao.NewListLoglinesRepository()
	listStoryPlansDAO := dao.NewListStoryPlansRepository()
	selectBeatsSheetDAO := dao.NewSelectBeatsSheetRepository()
	selectLoglineDAO := dao.NewSelectLoglineRepository()
	selectLoglineBySlugDAO := dao.NewSelectLoglineBySlugRepository()
	selectStoryPlanDAO := dao.NewSelectStoryPlanRepository()
	selectStoryPlanBySlugDAO := dao.NewSelectStoryPlanBySlugRepository()
	updateStoryPlanDAO := dao.NewUpdateStoryPlanRepository()

	expandBeatDAO := daoai.NewExpandBeatRepository()
	expandLoglineDAO := daoai.NewExpandLoglineRepository()
	generateBeatsSheetDAO := daoai.NewGenerateBeatsSheetRepository()
	generateLoglinesDAO := daoai.NewGenerateLoglinesRepository()
	regenerateBeatsDAO := daoai.NewRegenerateBeatsRepository()

	// SERVICES --------------------------------------------------------------------------------------------------------

	createBeatsSheetService := services.NewCreateBeatsSheetService(
		services.NewCreateBeatsSheetServiceSource(
			insertBeatsSheetDAO,
			selectStoryPlanDAO,
			selectLoglineDAO,
		),
	)
	createLoglineService := services.NewCreateLoglineService(
		services.NewCreateLoglineServiceSource(
			insertLoglineDAO,
			selectSlugIterationDAO,
		),
	)
	createStoryPlanService := services.NewCreateStoryPlanService(
		services.NewCreateStoryPlanServiceSource(
			insertStoryPlanDAO,
			selectSlugIterationDAO,
		),
	)
	expandBeatService := services.NewExpandBeatService(
		services.NewExpandBeatServiceSource(
			expandBeatDAO,
			selectBeatsSheetDAO,
			selectLoglineDAO,
			selectStoryPlanDAO,
		),
	)
	expandLoglineService := services.NewExpandLoglineService(expandLoglineDAO)
	generateBeatsSheetService := services.NewGenerateBeatsSheetService(
		services.NewGenerateBeatsSheetServiceSource(
			generateBeatsSheetDAO,
			selectLoglineDAO,
			selectStoryPlanDAO,
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
	listStoryPlansService := services.NewListStoryPlansService(listStoryPlansDAO)
	regenerateBeatsService := services.NewRegenerateBeatsService(
		services.NewRegenerateBeatsServiceSource(
			regenerateBeatsDAO,
			selectBeatsSheetDAO,
			selectLoglineDAO,
			selectStoryPlanDAO,
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
	selectStoryPlanService := services.NewSelectStoryPlanService(
		services.NewSelectStoryPlanServiceSource(
			selectStoryPlanDAO,
			selectStoryPlanBySlugDAO,
		),
	)
	updateStoryPlanService := services.NewUpdateStoryPlanService(updateStoryPlanDAO)

	// =================================================================================================================
	// SETUP ROUTER
	// =================================================================================================================

	router := chi.NewRouter()

	// MIDDLEWARES -----------------------------------------------------------------------------------------------------

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(config.API.Timeouts.Request))
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

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	router.Use(sentryHandler.Handle)

	// RUN -------------------------------------------------------------------------------------------------------------

	handler := &api.API{
		CreateBeatsSheetService: createBeatsSheetService,
		CreateLoglineService:    createLoglineService,
		CreateStoryPlanService:  createStoryPlanService,

		ExpandBeatService:    expandBeatService,
		ExpandLoglineService: expandLoglineService,

		GenerateBeatsSheetService: generateBeatsSheetService,
		GenerateLoglinesService:   generateLoglinesService,

		ListBeatsSheetsService: listBeatsSheetsService,
		ListLoglinesService:    listLoglinesService,
		ListStoryPlansService:  listStoryPlansService,

		RegenerateBeatsService: regenerateBeatsService,

		SelectBeatsSheetService: selectBeatsSheetService,
		SelectLoglineService:    selectLoglineService,
		SelectStoryPlanService:  selectStoryPlanService,

		UpdateStoryPlanService: updateStoryPlanService,
	}

	authSecurityService := authapiclient.NewSecurityHandlerService(config.API.ExternalAPIs.Auth)

	securityHandler, err := api.NewSecurity(config.Permissions, authSecurityService)
	if err != nil {
		logger.Fatalf(ctx, "start security handler: %v", err)
	}

	apiServer, err := codegen.NewServer(handler, securityHandler)
	if err != nil {
		logger.Fatalf(ctx, "start server: %v", err)
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

	logger.SetAttributes(attribute.Int("server.port", config.API.Port))
	logger.Info(ctx, "start http server")

	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Fatalf(ctx, "start http server: %v", err)
	}
}
