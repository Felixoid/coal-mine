/*
Copyright © 2021 Mikhail f. Shiryaev <mr.felixoid@gmail.com>

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
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Felixoid/coal-mine/cmd"
)

func main() {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	signal.Notify(cmd.CatchedSignals, signals...)
	cmd.Execute()
}
