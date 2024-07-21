genbuf:
	buf lint && buf generate

server:
	@go run main.go
	echo "🪖 Running Server"

client:
	@cd web && npm run dev
	echo "🪖 Running Client"