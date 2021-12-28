#!/bin/bash

for i in $(seq 0 1 12)
do
  gnome-terminal -e go run . i
done