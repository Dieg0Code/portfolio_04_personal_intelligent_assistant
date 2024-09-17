.PHONY: compile_lambda zip clean build-ServerlessApiScraper

compile_lambda:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap main.go

zip: compile_lambda
	zip rag_diary.zip bootstrap

clean:
	rm -f bootstrap rag_diary.zip

build-TestFunction: compile_lambda
	mkdir -p $(ARTIFACTS_DIR)
	cp bootstrap $(ARTIFACTS_DIR)