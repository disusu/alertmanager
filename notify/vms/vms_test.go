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
	"reflect"
	"sync"
	"testing"

	dyvmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dyvmsapi"
	"github.com/go-kit/log"
	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/alertmanager/template"
	"github.com/prometheus/alertmanager/types"
)

func TestNew(t *testing.T) {
	type args struct {
		c *config.VMSConfig
		t *template.Template
		l log.Logger
	}
	tests := []struct {
		name    string
		args    args
		want    *Notifier
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.c, tt.args.t, tt.args.l)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotifier_Notify(t *testing.T) {
	type fields struct {
		conf    *config.VMSConfig
		tmpl    *template.Template
		logger  log.Logger
		m       sync.Mutex
		client  *dyvmsapi.Client
		request *dyvmsapi.SingleCallByTtsRequest
	}
	type args struct {
		ctx context.Context
		as  []*types.Alert
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Notifier{
				conf:    tt.fields.conf,
				tmpl:    tt.fields.tmpl,
				logger:  tt.fields.logger,
				m:       tt.fields.m,
				client:  tt.fields.client,
				request: tt.fields.request,
			}
			got, err := n.Notify(tt.args.ctx, tt.args.as...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Notifier.Notify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Notifier.Notify() = %v, want %v", got, tt.want)
			}
		})
	}
}
