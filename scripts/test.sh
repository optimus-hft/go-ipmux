set -e
go test $(go list ./... | grep -v /vendor/) -v -race -coverprofile coverage.txt
go tool cover -func=coverage.txt | grep total