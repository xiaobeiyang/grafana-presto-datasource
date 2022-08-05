SHELL := /bin/bash
REV=$(shell git log --max-count=1 --pretty="format:%h")
GO_VER=$(shell go version|grep "go version"|cut -d' ' -f3|sed "s/[\s\t]*//"|sed "s/^go//")
VERSION = "undefined"
VER="-X main.Version=$(VERSION) -X main.Revision=$(REV) -X main.GoVersion=$(GO_VER) -X main.BuiltAt=`date -u '+%Y-%m-%d_%I:%M:%S%p'`"
SRC=$(shell find . -path ./ui -prune -or -type d | grep '^./')

dep:
	go mod vendor
	yarn install

build-go:
	cd pkg && GOOS=linux && GOARCH=amd64 && go build -ldflags $(VER) -v -o grafana-presto-datasource_linux_amd64
	cd pkg && GOOS=linux && GOARCH=arm64 && go build -ldflags $(VER) -v -o grafana-presto-datasource_linux_arm64
	cd pkg && GOOS=darwin && GOARCH=amd64 && go build -ldflags $(VER) -v -o grafana-presto-datasource_darwin_amd64
	cd pkg && GOOS=darwin && GOARCH=arm64 && go build -ldflags $(VER) -v -o grafana-presto-datasource_darwin_arm64

build-js:
	yarn build

package:
	rm -rf grafana-presto-datasource
	rm grafana-presto-datasource.tar.gz
	mkdir grafana-presto-datasource
	cp -r dist grafana-presto-datasource/ && cp pkg/grafana-presto-datasource_* grafana-presto-datasource/dist/
	tar vfcz grafana-presto-datasource.tar.gz grafana-presto-datasource
