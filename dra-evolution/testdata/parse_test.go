/*
Copyright 20224 The Kubernetes Authors.

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

package podspec

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"

	"github.com/kubernetes-sigs/wg-device-management/dra-evolution/pkg/api"
)

//go:embed *.yaml
var yamls embed.FS

func TestParse(t *testing.T) {
	files, err := yamls.ReadDir(".")
	require.NoError(t, err)

	scheme := runtime.NewScheme()
	require.NoError(t, api.AddToScheme(scheme))
	serializer := json.NewSerializerWithOptions(json.DefaultMetaFactory, scheme, scheme, json.SerializerOptions{Yaml: true, Pretty: true, Strict: true})

	for _, file := range files {
		t.Run(file.Name(), func(t *testing.T) {
			fh, err := yamls.Open(file.Name())
			require.NoError(t, err)
			content, err := io.ReadAll(fh)
			require.NoError(t, err)

			// Split at the "---" separator before working on
			// individual item. Only works for .yaml.
			items := bytes.Split(content, []byte("\n---"))
			if len(items) > 1 {
				for i, item := range items {
					if len(item) > 0 {
						t.Run(fmt.Sprintf("item_%d", i), func(t *testing.T) {
							testDecode(t, serializer, item)
						})
					}
				}
			} else {
				testDecode(t, serializer, content)
			}
		})
	}
}

func testDecode(t *testing.T, serializer *json.Serializer, content []byte) {
	obj, gvk, err := serializer.Decode(content, nil, nil)
	require.NoError(t, err)
	t.Logf("Got object %T = %s", obj, gvk)
}
