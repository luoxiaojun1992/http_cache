DARWIN_DIR = ../bin/darwin_amd64
LINUX_DIR = ../bin/linux_amd64
FREEBSD_DIR = ../bin/freebsd_amd64
WINDOWS_DIR = ../bin/windows_amd64

- install:
	cd src && GOOS=darwin GOARCH=amd64 go build -i -v -o http_cache_darwin && mkdir -p $(DARWIN_DIR) && \
	mv ./http_cache_darwin $(DARWIN_DIR)/http_cache && \
	tar -zcvf $(DARWIN_DIR)/http-cache.darwin-amd64.tar.gz $(DARWIN_DIR)/http_cache && \
	rm $(DARWIN_DIR)/http_cache && cd .. && \
	cd src && GOOS=linux GOARCH=amd64 go build -i -v -o http_cache_linux && mkdir -p $(LINUX_DIR) && \
	mv ./http_cache_linux $(LINUX_DIR)/http_cache && \
	tar -zcvf $(LINUX_DIR)/http-cache.linux-amd64.tar.gz $(LINUX_DIR)/http_cache && \
	rm $(LINUX_DIR)/http_cache && cd .. && \
	cd src && GOOS=freebsd GOARCH=amd64 go build -i -v -o http_cache_freebsd && mkdir -p $(FREEBSD_DIR) && \
	mv ./http_cache_freebsd $(FREEBSD_DIR)/http_cache && \
	tar -zcvf $(FREEBSD_DIR)/http-cache.freebsd-amd64.tar.gz $(FREEBSD_DIR)/http_cache && \
	rm $(FREEBSD_DIR)/http_cache && cd .. && \
	cd src && GOOS=windows GOARCH=amd64 go build -i -v -o http_cache_windows && mkdir -p $(WINDOWS_DIR) && \
	mv ./http_cache_windows $(WINDOWS_DIR)/http_cache && \
	tar -zcvf $(WINDOWS_DIR)/http-cache.windows-amd64.tar.gz $(WINDOWS_DIR)/http_cache && \
	rm $(WINDOWS_DIR)/http_cache && cd ..
