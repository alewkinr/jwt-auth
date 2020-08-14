dev:
	echo "Starting skaffold developers mode. For rebuild: https://skaffold.dev/docs/design/api/"
	skaffold dev -p dev --auto-build=false --port-forward

rebuild:
	echo "Starting to rebuild service"
	curl -POST \
		 -H 'Content-Type: application/x-www-form-urlencoded' \
		 -d '{"build": true}' \
		 'http://localhost:50052/v1/execute'

debug:
	echo "Starting skaffold debug mode. Add remote debug listener"
	skaffold debug -p debug --port-forward

lint:
	golangci-lint run --timeout 15m -v

release:lint
	go mod vendor
	go mod download
	go mod tidy
	git add -f ./vendor/*