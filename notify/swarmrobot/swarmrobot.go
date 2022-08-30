// Copyright 2019 Prometheus Team
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

package swarmrobot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
	commoncfg "github.com/prometheus/common/config"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/notify"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)

// Notifier implements a Notifier for swarmrobot notifications.
type Notifier struct {
	conf   *config.SwarmRobotConfig
	tmpl   *template.Template
	logger log.Logger
	client *http.Client
}

type weChatSwarmRobotMessage struct {
	Type     string                                 `yaml:"msgtype,omitempty" json:"msgtype,omitempty"`
	Text     weChatSwarmRobotTextMessageContent     `yaml:"text,omitempty" json:"text,omitempty"`
	Markdown weChatSwarmRobotMarkdownMessageContent `yaml:"markdown,omitempty" json:"markdown,omitempty"`
}

type weChatSwarmRobotMarkdownMessageContent struct {
	Content string `json:"content"`
}
type weChatSwarmRobotTextMessageContent struct {
	Content             string `json:"content"`
	MentionedList       string `json:"mentioned_list"`
	MentionedMobileList string `json:"mentioned_mobile_list"`
}

type weChatSwarmRobotResponse struct {
	Code  int    `json:"errcode"`
	Error string `json:"errmsg"`
}

// New returns a new SwarmRobot notifier.
func New(c *config.SwarmRobotConfig, t *template.Template, l log.Logger, httpOpts ...commoncfg.HTTPClientOption) (*Notifier, error) {
	client, err := commoncfg.NewClientFromConfig(*c.HTTPConfig, "swarmrobot", httpOpts...)
	if err != nil {
		return nil, err
	}

	return &Notifier{conf: c, tmpl: t, logger: l, client: client}, nil
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
	//TODO:
	// limit: 20/min
	// text-length: 2048 bytes
	// markdown-length: 4096 bytes
	if len(n.conf.APIKey) == 0 {
		return false, fmt.Errorf("invalid APIKey")
	}

	msg := &weChatSwarmRobotMessage{
		Type: n.conf.MessageType,
	}

	if msg.Type == "markdown" {
		msg.Markdown = weChatSwarmRobotMarkdownMessageContent{
			Content: tmpl(n.conf.Message),
		}
	} else {
		msg.Text = weChatSwarmRobotTextMessageContent{
			Content:             tmpl(n.conf.Message),
			MentionedList:       n.conf.MentionedList,
			MentionedMobileList: n.conf.MentionedMobileList,
		}
	}
	if err != nil {
		return false, fmt.Errorf("templating error: %s", err)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(msg); err != nil {
		return false, err
	}

	postMessageURL := n.conf.APIURL.Copy()
	postMessageURL.Path += "webhook/send"
	q := postMessageURL.Query()
	q.Set("key", string(n.conf.APIKey))
	postMessageURL.RawQuery = q.Encode()

	resp, err := notify.PostJSON(ctx, n.client, postMessageURL.String(), &buf)
	if err != nil {
		return true, notify.RedactURL(err)
	}
	defer notify.Drain(resp)

	if resp.StatusCode != 200 {
		return true, fmt.Errorf("unexpected status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return true, err
	}
	level.Debug(n.logger).Log("response", string(body), "incident", key)

	var weResp weChatSwarmRobotResponse
	if err := json.Unmarshal(body, &weResp); err != nil {
		return true, err
	}

	// https://work.weixin.qq.com/api/doc#10649
	if weResp.Code == 0 {
		return false, nil
	}

	return false, errors.New(weResp.Error)
}
