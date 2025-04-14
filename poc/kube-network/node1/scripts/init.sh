#!/bin/sh

echo "Initializing Node1..."

# Start Docker daemon
dockerd-entrypoint.sh &
sleep 5

# Keep container running
tail -f /dev/null