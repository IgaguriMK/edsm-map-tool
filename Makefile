ifeq ($(OS), Windows_NT)
	EXE:=.exe
else
	EXE:=""
endif

.PHONY: all
all: systemCoord$(EXE) transform$(EXE) imaging$(EXE)

systemCoord$(EXE):
	go build systemCoord.go

transform$(EXE):
	go build transform.go

imaging$(EXE):
	go build imaging.go

.PHONY: clean
clean:
	- rm systemCoord$(EXE)
	- rm transform$(EXE)
	- rm imaging$(EXE)
