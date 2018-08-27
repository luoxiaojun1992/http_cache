DARWIN_DIR = ../bin/darwin_amd64
LINUX_DIR = ../bin/linux_amd64
FREEBSD_DIR = ../bin/freebsd_amd64
WINDOWS_DIR = ../bin/windows_amd64

- install:
	cd src && GOOS=darwin GOARCH=amd64 go build -i -v -o http_cache && \
	tar -zcvf http-cache.darwin-amd64.tar.gz ./http_cache && rm ./http_cache && \
	mkdir -p $(DARWIN_DIR) && \
	mv ./http-cache.darwin-amd64.tar.gz $(DARWIN_DIR)/ && cd .. && \
	cd src && GOOS=linux GOARCH=amd64 go build -i -v -o http_cache && \
    tar -zcvf http-cache.linux-amd64.tar.gz ./http_cache && rm ./http_cache && \
    mkdir -p $(LINUX_DIR) && \
    mv ./http-cache.linux-amd64.tar.gz $(LINUX_DIR)/ && cd .. && \
    cd src && GOOS=freebsd GOARCH=amd64 go build -i -v -o http_cache && \
    tar -zcvf http-cache.freebsd-amd64.tar.gz ./http_cache && rm ./http_cache && \
    mkdir -p $(FREEBSD_DIR) && \
    mv ./http-cache.freebsd-amd64.tar.gz $(FREEBSD_DIR)/ && cd .. && \
    cd src && GOOS=windows GOARCH=amd64 go build -i -v -o http_cache && \
    tar -zcvf http-cache.windows-amd64.tar.gz ./http_cache && rm ./http_cache && \
    mkdir -p $(WINDOWS_DIR) && \
    mv ./http-cache.windows-amd64.tar.gz $(WINDOWS_DIR)/ && cd ..
