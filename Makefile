GO    		:= GO15VENDOREXPERIMENT=1 go
glide    	:= glide
release    	:= ./release.sh

all: build

deps: 
	@echo ">> getting dependencies"
	@$(glide) install

build: 
	@echo ">> building binaries"
	@$(release)

docker: 
	@echo ">> building micromdm-dev docker container"
	docker build -f Dockerfile.dev -t micromdm-dev .

