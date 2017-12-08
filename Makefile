ifeq ($(OS), Windows_NT)
	EXE=.exe
else
	EXE=
endif

.PHONY: all
all: coordExtract$(EXE) transform$(EXE) imaging$(EXE)

.PHONY: coordExtract
coordExtract$(EXE):
	go build coordExtract.go

.PHONY: transform
transform$(EXE):
	go build transform.go

.PHONY: imaging
imaging$(EXE):
	go build imaging.go

.PHONY: clean
clean:
	- rm coordExtract$(EXE)
	- rm transform$(EXE)
	- rm imaging$(EXE)
