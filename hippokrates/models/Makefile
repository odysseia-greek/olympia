PWD := $(CURDIR)
DOCKER_IMAGE := e2e-oapi-codegen:latest
NAME := alexandros

generate_from_api:
	#curl http://k3d-odysseia.greek:8080/ploutarchos/redoc/$(NAME)/yaml > $(NAME).yaml
	docker build -t $(DOCKER_IMAGE) .
	docker run --rm -v $(PWD):/models -w /models $(DOCKER_IMAGE) \
		oapi-codegen -package models -generate types $(NAME).yaml > $(NAME).gen.go