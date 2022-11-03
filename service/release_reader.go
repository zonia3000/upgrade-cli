package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const tagsUrl = "https://api.github.com/repos/entando/entando-releases/tags"

type TagData struct {
	Name string `json:"name"`
}

func GetLatestVersion() string {

	resp, err := http.Get(tagsUrl)
	if err != nil {
		fmt.Println("Unable to retrieve tags", err)
		os.Exit(1)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Unable to parse tags response", err)
		os.Exit(1)
	}

	tags := []TagData{}
	json.Unmarshal(bodyBytes, &tags)

	if len(tags) == 0 {
		fmt.Println("No tags found")
		os.Exit(1)
	}

	tag := tags[0].Name
	tag = strings.TrimPrefix(tag, "v")

	return tag
}
