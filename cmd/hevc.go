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
	"github.com/gammazero/workerpool"
	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var hevcCmd = &cobra.Command{
	Use:   "hevc",
	Short: "Convert videos to HEVC 10-bit.",
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

		sessions, err := cmd.Flags().GetInt("sessions")

		if err != nil {
			log.Fatal(err)
		}

		quality, err := cmd.Flags().GetInt("quality")

		if err != nil {
			log.Fatal(err)
		}

		cpu, err := cmd.Flags().GetBool("cpu")

		if err != nil {
			log.Fatal(err)
		}

		if cpu {
			sessions = 1
		}

		wp := workerpool.New(sessions)

		for _, file := range files {
			input := file

			wp.Submit(func() {
				log.Printf("encoding: '%s'", input)

				output, err := internal.NewPathFromDirOutput(input, Output, ".mp4")

				if err != nil {
					log.Println(err)
					return
				}

				var path string

				if Output == "" {
					path, err = internal.AddPrefix(output, "HEVC")

					if err != nil {
						log.Println(err)
						return
					}
				} else {
					path = output
				}

				err = encodeHEVC(input, path, cpu, quality)

				if err != nil {
					log.Println(err)
				}

				if sessions > 1 && !cpu && err != nil {
					log.Fatalf("error encoding '%s', if you are using a consumer NVIDIA GPU with support for NVENC, try lowering the session number in -s flag, also make sure no other application like Shadowplay or OBS is using a session", input)
				}

				log.Printf("finished: '%s'", output)
			})
		}

		wp.StopWait()
	},
}

func encodeHEVC(input string, output string, cpu bool, quality int) error {
	hevcArgs := ffmpeg.KwArgs{"map": "0", "map_metadata": "0", "qp": quality, "b:v": "0K", "movflags": "+faststart;use_metadata_tags"}

	if cpu {
		hevcArgs["pix_fmt"] = "yuv420p10le"
		hevcArgs["c:v"] = "libx265"
		hevcArgs["x265-params"] = "log-level=none"
	} else {
		hevcArgs["pix_fmt"] = "p010le"
		hevcArgs["c:v"] = "hevc_nvenc"
		hevcArgs["rc"] = "constqp"
	}

	err := ffmpeg.Input(input, ffmpeg.KwArgs{"hwaccel": "auto"}).Output(output, hevcArgs).GlobalArgs("-stats", "-loglevel", "warning").OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(hevcCmd)

	hevcCmd.Flags().IntP("sessions", "s", 1, "split the work into n NVENC sessions, most consumer NVIDIA GPUs supports a maximum of 3 sessions.")
	hevcCmd.Flags().IntP("quality", "q", 24, "the CQP value, 22 is lossless for HEVC.")
	hevcCmd.Flags().BoolP("cpu", "c", false, "use libx265 instead of hvec_nvenc.")
}
