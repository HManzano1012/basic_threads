serve:
	go run ./cmd/basicthreads/main.go

build:
	go build -o bin/main ./cmd/basicthreads/main.go

clean:
	rm -rf /bin/main
