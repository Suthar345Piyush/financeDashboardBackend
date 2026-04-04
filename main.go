// main file

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/suthar345Piyush/financedashboard/config"
	"github.com/suthar345Piyush/financedashboard/db"
	"github.com/suthar345Piyush/financedashboard/internal/handler"
	"github.com/suthar345Piyush/financedashboard/internal/models"
	"github.com/suthar345Piyush/financedashboard/internal/repository"
	"github.com/suthar345Piyush/financedashboard/internal/service"
	"github.com/suthar345Piyush/financedashboard/middleware"
	jwtpkg "github.com/suthar345Piyush/financedashboard/pkg/jwt"
)

func main() {
	cfg := config.Load()

	// database

	database, err := db.Connect(cfg.DBPath)

	if err != nil {
		log.Fatal("DB connect:", err)
	}

	defer database.Close()

	if err := db.RunMigrations(database, "./db/migrations"); err != nil {
		log.Fatal("Migrations:", err)
	}

	// jwt setup

	jwtManager := jwtpkg.NewManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	// repositories part

	userRepo := repository.NewUserRepository(database)
	recordRepo := repository.NewRecordRepository(database)
	dashboardRepo := repository.NewDashboardRepository(database)

	// services section

	authSvc := service.NewAuthService(userRepo, jwtManager)
	userSvc := service.NewUserService(userRepo)
	recordSvc := service.NewRecordService(recordRepo)
	dashboardSvc := service.NewDashboardService(dashboardRepo)

	// handler function

	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	recordH := handler.NewRecordHandler(recordSvc)
	dashboardH := handler.NewDashboardHandler(dashboardSvc)

	// setting up chi middleware

	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// routes (public)

	r.Post("/auth/register", authH.Register)
	r.Post("/auth/login", authH.Login)

	// routes (protected)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Authenticate(jwtManager))

		//users routes -  for only admin

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole(models.RoleAdmin))
			r.Get("/users", userH.List)
			r.Get("/users/{id}", userH.GetByID)
			r.Put("/users/{id}/role", userH.UpdateRole)
			r.Put("/users/{id}/status", userH.UpdateStatus)
		})

		// records routes

		r.Group(func(r chi.Router) {
			// analyst and admin can read

			r.With(middleware.RequireMinRole(models.RoleAnalyst)).Get("/records", recordH.List)
			r.With(middleware.RequireMinRole(models.RoleAnalyst)).Get("/records/{id}", recordH.GetByID)

			// only admin can write

			r.With(middleware.RequireRole(models.RoleAdmin)).Post("/records", recordH.Create)
			r.With(middleware.RequireRole(models.RoleAdmin)).Put("/records/{id}", recordH.Update)
			r.With(middleware.RequireRole(models.RoleAdmin)).Delete("/records/{id}", recordH.Delete)
		})

		// dashboard - all authenticated users

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireMinRole(models.RoleViewer))
			r.Get("/dashboard/summary", dashboardH.Summary)
			r.Get("/dashboard/categories", dashboardH.CategoryTotals)

			// monthly trends only for analyst

			r.With(middleware.RequireMinRole(models.RoleAnalyst)).Get("/dashboard/trends", dashboardH.MonthlyTrends)
			r.With(middleware.RequireMinRole(models.RoleAnalyst)).Get("/dashboard/recent", dashboardH.RecentActivity)

		})
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r); err != nil {
		log.Fatal(err)
	}

}
