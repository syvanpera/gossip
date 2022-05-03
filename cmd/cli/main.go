// Copyright © 2022 Tuomo Syvänperä <tuomo.syvanpera@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/syvanpera/gossip/pkg/cmd"
	"github.com/syvanpera/gossip/pkg/util"
)

var appName = "gossip"

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath(filepath.Clean(fmt.Sprintf("%s/%s", util.ConfigPath(), appName)))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("config.color", true)
	viper.SetDefault("config.browser", "default")
	viper.SetDefault("database.path", filepath.Clean(fmt.Sprintf("%s/%s/%s.db", util.DataPath(), appName, appName)))
}

func main() {
	cmd.Execute()
}
