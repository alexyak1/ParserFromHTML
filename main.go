package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
)

type TableData struct {
	name  string
	candy string
	eaten int
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

		pattern := regexp.MustCompile(`[a-zA-ZäöåÄÖÅ]*`)
		patternForInt := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

		name := pattern.FindAllString(textInRow[1], -4)[0]
		candy := pattern.FindAllString(textInRow[2], -1)[0]
		eatenCandies := patternForInt.FindAllString(textInRow[3], -1)[0]

		var eatenCandiesInt int
		fmt.Sscan(eatenCandies, &eatenCandiesInt) //convert string to int

		finalData = append(finalData, TableData{name, candy, eatenCandiesInt})
	}

	groupedByName := make(map[string]map[string]int)

	for _, data := range finalData {

		if el, ok := groupedByName[data.name]; ok {
			newCount := el["count"] + data.eaten

			groupedByName[data.name]["count"] = newCount

		} else {
			groupedByName[data.name] = map[string]int{"count": data.eaten}
		}

		if _, ok := groupedByName[data.name][data.candy]; ok {
			groupedByName[data.name][data.candy] = groupedByName[data.name][data.candy] + 1
		} else {
			groupedByName[data.name][data.candy] = 1
		}
	}
	// fmt.Println(groupedByName)

	// now sort by most eated candyes
	nameCount := make(map[string]int)

	for name, dataCollection := range groupedByName {
		nameCount[name] = dataCollection["count"]
	}

	valueKey := map[int][]string{}
	var numbersOfEatedCandys []int
	for key, v := range nameCount {
		valueKey[v] = append(valueKey[v], key)
	}

	for key := range valueKey {
		numbersOfEatedCandys = append(numbersOfEatedCandys, key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbersOfEatedCandys)))
	for _, key := range numbersOfEatedCandys {
		for _, name := range valueKey[key] {
			fmt.Printf("%s, %d\n", name, key)
		}
	}
}
