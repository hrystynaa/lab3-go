#!/bin/bash

curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "bgrect 200 200 600 600"
curl -X POST http://localhost:17000 -d "figure 400 400"
curl -X POST http://localhost:17000 -d "green"
curl -X POST http://localhost:17000 -d "move 100 100"
curl -X POST http://localhost:17000 -d "update"
curl -X POST http://localhost:17000 -d "reset"