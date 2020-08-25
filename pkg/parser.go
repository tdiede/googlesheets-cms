package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	gosheet "gopkg.in/Iwark/spreadsheet.v2"
)

//CLI is the NewServiceWithClient instance
type CLI struct {
	Service gosheet.Service
}

//Sheet defines a sheet in a spreadsheet
type Sheet *gosheet.Sheet

//Cell defines grid location and value
type Cell *gosheet.Cell

//ContentUpdate represents a single cell update value
type ContentUpdate struct {
	Row        int
	Column     int
	Value      string
	SheetTitle string
}

//SpreadsheetData is all data from spreadsheet
type SpreadsheetData []SheetData

//SheetData is all data structured from sheet
type SheetData []sheetRecord
type sheetRecord map[string]string

//NewCLI instantiates a service to connect with client secret
func NewCLI(secretPath string) (*CLI, error) {
	secret, err := ioutil.ReadFile(secretPath)
	if err != nil {
		return &CLI{}, errors.Errorf("%v", err)
	}
	conf, err := google.JWTConfigFromJSON(secret, gosheet.Scope)
	if err != nil {
		return &CLI{}, errors.Errorf("%v", err)
	}
	client := conf.Client(context.TODO())
	service := gosheet.NewServiceWithClient(client)

	fmt.Println("Instantiated service!")

	return &CLI{Service: *service}, nil
}

//GetSpreadsheet fetches the Google Spreadsheet with spreadsheetID
func (s *CLI) GetSpreadsheet(spreadsheetID string) (*gosheet.Spreadsheet, error) {
	spreadsheet, err := s.Service.FetchSpreadsheet(spreadsheetID)
	if err != nil {
		return &gosheet.Spreadsheet{}, errors.Errorf("error: %v", err)
	}

	return &spreadsheet, nil
}

//UpdateSpreadsheet wrapper
func (s *CLI) UpdateSpreadsheet(spreadsheetID string, update ContentUpdate) (bool, error) {
	spreadsheet, err := s.GetSpreadsheet(spreadsheetID)
	if err != nil {
		return false, errors.Errorf("error: %v", err)
	}
	sheet, err := spreadsheet.SheetByTitle(update.SheetTitle)
	if err != nil {
		return false, errors.Errorf("error: could not find sheet by tab name: %v", err)
	}
	sheet.Update(update.Row, update.Column, update.Value)

	// Make sure to call Synchronize to reflect the changes.
	err = sheet.Synchronize()
	if err != nil {
		return false, errors.Errorf("error: could not synchronize changes: %v", err)
	}

	return true, nil
}

//ReadSpreadsheet wrapper
func (s *CLI) ReadSpreadsheet(spreadsheetID string, outputJSON bool) (*SpreadsheetData, error) {
	var spreadsheetData SpreadsheetData

	spreadsheet, err := s.GetSpreadsheet(spreadsheetID)
	if err != nil {
		return nil, errors.Errorf("error: %v", err)
	}

	for _, sheet := range spreadsheet.Sheets {
		data := s.ParseSpreadsheet(sheet)
		spreadsheetData = append(spreadsheetData, *data)
		if outputJSON {
			// Call to write json for each sheet.
			s.ConvertToJSON(*data, sheet.Properties.Title)
		}
	}

	return &spreadsheetData, nil
}

//ParseSpreadsheet ...
func (s *CLI) ParseSpreadsheet(sheet gosheet.Sheet) *SheetData {
	var data SheetData
	var record = make(sheetRecord)
	var structFields []string

	for i, row := range sheet.Rows {
		for j, cell := range row {
			if i == 0 { //This is the struct field.
				structFields = append(structFields, cell.Value)
			} else {
				record[structFields[j]] = cell.Value
			}
		}
		if i != 0 {
			data = append(data, record)
		}
	}

	return &data
}

//ConvertToJSON takes sheetData map, outputs JSON file
func (s *CLI) ConvertToJSON(data SheetData, title string) (bool, error) {
	file, _ := json.MarshalIndent(data, "", " ")
	filename := strings.Join([]string{"data-", title, ".json"}, "")
	ioutil.WriteFile(filename, file, 0644)

	return true, nil
}
