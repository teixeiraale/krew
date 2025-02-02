// Copyright 2020 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"os"
	"strings"
	"testing"

	"sigs.k8s.io/krew/internal/environment"
	"sigs.k8s.io/krew/internal/testutil"
)

func TestIsBinDirInPATH_firstRun(t *testing.T) {
	tempDir, cleanup := testutil.NewTempDir(t)
	defer cleanup()

	paths := environment.NewPaths(tempDir.Path("does-not-exist"))
	res := IsBinDirInPATH(paths)
	if res == false {
		t.Errorf("expected positive result on first run")
	}
}

func TestIsBinDirInPATH_secondRun(t *testing.T) {
	tempDir, cleanup := testutil.NewTempDir(t)
	defer cleanup()
	paths := environment.NewPaths(tempDir.Root())
	res := IsBinDirInPATH(paths)
	if res == true {
		t.Errorf("expected negative result on second run")
	}
}

func TestSetupInstructions_windows(t *testing.T) {
	const instructionsContain = "add the\n\"%USERPROFILE%\\.krew\\bin\" directory to your PATH environment variable"
	os.Setenv("KREW_OS", "windows")
	defer func() { os.Unsetenv("KREW_OS") }()
	instructions := SetupInstructions()
	if !strings.Contains(instructions, instructionsContain) {
		t.Errorf("expected %q\nto contain %q", instructions, instructionsContain)
	}
}

func TestSetupInstructions_unix(t *testing.T) {
	tests := []struct {
		name                string
		shell               string
		instructionsContain string
	}{
		{
			name:                "When the shell is zsh",
			shell:               "/bin/zsh",
			instructionsContain: "~/.zshrc",
		},
		{
			name:                "When the shell is bash",
			shell:               "/bin/bash",
			instructionsContain: "~/.bash_profile or ~/.bashrc",
		},
		{
			name:                "When the shell is fish",
			shell:               "/bin/fish",
			instructionsContain: "config.fish",
		},
		{
			name:                "When the shell is unknown",
			shell:               "other",
			instructionsContain: "~/.bash_profile, ~/.bashrc, or ~/.zshrc",
		},
	}

	// always set KREW_OS, so that tests succeed on windows
	os.Setenv("KREW_OS", "linux")
	defer func(origShell string) {
		os.Unsetenv("KREW_OS")
		os.Setenv("SHELL", origShell)
	}(os.Getenv("SHELL"))

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			os.Setenv("SHELL", test.shell)
			instructions := SetupInstructions()
			if !strings.Contains(instructions, test.instructionsContain) {
				tt.Errorf("expected %q\nto contain %q", instructions, test.instructionsContain)
			}
		})
	}
}
