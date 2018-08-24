BINDIR = ../bin/`uname -s`_`uname -m`

- install:
	cd src && go build -i -v -o http_cache && mkdir -p $(BINDIR) && mv ./http_cache $(BINDIR)/http_cache & cd ..