git_short_hash=`git rev-parse --short HEAD`
project_name='odysseia'
image := $(shell docker images -q ${project_name}:${git_short_hash})
ROOT_DIR=$(shell pwd)

create-image:
ifeq ("${image}","")
	echo "creating base image"
	echo "docker build ${project_name}:$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) . --no-cache
else
	echo "${image}"
	echo "base image already present on this machine"
endif

create-image-force:
	echo "creating base image (forced)"
	echo "docker build ${project_name}:$(git_short_hash)"
	docker build -t $(project_name):$(git_short_hash) . --no-cache

generate_swagger:
	echo ${ROOT_DIR}
	docker run -it -v ${ROOT_DIR}:${ROOT_DIR} -e SWAGGER_GENERATE_EXTENSION=true --workdir ${ROOT_DIR}/alexandros quay.io/goswagger/swagger generate spec -o ./docs/swagger.json -m;
	curl -X 'POST' \
      'https://converter.swagger.io/api/convert' \
	  -H 'accept: application/yaml' \
	  -H 'Content-Type: application/json' \
      -d '@./alexandros/docs/swagger.json' > ./alexandros/docs/openapi.yaml