/*
Copyright © 2024 Dragon Security

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/dragonsecurity/dev-setup/efs"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime/debug"

	"github.com/leaanthony/debme"
	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "install and setup git on your machine",
	Long:  `install and setup git on your machine.`,
	Run: func(cmd *cobra.Command, args []string) {
		isVerbose := cmd.Flag("verbose").Value.String() == "true"

		logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
		err := run(logger, isVerbose)
		if err != nil {
			trace := string(debug.Stack())
			logger.Error(err.Error(), "trace", trace)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)
}

func run(logger *slog.Logger, isVerbose bool) error {
	scriptContent, err := efs.EmbeddedFiles.ReadFile("scripts/git.sh")
	if err != nil {
		logger.Error("Failed to read embedded script", slog.Any("error", err))
	}

	// Create a temporary file to write the script
	tmpFile, err := os.CreateTemp("", "script-*.sh")
	if err != nil {
		logger.Error("Failed to create temp file", slog.Any("error", err))
	}

	// Clean up the file later
	defer os.Remove(tmpFile.Name())

	// Write the script to the temporary file
	if _, err := tmpFile.Write(scriptContent); err != nil {
		logger.Error("Failed to write script to temp file", slog.Any("error", err))
	}
	if err := tmpFile.Close(); err != nil {
		logger.Error("Failed to close temp file", slog.Any("error", err))
	}

	// Make the script executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		logger.Error("Failed to chmod temp file", slog.Any("error", err))
	}

	fmt.Println(tmpFile.Name())
	localCmd := exec.Command("/bin/bash", tmpFile.Name())

	if isVerbose {
		localCmd.Stdout = os.Stdout
	}
	localCmd.Stderr = os.Stderr

	if err := localCmd.Run(); err != nil {
		logger.Error("Failed to execute script", slog.Any("error", err))
	}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(usr.HomeDir)

	root, _ := debme.FS(efs.EmbeddedFiles, "dotfiles")
	dstFile := path.Join(usr.HomeDir, ".gitconfig")
	err = root.CopyFile("gitconfig.txt", dstFile, 0644)
	if err != nil {
		logger.Error("Failed to copy .gitconfig file", slog.Any("error", err))
	}
	dstFile = path.Join(usr.HomeDir, ".gitignore")
	err = root.CopyFile("gitignore.txt", dstFile, 0644)
	return nil
}
