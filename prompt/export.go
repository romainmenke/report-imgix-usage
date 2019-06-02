package prompt

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func promptExport(sources *sources.Sources) bool {
	items := []string{
		"csv",
		backToMainSelect,
		exitSelect,
	}

	prompt := promptui.Select{
		Label: "Export?",
		Items: items,
	}

	_, csv, err := prompt.Run()
	if err != nil {
		fmt.Printf("Export entry failed %v\n", err)
		return false
	}

	if v, ok := handleDefaultOptions(csv); ok {
		return v
	}

	items = []string{
		"cost",
		"bandwidth",
		"images",
	}

	prompt = promptui.Select{
		Label: "Field?",
		Items: items,
	}

	_, field, err := prompt.Run()
	if err != nil {
		fmt.Printf("Export entry failed %v\n", err)
		return false
	}

	items = []string{
		"normal",
		"cumulative",
	}

	prompt = promptui.Select{
		Label: "Data?",
		Items: items,
	}

	_, valueAggregation, err := prompt.Run()
	if err != nil {
		fmt.Printf("Export entry failed %v\n", err)
		return false
	}

	promptExportCsv(sources, field, valueAggregation == "cumulative")

	return true
}

func promptExportCsv(sources *sources.Sources, field string, cumulative bool) {
	validate := func(input string) error {
		if !isValidFilePath(input) {
			return errors.New("invalid export path")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Export File",
		Validate: validate,
	}

	filePath, err := prompt.Run()
	if err != nil {
		promptExportCsv(sources, field, cumulative)
		return
	}

	// csv file
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	csv := csv.NewWriter(f)
	defer csv.Flush()

	// csv header
	header := []string{
		"imgix - " + field,
	}

	timesMap := map[string]int{}
	for _, sourceData := range sources.Data {
		for t := range sourceData.Counters {
			timesMap[t] = 0
		}
	}

	times := []string{}
	for t := range timesMap {
		times = append(times, t)
	}

	sort.Strings(times)

	for i, t := range times {
		timesMap[t] = i
	}

	var sourceDataMatrix matrix
	switch field {
	case "cost":
		sourceDataMatrix = floatMatrix{}
	case "bandwidth":
		sourceDataMatrix = intMatrix{}
	case "images":
		sourceDataMatrix = intMatrix{}
	}

	sourceDataMatrix = sourceDataMatrix.size(len(times), len(sources.Data))

	for y, sourceData := range sources.Data {

		for t, c := range sourceData.Counters {
			timeIndex, ok := timesMap[t]
			if !ok {
				continue // should not happen
			}

			c.SetSum()

			switch field {
			case "cost":

				v := cost(c.Sum.Bandwidth, c.Sum.Images)
				sourceDataMatrix.insert(timeIndex, y, v)

			case "bandwidth":

				v := c.Sum.Bandwidth
				sourceDataMatrix.insert(timeIndex, y, v/(1024*1024*1024))

			case "images":

				v := c.Sum.Images
				sourceDataMatrix.insert(timeIndex, y, v)
			}
		}
	}

	if cumulative {
		sourceDataMatrix = sourceDataMatrix.cumulative()
	}

	header = append(header, times...)

	err = csv.Write(header)
	if err != nil {
		panic(err)
	}

	for y, row := range sourceDataMatrix.toString() {
		row = append([]string{sources.Data[y].Attributes.Name}, row...)

		err := csv.Write(row)
		if err != nil {
			panic(err)
		}
	}
}

func isValidFilePath(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}
