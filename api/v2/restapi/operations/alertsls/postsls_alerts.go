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
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostslsAlertsHandlerFunc turns a function with the right signature into a postsls alerts handler
type PostslsAlertsHandlerFunc func(PostslsAlertsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostslsAlertsHandlerFunc) Handle(params PostslsAlertsParams) middleware.Responder {
	return fn(params)
}

// PostslsAlertsHandler interface for that can handle valid postsls alerts params
type PostslsAlertsHandler interface {
	Handle(PostslsAlertsParams) middleware.Responder
}

// NewPostslsAlerts creates a new http.Handler for the postsls alerts operation
func NewPostslsAlerts(ctx *middleware.Context, handler PostslsAlertsHandler) *PostslsAlerts {
	return &PostslsAlerts{Context: ctx, Handler: handler}
}

/*PostslsAlerts swagger:route POST /alerts/sls alertsls postslsAlerts

Create new sls Alerts

*/
type PostslsAlerts struct {
	Context *middleware.Context
	Handler PostslsAlertsHandler
}

func (o *PostslsAlerts) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostslsAlertsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
