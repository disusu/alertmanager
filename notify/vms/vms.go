// Copyright 2022 Prometheus Team
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vms

import (
	"context"
	"fmt"
	"sync"

	dyvmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/go-kit/log"

	"github.com/go-kit/log/level"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)

// Notifier implements a Notifier for vms notifications.
type Notifier struct {
	conf    *config.VMSConfig
	tmpl    *template.Template
	logger  log.Logger
	m       sync.Mutex
	client  *dyvmsapi.Client
	request *dyvmsapi.SingleCallByTtsRequest
}

// New returns a new vms notifier.
func New(c *config.VMSConfig, t *template.Template, l log.Logger) (*Notifier, error) {
	var (
		client *dyvmsapi.Client
		err    error
	)
	if string(c.AccessKeyID) == "" {
		client, err = dyvmsapi.NewClientWithEcsRamRole(c.RegionID, string(c.RoleName))
		if err != nil {
			return nil, err
		}
	} else {
		client, err = dyvmsapi.NewClientWithAccessKey(c.RegionID, string(c.AccessKeyID), string(c.AccessKeySecret))
		if err != nil {
			return nil, err
		}
	}
	req := dyvmsapi.CreateSingleCallByTtsRequest()
	return &Notifier{conf: c, tmpl: t, logger: l, client: client, request: req}, nil
}

// Notify implements the Notifier interface.
func (n *Notifier) Notify(ctx context.Context, as ...*types.Alert) (bool, error) {
	n.m.Lock()
	defer n.m.Unlock()

	n.request.Scheme = "https"
	n.request.TtsCode = n.conf.TtsCode

	key, err := notify.ExtractGroupKey(ctx)
	if err != nil {
		return false, err
	}

	level.Debug(n.logger).Log("incident", key)

	var oas []*types.Alert
	first := true
	for k := range as {
		oas = append(oas, as[k])
		data := notify.GetTemplateData(ctx, n.tmpl, oas, n.logger)

		tmpl := notify.TmplText(n.tmpl, data, &err)
		if err != nil {
			return false, err
		}

		n.request.TtsParam = tmpl(n.conf.TtsParam)
		//whether all alarms
		if first || n.conf.IsBomb {
			first = false
			for _, v := range n.conf.PhoneNumber {
				n.request.CalledNumber = v
				resp, err := n.client.SingleCallByTts(n.request)
				if err != nil {
					return false, err
				}
				level.Debug(n.logger).Log("msg", "response info", "num", k, "code", resp.Code, "message", resp.Message, "callid", resp.CallId, "requestid", resp.RequestId, "incident", key)

				if resp.Code != "OK" {
					return false, fmt.Errorf("unexpected resp code %v, response message: %v, callId: %v, requestId: %v", resp.Code, resp.Message, resp.CallId, resp.RequestId)
				}
			}
		}
		//Only one vms alert at a time
		oas = nil
	}

	return false, nil
}
