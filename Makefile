build: clean
	go build

clean:
	rm -f ./ws

test:
	printf '{"hello":"world"}' | ./ws ws://echo.websocket.org