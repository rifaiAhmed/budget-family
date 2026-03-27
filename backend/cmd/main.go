package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"budget-family/internal/config"
	"budget-family/internal/entity"
	"budget-family/internal/handler"
	"budget-family/internal/middleware"
	"budget-family/internal/repository"
	"budget-family/internal/service"
	"budget-family/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	logger := config.NewLogger(cfg)
	defer func() { _ = logger.Sync() }()

	db, err := config.NewDatabase(cfg)
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	if err := config.RunMigrations(cfg, logger); err != nil {
		logger.Fatal("failed to run migrations", zap.Error(err))
	}

	if err := seedIfEmpty(db, logger); err != nil {
		logger.Fatal("failed to seed initial data", zap.Error(err))
	}

	rdb, err := config.NewRedis(cfg)
	if err != nil {
		logger.Fatal("failed to connect redis", zap.Error(err))
	}
	_ = rdb

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.ZapLogger(logger))
	r.Use(middleware.CORS(middleware.CORSConfig{AllowedOrigins: cfg.CORS.AllowedOrigins}))
	r.Use(gin.Recovery())

	validator := utils.NewValidator()

	userRepo := repository.NewUserRepository(db)
	familyRepo := repository.NewFamilyRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)
	goalRepo := repository.NewGoalRepository(db)
	billRepo := repository.NewBillRepository(db)
	invitationRepo := repository.NewInvitationRepository(db)

	jwtManager := utils.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.AccessTokenTTL, cfg.Auth.RefreshTokenTTL)

	authSvc := service.NewAuthService(cfg, logger, userRepo, jwtManager)
	familySvc := service.NewFamilyService(cfg, logger, familyRepo, invitationRepo)
	walletSvc := service.NewWalletService(cfg, logger, walletRepo)
	categorySvc := service.NewCategoryService(cfg, logger, categoryRepo)
	transactionSvc := service.NewTransactionService(cfg, logger, db, transactionRepo, walletRepo, budgetRepo)
	budgetSvc := service.NewBudgetService(cfg, logger, budgetRepo, transactionRepo)
	goalSvc := service.NewGoalService(cfg, logger, goalRepo)
	billSvc := service.NewBillService(cfg, logger, billRepo)

	jwtMW := middleware.NewJWTMiddleware(jwtManager, cfg.Auth.Issuer)

	authH := handler.NewAuthHandler(logger, validator, authSvc)
	familyH := handler.NewFamilyHandler(logger, validator, familySvc)
	walletH := handler.NewWalletHandler(logger, validator, walletSvc)
	categoryH := handler.NewCategoryHandler(logger, validator, categorySvc)
	transactionH := handler.NewTransactionHandler(logger, validator, transactionSvc)
	budgetH := handler.NewBudgetHandler(logger, validator, budgetSvc)
	goalH := handler.NewGoalHandler(logger, validator, goalSvc)
	billH := handler.NewBillHandler(logger, validator, billSvc)

	r.GET("/health", func(c *gin.Context) { utils.Success(c, http.StatusOK, "ok", gin.H{"status": "up"}) })

	auth := r.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.GET("/me", jwtMW.RequireAuth(), authH.Me)
	}

	api := r.Group("/")
	api.Use(jwtMW.RequireAuth())
	{
		api.POST("/families", familyH.Create)
		api.GET("/families", familyH.List)
		api.POST("/families/invite", familyH.Invite)

		api.GET("/wallets", walletH.List)
		api.POST("/wallets", walletH.Create)
		api.PUT("/wallets/:id", walletH.Update)
		api.DELETE("/wallets/:id", walletH.Delete)

		api.GET("/categories", categoryH.List)
		api.POST("/categories", categoryH.Create)

		api.POST("/transactions", transactionH.Create)
		api.GET("/transactions", transactionH.List)
		api.GET("/transactions/summary", transactionH.Summary)
		api.GET("/transactions/:id", transactionH.Get)
		api.PUT("/transactions/:id", transactionH.Update)
		api.DELETE("/transactions/:id", transactionH.Delete)

		api.POST("/budgets", budgetH.Create)
		api.GET("/budgets", budgetH.List)
		api.GET("/budgets/usage", budgetH.Usage)

		api.POST("/goals", goalH.Create)
		api.GET("/goals", goalH.List)
		api.PUT("/goals/:id", goalH.Update)

		api.POST("/bills", billH.Create)
		api.GET("/bills", billH.List)
	}

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server started", zap.String("addr", cfg.Server.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown failed", zap.Error(err))
	}
}

func seedIfEmpty(db *gorm.DB, logger *zap.Logger) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var userCount int64
		if err := tx.Model(&entity.User{}).Count(&userCount).Error; err != nil {
			return err
		}
		if userCount == 0 {
			return nil
		}

		var familyCount int64
		if err := tx.Model(&entity.Family{}).Count(&familyCount).Error; err != nil {
			return err
		}

		var fam entity.Family
		seededFamily := false
		if familyCount == 0 {
			var firstUser entity.User
			if err := tx.Order("created_at ASC").First(&firstUser).Error; err != nil {
				return err
			}

			familyName := strings.TrimSpace(firstUser.Name)
			if familyName == "" {
				familyName = "My Family"
			}

			fam = entity.Family{ID: uuid.New(), Name: familyName, OwnerID: firstUser.ID}
			owner := &entity.FamilyMember{ID: uuid.New(), FamilyID: fam.ID, UserID: firstUser.ID, Role: "owner"}
			if err := tx.Create(&fam).Error; err != nil {
				return err
			}
			if err := tx.Create(owner).Error; err != nil {
				return err
			}
			seededFamily = true
		} else {
			if err := tx.Order("created_at ASC").First(&fam).Error; err != nil {
				return err
			}
		}

		seededWallets := false
		var walletCount int64
		if err := tx.Model(&entity.Wallet{}).Where("family_id = ?", fam.ID).Count(&walletCount).Error; err != nil {
			return err
		}
		if walletCount == 0 {
			wallets := []*entity.Wallet{
				{ID: uuid.New(), FamilyID: fam.ID, Name: "Cash", Type: "cash", Balance: decimal.NewFromInt(0)},
				{ID: uuid.New(), FamilyID: fam.ID, Name: "Bank", Type: "bank", Balance: decimal.NewFromInt(0)},
			}
			if err := tx.Create(&wallets).Error; err != nil {
				return err
			}
			seededWallets = true
		}

		seededCategories := false
		var categoryCount int64
		if err := tx.Model(&entity.Category{}).Where("family_id = ?", fam.ID).Count(&categoryCount).Error; err != nil {
			return err
		}
		if categoryCount == 0 {
			cats := []*entity.Category{
				{ID: uuid.New(), FamilyID: fam.ID, Name: "Salary", Type: "income", Icon: "payments"},
				{ID: uuid.New(), FamilyID: fam.ID, Name: "Food", Type: "expense", Icon: "restaurant"},
				{ID: uuid.New(), FamilyID: fam.ID, Name: "Transport", Type: "expense", Icon: "directions_car"},
				{ID: uuid.New(), FamilyID: fam.ID, Name: "Bills", Type: "expense", Icon: "receipt_long"},
			}
			if err := tx.Create(&cats).Error; err != nil {
				return err
			}
			seededCategories = true
		}

		if seededFamily || seededWallets || seededCategories {
			logger.Info(
				"seeded initial data",
				zap.String("family_id", fam.ID.String()),
				zap.Bool("family", seededFamily),
				zap.Bool("wallets", seededWallets),
				zap.Bool("categories", seededCategories),
			)
		}

		return nil
	})
}
