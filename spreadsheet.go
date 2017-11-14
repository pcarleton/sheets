package sheets

import (
  "fmt"
  "bufio"
  "io"
  "strings"

	"google.golang.org/api/sheets/v4"
)


type Spreadsheet struct {
  Client *Client
  *sheets.Spreadsheet
}

type Sheet struct {
  *sheets.Sheet
  Spreadsheet *Spreadsheet
  Client *Client
}


func (s *Spreadsheet) Id() string {
  return s.SpreadsheetId
}

func (s *Spreadsheet) Url() string {
  return s.SpreadsheetUrl
}

func (s *Spreadsheet) Import(sheetName string, data [][]string, cellRange CellRange) error {
  // Convert to interfaces to satisfy the Google API
  converted := make([][]interface{}, 0)

  for _, row := range(data) {
    converted = append(converted, strToInterface(row))
  }

  // TODO: Check if sheet exists already
  vRange := &sheets.ValueRange{
    Range: cellRange.String(),
    Values: converted,
  }

  req := s.Client.Sheets.Spreadsheets.Values.Update(s.Id(), cellRange.String(), vRange)

  req.ValueInputOption("USER_ENTERED")
  _, err := req.Do()

  return err
}

func (s *Spreadsheet) GetSheet(title string) *Sheet {
  query := strings.ToLower(title)
  for _, sheet := range s.Sheets {
    lowerTitle := strings.ToLower(sheet.Properties.Title)
    if lowerTitle == query {
      return &Sheet{sheet, s, s.Client}
    }
  }
  return nil
}


func (s *Sheet) Resize(rows, cols int) error {
  return nil
}

func (s *Sheet) Update(data [][]string, start CellPos) error {
  return nil
}

func (s *Spreadsheet) DoBatch(reqs ...*sheets.Request) (*sheets.BatchUpdateSpreadsheetResponse, error) {
  batchUpdateReq := sheets.BatchUpdateSpreadsheetRequest{
    Requests: reqs,
    IncludeSpreadsheetInResponse: true,
  }

  resp, err := s.Client.Sheets.Spreadsheets.BatchUpdate(s.Id(), &batchUpdateReq).Do()

  if err != nil {
    return nil, err
  }

  s.Spreadsheet = resp.UpdatedSpreadsheet

  return resp, nil
}



func (s *Spreadsheet) AddSheet(title string) (*Sheet, error) {
  sheet := s.GetSheet(title)

  if sheet != nil {
    return sheet, nil
  }

  props := sheets.SheetProperties{Title: title}
  addReq := sheets.Request{AddSheet: &sheets.AddSheetRequest{Properties: &props}}

  _, err := s.DoBatch(&addReq)
  if err != nil {
    return nil, err
  }

  sheet = s.GetSheet(title)

  if sheet == nil {
    return nil, fmt.Errorf("Unable to get sheet after adding it: %s", title)
  }

  return sheet, nil
}

func (s *Spreadsheet) Share(email string) error {
  return s.Client.ShareFile(s.Id(), email)
}


func TsvToArr(reader io.Reader) [][]string {
    delimiter := "\t"
    scanner := bufio.NewScanner(reader)

    data := make([][]string, 0)

    for scanner.Scan() {
      pieces := strings.Split(scanner.Text(), delimiter)
      data = append(data, pieces)
    }

    return data
}

func strToInterface(strs []string) []interface{} {
      arr := make([]interface{}, len(strs))

      for i, s := range(strs) {
        arr[i] = s
      }
      return arr
}

