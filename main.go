package main

import (
	"log"
	"time"

	"nfs-api/agreement"
	"nfs-api/appfeature"
	"nfs-api/auth"
	"nfs-api/checklistnote"
	"nfs-api/config"
	"nfs-api/database"
	"nfs-api/docs"
	"nfs-api/invoice"
	"nfs-api/nlachecklist"
	"nfs-api/nlafile"
	"nfs-api/notary"
	"nfs-api/office"
	"nfs-api/officefeature"
	"nfs-api/officeform"
	"nfs-api/officepayment"
	"nfs-api/officeservice"
	"nfs-api/permission"
	"nfs-api/userpermission"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           NFS API
// @version         1.0
// @description     NFS API Server

// @host            nfs-api.andasy.dev
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := config.Load()

	database.Connect(cfg)
	database.Migrate(
		&office.Office{},
		&auth.User{},
		&permission.Permission{},
		&userpermission.UserPermission{},
		&officeservice.OfficeService{},
		&officeform.OfficeForm{},
		&notary.Notary{},
		&agreement.Agreement{},
		&nlafile.NlaFile{},
		&invoice.Invoice{},
		&invoice.InvoiceService{},
		&appfeature.AppFeature{},
		&officefeature.OfficeFeature{},
		&officepayment.OfficePayment{},
		&nlachecklist.NlaChecklist{},
		&nlachecklist.ChecklistDossier{},
		&checklistnote.ChecklistNote{},
	)

	// Clean up expired tokens from blacklist every hour
	auth.CleanupBlacklist(1 * time.Hour)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "https://nfs.andasy.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Connected"})
	})

	if cfg.AppEnv != "production" {
		docs.SwaggerInfo.Host = "nfs-api-dev.andasy.dev"
		r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	api := r.Group("/api/v1")
	{
		authService := auth.NewService(cfg)
		authHandler := auth.NewHandler(authService)
		authHandler.RegisterRoutes(api)

		protected := api.Group("", auth.JWTMiddleware(cfg))
		authHandler.RegisterProtectedRoutes(protected)

		officeService := office.NewService()
		officeHandler := office.NewHandler(officeService)
		officeHandler.RegisterRoutes(protected)

		permService := permission.NewService()
		permHandler := permission.NewHandler(permService)
		permHandler.RegisterRoutes(protected)

		upService := userpermission.NewService()
		upHandler := userpermission.NewHandler(upService)
		upHandler.RegisterRoutes(protected)

		osSvc := officeservice.NewService()
		osHandler := officeservice.NewHandler(osSvc)
		osHandler.RegisterRoutes(protected)

		ofSvc := officeform.NewService()
		ofHandler := officeform.NewHandler(ofSvc)
		ofHandler.RegisterRoutes(protected)

		notarySvc := notary.NewService()
		notaryHandler := notary.NewHandler(notarySvc)
		notaryHandler.RegisterRoutes(protected)

		agreementSvc := agreement.NewService()
		agreementHandler := agreement.NewHandler(agreementSvc)
		agreementHandler.RegisterRoutes(protected)

		nlaSvc := nlafile.NewService()
		nlaHandler := nlafile.NewHandler(nlaSvc)
		nlaHandler.RegisterRoutes(protected)

		invoiceSvc := invoice.NewService()
		invoiceHandler := invoice.NewHandler(invoiceSvc)
		invoiceHandler.RegisterRoutes(protected)

		appFeatureSvc := appfeature.NewService()
		appFeatureHandler := appfeature.NewHandler(appFeatureSvc)
		appFeatureHandler.RegisterRoutes(protected)

		officeFeatureSvc := officefeature.NewService()
		officeFeatureHandler := officefeature.NewHandler(officeFeatureSvc)
		officeFeatureHandler.RegisterRoutes(protected)

		officePaymentSvc := officepayment.NewService()
		officePaymentHandler := officepayment.NewHandler(officePaymentSvc)
		officePaymentHandler.RegisterRoutes(protected)

		nlaChecklistSvc := nlachecklist.NewService()
		nlaChecklistHandler := nlachecklist.NewHandler(nlaChecklistSvc)
		nlaChecklistHandler.RegisterRoutes(protected)

		checklistNoteSvc := checklistnote.NewService()
		checklistNoteHandler := checklistnote.NewHandler(checklistNoteSvc)
		checklistNoteHandler.RegisterRoutes(protected)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
