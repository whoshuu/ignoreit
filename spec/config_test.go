package spec

import (
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

const (
	repoName     = "github/gitignore"
	branchName   = "master"
	testFilename = ".ignoreit.test.yml"
	rawConfig    = `sources:
- repo: github/gitignore
  branch: ghfw
  entries:
  - Ada
  - Python
- repo: github/gitignore
  branch: master
  entries:
  - C++
  - CMake
  - Go
custom:
- .custompattern
- .anothercustompattern
schema_version: 1
`
	rawConfigBadSchema = `sources: []
custom: []
schema_version: 2
`
)

type sourceInput struct {
	repo   string
	branch string
}

var (
	allInputs = []sourceInput{{"", ""}, {repoName, branchName}, {"", branchName}, {repoName, ""}}
)

func TestCreateSourceSingle(t *testing.T) {
	config := Config{}

	source := config.CreateSource(repoName, branchName)
	if source == nil {
		t.Errorf("Source should not be nil for repo: %s, branch: %s", repoName, branchName)
	}
}

func TestCreateSourceSingleFromMany(t *testing.T) {
	config := Config{}

	for _, input := range allInputs {
		config.CreateSource(input.repo, input.branch)
	}

	if len(config.Sources) != 1 {
		t.Errorf("Should have only created 1 source, but created %d instead", len(config.Sources))
	}
}

func TestCreateSourceAlreadyExists(t *testing.T) {
	config := Config{}

	config.CreateSource(repoName, branchName)
	source := config.CreateSource(repoName, branchName)
	if source == nil {
		t.Errorf("Source should not be nil for repo: %s, branch: %s", repoName, branchName)
	}

	if len(config.Sources) != 1 {
		t.Errorf("Should have only created 1 source, but created %d instead", len(config.Sources))
	}
}

func TestCreateSourceMany(t *testing.T) {
	config := Config{}

	for i := 0; i < 1000; i++ {
		config.CreateSource(randSeq(20), randSeq(20))
	}

	if len(config.Sources) != 1000 {
		t.Errorf("Should have only created 1000 sources, but created %d instead", len(config.Sources))
	}
}

func TestGetSourceEmpty(t *testing.T) {
	config := Config{}

	for _, input := range allInputs {
		source := config.GetSource(input.repo, input.branch)
		if source != nil {
			t.Errorf("Source should be nil for repo: %s, branch: %s", input.repo, input.branch)
		}
	}
}

func TestGetSourceSingle(t *testing.T) {
	config := Config{}

	config.CreateSource(repoName, branchName)
	source := config.GetSource(repoName, branchName)
	if source == nil {
		t.Errorf("Source should not be nil for repo: %s, branch: %s", repoName, branchName)
	}
}

func TestLoadConfigEmptyFilename(t *testing.T) {
	_, err := LoadConfig("")

	if err == nil {
		t.Error("Error should be returned")
	}
}

func TestLoadConfigNonexistentFilename(t *testing.T) {
	config, err := LoadConfig(randSeq(50))

	if err != nil {
		t.Errorf("Error should not be returned: %s", err)
	}

	if len(config.Sources) != 0 {
		t.Errorf("Length of Sources should be 0, got %d instead", len(config.Sources))
	}

	if len(config.Custom) != 0 {
		t.Errorf("Length of Custom patterns should be 0, got %d instead", len(config.Custom))
	}

	if config.SchemaVersion != 1 {
		t.Errorf("SchemaVersion should be 1, got %d instead", config.SchemaVersion)
	}
}

func TestLoadConfigExistingFilename(t *testing.T) {
	if err := ioutil.WriteFile(testFilename, []byte(rawConfig), 0644); err != nil {
		panic(err)
	}
	defer os.Remove(testFilename)

	config, err := LoadConfig(testFilename)

	if err != nil {
		t.Errorf("Error should not be returned: %s", err)
	}

	if len(config.Sources) != 2 {
		t.Errorf("Length of Sources should be 2, got %d instead", len(config.Sources))
	}

	expectedCustomPatterns := []string{".custompattern", ".anothercustompattern"}

	actualSource := config.Sources[0]
	if actualSource.Repo != repoName || actualSource.Branch != "ghfw" {
		t.Errorf("Source should be repo: %s, branch: %s, got repo: %s, branch: %s instead", repoName, "ghfw", actualSource.Repo, actualSource.Branch)
	}

	expectedEntries := []string{"Ada", "Python"}

	if len(actualSource.Entries) != len(expectedEntries) {
		t.Errorf("Length of Source Entries should be %d, got %d instead", len(expectedEntries), len(actualSource.Entries))
	}
	for i := range actualSource.Entries {
		if actualSource.Entries[i] != expectedEntries[i] {
			t.Errorf("Source should have %s at index %d, got %s instead", expectedEntries[i], i, actualSource.Entries[i])
		}
	}

	actualSource = config.Sources[1]
	if actualSource.Repo != repoName || actualSource.Branch != branchName {
		t.Errorf("Source should be repo: %s, branch: %s, got repo: %s, branch: %s instead", repoName, branchName, actualSource.Repo, actualSource.Branch)
	}

	expectedEntries = []string{"C++", "CMake", "Go"}

	if len(actualSource.Entries) != len(expectedEntries) {
		t.Errorf("Length of Source Entries should be %d, got %d instead", len(expectedEntries), len(actualSource.Entries))
	}
	for i := range actualSource.Entries {
		if actualSource.Entries[i] != expectedEntries[i] {
			t.Errorf("Source should have %s at index %d, got %s instead", expectedEntries[i], i, actualSource.Entries[i])
		}
	}

	if len(config.Custom) != len(expectedCustomPatterns) {
		t.Errorf("Length of Custom patterns should be %d, got %d instead", len(expectedCustomPatterns), len(config.Custom))
	}

	for i := range config.Custom {
		if config.Custom[i] != expectedCustomPatterns[i] {
			t.Errorf("Custom should have custom pattern %s at index %d, got %s instead", expectedCustomPatterns[i], i, config.Custom[i])
		}
	}

	if config.SchemaVersion != 1 {
		t.Errorf("SchemaVersion should be 1, got %d instead", config.SchemaVersion)
	}
}

func TestLoadConfigBadSchema(t *testing.T) {
	if err := ioutil.WriteFile(testFilename, []byte(rawConfigBadSchema), 0644); err != nil {
		panic(err)
	}
	defer os.Remove(testFilename)

	_, err := LoadConfig(testFilename)

	if err.Error() != "Schema version 2 does not match expected version 1" {
		t.Errorf("Schema check should have failed for bad schema: %s", err)
	}
}

//- test saving config
//- test dedupe in save
//- test empty clean in save
//- test unreadable file in loadconfig
//- test yml unmarshal error
func TestSave(t *testing.T) {
}
