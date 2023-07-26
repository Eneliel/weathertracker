APP_BIN = cmd/app/build/myapp

build-and-run: clean $(APP_BIN)
	$(APP_BIN)

$(APP_BIN):
	go build -o $(APP_BIN) ./cmd/app/main.go

clean: 
	rmdir /s /q cmd\app\build || exit 0