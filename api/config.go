package api

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Config - environment variables
type Config struct {
	Port     string
	MongoURI string
	MongoDB  string
}

// GetLocals - load config variables
func GetLocals(logger *zap.SugaredLogger) *Config {
	if err := godotenv.Load(); err != nil {
		logger.Errorf(".env variables failed to load: %+v", err)
	}

	config := &Config{
		Port:     os.Getenv("PORT"),
		MongoURI: os.Getenv("MONGO_URI"),
		MongoDB:  os.Getenv("MONGO_DB"),
	}

	return config
}

// ProvideDatabase - connect to database
func ProvideDatabase(logger *zap.SugaredLogger, config *Config) *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		logger.Panic(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logger.Panic(err)
	}

	logger.Info("Successfully connected to MongoDB...")

	db := client.Database(config.MongoDB)

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
