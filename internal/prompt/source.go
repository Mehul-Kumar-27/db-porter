package prompt

import (
	"github.com/AlecAivazis/survey/v2"
)

func SelectSourceType() (string, error) {
	var sourceType string
	prompt := &survey.Select{
		Message: "Select the source type",
		Options: []string{"gsheet", "snowflake"},
	}

	err := survey.AskOne(prompt, &sourceType)
	if err != nil {
		return "", err
	}

	return sourceType, nil
}
