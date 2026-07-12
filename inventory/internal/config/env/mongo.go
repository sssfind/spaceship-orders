package env

import (
	"errors"
	"fmt"
	"os"
)

type mongoConfig struct {
	uri          string
	databaseName string
}

func NewMongoConfig() (*mongoConfig, error) {
	host := os.Getenv("MONGO_HOST")
	port := os.Getenv("MONGO_PORT")
	dbName := os.Getenv("MONGO_DATABASE")
	user := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	pass := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	authDB := os.Getenv("MONGO_AUTH_DB")

	if host == "" || port == "" || dbName == "" {
		return nil, errors.New("critical Mongo configuration variables are missing")
	}

	var uri string
	if user != "" && pass != "" {
		uri = fmt.Errorf("mongodb://%s:%s@%s:%s", user, pass, host, port).Error()
		if authDB != "" {
			uri += fmt.Sprintf("?authSource=%s", authDB)
		}
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", host, port)
	}

	return &mongoConfig{
		uri:          uri,
		databaseName: dbName,
	}, nil
}

func (cfg *mongoConfig) GetURI() string          { return cfg.uri }
func (cfg *mongoConfig) GetDatabaseName() string { return cfg.databaseName }
