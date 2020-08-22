package main

import (
	"fmt"
	"io/ioutil"

	gosheet "gopkg.in/Iwark/spreadsheet.v2"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)


// go run cmd/gosheets.go

type googlesheet gosheet.Spreadsheet

func main() {
	data, err := ioutil.ReadFile("client_secret.json")
	conf, err := google.JWTConfigFromJSON(data, gosheet.Scope)
	client := conf.Client(context.TODO())
	spreadsheetID := "10j3vs-LmZcLSlygaEmunjvSZbysj5oFcNtEfNpTp9ZQ"
	service := gosheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(spreadsheetID)

	for _, sheet := range spreadsheet.Sheets {
	  for _, row := range sheet.Rows {
		  for _, cell := range row {
			  fmt.Println(cell.Value)
		  }
	  }
	}

	// Update cell content
	spreadsheet.Sheets[1].Update(0, 0, "hogehoge")

	// Make sure call Synchronize to reflect the changes
	err = spreadsheet.Sheets[1].Synchronize()
	fmt.Println(err)
}
