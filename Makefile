IMAGE_NAME = "buzeus/sample-app"
IMAGE_TAG = "latest"

check_install:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	swagger generate spec -o ./swagger/swagger.yml --scan-models

lint:
	go fmt
	go vet

vendor:
	go mod vendor

build_image: lint vendor swagger
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

run_local: build_image
	docker run -it ${IMAGE_NAME}:${IMAGE_TAG}

push_image: build_image
	docker push ${IMAGE_NAME}:${IMAGE_TAG}