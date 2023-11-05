test:
	go test -v -race ./

coverage:
	go test ./ -race -shuffle=on -coverprofile=coverage.out -covermode=atomic
