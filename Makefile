EXENAME = contacts

.PHONY: exe.linux
exe.linux:
	GOOS=linux GOARCH=amd64 go build -o $(EXENAME).linux

.PHONY: exe
exe: 
	go build -o $(EXENAME)

.PHONY: clean
clean:
	rm $(EXENAME) $(EXENAME).linux

.PHONY: ping
ping:
	@echo "PONG"
