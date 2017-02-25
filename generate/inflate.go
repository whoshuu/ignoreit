package generate

import (
	"bufio"
	"fmt"
	"os"

	"github.com/whoshuu/ignoreit/network"
	"github.com/whoshuu/ignoreit/spec"
)

// Inflate generates a .gitignore file from the input config.
// Each source specified in the config will be given its own section in the output file.
// Custom ignore patterns are appended at the end of the file in their own section.
func Inflate(config spec.Config, ignoreFilename string) error {
	var generatedLines []string

	generatedLines = append(generatedLines, fmt.Sprintf("#### Auto-generated .gitignore by ignoreit tool (schema version: %d) ####\n", config.SchemaVersion))

	for _, source := range config.Sources {
		sourceLines, err := inflatSource(source)
		if err != nil {
			return fmt.Errorf("Error inflating source [%s - %s]: %s", source.Repo, source.Branch, err)
		}
		generatedLines = append(generatedLines, sourceLines...)
	}

	if len(config.Custom) > 0 {
		generatedLines = append(generatedLines, fmt.Sprint("\n### Custom Patterns ###\n\n"))
		for _, pattern := range config.Custom {
			generatedLines = append(generatedLines, fmt.Sprintln(pattern))
		}
	}

	return writeToFile(ignoreFilename, generatedLines)
}

func inflatSource(source spec.Source) ([]string, error) {
	var sourceLines []string
	if len(source.Entries) > 0 {
		sourceLines = append(sourceLines, fmt.Sprintln("\n### Source:", source.Repo, "-", source.Branch, "###"))
		for _, entry := range source.Entries {
			contents := network.EntryContents(source.GetDownloadLink(entry))
			if contents != "" {
				sourceLines = append(sourceLines, fmt.Sprintln("\n## Entry:", entry, "##"))
				sourceLines = append(sourceLines, fmt.Sprint(contents))
			}
		}
	}

	return sourceLines, nil
}

func writeToFile(filename string, lines []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprint(w, line)
	}

	return w.Flush()
}
