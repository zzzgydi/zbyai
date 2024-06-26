HOMEDIR := $(shell pwd)
PROJECTNAME := $(shell basename $(HOMEDIR))
OUTDIR  := $(HOMEDIR)/output

GOBUILD := go build -ldflags "-X google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn"

build:
	mkdir -p $(OUTDIR)
	$(GOBUILD) -o $(OUTDIR)/$(PROJECTNAME) ./cmd/main.go


dev: build
	mkdir -p $(OUTDIR)/log
	cp -r $(HOMEDIR)/config  $(OUTDIR)/
	cp -r $(HOMEDIR)/prompt  $(OUTDIR)/

	CONFIG_ENV=dev $(OUTDIR)/$(PROJECTNAME)
