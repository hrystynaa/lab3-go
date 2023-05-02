#!/bin/bash

index=0
defX=600
defY=150
steps=20
interval=0.1

dx=$((defX/steps))
dy=$((defY/steps))

sleep $interval
curl -X POST http://localhost:17000 -d "reset"
curl -X POST http://localhost:17000 -d "white"
curl -X POST http://localhost:17000 -d "figure 100 100"

while true; do
  for i in {1..20}
  do
      curl -X POST http://localhost:17000 -d "move $dx $dy"
      curl -X POST http://localhost:17000 -d "update"
      sleep $interval
  done

  index=$((index+1))
  dx=$((dx*-1))

  if ((index > 3)); then
    dy=$((dy*-1))
    index=0
  fi

done 
