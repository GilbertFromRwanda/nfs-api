package main

import (
	"log"

	"nfs-api/agreement"
	"nfs-api/auth"
	"nfs-api/config"
	"nfs-api/database"
	"nfs-api/nlafile"
	"nfs-api/notary"
	"nfs-api/office"
	"nfs-api/officeform"
	"nfs-api/officeservice"
	"nfs-api/permission"
	"nfs-api/userpermission"

	"github.com/gin-gonic/gin"
)

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
	)

	r := gin.Default()

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
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
