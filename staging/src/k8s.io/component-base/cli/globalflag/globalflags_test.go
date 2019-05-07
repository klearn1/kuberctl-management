/*
Copyright 2018 The Kubernetes Authors.

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

package globalflag

import (
	"flag"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/diff"
	cliflag "k8s.io/component-base/cli/flag"
)

func TestAddGlobalFlags(t *testing.T) {
	namedFlagSets := &cliflag.NamedFlagSets{}
	nfs := namedFlagSets.FlagSet("global")
	AddGlobalFlags(nfs, "test-cmd")

	actualFlag := []string{}
	nfs.VisitAll(func(flag *pflag.Flag) {
		actualFlag = append(actualFlag, flag.Name)
	})

	// Get all flags from flags.CommandLine, except flag `test.*`.
	wantedFlag := []string{"help"}
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.VisitAll(func(flag *pflag.Flag) {
		if !strings.Contains(flag.Name, "test.") {
			wantedFlag = append(wantedFlag, normalize(flag.Name))
		}
	})
	sort.Strings(wantedFlag)

	if !reflect.DeepEqual(wantedFlag, actualFlag) {
		t.Errorf("[Default]: expected %+v, got %+v", wantedFlag, actualFlag)
	}

	tests := []struct {
		expectedFlag  []string
		matchExpected bool
	}{
		{
			// Happy case
			expectedFlag:  []string{"alsologtostderr", "help", "log-backtrace-at", "log-dir", "log-file", "log-file-max-size", "log-flush-frequency", "logtostderr", "skip-headers", "skip-log-headers", "stderrthreshold", "v", "vmodule"},
			matchExpected: false,
		},
		{
			// Missing flag
			expectedFlag:  []string{"logtostderr", "log-dir"},
			matchExpected: true,
		},
		{
			// Empty flag
			expectedFlag:  []string{},
			matchExpected: true,
		},
		{
			// Invalid flag
			expectedFlag:  []string{"foo"},
			matchExpected: true,
		},
	}

	for i, test := range tests {
		if reflect.DeepEqual(test.expectedFlag, actualFlag) == test.matchExpected {
			t.Errorf("[%d]: expected %+v, got %+v", i, test.expectedFlag, actualFlag)
		}
	}
}

func TestSideEffect(t *testing.T) {
	// construct an initial flagset
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	AddGlobalFlags(fs, "")

	// change some of the values from their defaults and record the new values
	want := map[string]string{}
	fs.VisitAll(func(f *pflag.Flag) {
		// change values for a few basic types
		switch vt := f.Value.Type(); {
		case vt == "string":
			if f.Value.String() == "foo" {
				f.Value.Set("bar")
			} else {
				f.Value.Set("foo")
			}
		case vt == "bool":
			if f.Value.String() == "true" {
				f.Value.Set("false")
			} else {
				f.Value.Set("true")
			}
		case strings.Contains(vt, "int"):
			if f.Value.String() == "1" {
				f.Value.Set("2")
			} else {
				f.Value.Set("1")
			}
		}
		// record the values as set
		want[f.Name] = f.Value.String()
	})

	// construct another flagset
	fs2 := pflag.NewFlagSet("", pflag.ExitOnError)
	AddGlobalFlags(fs2, "")

	// check if the values in the original flagset still match what we recorded
	got := map[string]string{}
	fs.VisitAll(func(f *pflag.Flag) {
		got[f.Name] = f.Value.String()
	})

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("global flag registration has side effects. Diff:\n%s",
			diff.ObjectDiff(want, got))
	}
}
