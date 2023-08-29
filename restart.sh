#!/bin/bash

. ~/.bash_profile
go build
mv gpt-wework gpt-love 
pkill -9 gpt-love 
./gpt-love &
