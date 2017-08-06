package name

import (
	"testing"
)

type goldenCase struct {
	snake, lowerCamel, upperCamel string
}

var goldenCases = []goldenCase{
	{"", "", ""},
	{"i", "i", "I"},
	{"name", "name", "Name"},
	{"ID", "ID", "ID"},
	{"wi_fi", "wiFi", "WiFi"},

	// single outer abbreviation
	{"vitamin_C", "vitaminC", "VitaminC"},
	{"T_cell", "TCell", "TCell"},

	// double outer abbreviation
	{"master_DB", "masterDB", "MasterDB"},
	{"IO_bounds", "IOBounds", "IOBounds"},

	// tripple outer abbreviation
	{"main_API", "mainAPI", "MainAPI"},
	{"TCP_conn", "TCPConn", "TCPConn"},

	// inner abbreviation
	{"raw_URL_query", "rawURLQuery", "RawURLQuery"},

	// numbers
	{"4x4", "4x4", "4x4"},
	{"no5", "no5", "No5"},
	{"DB2", "DB2", "DB2"},
	{"3M", "3M", "3M"},
	{"7_up", "7Up", "7Up"},
	{"20th", "20th", "20th"},
}

func TestSnakeToSnake(t *testing.T) {
	for _, golden := range goldenCases {
		s := golden.snake
		if got := SnakeCase(s); got != s {
			t.Errorf("%q: got %q", s, got)
		}
	}
}

func TestLowerCamelToLowerCamel(t *testing.T) {
	for _, golden := range goldenCases {
		s := golden.lowerCamel
		if got := CamelCase(s, false); got != s {
			t.Errorf("%q: got %q", s, got)
		}
	}
}

func TestUpperCamelToUpperCamel(t *testing.T) {
	for _, golden := range goldenCases {
		s := golden.upperCamel
		if got := CamelCase(s, true); got != s {
			t.Errorf("%q: got %q", s, got)
		}
	}
}

func TestSnakeToLowerCamel(t *testing.T) {
	for _, golden := range goldenCases {
		snake, want := golden.snake, golden.lowerCamel
		if got := CamelCase(snake, false); got != want {
			t.Errorf("%q: got %q, want %q", snake, got, want)
		}
	}
}

func TestSnakeToUpperCamel(t *testing.T) {
	for _, golden := range goldenCases {
		snake, want := golden.snake, golden.upperCamel
		if got := CamelCase(snake, true); got != want {
			t.Errorf("%q: got %q, want %q", snake, got, want)
		}
	}
}

func TestLowerCamelToSnake(t *testing.T) {
	for _, golden := range goldenCases {
		camel, want := golden.lowerCamel, golden.snake
		if got := SnakeCase(camel); got != want {
			t.Errorf("%q: got %q, want %q", camel, got, want)
		}
	}
}

func TestUpperCamelToSnake(t *testing.T) {
	for _, golden := range goldenCases {
		camel, want := golden.upperCamel, golden.snake
		if got := SnakeCase(camel); got != want {
			t.Errorf("%q: got %q, want %q", camel, got, want)
		}
	}
}
