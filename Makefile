# init project path
HOMEDIR := $(shell pwd)
OUTDIR  := $(HOMEDIR)/output

# 应用名称/二进制文件名称
APPNAME = pod-service-relations

# 镜像仓库 & kubernetes命名空间
IMG_HUB = registry.xxx.com

ifdef GOROOT
	GOROOT := $(GOROOT)
	GO     := $(GOROOT)/bin/go
endif

TAG=${shell git rev-parse --short HEAD}

# init command params
GO      := go
GOROOT  := $(shell $(GO) env GOROOT)
GOPATH  := $(shell $(GO) env GOPATH)
GOMOD   := $(GO) mod
GOBUILD := $(GO) build
GOTEST  := $(GO) test -race -timeout 30s -gcflags="-N -l"
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")


# test cover files
COVPROF := $(HOMEDIR)/covprof.out  # coverage profile
COVFUNC := $(HOMEDIR)/covfunc.txt  # coverage profile information for each function
COVHTML := $(HOMEDIR)/covhtml.html # HTML representation of coverage profile

# make, make all
all: prepare compile

# set proxy env
set-env:
	$(GO) env -w GO111MODULE=on
	$(GO) env -w GONOPROXY=\*.xunlei.com\*
	$(GO) env -w GOPROXY=https://goproxy.cn,direct
	$(GO) env -w GONOSUMDB=\*
	$(GO) env -w CC=/usr/bin/gcc
	$(GO) env -w CXX=/usr/bin/g++

#make prepare, download dependencies
prepare: gomod

gomod: set-env
	$(GOMOD) download -x || $(GOMOD) download -x

#make compile
compile: build

build:
	$(GOBUILD) -o $(HOMEDIR)/$(APPNAME)

# make test, test your code
test: prepare test-case
test-case:
	$(GOTEST) -v -cover $(GOPKGS)

# make package
package: package-bin
package-bin:
	mkdir -p $(OUTDIR)/bin
	cp $(APPNAME) $(OUTDIR)/bin

image:package-bin
	docker build -t $(IMG_HUB)/$(APPNAME)-prod:$(TAG) .

push:image
	docker push $(IMG_HUB)/$(APPNAME)-prod:$(TAG)

run:image
	@echo "Locally run service in the form of swarm..."
	@-docker service rm $(APPNAME) > /dev/null 2>&1  || true
	@-docker service create --name $(APPNAME) \
	--mount type=bind,source=/home/work/$(APPNAME)/log,destination=/usr/src/app/log \
	--no-resolve-image \
	--network devel -p 8080:8080 \
	$(IMG_HUB)/$(APPNAME):$(TAG)

push-test:package-test
	docker build -t $(IMG_HUB)/$(APPNAME)-test:$(TAG) .
	docker push $(IMG_HUB)/$(APPNAME)-test:$(TAG)

# make clean
clean:
	$(GO) clean
	rm -rf $(OUTDIR)
	rm -rf $(HOMEDIR)/pod-service-relations
	rm -rf $(GOPATH)/pkg/darwin_amd64

# avoid filename conflict and speed up build 
.PHONY: all prepare compile test package clean build
