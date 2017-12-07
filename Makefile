
.PHONY: all
all: coordExtract transform imaging

.PHONY: coordExtract
coordExtract:
	go build coordExtract.go

.PHONY: transform
transform:
	go build transform.go

.PHONY: imaging
imaging:
	go build imaging.go

.PHONY: clean
clean:
	- rm coordExtract
	- rm transform
	- rm imaging
	- rm *.exe
