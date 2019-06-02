package prompt

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
)

func promptStartEnd() (time.Time, time.Time, error) {
	validate := func(input string) error {
		_, err := time.Parse("2006-01", input)
		if err != nil {
			return err
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Start Date : 2006-01",
		Validate: validate,
	}

	startStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("Start Date entry failed %v\n", err)
		return time.Now(), time.Now(), err
	}

	prompt = promptui.Prompt{
		Label:    "End Date : 2006-01",
		Validate: validate,
	}

	endStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("End Date entry failed %v\n", err)
		return time.Now(), time.Now(), err
	}

	start, _ := time.Parse("2006-01", startStr)
	end, _ := time.Parse("2006-01", endStr)

	return start, end, nil
}
