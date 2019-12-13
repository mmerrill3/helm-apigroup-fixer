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

package main // import "github.com/mmerrill3/helm-apigroup-fixer/cmd"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/mmerrill3/helm-apigroup-fixer/pkg/kube"
	"github.com/mmerrill3/helm-apigroup-fixer/pkg/storage/driver"
	"github.com/mmerrill3/helm-apigroup-fixer/pkg/version"
	//"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	//typev1beta1 "k8s.io/api/apps/v1beta1"
	//"k8s.io/apimachinery/pkg/runtime"
)

var (
	ingest   = flag.Bool("ingest", false, "whether to read manifest.yaml an apply it, otherwise, produce manifest.yaml")
	release  = flag.String("release", "", "tiller release name to manage")
	tillerns = flag.String("tiller-namespace", "kube-system", "tiller namespace that holds the release name as a config map")
)

//Main is the starting point for our app
func main() {
	flag.Parse()

	logger := newLogger("main")

	logger.Printf("Starting helm configmap fixer version %s", version.GetVersion())
	logger.Printf("Ingesting: %s", *ingest)
	if *release == "" {
		logger.Fatal("need to pass in a release, like dev-wildcard-cert.v25")
	}

	//get the configmap for the release
	//get the configmap client
	clientset, err := kube.New(nil).KubernetesClientSet()
	if err != nil {
		logger.Fatalf("Cannot initialize Kubernetes connection: %s", err)
	}
	fmt.Printf("got k8s client: %+v\n", clientset)

	cfgmaps := driver.NewConfigMaps(clientset.CoreV1().ConfigMaps(*tillerns))
	//fmt.Printf("got k8s cfgmaps client: %+v\n", cfgmaps)

	storedRelease, err := cfgmaps.Get(*release)
	if err != nil {
		logger.Fatalf("Cannot find release: %s", err)
	}
	//update the manifest for deployments in memory

	if !(*ingest) {
		manifest := storedRelease.Manifest
		fmt.Printf("got manifest: %+v\n", manifest)
		//write to file
		err = ioutil.WriteFile("manifest.yaml", []byte(manifest), 0644)
	} else {
		bmanifest, err := ioutil.ReadFile("manifest.yaml")
		if err != nil {
			logger.Fatalf("Cannot find manifest.yaml: %s", err)
		}
		storedRelease.Manifest = string(bmanifest)
		err = cfgmaps.Update(*release, storedRelease)
		if err != nil {
			logger.Fatalf("Cannot find update the stored configmap: %s", err)
		}
	}
	/*
		options := genericclioptions.NewConfigFlags()
		ddConfig := "/Users/mmerrill/.kube/configDockerDesktop"
		options.KubeConfig = &ddConfig
		kubeClient := kube.New(options)
		fmt.Printf("got kubeClient: %+v\n", kubeClient)
		kubeClient.Log = newLogger("kube").Printf

		result, err := kubeClient.BuildUnstructured(tillerNamespace, bytes.NewBufferString(manifest))
		if err != nil {
			logger.Fatalf("Cannot parse manifest: %s", err)
		}
		fmt.Printf("got result: %+v\n", result)

		//dev-wildcard-cert.v25
		//kube-system

		info := &resource.Info{
			Namespace: tillerNamespace,
			Name:      "wildcard-cert-manager",
			Mapping: &meta.RESTMapping{
				GroupVersionKind: schema.GroupVersionKind{
					Group:   "extensions",
					Version: "v1beta1",
					Kind:    "Deployment",
				},
			},
		}

		filteredInfo := result.Get(info)
		if filteredInfo == nil {
			logger.Fatalf("Cannot find filterd Info: %+v\n", *info)
		}
		fmt.Printf("got filtered info: %+v\n", filteredInfo.Object)

	*/
	//convert to apps/v1
	//v1beta1.Convert_v1beta1_DeploymentSpec_To_apps_DeploymentSpec()
	//write back the configmap
	//voila!
}

func newLogger(prefix string) *log.Logger {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("[%s] ", prefix)
	}
	return log.New(os.Stderr, prefix, log.Flags())
}
