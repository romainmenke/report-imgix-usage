package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func PromptAuthAndGetSources() *sources.Sources {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Auth",
		Validate: validate,
		Mask:     '*',
	}

	auth, err := prompt.Run()
	if err != nil {
		fmt.Printf("Auth entry failed %v\n", err)
		return nil
	}

	return getAllData(auth)
}
