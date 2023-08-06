package cmd

import (
	"bufio"
	"dat/contents"
	"dat/highlight"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
	"github.com/savioxavier/termlink"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "dat",
	Short: "Dat is the cat command with super-power!",
	Long:  `User-friendly and colored cat replacement.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		currentPath, _ := os.Getwd()
		filePath := currentPath + "/" + fileName
		fileToRead, err := os.Open(filePath)
		if err != nil {
			log.Panicln(err)
		}
		fileScanner := bufio.NewScanner(fileToRead)
		fileScanner.Split(bufio.ScanLines)

		lineNumber := 0
		parseNumberFlag(fileScanner, &lineNumber)
		fileRows := readFile(fileScanner, &lineNumber)

		highlight.Highlight(fileName, &fileRows)

		printFileInfo(filePath)
		for _, line := range fileRows {
			if showLineNumbers {
				color.Set(color.FgHiBlack)
				fmt.Printf("%d | ", line.Number)
			}

			for i, char := range line.Chars {
				if lineColumn > 0 && i < (lineColumn-1) {
					continue
				}

				color.Set(char.Color)
				fmt.Print(char.Char)
			}
			fmt.Print("\n")
		}
		color.Unset()
	},
}

// --lineNumbers flag
var showLineNumbers bool

// --line flag
var line string
var lineRowNumber int
var lineColumn int

func Execute() {
	rootCmd.SetVersionTemplate("1.0.0")
	rootCmd.PersistentFlags().BoolVarP(&showLineNumbers, "lineNumbers", "n", false, "Print line numbers")
	rootCmd.PersistentFlags().StringVarP(&line, "line", "l", "", "Show line (and column) at this position. Ex. 5:3 or 5")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseNumberFlag(fileScanner *bufio.Scanner, lineNumber *int) {
	if line == "" {
		return
	}

	if !strings.Contains(line, ":") {
		lineRowNumber, _ = strconv.Atoi(line)
	} else {
		lineSplit := strings.Split(line, ":")
		lineRowNumber, _ = strconv.Atoi(lineSplit[0])
		lineColumn, _ = strconv.Atoi(lineSplit[1])
	}

	for i := 0; i < lineRowNumber-1; i++ {
		*lineNumber++
		fileScanner.Scan()
	}
}

func readFile(fileScanner *bufio.Scanner, lineNumber *int) []contents.Row {
	var fileRows []contents.Row
	for fileScanner.Scan() {
		*lineNumber++

		var chars []contents.Char
		for _, char := range fileScanner.Text() {
			chars = append(chars, contents.Char{
				Char:  string(char),
				Color: color.Reset,
			})
		}

		row := contents.Row{
			Number: *lineNumber,
			Chars:  chars,
		}
		fileRows = append(fileRows, row)
		if lineRowNumber > 0 && *lineNumber >= lineRowNumber {
			break
		}
	}

	return fileRows
}

func printFileInfo(filePath string) {
	fileStat, _ := os.Stat(filePath)

	color.Set(color.FgWhite)
	fmt.Print(termlink.Link(fileStat.Name(), "file://"+filePath))

	color.Set(color.FgHiBlack)
	fmt.Print(" | ")

	color.Set(color.FgWhite)
	fmt.Printf("%s", humanize.Bytes(uint64(fileStat.Size())))

	color.Set(color.FgHiBlack)
	fmt.Print(" | ")

	color.Set(color.FgWhite)
	fmt.Print(humanize.Time(fileStat.ModTime()))

	color.Set(color.FgHiBlack)
	fmt.Print(" | ")

	color.Set(color.FgWhite)
	fmt.Print(fileStat.Mode())

	fmt.Print("\n")
	color.Unset()
}
