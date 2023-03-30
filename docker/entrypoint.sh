#!/bin/bash

# start pigeon
/curve-manager/pigeon start

# start nginx
nginx -g "daemon off;"
