test:
	go test -v -coverprofile coverage/cover.out
	go tool cover -html=coverage/cover.out