#!/bin/bash

cargo build
echo $(cd client && go run .)
