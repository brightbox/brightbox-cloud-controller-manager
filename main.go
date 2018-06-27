// Copyright 2018 Brightbox Systems Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
The external controller manager is responsible for running controller
loops that are cloud provider dependent. It uses the API to listen to
new events on resources.
*/

package main

import (
	goflag "flag"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	_ "github.com/brightbox/brightbox-cloud-controller-manager/brightbox"
	utilflag "k8s.io/apiserver/pkg/util/flag"
	"k8s.io/apiserver/pkg/util/logs"
	"k8s.io/kubernetes/cmd/cloud-controller-manager/app"
	_ "k8s.io/kubernetes/pkg/client/metrics/prometheus" // for client metric registration
	_ "k8s.io/kubernetes/pkg/version/prometheus"        // for version metric registration

	"github.com/golang/glog"
	"github.com/spf13/pflag"
)

/*
These are exogenous constants initialised at compile time by the
compiler command line.

Yes the syntax for these should be better. Naughty golang.
*/
var version string
var build string

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	command := app.NewCloudControllerManagerCommand()

	/* TODO: once we switch everything over to Cobra commands, we
	can go back to calling utilflag.InitFlags() (by removing its
	pflag.Parse() call). For now, we have to set the normalize func
	and add the go flag set by hand.
	*/
	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	// utilflag.InitFlags()

	// Workaround for this issue:
	// https://github.com/kubernetes/kubernetes/issues/17162
	goflag.CommandLine.Parse([]string{})

	logs.InitLogs()
	defer logs.FlushLogs()

	if glog.V(1) {
		glog.Infof("%s version: %s (%s)", filepath.Base(os.Args[0]), version, build)
	}

	if err := command.Execute(); err != nil {
		glog.Exitf("error: %v", err)
	}
}
