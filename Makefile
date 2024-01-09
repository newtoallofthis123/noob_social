BINARY_NAME=noob_social

build: 
	@templ generate && cd cmd && go build -o ../bin/$(BINARY_NAME)

run: build
	@./bin/$(BINARY_NAME)

clean:
	@rm -f bin/$(BINARY_NAME)
