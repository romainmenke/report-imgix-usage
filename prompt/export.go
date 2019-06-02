package prompt

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/romainmenke/report-imgix-usage/counters"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func promptExport(sources *sources.Sources) bool {
	items := []string{
		"cost      -> csv",
		"bandwidth -> csv",
		"images    -> csv",
		backToMainSelect,
		exitSelect,
	}

	prompt := promptui.Select{
		Label: "Export?",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Export entry failed %v\n", err)
		return false
	}

	if v, ok := handleDefaultOptions(result); ok {
		return v
	}

	switch result {
	case "cost      -> csv":
		promptExportCsv(sources, "cost", func(c *counters.Counters) string {
			c.SetCumulative()
			return fmt.Sprintf("%.2f", cost(c.Cumulative.Bandwidth, c.Cumulative.Images))
		})

	case "bandwidth -> csv":
		promptExportCsv(sources, "bandwidth", func(c *counters.Counters) string {
			c.SetCumulative()
			return fmt.Sprint(c.Cumulative.Bandwidth / (1024 * 1024 * 1024))
		})

	case "images    -> csv":
		promptExportCsv(sources, "images", func(c *counters.Counters) string {
			c.SetCumulative()
			return fmt.Sprint(c.Cumulative.Images)
		})
	}

	return true
}

func promptExportCsv(sources *sources.Sources, label string, valueFunc func(c *counters.Counters) string) {
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
		promptExportCsv(sources, label, valueFunc)
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
		"imgix - " + label,
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

	sourceDataMatrix := [][]string{}
	for _, sourceData := range sources.Data {
		row := make([]string, len(times), len(times))

		for t, c := range sourceData.Counters {
			timeIndex, ok := timesMap[t]
			if !ok {
				continue // should not happen
			}

			row[timeIndex] = valueFunc(c)
		}

		row = append([]string{sourceData.Attributes.Name}, row...)
		sourceDataMatrix = append(sourceDataMatrix, row)
	}

	header = append(header, times...)

	err = csv.Write(header)
	if err != nil {
		panic(err)
	}

	for _, row := range sourceDataMatrix {
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
