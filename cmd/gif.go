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

var gifCmd = &cobra.Command{
	Use:   "gif",
	Short: "Converts videos to a gif.",
	Run: func(cmd *cobra.Command, args []string) {
		var files []string

		inputIsDir, err := internal.CheckOutputCompatibility(Input, Output, Recursive)

		if err != nil {
			log.Fatal(err)
		}

		if inputIsDir {
			paths, err := internal.GetFiles(Input, Recursive, internal.VideoExt[:])

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

		from, err := cmd.Flags().GetInt("from")

		if err != nil {
			log.Fatal(err)
		}

		gifArgs := ffmpeg.KwArgs{"vf": "fps=10,scale=-1:240:flags=lanczos,split[s0][s1];[s0]palettegen=max_colors=64[p];[s1][p]paletteuse=dither=bayer"}

		if !inputIsDir {
			if from != 0 {
				gifArgs["ss"] = from
			}

			to, _ := cmd.Flags().GetInt("to")

			if to != 0 {
				gifArgs["to"] = to
			}
		}

		for _, file := range files {
			log.Printf("encoding: '%s'", file)

			output, err := internal.NewPathFromDirOutput(file, Output, ".gif")

			if err != nil {
				log.Println(err)
				return
			}

			err = converToGIF(file, output, gifArgs)

			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("finished: '%s'", output)
		}
	},
}

func converToGIF(input string, output string, gifArgs ffmpeg.KwArgs) error {
	err := ffmpeg.Input(input).Output(output, gifArgs).GlobalArgs("-stats", "-loglevel", "warning").OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(gifCmd)

	gifCmd.Flags().IntP("from", "f", 0, "decodes but discards input until the timestamps reaches the specified value.")
	gifCmd.Flags().IntP("to", "t", 0, "stop writing the output after its duration reaches the specified value.")
}
