package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/whoshuu/ignoreit/generate"
	"github.com/whoshuu/ignoreit/spec"
)

const (
	configFilename = ".ignoreit.yml"
	ignoreFilename = ".gitignore"
	defaultRepo    = "github/gitignore"
	defaultBranch  = "master"
)

func main() {
	config, err := spec.LoadConfig(configFilename)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	app := cli.NewApp()
	app.Name = "ignoreit"
	app.Usage = "Manage .gitignore templates declaratively"

	var repo string
	var branch string
	addAndRemoveFlags := []cli.Flag{
		cli.StringFlag{
			Name:        "repo, r",
			Value:       defaultRepo,
			Usage:       "uses .gitignore files from https://github.com/`REPO`",
			Destination: &repo,
		},
		cli.StringFlag{
			Name:        "branch, b",
			Value:       defaultBranch,
			Usage:       "git `BRANCH` of the REPO",
			Destination: &branch,
		},
	}
	app.Commands = []cli.Command{
		{
			// Add source and branch flags
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add entries to .ignoreit.yml",
			Flags:   addAndRemoveFlags,
			Action: func(c *cli.Context) error {
				source := config.CreateSource(repo, branch)
				var err error
				if source != nil {
					for _, entry := range c.Args() {
						if err = source.AddEntry(entry); err != nil {
							log.Fatalf("Error adding entry: %v", err)
						}
					}
					return config.Save(configFilename)
				}
				return err
			},
		},
		{
			Name:    "remove",
			Aliases: []string{"rm"},
			Usage:   "remove entries to .ignoreit.yml",
			Flags:   addAndRemoveFlags,
			Action: func(c *cli.Context) error {
				source := config.GetSource(repo, branch)
				var err error
				if source != nil {
					for _, entry := range c.Args() {
						if err = source.RemoveEntry(entry); err != nil {
							log.Fatalf("Error removing entry: %v", err)
						}
					}
					return config.Save(configFilename)
				}
				return err
			},
		},
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate a .gitignore from .ignoreit.yml",
			Action: func(c *cli.Context) error {
				return generate.Inflate(config, ignoreFilename)
			},
		},
	}

	app.Run(os.Args)
}
