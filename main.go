package main

import (
	"log"

	"nfs-api/agreement"
	"nfs-api/appfeature"
	"nfs-api/auth"
	"nfs-api/checklistnote"
	"nfs-api/config"
	"nfs-api/database"
	_ "nfs-api/docs"
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

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           NFS API
// @version         1.0
// @description     NFS API Server

// @host            localhost:8080
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

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		authService := auth.NewService(cfg)
		authHandler := auth.NewHandler(authService)
		authHandler.RegisterRoutes(api)

		officeService := office.NewService()
		officeHandler := office.NewHandler(officeService)
		officeHandler.RegisterRoutes(api)

		permService := permission.NewService()
		permHandler := permission.NewHandler(permService)
		permHandler.RegisterRoutes(api)

		upService := userpermission.NewService()
		upHandler := userpermission.NewHandler(upService)
		upHandler.RegisterRoutes(api)

		osSvc := officeservice.NewService()
		osHandler := officeservice.NewHandler(osSvc)
		osHandler.RegisterRoutes(api)

		ofSvc := officeform.NewService()
		ofHandler := officeform.NewHandler(ofSvc)
		ofHandler.RegisterRoutes(api)

		notarySvc := notary.NewService()
		notaryHandler := notary.NewHandler(notarySvc)
		notaryHandler.RegisterRoutes(api)

		agreementSvc := agreement.NewService()
		agreementHandler := agreement.NewHandler(agreementSvc)
		agreementHandler.RegisterRoutes(api)

		nlaSvc := nlafile.NewService()
		nlaHandler := nlafile.NewHandler(nlaSvc)
		nlaHandler.RegisterRoutes(api)

		invoiceSvc := invoice.NewService()
		invoiceHandler := invoice.NewHandler(invoiceSvc)
		invoiceHandler.RegisterRoutes(api)

		appFeatureSvc := appfeature.NewService()
		appFeatureHandler := appfeature.NewHandler(appFeatureSvc)
		appFeatureHandler.RegisterRoutes(api)

		officeFeatureSvc := officefeature.NewService()
		officeFeatureHandler := officefeature.NewHandler(officeFeatureSvc)
		officeFeatureHandler.RegisterRoutes(api)

		officePaymentSvc := officepayment.NewService()
		officePaymentHandler := officepayment.NewHandler(officePaymentSvc)
		officePaymentHandler.RegisterRoutes(api)

		nlaChecklistSvc := nlachecklist.NewService()
		nlaChecklistHandler := nlachecklist.NewHandler(nlaChecklistSvc)
		nlaChecklistHandler.RegisterRoutes(api)

		checklistNoteSvc := checklistnote.NewService()
		checklistNoteHandler := checklistnote.NewHandler(checklistNoteSvc)
		checklistNoteHandler.RegisterRoutes(api)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
