tidy:
	@go mod tidy
	
run:
	@go run .

build:
	@go build \
		-ldflags "-X main.buildName=document-service -X main.buildVersion=`git rev-parse --short HEAD`" \
		-o document-service.app main.go

docker-build:
	@docker build -f Dockerfile -t "document-service:1.0" --build-arg BUILD_DATE="$(date -u +"%Y-%m-%dT%H:%M:%SZ")" .

docker-run-container:
	@docker container run -ti --init --rm \
		--name document-service \
		--memory="1g" --cpus="1" \
		--net host \
		-v $$HOME/go/docker:/go \
		-v $(PWD):/data \
		-w /data \
		redhat/ubi8-micro:8.5-437 ./document-service.app