package prompt

import (
	"fmt"
	"sort"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/romainmenke/report-imgix-usage/counters"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func promptSources(sources *sources.Sources) bool {
	items := []string{}

	for _, sourceData := range sources.Data {
		items = append(items, sourceData.Attributes.Name)
	}

	sort.Strings(items)

	items = append([]string{"all sources"}, items...)
	items = append(items, []string{backToMainSelect, exitSelect}...)

	prompt := promptui.Select{
		Label: "Select Source",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Source selection failed %v\n", err)
		return false
	}

	if v, ok := handleDefaultOptions(result); ok {
		return v
	}

	if result == "all sources" {
		return promptListForAllSource(sources)
	}

	for _, sourceData := range sources.Data {
		if result == sourceData.Attributes.Name {
			return promptListForSingleSource(sourceData)
		}
	}

	return false
}

func promptListForSingleSource(sourceData *sources.Data) bool {
	start, end, err := promptStartEnd()
	if err != nil {
		return false
	}

	foundCounters := counters.MultipleCounters{}
	for t, counters := range sourceData.Counters {
		compare, err := time.Parse("2006-01", t)
		if err != nil {
			panic(err)
		}

		if (compare.After(start) || compare.Equal(start)) && (compare.Before(end) || compare.Equal(end)) {
			foundCounters = append(foundCounters, counters)
		}
	}

	selectedCounters := foundCounters.Sum()

	fmt.Printf(
		"%-20s%-20s%-20s\n",
		"cost",
		"bandwidth",
		"images",
	)

	fmt.Printf(
		"%-20s%-20s%-20s\n",
		fmt.Sprintf("%.2f", cost(selectedCounters.Sum.Bandwidth, selectedCounters.Sum.Images)),
		fmt.Sprintf("%d gb", selectedCounters.Sum.Bandwidth/(1024*1024*1024)),
		fmt.Sprint(selectedCounters.Sum.Images),
	)

	return true
}

func promptListForAllSource(allSources *sources.Sources) bool {
	start, end, err := promptStartEnd()
	if err != nil {
		return false
	}

	foundCounters := counters.MultipleCounters{}
	for _, sourceData := range allSources.Data {
		for t, counters := range sourceData.Counters {
			compare, err := time.Parse("2006-01", t)
			if err != nil {
				panic(err)
			}

			if (compare.After(start) || compare.Equal(start)) && (compare.Before(end) || compare.Equal(end)) {
				foundCounters = append(foundCounters, counters)
			}
		}
	}

	selectedCounters := foundCounters.Sum()

	fmt.Printf(
		"%-20s%-20s%-20s\n",
		"cost",
		"bandwidth",
		"images",
	)

	fmt.Printf(
		"%-20s%-20s%-20s\n",
		fmt.Sprintf("%.2f", cost(selectedCounters.Sum.Bandwidth, selectedCounters.Sum.Images)),
		fmt.Sprintf("%d gb", selectedCounters.Sum.Bandwidth/(1024*1024*1024)),
		fmt.Sprint(selectedCounters.Sum.Images),
	)

	return true
}
