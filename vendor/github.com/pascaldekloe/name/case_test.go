package name

import (
	"strings"
	"testing"
)

var goldenCamelSnakes = map[string]string{
	"":      "",
	"name":  "name",
	"Label": "label",
	"ID":    "ID",

	"loFi": "lo_fi",
	"HiFi": "hi_fi",

	// inner abbreviation
	"rawURLQuery": "raw_URL_query",

	// single outer abbreviation
	"vitaminC": "vitamin_C",
	"TCell":    "T_cell",

	// double outer abbreviation
	"masterDB": "master_DB",
	"IOBounds": "IO_bounds",

	// tripple outer abbreviation
	"mainAPI": "main_API",
	"TCPConn": "TCP_conn",

	// numbers
	"b2b":  "b2b",
	"4x4":  "4x4",
	"No5":  "no5",
	"DB2":  "DB2",
	"3M":   "3M",
	"7Up":  "7_up",
	"20th": "20th",
}

func TestCamelToSnake(t *testing.T) {
	for camel, snake := range goldenCamelSnakes {
		if got := SnakeCase(camel); got != snake {
			t.Errorf("snake case %q got %q, want %q", camel, got, snake)
		}
	}
}

func TestSnakeToSnake(t *testing.T) {
	for _, s := range goldenCamelSnakes {
		if got := SnakeCase(s); got != s {
			t.Errorf("snake case %q got %q", s, got)
		}
	}
}

func TestSnakeToCamel(t *testing.T) {
	for camel, snake := range goldenCamelSnakes {
		want := strings.Title(camel)
		got := CamelCase(snake, true)
		if got != want {
			t.Errorf("camel case %q titled got %q, want %q", snake, got, camel)
		}
	}
}
