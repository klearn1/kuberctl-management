// Copyright 2017 The etcd Authors
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

package clientv3

import (
	"fmt"
	"math"
	"time"

	"google.golang.org/grpc"
)

var (
	// client-side handling retrying of request failures where data was not written to the wire or
	// where server indicates it did not process the data. gRPC default is default is "WaitForReady(false)"
	// but for etcd we default to "WaitForReady(true)" to minimize client request error responses due to
	// transient failures.
	defaultWaitForReady = grpc.WaitForReady(true)

	// client-side request send limit, gRPC default is math.MaxInt32
	// Make sure that "client-side send limit < server-side default send/recv limit"
	// Same value as "embed.DefaultMaxRequestBytes" plus gRPC overhead bytes
	defaultMaxCallSendMsgSize = grpc.MaxCallSendMsgSize(2 * 1024 * 1024)

	// client-side response receive limit, gRPC default is 4MB
	// Make sure that "client-side receive limit >= server-side default send/recv limit"
	// because range response can easily exceed request send limits
	// Default to math.MaxInt32; writes exceeding server-side send limit fails anyway
	defaultMaxCallRecvMsgSize = grpc.MaxCallRecvMsgSize(math.MaxInt32)

	// client-side non-streaming retry limit, only applied to requests where server responds with
	// a error code clearly indicating it was unable to process the request such as codes.Unavailable.
	// If set to 0, retry is disabled.
	defaultUnaryMaxRetries uint = 100

	// client-side streaming retry limit, only applied to requests where server responds with
	// a error code clearly indicating it was unable to process the request such as codes.Unavailable.
	// If set to 0, retry is disabled.
	defaultStreamMaxRetries = ^uint(0) // max uint

	// client-side retry backoff wait between requests.
	defaultBackoffWaitBetween = 25 * time.Millisecond

	// client-side retry backoff default jitter fraction.
	defaultBackoffJitterFraction = 0.10
)

// defaultCallOpts defines a list of default "gRPC.CallOption".
// Some options are exposed to "clientv3.Config".
// Defaults will be overridden by the settings in "clientv3.Config".
var defaultCallOpts = []grpc.CallOption{
	defaultWaitForReady,
	defaultMaxCallSendMsgSize,
	defaultMaxCallRecvMsgSize,
}

// MaxLeaseTTL is the maximum lease TTL value
const MaxLeaseTTL = 9000000000

func Options(cfg *Config) ([]grpc.CallOption, error) {
	callOpts := defaultCallOpts
	if cfg.MaxCallSendMsgSize > 0 || cfg.MaxCallRecvMsgSize > 0 {
		if cfg.MaxCallRecvMsgSize > 0 && cfg.MaxCallSendMsgSize > cfg.MaxCallRecvMsgSize {
			return nil, fmt.Errorf("gRPC message recv limit (%d bytes) must be greater than send limit (%d bytes)", cfg.MaxCallRecvMsgSize, cfg.MaxCallSendMsgSize)
		}
		callOpts = []grpc.CallOption{
			defaultWaitForReady,
			defaultMaxCallSendMsgSize,
			defaultMaxCallRecvMsgSize,
		}
		if cfg.MaxCallSendMsgSize > 0 {
			callOpts[1] = grpc.MaxCallSendMsgSize(cfg.MaxCallSendMsgSize)
		}
		if cfg.MaxCallRecvMsgSize > 0 {
			callOpts[2] = grpc.MaxCallRecvMsgSize(cfg.MaxCallRecvMsgSize)
		}
	}
	return callOpts, nil
}
