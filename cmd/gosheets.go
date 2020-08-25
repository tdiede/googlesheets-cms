package main

import (
	"flag"
	"fmt"

	"github.com/tdiede/gosheets/pkg"
)

const (
	_clientSecretPath = "client_secret.json"
)

// go run cmd/gosheets.go -s 10j3vs-LmZcLSlygaEmunjvSZbysj5oFcNtEfNpTp9ZQ -method read

func main() {
	cli, err := parser.NewCLI(_clientSecretPath)
	if err != nil {
		fmt.Printf("error: %v", err)
	}

	spreadsheetID, method := parseFlags()

	switch method {
	case "read":
		cli.ReadSpreadsheet(spreadsheetID, false)
	case "update":
		cellUpdate := parser.ContentUpdate{
			Row:        2,
			Column:     2,
			Value:      "hello",
			SheetTitle: "tab1",
		}
		cli.UpdateSpreadsheet(spreadsheetID, cellUpdate)
	case "convert":
		fmt.Println("Converting spreadsheet contents to JSON file.")
		cli.ReadSpreadsheet(spreadsheetID, true)
	default:
		fmt.Printf("%s.\n", method)
	}
}

func parseFlags() (string, string) {
	s := flag.String("s", "", "Spreadsheet ID from Google Sheets, e.g. https://docs.google.com/spreadsheets/d/<Spreadsheet ID>")
	m := flag.String("method", "", "The command used for desired action: 'read', 'update', 'convert' (to JSON)")
	flag.Parse()
	return *s, *m
}
