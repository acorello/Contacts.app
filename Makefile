EXENAME = contacts

exe.linux:
	GOOS=linux GOARCH=amd64 go build -o $(EXENAME).linux

exe: 
	go build -o $(EXENAME)

clean:
	rm $(EXENAME) $(EXENAME).linux

.PHONY: exe exe.linux