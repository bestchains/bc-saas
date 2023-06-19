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

	"github.com/bestchains/bc-saas/pkg/events"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"k8s.io/klog/v2"
)

// listener listens chaincode events
type listener struct {
	// event registered and its handler
	registeredEvents map[events.Event]events.EventHandler

	eventsSub <-chan *client.ChaincodeEvent
}

func NewListener(eventsSub <-chan *client.ChaincodeEvent, registeredEvents map[events.Event]events.EventHandler) (Listener, error) {
	l := &listener{
		eventsSub: eventsSub,
	}

	l.eventsSub = eventsSub

	// registeredEvents
	if registeredEvents == nil {
		registeredEvents = make(map[events.Event]events.EventHandler)
	}
	l.registeredEvents = registeredEvents

	return l, nil
}

func (l *listener) Events(ctx context.Context) {
	klog.Info("starting fetch events")
	var err error
	for {
		select {
		case e := <-l.eventsSub:
			// check whether event registered
			eventHandler, ok := l.registeredEvents[events.Event(e.EventName)]
			if !ok {
				klog.Warningf("Event %s not registered, skip", e.EventName)
				continue
			}
			err = eventHandler(e)
			if err != nil {
				klog.Errorf("[Error] handle event %s erorr %s", e.EventName, err.Error())
				continue
			}

			klog.V(5).Infof("[Debug] event %+v", *e)

		case <-ctx.Done():
			klog.Info("context break down")
			return
		}
	}
}
