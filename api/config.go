package api

import (
	"context"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LoadVariables - load config variables
func LoadVariables(logger *zap.SugaredLogger) {
	if err := godotenv.Load(); err != nil {
		logger.Errorf(".env variables failed to load (only applicable in a development environment): %+v", err)
	}
}

// ProvideDatabase - connect to database
func ProvideDatabase(logger *zap.SugaredLogger) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		logger.Panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Panic(err)
	}

	logger.Info("Successfully connected to MongoDB...")

	db := client.Database(os.Getenv("MONGO_DB"))

	return db
}

// ProvideLogger to fx
func ProvideLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	logger, _ := config.Build()
	slogger := logger.Sugar()

	return slogger
}

// LoggerFXModule provided to fx
var LoggerFXModule = fx.Options(
	fx.Provide(ProvideLogger),
)
