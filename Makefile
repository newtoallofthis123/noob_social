BINARY_NAME=noob_social

build: 
	@templ generate && cd cmd && go build -o ../bin/$(BINARY_NAME)

run: build
	@./bin/$(BINARY_NAME)

clean:
	@rm -f bin/$(BINARY_NAME)

tailwind:
	npx tailwindcss -i static/input.css -o static/output_prod_styles.css --watch

css:
	bunx tailwindcss -i static/input.css -o static/output_prod_styles.css --watch
