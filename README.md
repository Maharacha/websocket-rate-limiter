Start a rpc node that will not sync but be able to respond on requests.

    ./polkadot --chain westend --rpc-external --rpc-port=9933 --rpc-cors all --rpc-methods unsafe --reserved-only

Start the websocket application

    go run main.go

Use websocat to try it

    websocat --jsonrpc ws://localhost:8080
