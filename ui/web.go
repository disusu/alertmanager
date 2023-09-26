// Copyright 2015 Prometheus Team
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

package ui

import (
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof" // Comment this line to disable pprof endpoint.
	"path"

	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/route"
	"github.com/prometheus/common/server"

	"github.com/prometheus/alertmanager/asset"
	reactApp "github.com/prometheus/alertmanager/ui/react-app"
)

var reactRouterPaths = []string{
	"/",
	"/status",
}

// Register registers handlers to serve files for the web interface.
func Register(r *route.Router, reloadCh chan<- chan error, promlogConfig promlog.Config, logger log.Logger) {
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		disableCaching(w)

		req.URL.Path = "/static/"
		fs := http.FileServer(asset.Assets)
		fs.ServeHTTP(w, req)
	})

	r.Get("/script.js", func(w http.ResponseWriter, req *http.Request) {
		disableCaching(w)

		req.URL.Path = "/static/script.js"
		fs := http.FileServer(asset.Assets)
		fs.ServeHTTP(w, req)
	})

	r.Get("/favicon.ico", func(w http.ResponseWriter, req *http.Request) {
		disableCaching(w)

		req.URL.Path = "/static/favicon.ico"
		fs := http.FileServer(asset.Assets)
		fs.ServeHTTP(w, req)
	})

	r.Get("/lib/*path", func(w http.ResponseWriter, req *http.Request) {
		disableCaching(w)

		req.URL.Path = path.Join("/static/lib", route.Param(req.Context(), "path"))
		fs := http.FileServer(asset.Assets)
		fs.ServeHTTP(w, req)
	})

	serveReactApp := func(w http.ResponseWriter, r *http.Request) {
		f, err := reactApp.Assets.Open("/dist/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error opening React index.html: %v", err)
			return
		}
		defer func() { _ = f.Close() }()
		idx, err := io.ReadAll(f)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error reading React index.html: %v", err)
			return
		}
		w.Write(idx)
	}

	// Static files required by the React app.
	r.Get("/react-app/*filepath", func(w http.ResponseWriter, r *http.Request) {
		for _, rt := range reactRouterPaths {
			if r.URL.Path != "/react-app"+rt {
				continue
			}
			serveReactApp(w, r)
			return
		}
		r.URL.Path = path.Join("/dist", route.Param(r.Context(), "filepath"))
		fs := server.StaticFileServer(reactApp.Assets)
		fs.ServeHTTP(w, r)
	})

	r.Post("/-/reload", func(w http.ResponseWriter, req *http.Request) {
		errc := make(chan error)
		defer close(errc)

		reloadCh <- errc
		if err := <-errc; err != nil {
			http.Error(w, fmt.Sprintf("failed to reload config: %s", err), http.StatusInternalServerError)
		}
	})

	r.Get("/-/healthy", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	r.Head("/-/healthy", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r.Get("/-/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	r.Head("/-/ready", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if promlogConfig.Level.String() == "debug" {
		r.Post("/-/reload", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			errc := make(chan error)
			defer close(errc)

			reloadCh <- errc
			if err := <-errc; err != nil {
				http.Error(w, fmt.Sprintf("failed to reload config: %s", err), http.StatusInternalServerError)
			}
		}))

		r.Get("/debug/*subpath", http.DefaultServeMux.ServeHTTP)
		r.Post("/debug/*subpath", http.DefaultServeMux.ServeHTTP)
	}
}

func disableCaching(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0") // Prevent proxies from caching.
}
