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

package sms

import (
	"context"
	"fmt"

	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-kit/log"

	"github.com/go-kit/log/level"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)

// Notifier implements a Notifier for SMS notifications.
type Notifier struct {
	conf    *config.SMSConfig
	tmpl    *template.Template
	logger  log.Logger
	client  *dysmsapi.Client
	request *dysmsapi.SendSmsRequest
}

// New returns a new SMS notifier.
func New(c *config.SMSConfig, t *template.Template, l log.Logger) (*Notifier, error) {
	var (
		client *dysmsapi.Client
		err    error
	)
	if string(c.AccessKeyID) == "" {
		client, err = dysmsapi.NewClientWithEcsRamRole(c.RegionID, string(c.RoleName))
		if err != nil {
			return nil, err
		}
	} else {
		client, err = dysmsapi.NewClientWithAccessKey(c.RegionID, string(c.AccessKeyID), string(c.AccessKeySecret))
		if err != nil {
			return nil, err
		}
	}
	req := dysmsapi.CreateSendSmsRequest()
	return &Notifier{conf: c, tmpl: t, logger: l, client: client, request: req}, nil
}

// Notify implements the Notifier interface.
func (n *Notifier) Notify(ctx context.Context, as ...*types.Alert) (bool, error) {
	key, err := notify.ExtractGroupKey(ctx)
	if err != nil {
		return false, err
	}

	level.Debug(n.logger).Log("incident", key)
	data := notify.GetTemplateData(ctx, n.tmpl, as, n.logger)

	tmpl := notify.TmplText(n.tmpl, data, &err)
	if err != nil {
		return false, err
	}

	n.request.Scheme = "https"
	n.request.PhoneNumbers = n.conf.PhoneNumber
	n.request.SignName = n.conf.SignName
	n.request.TemplateCode = n.conf.TemplateCode
	n.request.TemplateParam = tmpl(n.conf.TemplateParam)

	resp, err := n.client.SendSms(n.request)
	if err != nil {
		return false, err
	}

	level.Debug(n.logger).Log("msg", "response info", "code", resp.Code, "message", resp.Message, "bizid", resp.BizId, "requestid", resp.RequestId, "incident", key)

	if resp.Code != "OK" {
		return false, fmt.Errorf("unexpected resp code %v, response message: %v, bizId: %v, requestId: %v", resp.Code, resp.Message, resp.BizId, resp.RequestId)
	}

	return false, nil
}
