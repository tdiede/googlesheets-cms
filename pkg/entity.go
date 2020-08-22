package entity

import gosheet "gopkg.in/Iwark/spreadsheet.v2"

type googlesheet gosheet.Spreadsheet

type CellContents struct {
	x       int
	y       int
	content string
}
