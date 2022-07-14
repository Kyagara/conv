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
	"fmt"
	"log"

	"github.com/Kyagara/conv/internal"
	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var remuxCmd = &cobra.Command{
	Use:   "remux",
	Short: "Remux videos to MP4.",
	Run: func(cmd *cobra.Command, args []string) {
		var files []string

		inputIsDir, err := internal.CheckOutputCompatibility(Input, Output, Recursive)

		if err != nil {
			log.Fatal(err)
		}

		if inputIsDir {
			paths, err := internal.GetFiles(Input, Recursive, internal.ContainersExt[:])

			if err != nil {
				log.Fatal(err)
			}

			files = append(files, paths...)
		} else {
			files = append(files, Input)
		}

		if len(files) == 0 {
			log.Fatal(fmt.Errorf("no files found"))
		}

		for _, file := range files {
			log.Printf("encoding: '%s'", file)

			output, err := internal.NewPathFromDirOutput(file, Output, ".mp4")

			if err != nil {
				log.Println(err)
				continue
			}

			err = remuxToMP4(file, output)

			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("finished: '%s'", output)
		}
	},
}

func remuxToMP4(input string, output string) error {
	err := ffmpeg.Input(input).Output(output, ffmpeg.KwArgs{"map": 0, "codec": "copy"}).GlobalArgs("-stats", "-loglevel", "warning").OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(remuxCmd)
}
