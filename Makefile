build:
	go build

deps:
	go get -u "github.com/gonum/matrix/mat64"
	go get -u "github.com/DavidGamba/go-getoptions"
	go get -u "github.com/montanaflynn/stats"
	go get -u "github.com/gonum/plot"
	go get -u "github.com/gonum/plot/plotter"
	go get -u "github.com/gonum/plot/vg"

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
	asciidoctor -b manpage csv-analysis.adoc

open:
	open README.html
