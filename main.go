package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type TableData struct {
	name  string
	candy string
	eaten string
}

func main() {

	response, err := http.Get("https://candystore.zimpler.net/")
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	dataFromDOM, _ := ioutil.ReadAll(response.Body)
	pageContent := string(dataFromDOM)

	startIndex := strings.Index(pageContent, "<table id=\"top.customers\" class=\"top.customers details\">")
	if startIndex == -1 {
		fmt.Println("No table with id top.customers found")
		os.Exit(0)
	}

	startIndex += 8

	endIndex := strings.Index(pageContent, "<footer>")
	if endIndex == -1 {
		fmt.Println("No closing tag")
		os.Exit(0)
	}

	tableData := pageContent[startIndex:endIndex]

	splitted := strings.Split(tableData, "<tr>")

	var finalData []TableData

	for i, tableRowData := range splitted {
		if i == 0 {
			continue
		}

		textInRow := strings.Split(tableRowData, "<td>")

		pattern := regexp.MustCompile(`([^"]*) *</td>`)
		patternForInt := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

		submatchall := patternForInt.FindAllString(textInRow[3], -1)

		finalData = append(finalData, TableData{pattern.ReplaceAllString(textInRow[1], "${1}"),
			pattern.ReplaceAllString(textInRow[2], "${1}"),
			submatchall[0]})
	}
	fmt.Println(finalData)

}
