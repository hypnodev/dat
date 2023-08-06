package highlight

import (
	"dat/contents"
	"fmt"
	"github.com/fatih/color"
	"github.com/zyedidia/highlight"
	"os"
	"strings"
)

func Highlight(fileName string, rows *[]contents.Row) {
	fileType := getFileType(fileName, *rows)

	highlightPath := os.Getenv("DAT_HIGHLIGHT_FILE")
	syntaxFile, _ := os.ReadFile(highlightPath + "/syntax_files/" + fileType + ".yaml")
	syntaxDef, err := highlight.ParseDef(syntaxFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	h := highlight.NewHighlighter(syntaxDef)
	matches := h.HighlightString(contents.Join(*rows))

	for lineNumber, row := range *rows {
		var consoleColor color.Attribute
		for columnNumber := range row.Chars {
			if group, ok := matches[lineNumber][columnNumber]; ok {
				if group == highlight.Groups["statement"] {
					consoleColor = color.FgGreen
				} else if group == highlight.Groups["preproc"] {
					consoleColor = color.FgHiRed
				} else if group == highlight.Groups["special"] {
					consoleColor = color.FgBlue
				} else if group == highlight.Groups["constant.string"] {
					consoleColor = color.FgCyan
				} else if group == highlight.Groups["constant.specialChar"] {
					consoleColor = color.FgHiMagenta
				} else if group == highlight.Groups["type"] {
					consoleColor = color.FgYellow
				} else if group == highlight.Groups["constant.number"] {
					consoleColor = color.FgCyan
				} else if group == highlight.Groups["comment"] {
					consoleColor = color.FgHiGreen
				}
			}
			(*rows)[lineNumber].Chars[columnNumber].Color = consoleColor
		}

		if group, ok := matches[lineNumber][len(row.Chars)]; ok {
			if group == highlight.Groups["default"] || group == highlight.Groups[""] {
				color.Unset()
			}
		}
	}
}

func preloadDefinitions() []*highlight.Def {
	var defs []*highlight.Def
	highlightPath := os.Getenv("DAT_HIGHLIGHT_FILE")
	files, _ := os.ReadDir(highlightPath + "/syntax_files")

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yaml") {
			input, _ := os.ReadFile(highlightPath + "/syntax_files/" + f.Name())
			d, err := highlight.ParseDef(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			defs = append(defs, d)
		}
	}

	return defs
}

func getFileType(fileName string, rows []contents.Row) string {
	defs := preloadDefinitions()
	firstLine := []byte(rows[0].Text())
	def := highlight.DetectFiletype(defs, fileName, firstLine)

	return def.FileType
}
