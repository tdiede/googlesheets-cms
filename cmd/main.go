package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/tdiede/gosheets/pkg"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	gosheet "gopkg.in/Iwark/spreadsheet.v2"
)

const (
	_clientSecretPath = "client_secret.json"
)

// go run cmd/gosheets.go -s 10j3vs-LmZcLSlygaEmunjvSZbysj5oFcNtEfNpTp9ZQ

func main() {
	service, err := initializeService()
	if err != nil {
		return
	}

	spreadsheetID := parseFlags()

	cellUpdate := entity.CellContents{
		x:       2,
		y:       2,
		content: "hello",
	}
	readSpreadsheetContents(service, spreadsheetID)
	updateSpreadsheetContents(service, spreadsheetID, cellUpdate)
}

func initializeService() (*gosheet.Service, error) {
	secret, err := ioutil.ReadFile(_clientSecretPath)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(secret, gosheet.Scope)
	if err != nil {
		return nil, err
	}
	client := conf.Client(context.TODO())
	return gosheet.NewServiceWithClient(client), nil
}

func updateSpreadsheetContents(service *gosheet.Service, spreadsheetID string, update entity.CellContents) error {
	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	spreadsheet.Sheets[1].Update(update.x, 0, "hogehoge")
}

func readSpreadsheetContents(service *gosheet.Service, spreadsheetID string) error {
	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for _, sheet := range spreadsheet.Sheets {
		for _, row := range sheet.Rows {
			for _, cell := range row {
				fmt.Println(cell.Value)
			}
		}
	}

	// Make sure call Synchronize to reflect the changes
	err = spreadsheet.Sheets[1].Synchronize()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

func parseFlags() string {
	s := flag.String("s", "", "Spreadsheet ID from Google Sheets, e.g. https://docs.google.com/spreadsheets/d/<Spreadsheet ID>")
	flag.Parse()
	return *s
}
