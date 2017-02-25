package spec

import (
	"sort"

	"github.com/whoshuu/ignoreit/network"
)

// Source represents a collection of .gitignore resources.
// Repo and Branch uniquely identify a remote repository of .gitignore files.
// Entries is a list of files to sync with, exluding the .gitignore suffix. Ex: Go is a valid entry.
type Source struct {
	Repo    string   `yaml:"repo"`
	Branch  string   `yaml:"branch"`
	Entries []string `yaml:"entries"`
}

// Sources is a collection of Source structs
type Sources []Source

func (sources Sources) Len() int {
	return len(sources)
}

func (sources Sources) Less(i, j int) bool {
	if sources[i].Repo == sources[j].Repo {
		return sources[i].Branch < sources[j].Branch
	}

	return sources[i].Repo < sources[j].Repo
}

func (sources Sources) Swap(i, j int) {
	sources[i], sources[j] = sources[j], sources[i]
}

// GetDownloadLink returns the link to download a raw form of the entry from the source.
func (source Source) GetDownloadLink(entry string) string {
	return "https://raw.githubusercontent.com/" + source.Repo + "/" + source.Branch + "/" + entry + ".gitignore"
}

// AddEntry adds the entry to the source.Entries slice.
// If the entry already exists, nothing is modified and this method returns early.
func (source *Source) AddEntry(entry string) error {
	for _, existingEntry := range source.Entries {
		if existingEntry == entry {
			return nil
		}
	}

	if network.EntryExists(source.GetDownloadLink(entry)) {
		source.Entries = append(source.Entries, entry)
	}

	return nil
}

// RemoveEntry removes the entry from the source.Entries slice.
// If the entry is removed, every following entry is pushed back to keep the slice tight.
// Calling this multiple times in a row may be inefficient.
func (source *Source) RemoveEntry(entry string) error {
	for i, existingEntry := range source.Entries {
		if existingEntry == entry {
			source.Entries = append(source.Entries[:i], source.Entries[i+1:]...)
			break
		}
	}

	return nil
}

// Clean sorts and dedupes source.Entries, where deduping is the equivalent of removing an entry.
// The resulting source.Entries should be a tightly packed, sorted, and unique slice of strings.
func (source *Source) Clean() error {
	sort.Strings(source.Entries)
	for i := 0; i < len(source.Entries)-1; {
		if source.Entries[i] == source.Entries[i+1] {
			source.Entries = append(source.Entries[:i], source.Entries[i+1:]...)
		} else {
			i++
		}
	}
	return nil
}
