// Code generated by go-swagger; DO NOT EDIT.

// Copyright Prometheus Team
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
//

package alertsls

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/prometheus/alertmanager/api/v2/models"
)

// NewPostslsAlertsParams creates a new PostslsAlertsParams object
// with the default values initialized.
func NewPostslsAlertsParams() *PostslsAlertsParams {
	var ()
	return &PostslsAlertsParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewPostslsAlertsParamsWithTimeout creates a new PostslsAlertsParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewPostslsAlertsParamsWithTimeout(timeout time.Duration) *PostslsAlertsParams {
	var ()
	return &PostslsAlertsParams{

		timeout: timeout,
	}
}

// NewPostslsAlertsParamsWithContext creates a new PostslsAlertsParams object
// with the default values initialized, and the ability to set a context for a request
func NewPostslsAlertsParamsWithContext(ctx context.Context) *PostslsAlertsParams {
	var ()
	return &PostslsAlertsParams{

		Context: ctx,
	}
}

// NewPostslsAlertsParamsWithHTTPClient creates a new PostslsAlertsParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewPostslsAlertsParamsWithHTTPClient(client *http.Client) *PostslsAlertsParams {
	var ()
	return &PostslsAlertsParams{
		HTTPClient: client,
	}
}

/*PostslsAlertsParams contains all the parameters to send to the API endpoint
for the postsls alerts operation typically these are written to a http.Request
*/
type PostslsAlertsParams struct {

	/*Alerts
	  The sls alerts to create

	*/
	Alerts *models.PostableSlsAlerts

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the postsls alerts params
func (o *PostslsAlertsParams) WithTimeout(timeout time.Duration) *PostslsAlertsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the postsls alerts params
func (o *PostslsAlertsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the postsls alerts params
func (o *PostslsAlertsParams) WithContext(ctx context.Context) *PostslsAlertsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the postsls alerts params
func (o *PostslsAlertsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the postsls alerts params
func (o *PostslsAlertsParams) WithHTTPClient(client *http.Client) *PostslsAlertsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the postsls alerts params
func (o *PostslsAlertsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithAlerts adds the alerts to the postsls alerts params
func (o *PostslsAlertsParams) WithAlerts(alerts *models.PostableSlsAlerts) *PostslsAlertsParams {
	o.SetAlerts(alerts)
	return o
}

// SetAlerts adds the alerts to the postsls alerts params
func (o *PostslsAlertsParams) SetAlerts(alerts *models.PostableSlsAlerts) {
	o.Alerts = alerts
}

// WriteToRequest writes these params to a swagger request
func (o *PostslsAlertsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Alerts != nil {
		if err := r.SetBodyParam(o.Alerts); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
