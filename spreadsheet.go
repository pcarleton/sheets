package sheets

import (
	"google.golang.org/api/sheets/v4"
)

const (
  Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)


type Spreadsheet struct {
  client *Client
  info *sheets.Spreadsheet
}


func aRangeLetter(idx int) string {
  secondLet := idx % len(Alphabet)

  if idx > len(Alphabet) {
    firstLet := idx / len(Alphabet)
    return fmt.Sprintf("%s%s",
    Alphabet[firstLet:firstLet + 1], Alphabet[secondLet:secondLet + 1])
  }

  return fmt.Sprintf("%s", Alphabet[secondLet:secondLet+1])
}

type CellPos struct {
  Row int
  Col int
}

func (c CellPos) A1Notation() string {
  return fmt.Sprintf("%s%d", aRangeLetter(c.Col), c.Row + 1)
}

type CellRange struct {
  Start CellPos
  End CellPos
}

func (a CellRange) String() string {
  return fmt.Sprintf("%s:%s", a.Start.A1Notation(), a.End.A1Notation())
}

type SheetRange struct {
  SheetName string
  Range CellRange
}

func (s *SheetRange) String string {
  return fmt.Sprintf("%s!%s", s.SheetName, s.Range.String())
}

func DefaultRange(data [][]interface{}) CellRange {
  bottomLeft := CellPos{len(data), len(data[0])}

  return CellRange(CellPos{}, bottomLeft)
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

func strToInterface(strs []string) []interface{} {
      arr := make([]interface{}, len(strs))

      for i, s := range(strs) {
        arr[i] = s
      }
      return arr
}


func (s *Spreadsheet) Id() string {
  return s.spreadsheet.SpreadsheetId
}

func (s *Spreadsheet) Import(sheetName, data [][]string, cellRange CellRange) error {
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

  req := s.client.Sheets.Spreadsheets.Values.Update(s.Id(), aRange, vRange)

  req.ValueInputOption("USER_ENTERED")
  _, err = req.Do()

  return err
}


