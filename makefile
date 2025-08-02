include .env

genbuf:
	buf lint && buf generate

server:
	@go run main.go
	echo "🪖 Running Server"

d-server:
	@echo "🔨 Running Server with no game timeout"
	@go run main.go -timeout=86400

client:
	@cd web && npm run dev
	echo "🪖 Running Client"

deploy:
	gcloud builds submit --region=${LOCATION} --tag ${LOCATION}-docker.pkg.dev/${PROJECT_ID}/submarine-game/${IMAGE_NAME}:${TAG} \
	--project ${PROJECT_ID}

login:
	gcloud auth login

