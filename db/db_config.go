package db

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/supabase-community/supabase-go"
)

func NewDBConnection(supabaseURL, supabaseKey string) (*supabase.Client, error) {
	dbClient, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		logrus.WithError(err).Error("cannot initialize client")
		return nil, fmt.Errorf("cannot initialize client: %v", err)
	}

	logrus.Info("Database connection initialized successfully")
	return dbClient, nil
}
