package spec

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

const (
	schemaVersion = 1
)

// Config encapsulates a specification of .gitignore entries and their sources.
// It includes a list of custom strings that can be used as additional .gitignore patterns.
// The schema is versioned to enable forward and backward compatibility.
type Config struct {
	Sources       Sources  `yaml:"sources"`
	Custom        []string `yaml:"custom"`
	SchemaVersion uint     `yaml:"schema_version"`
}

// GetSource grabs a modifiable reference to a Source if it exists in the input config.
// If a Source of the repo and branch name doesn't exist, nil is returned instead.
func (config Config) GetSource(repo, branch string) *Source {
	if repo == "" || branch == "" {
		return nil
	}
	var source *Source
	for i := range config.Sources {
		if config.Sources[i].Repo == repo && config.Sources[i].Branch == branch {
			source = &config.Sources[i]
			break
		}
	}
	return source
}

// CreateSource creates a modifiable reference to a Source if it doesn't yet exist.
// Otherwise, it returns that reference without modifying config.
func (config *Config) CreateSource(repo, branch string) *Source {
	if repo == "" || branch == "" {
		return nil
	}
	source := config.GetSource(repo, branch)
	if source == nil {
		config.Sources = append(config.Sources, Source{repo, branch, []string{}})
		source = &config.Sources[len(config.Sources)-1]
	}
	return source
}

// Save will write the config to disk in YAML format for readability.
// Prior to the write, the config is deduped and scrubbed.
// Sources with no Entries will be removed from config.
// Custom patterns are left unmodified as users are responsible for proper maintenance of that array.
func (config *Config) Save(configFilename string) error {
	config.clean()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshalling config [ %v ] to yaml: %s", config, err)
	}

	return ioutil.WriteFile(configFilename, data, 0644)
}

// LoadConfig will unmarshal a Config struct from a config file in the current working directory.
// If the schema version of the loaded file is different from the tool's schema version, it is rejected.
func LoadConfig(configFilename string) (Config, error) {
	config := Config{}
	config.SchemaVersion = schemaVersion

	if configFilename == "" {
		return config, fmt.Errorf("cannot specify empty string for configFilename")
	}

	contents, err := ioutil.ReadFile(configFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}

		return config, err
	}

	err = yaml.Unmarshal(contents, &config)
	if err == nil {
		err = config.checkSchema()
	}

	return config, err
}

func (config *Config) clean() {
	sort.Sort(config.Sources)
	for i := 0; i < len(config.Sources); {
		config.Sources[i].Clean()
		if len(config.Sources[i].Entries) == 0 {
			config.Sources = append(config.Sources[:i], config.Sources[i+1:]...)
			continue
		}
		i++
	}
}

func (config Config) checkSchema() error {
	if config.SchemaVersion != schemaVersion {
		return fmt.Errorf("Schema version %d does not match expected version %d", config.SchemaVersion, schemaVersion)
	}

	return nil
}
