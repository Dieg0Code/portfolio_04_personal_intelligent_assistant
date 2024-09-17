package db

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

var (
	supabaseURL = os.Getenv("SUPABASE_URL")
	supabaseKey = os.Getenv("SUPABASE_KEY")
)

func NewDBConnection() (*supabase.Client, error) {
	dbClient, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		logrus.WithError(err).Error("cannot initialize client")
		return nil, fmt.Errorf("cannot initialize client: %v", err)
	}

	logrus.Info("Database connection initialized successfully")
	return dbClient, nil
}
