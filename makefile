
all: website backend

website:
	hugo --minify

backend:
	go mod tidy
	go build .

.PHONY all
