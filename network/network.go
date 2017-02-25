package network

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// EntryExists checks if the input url points to a valid hosted .gitignore file.
// An HTTP request with method HEAD expects to return 200 OK in the response.
// Any other response is interpreted to mean that the .gitignore entry does not exist.
func EntryExists(url string) bool {
	resp, err := http.Head(url)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return resp.StatusCode == http.StatusOK
}

// EntryContents gets the contents of the .gitignore file pointed to be the input url.
// If the response is not 200 OK, an empty string is returned instead.
func EntryContents(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
	}

	return string(body)
}
