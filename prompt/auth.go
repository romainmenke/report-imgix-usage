package prompt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manifoldco/promptui"
	"github.com/romainmenke/report-imgix-usage/sources"
)

func PromptAuthAndGetSources() *sources.Sources {
	client := http.DefaultClient
	client.Transport = http.DefaultTransport

	if promptAuthSelect(client) {
		return getAllData(client)
	}

	return nil
}

func promptAuthSelect(client *http.Client) bool {
	items := []string{
		"email and password",
		"token",
	}

	prompt := promptui.Select{
		Label: "Authentication method?",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Authentication selection failed %v\n", err)
		return false
	}

	if result == "token" {
		return promptAuthToken(client)
	}

	if result == "email and password" {
		return promptAuthEmailPassword(client)
	}

	return false
}

func promptAuthToken(client *http.Client) bool {
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
		return false
	}

	client.Transport = AuthRoundTripper(client, auth)

	return true
}

func promptAuthEmailPassword(client *http.Client) bool {
	validate := func(input string) error {
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Email",
		Validate: validate,
	}

	email, err := prompt.Run()
	if err != nil {
		fmt.Printf("Email entry failed %v\n", err)
		return false
	}

	prompt = promptui.Prompt{
		Label:    "Password",
		Validate: validate,
		Mask:     '*',
	}

	password, err := prompt.Run()
	if err != nil {
		fmt.Printf("Password entry failed %v\n", err)
		return false
	}

	payload := map[string]map[string]interface{}{
		"data": map[string]interface{}{
			"type": "users",
			"attributes": map[string]string{
				"email":    email,
				"password": password,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Authentication failed")
		return false
	}

	resp, err := http.Post("https://api.imgix.com/v4/users/login", "application/vnd.api+json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Authentication failed")
		return false
	}

	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Authentication failed")
		return false
	}

	respPayload := map[string]interface{}{}
	err = json.Unmarshal(respData, &respPayload)
	if err != nil {
		fmt.Println("Authentication failed")
		return false
	}

	if data, ok := respPayload["data"].(map[string]interface{}); !ok {
		fmt.Println("Authentication failed")
		return false
	} else if attributes, ok := data["attributes"].(map[string]interface{}); !ok {
		fmt.Println("Authentication failed")
		return false
	} else if auth, ok := attributes["api_key"].(string); !ok {
		fmt.Println("Authentication failed")
		return false
	} else {
		client.Transport = AuthRoundTripper(client, auth)
	}

	return true
}

type RoundTripper func(req *http.Request) (*http.Response, error)

func (x RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return x(req)
}

func AuthRoundTripper(client *http.Client, auth string) http.RoundTripper {
	next := client.Transport

	return RoundTripper(func(req *http.Request) (*http.Response, error) {
		req.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth+":"))))

		return next.RoundTrip(req)
	})
}
