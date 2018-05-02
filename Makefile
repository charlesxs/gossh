default: gossh.go
	glide install
	go build -o gossh gossh.go

install:
	install gossh $(GOPATH)/bin/

clean:
	rm -f gossh

