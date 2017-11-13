package sheets

import (
  "testing"
)

var posTests = []struct {
  pos CellPos
  expected string
}{
  {CellPos{0, 0}, "A1"},
  {CellPos{1, 0}, "A2"},
  {CellPos{0, 1}, "B1"},
  {CellPos{1, 1}, "B2"},
  {CellPos{10, 10}, "K11"},
  {CellPos{0, 25}, "Z1"},
  {CellPos{0, 26}, "AA1"},
  {CellPos{0, 27}, "AB1"},
  {CellPos{0, 52}, "BA1"},

  {CellPos{0, 624}, "XA1"},
  {CellPos{0, 650}, "YA1"},
  {CellPos{0, 675}, "YZ1"},
  {CellPos{0, 701}, "ZZ1"},

}

func TestCellPosA1Notation(t *testing.T) {

  for _, tt := range posTests {
    got := tt.pos.A1Notation()
    if  got != tt.expected {
      t.Errorf("Wanted %s, but got %s for %v", tt.expected, got, tt.pos)
    }
  }
}

