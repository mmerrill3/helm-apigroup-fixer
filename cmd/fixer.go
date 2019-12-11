/*
Copyright Mike Merrill.

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

package main // import "github.com/mmerrill3/helm-agpigroup-fixer/cmd"

import (
	"log"
	"github.com/mmerrill3/helm-apigroup-fixer/pkg/version"
	"github.com/mmerrill3/helm-apigroup-fixer/pkg/kube"
	"bufio"
	"os"
	"fmt"
)

//Main is the starting point for our app
func main() {

	logger := newLogger("main")

	logger.Printf("Starting helm configmap fixer version %s", version.GetVersion())

	//get the chart name and the tiller namespace to look at
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the chart: ")
	chart, err := reader.ReadString("\n")
	if err != nil {
		logger.Fatal("did not understand the chart")
	}
	fmt.Print("Enter the tiller namespace: ")
	tillerNamespace, err := reader.ReadString("\n")
	if err != nil {
		logger.Fatal("did not understand the tiller namespace")
	}
	//get the configmap for the release
	//get the configmap client
	clientset, err := kube.New(nil).KubernetesClientSet()
	if err != nil {
		logger.Fatalf("Cannot initialize Kubernetes connection: %s", err)
	}

	//update the manifest for deployments in memory

	//write back the configmap

	//voila!
}

func newLogger(prefix string) *log.Logger {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("[%s] ", prefix)
	}
	return log.New(os.Stderr, prefix, log.Flags())
}
