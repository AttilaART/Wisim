#! /usr/bin/bash

trap 'echo "\nRecieved termination signal"; pkill "server"; pkill "npm run dev"; exit' SIGINT

echo "Starting both servers"
$(
	cd wisim-backend/server
	go run . 8000 1 1
) &

$(
	cd wisim-nouveau
	npm run dev &
) &

while true; do
	sleep 0.1s
done
