package gsheet

import (
	"google.golang.org/api/sheets/v4"
)

type GsheetAdapter struct {
	service         *sheets.Service
	CredentialsPath string
}

func NewGsheetAdapter(credentialsPath string) (*GsheetAdapter, error) {
	client := &GsheetAdapter{
		CredentialsPath: credentialsPath,
	}

	
	return client, nil
}
