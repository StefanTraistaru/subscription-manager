#!/bin/bash

# Remove all started containers
docker ps -a -q | xargs docker rm 2> /dev/null || true

# Remove all dangling docker images
docker images -qa -f 'dangling=true' | xargs docker rmi 2> /dev/null || true