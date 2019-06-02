package prompt

const backToMainSelect = "back to main"
const exitSelect = "exit"

func handleDefaultOptions(input string) (bool, bool) {
	if input == exitSelect {
		return false, true
	}

	if input == backToMainSelect {
		return true, true
	}

	return false, false
}
