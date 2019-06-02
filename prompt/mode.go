package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func PromptMode(sources *sources.Sources) bool {
	items := []string{
		"drilldown",
		"export",
		exitSelect,
	}

	prompt := promptui.Select{
		Label: "Mode?",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Mode selection failed %v\n", err)
		return false
	}

	if result == exitSelect {
		return false
	}

	if result == "drilldown" {
		return promptSources(sources)
	}

	if result == "export" {
		return promptExport(sources)
	}

	return false
}
