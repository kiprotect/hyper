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
	"os"
	"runtime"
	"runtime/pprof"
)

func runWithProfiler(name string, runner func() error) error {

	hyper.Log.Info("Running with profiler...")

	fc, err := os.Create(fmt.Sprintf("%s-cpu.pprof", name))

	if err != nil {
		return fmt.Errorf("error creating CPU profile file: %w", err)
	}

	if err := pprof.StartCPUProfile(fc); err != nil {
		return fmt.Errorf("error starting CPU profiling: %w", err)
	}

	defer pprof.StopCPUProfile()

	runnerErr := runner()

	fm, err := os.Create(fmt.Sprintf("%s-mem.pprof", name))

	if err != nil {
		return fmt.Errorf("error creating MEM profile file: %w", err)
	}

	runtime.GC() // get up-to-date statistics

	if err := pprof.WriteHeapProfile(fm); err != nil {
		return fmt.Errorf("error writing heap profile: %w", err)
	}

	fm.Close()

	return runnerErr
}
