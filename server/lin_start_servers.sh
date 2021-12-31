#!/bin/bash

for i in $(seq 0 1 7)
do
  gnome-terminal -e go run . i
done