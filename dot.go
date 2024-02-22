// dotenv - a wrapper for joho/godotenv
// Copyright (c) 2024 Michael D Henderson. All rights reserved.

// Package dotenv implements a wrapper around John Barton's package for loading dot files.
package dotenv

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io/fs"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	// environments must be sorted by ascending priority.
	environments []string = []string{"development", "test", "production"}
)

// Load tries to emulate the priority list from the dotenv page at
// https://github.com/bkeepers/dotenv#what-other-env-files-can-i-use.
// This is essentially .env.{environment}.local, .env.local,
// .env.{environment}, then .env. Local files take priority over
// environment files, which take priority over the global .env file.
func Load(prefix string, show, verbose bool) error {
	// Pri  Filename______________  Env__  .gitignore?
	// 1st  .env.development.local  dev    yes
	// 1st  .env.test.local         test   yes
	// 1st  .env.production.local   prod   yes
	// 2nd  .env.local              all    yes
	// 3rd  .env.development        dev    no, but be wary of secrets
	// 3rd  .env.test               test   no, but be wary of secrets
	// 3rd  .env.production         prod   no, but be wary of secrets
	// 4th  .env                    all    no, but be wary of secrets
	//
	// Notes:
	//   - The .env.*.local files are for local overrides of environment-specific settings.
	//     We assume that they're created by the deployment process.
	//     They can contain sensitive information like credentials or tokens.
	//     They are loaded first, so they will override settings in the shared files.
	//     They should never be checked into the repository.
	//   - The .env.local file is loaded in development and production; never in test.
	//     It should never be checked into the repository.
	//   - The .env.* files are shared environment-specific settings.
	//     They should not contain sensitive information like credentials or tokens.
	//     They should always be checked into the repository.
	//   - The .env file is loaded in all environments.
	//     It should not contain sensitive information like credentials or tokens.
	//     It is loaded last, so all other files will override any settings.
	//     It should always be checked into the repository.

	if verbose {
		log.Printf("[dot] %-30s == %q\n", "env var prefix", prefix)
	}
	envvar := "ENV"
	if prefix != "" {
		envvar = prefix + "_ENV"
	}
	env := os.Getenv(envvar)
	if verbose {
		log.Printf("[dot] %-30s == %q\n", envvar, env)
	}
	found := false
	for _, environment := range environments {
		if environment == env {
			found = true
			break
		}
	}
	if !found {
		log.Printf("[dot] error: env should be in %v\n", environments)
		return fmt.Errorf("unknown environment")
	}

	// local environment files are the highest priority
	for _, local := range environments {
		if local == env {
			if err := godotenv.Load(".env." + local + ".local"); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					return err
				}
			} else if verbose {
				log.Printf("[dot] loaded %q\n", ".env."+local+".local")
			}
		}
	}

	// .env.local is loaded for all environments except test.
	if env != "test" {
		if err := godotenv.Load(".env.local"); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		} else if verbose {
			log.Printf("[dot] loaded %q\n", ".env.local")
		}
	}

	// shared environment specific settings
	for _, shared := range environments {
		if shared == env {
			if err := godotenv.Load(".env." + shared); err != nil {
				if !errors.Is(err, fs.ErrNotExist) {
					return err
				}
			} else if verbose {
				log.Printf("[dot] loaded %q\n", ".env."+shared)
			}
		}
	}

	// .env is the lowest priority
	if err := godotenv.Load(".env"); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	} else if verbose {
		log.Printf("[dot] loaded %q\n", ".env")
	}

	if show && prefix != "" {
		type kv struct {
			key, value string
		}
		var vars []kv
		for _, v := range os.Environ() {
			if strings.HasPrefix(v, prefix+"_") {
				key, val, _ := strings.Cut(v, "=")
				vars = append(vars, kv{key, val})
			}
		}
		sort.Slice(vars, func(i, j int) bool {
			return vars[i].key < vars[j].key
		})
		for _, v := range vars {
			log.Printf("[dot] %-30s == %q\n", v.key, v.value)
		}
	}

	return nil
}
