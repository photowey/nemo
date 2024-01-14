/*
 * Copyright © 2023 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package filez

import (
	"os"
	"path/filepath"
)

func ToAbsIfNecessary(path string) (string, error) {
	if IsAbs(path) {
		return path, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func IsFile(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	if IsDir(path) {
		return false
	}

	return true
}

func IsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}

	return file.IsDir()
}
