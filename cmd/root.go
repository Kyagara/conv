/*
Copyright © 2022 Kyagara

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

var Input string
var Output string
var Recursive bool

var rootCmd = &cobra.Command{
	Use: "conv",
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cmd := exec.Command("ffmpeg", "-version")

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	rootCmd.PersistentFlags().StringVarP(&Input, "input", "i", "", "input file/path.")
	rootCmd.PersistentFlags().StringVarP(&Output, "output", "o", "", "output file/path.")
	rootCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "search for files inside folders in the directory provided.")

	err = rootCmd.MarkPersistentFlagRequired("input")

	if err != nil {
		log.Fatal(err)
	}
}
