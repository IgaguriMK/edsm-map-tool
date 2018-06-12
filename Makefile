ifeq ($(OS), Windows_NT)
	EXE:=.exe
else
	EXE:=""
endif

.PHONY: all
all: systemCoord transform fromtext imaging

systemCoord:
	go build systemCoord.go

transform:
	go build transform.go

fromtext:
	go build fromtext.go

imaging:
	go build imaging.go

.PHONY: clean
clean:
	- rm systemCoord$(EXE)
	- rm transform$(EXE)
	- rm imaging$(EXE)
