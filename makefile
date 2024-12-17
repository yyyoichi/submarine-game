include .env

genbuf:
	buf lint && buf generate

server:
	@go run main.go
	echo "🪖 Running Server"

client:
	@cd web && npm run dev
	echo "🪖 Running Client"

deploy:
	gcloud builds submit --region=${LOCATION} --tag ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/submarine-game/${IMAGE_NAME}:${TAG} \
	--project ${PROJECT_ID}

login:
	gcloud auth login

dev: build
	@echo "🏗️  Building WebAssembly"
	@echo "🚀 Running Server"
	@go run cmd/app/main.go

build:
	env GOOS=js GOARCH=wasm go build -o internal/app/public/dist/game.wasm cmd/screen/main.go

initwasm:
	cp $(shell go env GOROOT)/misc/wasm/wasm_exec.js internal/app/public/dist/


