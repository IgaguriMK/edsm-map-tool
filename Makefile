GOLINT_OPTS=-min_confidence 0.8 -set_exit_status

.PHONY: all
all: build lint


.PHONY: build
build: systemCoord transform fromtext imaging

systemCoord:
	go build systemCoord.go
	- golint $(GOLINT_OPTS) systemCoord.go

transform:
	go build transform.go
	- golint $(GOLINT_OPTS) transform.go

fromtext:
	go build fromtext.go
	- golint $(GOLINT_OPTS) fromtext.go

imaging:
	go build imaging.go
	- golint $(GOLINT_OPTS) imaging.go


.PHONY: lint
lint:
	golint $(GOLINT_OPTS) sysCoord/*.go
	golint $(GOLINT_OPTS) *.go


.PHONY: clean
clean:
	- rm *.exe
	- rm systemCoord
	- rm transform
	- rm imaging
