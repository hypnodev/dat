package contents

import "github.com/fatih/color"

type Char struct {
	Char  string
	Color color.Attribute
}

type Row struct {
	Number int
	Chars  []Char
}

func (row Row) Text() string {
	line := ""

	for _, char := range row.Chars {
		line = line + char.Char
	}

	return line
}

func Join(rows []Row) string {
	lines := ""
	for _, row := range rows {
		lines = lines + row.Text() + "\n"
	}

	return lines
}
