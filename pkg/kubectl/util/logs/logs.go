/*
Copyright 2014 The Kubernetes Authors.

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

package logs

import (
	"log"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/util/wait"
)

var logFlushFreq = pflag.Duration("log-flush-frequency", 5*time.Second, "Maximum number of seconds between log flushes")

// ErrorLogWriter serves as a bridge between the standard log package and the glog package.
type ErrorLogWriter struct{}

// Write implements the io.Writer interface.
func (writer ErrorLogWriter) Write(data []byte) (n int, err error) {
	return os.Stderr.Write(data)
}

// InitLogs initializes logs the way we want for kubernetes.
func InitLogs() {
	errWriter := ErrorLogWriter{}
	log.SetOutput(errWriter)
	log.SetFlags(0)
	glog.SetOutput(errWriter)
	// The default glog flush interval is 5 seconds.
	go wait.Until(glog.Flush, *logFlushFreq, wait.NeverStop)
}

// FlushLogs flushes logs immediately.
func FlushLogs() {
	glog.Flush()
}

// NewLogger creates a new log.Logger which sends logs to glog.Info.
func NewLogger(prefix string) *log.Logger {
	return log.New(ErrorLogWriter{}, prefix, 0)
}
