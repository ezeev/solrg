

coverage:
	go test -coverprofile=coverage.out -v
	go tool cover -html=coverage.out

