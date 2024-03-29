// KIProtect Hyper
// Copyright (C) 2021-2023 KIProtect GmbH
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package helpers

import (
	"fmt"
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/helpers"
	"github.com/urfave/cli"
)

type Decorator func(f func(c *cli.Context) error, service string) func(c *cli.Context) error

func Decorate(commands []cli.Command, decorator Decorator, service string) []cli.Command {
	newCommands := make([]cli.Command, len(commands))
	for i, command := range commands {
		if command.Action != nil {
			command.Action = decorator(command.Action.(func(c *cli.Context) error), service)
		}
		if command.Subcommands != nil {
			command.Subcommands = Decorate(command.Subcommands, decorator, service)
		}
		newCommands[i] = command
	}
	return newCommands
}

func Settings(definitions *hyper.Definitions) (*hyper.Settings, error) {
	if settingsPaths, fs, err := helpers.SettingsPaths("HYPER_SETTINGS"); err != nil {
		return nil, err
	} else {
		return helpers.Settings(settingsPaths, fs, definitions)
	}
}

var CommonCommands = []cli.Command{
	{
		Name:   "version",
		Usage:  "Print the software version",
		Action: func(c *cli.Context) error { fmt.Println(hyper.Version); return nil },
	},
}

var CommonFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "level",
		Value: "info",
		Usage: "The desired log level",
	},
	cli.StringFlag{
		Name:  "format",
		Value: "",
		Usage: "The desired log format (possible values: iris)",
	},
	cli.StringFlag{
		Name:  "profile",
		Value: "",
		Usage: "enable profiler and store results to given filename",
	},
}

func InitCLI(f func(c *cli.Context) error, service string) func(c *cli.Context) error {
	return func(c *cli.Context) error {

		level := c.GlobalString("level")
		logLevel, err := hyper.ParseLevel(level)
		if err != nil {
			return fmt.Errorf("error parsing flag: %w", err)
		}
		hyper.Log.SetLevel(logLevel)

		hyper.Log.Debugf("%s version: %s", service, hyper.Version)

		format := c.GlobalString("format")
		if format != "" {
			if err := hyper.SetLogFormat(format, service); err != nil {
				return fmt.Errorf("error setting log formatter: %w", err)
			}
		}

		runner := func() error { return f(c) }
		profiler := c.GlobalString("profile")
		if profiler != "" {
			return runWithProfiler(profiler, runner)
		}

		return f(c)
	}
}
