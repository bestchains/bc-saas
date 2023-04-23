/*
Copyright 2023 The Bestchains Authors.

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

package listener

import (
	"context"
	"time"

	"k8s.io/klog/v2"
)

type logListner struct {
	x int
}

func NewLogListener() Listener {
	return &logListner{x: 0}
}
func (ll *logListner) Events(ctx context.Context) {
	klog.Infoln("starting logListener")
	for {
		select {
		case <-ctx.Done():
			klog.Infof("logListener done")
			return
		default:
			ll.x = (ll.x + 1) % 10000
			klog.Infof("default logListener generate a number %d", ll.x)
			<-time.After(5 * time.Second)
		}
	}
}
