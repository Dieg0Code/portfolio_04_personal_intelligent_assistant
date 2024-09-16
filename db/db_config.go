package db

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func NewDBConnection() *supabase.Client {
	api_url := os.Getenv("SUPABASE_URL")
	api_key := os.Getenv("SUPABASE_KEY")

	dbClient, err := supabase.NewClient(api_url, api_key, nil)
	if err != nil {
		logrus.WithError(err).Error("cannot initalize client")
		panic(fmt.Sprintf("cannot initalize client: %v", err))
	}

	return dbClient
}
