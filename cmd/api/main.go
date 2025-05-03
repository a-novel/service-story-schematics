package main

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	authapiclient "github.com/a-novel/service-authentication/api/apiclient"
	apicodegen "github.com/a-novel/service-authentication/api/codegen"
	apimodels "github.com/a-novel/service-authentication/models"

	"github.com/a-novel-kit/context"
	"github.com/a-novel-kit/golm/bindings/groq"
	golmmiddleware "github.com/a-novel-kit/middlewares/golm"
	zeromiddleware "github.com/a-novel-kit/middlewares/zerolog"

	"github.com/a-novel/service-story-schematics/api"
	"github.com/a-novel/service-story-schematics/api/codegen"
	"github.com/a-novel/service-story-schematics/config"
	"github.com/a-novel/service-story-schematics/internal/dao"
	"github.com/a-novel/service-story-schematics/internal/daoai"
	"github.com/a-novel/service-story-schematics/internal/lib"
	"github.com/a-novel/service-story-schematics/internal/services"
)

func main() {
	// =================================================================================================================
	// LOAD DEPENDENCIES (EXTERNAL)
	// =================================================================================================================
	logger := zerolog.New(os.Stderr).With().
		Str("app", "story-schematics").
		Timestamp().
		Logger()

	if config.LoggerColor {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger.Info().Msg("starting application...")

	// Initialize all contexts at once.
	ctx, err := lib.NewAgoraContext(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Msg("initialize agora context")
	}

	chatGroqBinding := groq.New(config.Groq.APIKey, config.Groq.Model)

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
	router.Use(zeromiddleware.ZeroLog(&logger))
	router.Use(golmmiddleware.Golm(chatGroqBinding))
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

	if config.Sentry.DSN != "" {
		sentryOptions := sentry.ClientOptions{
			Dsn: config.Sentry.DSN,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					// Add IP Address to user information.
					event.User.IPAddress = req.RemoteAddr
				}

				return event
			},
		}

		if err = sentry.Init(sentryOptions); err != nil {
			logger.Fatal().Err(err).Msg("initialize sentry")
		}

		sentryHandler := sentryhttp.New(sentryhttp.Options{})
		router.Use(sentryHandler.Handle)
	}

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

	apiPermissions := map[apicodegen.OperationName][]apimodels.Permission{
		codegen.PingOperation:        {},
		codegen.HealthcheckOperation: {},

		codegen.CreateBeatsSheetOperation: {"beats-sheet:create"},
		codegen.CreateLoglineOperation:    {"logline:create"},
		codegen.CreateStoryPlanOperation:  {"story-plan:create"},

		codegen.ExpandBeatOperation:    {"beat:expand"},
		codegen.ExpandLoglineOperation: {"logline:expand"},

		codegen.GenerateBeatsSheetOperation: {"beats-sheet:generate"},
		codegen.GenerateLoglinesOperation:   {"loglines:generate"},

		codegen.GetBeatsSheetsOperation: {"beats-sheets:read"},
		codegen.GetLoglinesOperation:    {"loglines:read"},
		codegen.GetStoryPlansOperation:  {"story-plans:read"},

		codegen.RegenerateBeatsOperation: {"beats:regenerate"},

		codegen.GetBeatsSheetOperation: {"beats-sheet:read"},
		codegen.GetLoglineOperation:    {"logline:read"},
		codegen.GetStoryPlanOperation:  {"story-plan:read"},

		codegen.UpdateStoryPlanOperation: {"story-plan:update"},
	}

	authSecurityService := authapiclient.NewSecurityHandlerService(config.API.ExternalAPIs.Auth)

	securityHandler, err := api.NewSecurity(apiPermissions, config.Permissions, authSecurityService)
	if err != nil {
		logger.Fatal().Err(err).Msg("initialize security handler")
	}

	apiServer, err := codegen.NewServer(handler, securityHandler)
	if err != nil {
		logger.Fatal().Err(err).Msg("initialize server")
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

	logger.Info().
		Str("address", httpServer.Addr).
		Msg("application started!")

	if err = httpServer.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("application stopped")
	}
}
