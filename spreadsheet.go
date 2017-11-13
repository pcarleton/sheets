package sheets

import (
  "bufio"
  "io"
  "strings"

	"google.golang.org/api/sheets/v4"
)


type Spreadsheet struct {
  client *Client
  info *sheets.Spreadsheet
}

func strToInterface(strs []string) []interface{} {
      arr := make([]interface{}, len(strs))

      for i, s := range(strs) {
        arr[i] = s
      }
      return arr
}


func (s *Spreadsheet) Id() string {
  return s.info.SpreadsheetId
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

  req := s.client.Sheets.Spreadsheets.Values.Update(s.Id(), cellRange.String(), vRange)

  req.ValueInputOption("USER_ENTERED")
  _, err := req.Do()

  return err
}


func TsvToArr(reader io.Reader) ([][]string, error) {
    delimiter := "\t"
    scanner := bufio.NewScanner(reader)

    data := make([][]string, 0)

    for scanner.Scan() {
      pieces := strings.Split(scanner.Text(), delimiter)
      data = append(data, pieces)
    }

    return data, nil
}

