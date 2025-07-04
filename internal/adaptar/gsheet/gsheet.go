package gsheet

import (
	"context"
	"fmt"

	log "Mehul-Kumar-27/dbporter/logger"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	scope = []string{
		"https://spreadsheets.google.com/feeds",
		"https://www.googleapis.com/auth/drive",
	}
)

type GsheetAdapter struct {
	service         *sheets.Service
	CredentialsPath string
	logger          *log.Logger
	Config          GsheetConfig
}

func (g *GsheetAdapter) initializeGsheetAdapter() error {
	ctx := context.Background()
	service, err := sheets.NewService(ctx, option.WithCredentialsFile(g.CredentialsPath), option.WithScopes(scope...))
	if err != nil {
		return fmt.Errorf("error creating gsheet client %s", err)
	}

	g.service = service
	return nil
}

func NewGsheetAdapter(credentialsPath string) (*GsheetAdapter, error) {
	l := log.New("GSHEET", nil)
	client := &GsheetAdapter{
		CredentialsPath: credentialsPath,
		logger:          l,
	}

	err := client.initializeGsheetAdapter()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (g *GsheetAdapter) GetData() error {
	spreadsheet, err := g.service.Spreadsheets.Get(g.Config.SheetID).Do()
	if err != nil {
		return fmt.Errorf("error getting spreadsheet %s", err)
	}

	g.logger.Info("Spreadsheet retrieved successfully", "spreadsheet", spreadsheet)
	return nil
}
