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
	cmdHelpers "github.com/kiprotect/hyper/cmd/helpers"
	"github.com/kiprotect/hyper/definitions"
	"github.com/kiprotect/hyper/helpers"
	"github.com/kiprotect/hyper/metrics"
	"github.com/kiprotect/hyper/proxy"
	proxyHelpers "github.com/kiprotect/hyper/proxy/helpers"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

func CLI(settings *proxy.Settings) {

	var err error

	app := cli.NewApp()
	app.Name = "Proxy Server"
	app.Usage = "Run all proxy server commands"
	app.Flags = cmdHelpers.CommonFlags

	bareCommands := append([]cli.Command{
		{
			Name:  "run",
			Usage: "Run the public or private proxy",
			Subcommands: []cli.Command{
				{
					Name:  "public",
					Flags: []cli.Flag{},
					Usage: "Run the public TLS proxy.",
					Action: func(c *cli.Context) error {
						hyper.Log.Info("Starting public proxy...")

						if settings.Public == nil {
							hyper.Log.Fatal("Public settings undefined!")
						}

						server, err := proxy.MakePublicServer(settings.Public, settings.Definitions)

						if err != nil {
							hyper.Log.Fatal(err)
						}

						if err := server.Start(); err != nil {
							hyper.Log.Fatal(err)
						}

						metricsServer := metrics.MakePrometheusMetricsServer(settings.Metrics)

						// we wait for CTRL-C / Interrupt
						sigchan := make(chan os.Signal, 1)
						signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

						hyper.Log.Info("Waiting for CTRL-C...")

						<-sigchan

						hyper.Log.Info("Stopping proxy...")

						if err := server.Stop(); err != nil {
							hyper.Log.Fatal(err)
						}

						if metricsServer != nil {
							if err := metricsServer.Stop(); err != nil {
								hyper.Log.Fatal(err)
							}
						}

						return nil
					},
				},
				{
					Name:  "private",
					Flags: []cli.Flag{},
					Usage: "Run the private TLS proxy.",
					Action: func(c *cli.Context) error {
						hyper.Log.Info("Starting private proxy...")

						if settings.Private == nil {
							hyper.Log.Fatal("Private settings undefined!")
						}

						server, err := proxy.MakePrivateServer(settings.Private, settings.Definitions)

						if err != nil {
							hyper.Log.Fatal(err)
						}

						if err := server.Start(); err != nil {
							hyper.Log.Fatal(err)
						}

						metricsServer := metrics.MakePrometheusMetricsServer(settings.Metrics)

						// we wait for CTRL-C / Interrupt
						sigchan := make(chan os.Signal, 1)
						signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

						hyper.Log.Info("Waiting for CTRL-C...")

						<-sigchan

						hyper.Log.Info("Stopping proxy...")

						if err := server.Stop(); err != nil {
							hyper.Log.Fatal(err)
						}

						if metricsServer != nil {
							if err := metricsServer.Stop(); err != nil {
								hyper.Log.Fatal(err)
							}
						}

						return nil
					},
				},
			},
		},
	}, cmdHelpers.CommonCommands...)

	app.Commands = cmdHelpers.Decorate(bareCommands, cmdHelpers.InitCLI, "PRO")

	err = app.Run(os.Args)

	if err != nil {
		hyper.Log.Error(err)
	}

}

func main() {
	if settingsPaths, fs, err := helpers.SettingsPaths("PROXY_SETTINGS"); err != nil {
		hyper.Log.Error(err)
		return
	} else if settings, err := proxyHelpers.Settings(settingsPaths, fs, &definitions.Default); err != nil {
		hyper.Log.Error(err)
		return
	} else {
		CLI(settings)
	}
}
