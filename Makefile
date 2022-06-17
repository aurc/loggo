debug:
	go build -gcflags='-N -l' cmd/main.go \
        && dlv --listen=:2345 --headless=true --api-version=2 exec ./main