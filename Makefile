test:
	go test ./...

cover:
	go test -coverprofile=c.out -covermode=atomic
	go test -coverprofile=c.out -covermode=atomic ./csvutil/...
	go test -coverprofile=c.out -covermode=atomic ./regression/...

view:
	go tool cover -html=c.out

doc:
	asciidoctor README.adoc

open:
	open README.html
