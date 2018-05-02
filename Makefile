default: gossh.go
	@glide install
	@go build -o gossh gossh.go

install:
	install gossh $(GOPATH)/bin/

clean:
	rm -f gossh

help:
	@echo "make             # compile"
	@echo "make install     # install gossh to $(GOPATH)/bin"
	@echo "make clean       # clean compile information"
