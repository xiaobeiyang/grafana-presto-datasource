package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Revision  = ""
	GoVersion = ""
	BuiltAt   = ""
	Version   = ""
)

func main() {
	v := flag.Bool("v", false, "show version")
	flag.Parse()
	if *v {
		fmt.Printf("Version: %v, Revision: %v, GoVersion: %v, BuiltAt: %v\n", Version, Revision, GoVersion, BuiltAt)
		return
	}

	backend.Logger.Info("Starting presto datasource backend...")
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":3100", nil)
	}()

	if err := datasource.Manage("grafana-presto-datasource", NewDatasourceInstance, datasource.ManageOpts{}); err != nil {
		log.DefaultLogger.Error(err.Error())
		os.Exit(1)
	}
}
