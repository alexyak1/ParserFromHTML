package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type TableData struct {
	name  string
	candy string
	eaten int
}

type StructedResult struct {
	Name           string
	FavouriteSnack string
	TotalSnacks    string
}

var globalAggregatedData []StructedResult

func main() {
	tableData := getExtractedTable("https://candystore.zimpler.net/")
	finalTableData := getTableData(strings.Split(tableData, "<tr>"))

	candiesByName := make(map[string]map[string]int)
	totalSnacksPerName := make(map[string]int)

	for _, data := range finalTableData {
		if _, ok := totalSnacksPerName[data.name]; ok {
			newCount := totalSnacksPerName[data.name] + data.eaten
			totalSnacksPerName[data.name] = newCount
		} else {
			totalSnacksPerName[data.name] = data.eaten
		}

		if _, ok := candiesByName[data.name]; ok {
		} else {
			candiesByName[data.name] = map[string]int{data.candy: 0}
		}

		if _, ok := candiesByName[data.name][data.candy]; ok {
			candiesByName[data.name][data.candy] = candiesByName[data.name][data.candy] + data.eaten
		} else {
			candiesByName[data.name][data.candy] = data.eaten
		}
	}
	favoriteCandy := getFavoriteCandy(candiesByName)

	setSortedData(totalSnacksPerName, favoriteCandy)
	fmt.Println("Data from table to show: ", globalAggregatedData)

	handleRequests()
}
func handleRequests() {
	http.HandleFunc("/", returnAllData)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
func returnAllData(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(globalAggregatedData)
}

func setSortedData(totalSnacksPerName map[string]int, favoriteCandy map[string]string) {
	var result []StructedResult

	valueKey := map[int][]string{}
	var numbersOfEatedCandys []int
	for key, v := range totalSnacksPerName {
		valueKey[v] = append(valueKey[v], key)
	}

	for key := range valueKey {
		numbersOfEatedCandys = append(numbersOfEatedCandys, key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(numbersOfEatedCandys)))
	for _, totalSnacks := range numbersOfEatedCandys {
		for _, name := range valueKey[totalSnacks] {
			result = append(result, StructedResult{name, favoriteCandy[name], strconv.Itoa(totalSnacks)})
		}
	}
	globalAggregatedData = result
}

func getExtractedTable(url string) string {
	response, err := http.Get(url)
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

	return pageContent[startIndex:endIndex]
}

func getTableData(splitted []string) []TableData {
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
	return finalData
}

func getFavoriteCandy(candiesByName map[string]map[string]int) map[string]string {
	candyCount := make(map[string]string)

	for name, dataCollection := range candiesByName {
		famostCandy := make(map[string]int)
		candyName := ""
		for candy, count := range dataCollection {
			if len(famostCandy) == 0 {
				famostCandy[candy] = count
				candyName = candy
			}
			if count > famostCandy[candyName] {
				famostCandy[candy] = count
				candyName = candy
			}

			candyCount[name] = candyName
		}
	}
	return candyCount
}
