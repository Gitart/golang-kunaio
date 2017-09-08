// Copyright 2017 Aleksey Morarash <tuxofil@gmail.com>
//
// Licensed under the BSD 2 Clause License (the "License");
// you may not use the file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://opensource.org/licenses/BSD-2-Clause
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kunaio

import (
	"log"
	"os"
)

var (
	gDebug bool
)

func init() {
	gDebug = os.Getenv("KUNAIO_DEBUG") != ""
}

func debugLog(format string, args ...interface{}) {
	if gDebug {
		log.Printf("DEBUG: "+format, args...)
	}
}
