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
package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var VideoExt = [8]string{".mp4", ".avi", ".mov", ".wmv", ".ts", ".m2ts", ".mkv", ".mts"}
var ContainersExt = [2]string{".mkv", ".flv"}

func CheckOutputCompatibility(input string, output string, recursive bool) (bool, error) {
	pathInput, err := os.Stat(input)

	if err != nil {
		return false, err
	}

	var pathOutput fs.FileInfo

	if output != "" {
		file, err := os.Stat(output)

		if err != nil {
			return false, err
		}

		pathOutput = file
	}

	// If the input is not a directory and the recursive flag is set, return an error
	if !pathInput.IsDir() && recursive {
		return false, fmt.Errorf("recursive flag is set but the input is not a directory, try changing the input to a directory")
	}

	if pathOutput != nil {
		// If the input is a directory but the output is not, return an error
		if pathInput.IsDir() && !pathOutput.IsDir() {
			return false, fmt.Errorf("the output provided is a file but the input was a directory, try changing the output to a directory")
		}
	}

	return pathInput.IsDir(), nil
}

func GetFiles(path string, recursive bool, extensions []string) ([]string, error) {
	f, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	files, err := f.Readdir(0)

	if err != nil {
		return nil, err
	}

	var paths []string

	for _, file := range files {
		if file.IsDir() && recursive {
			files, err := GetFiles(fmt.Sprintf("%s%s\\", path, file.Name()), recursive, extensions)

			if err != nil {
				continue
			}

			paths = append(paths, files...)

			continue
		}

		if file.IsDir() && !recursive {
			continue
		}

		if !checkExt(file.Name(), extensions) {
			continue
		}

		paths = append(paths, fmt.Sprintf("%s%s", path, file.Name()))
	}

	return paths, nil
}

func AddPrefix(input string, prefix string) (string, error) {
	file, err := os.Stat(input)

	if err != nil {
		return "", err
	}

	filename := file.Name()

	newName := fmt.Sprintf("%s-%s", prefix, filename)

	return fmt.Sprintf("%s%s", strings.TrimSuffix(input, filename), newName), nil
}

func ReplaceExtension(input string, extension string) string {
	return strings.TrimSuffix(input, filepath.Ext(input)) + extension
}

func NewPathFromDirOutput(input string, output string, extension string) (string, error) {
	if output == "" {
		file, err := os.Stat(input)

		if err != nil {
			return "", err
		}

		if !file.IsDir() {
			return ReplaceExtension(input, extension), nil
		}

		name := ReplaceExtension(input, extension)

		return fmt.Sprintf("%s%s", input, name), nil
	}

	file, err := os.Stat(output)

	if err != nil {
		return "", err
	}

	if !file.IsDir() {
		return output, nil
	}

	file, err = os.Stat(input)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\\%s", output, file.Name()), nil
}

func checkExt(file string, extensions []string) bool {
	ext := filepath.Ext(file)

	for _, extension := range extensions {
		if extension == ext {
			return true
		}
	}

	return false
}
