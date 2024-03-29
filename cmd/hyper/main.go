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

package main

import (
	"github.com/kiprotect/hyper"
	"github.com/kiprotect/hyper/cmd/helpers"
	"github.com/kiprotect/hyper/definitions"
	"github.com/urfave/cli"
	"os"
)

func main() {
	if settings, err := helpers.Settings(&definitions.Default); err != nil {
		hyper.Log.Error(err)
		return
	} else {
		CLI(settings)
	}
}

func CLI(settings *hyper.Settings) {

	var err error

	app := cli.NewApp()
	app.Name = "Endpoint Server"
	app.Usage = "Run all server commands"
	app.Flags = helpers.CommonFlags

	bareCommands := helpers.CommonCommands

	// we add commands from the definitions
	for _, commandsDefinition := range settings.Definitions.CommandsDefinitions {
		if commands, err := commandsDefinition.Maker(settings); err != nil {
			hyper.Log.Fatal(err)
		} else {
			bareCommands = append(bareCommands, commands...)
		}
	}

	app.Commands = helpers.Decorate(bareCommands, helpers.InitCLI, "Hyper")

	err = app.Run(os.Args)

	if err != nil {
		hyper.Log.Error(err)
	}

}
