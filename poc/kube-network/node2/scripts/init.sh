#!/bin/sh

# Since we're using docker:dind image, Docker and iptables are already installed
# We just need Node2-specific configuration

echo "Initializing Node2..."

# Start Docker daemon (it might not be started in the dind image)
dockerd-entrypoint.sh &
sleep 5

# Create directory for web content
mkdir -p /tmp/webroot
mkdir -p /tmp/webroot2

# Create index.html with some text
cat > /tmp/webroot/index.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Node2 Web Server</title>
</head>
<body>
    <h1>Hello from Pod 1</h1>
    <h3>pod 1 is using vue3</h3>
</body>
</html>
EOF

cat > /tmp/webroot2/index.html << EOF
<!DOCTYPE html>
<html>
<head>
    <title>Node2 Web Server</title>
</head>
<body>
    <h1>Hello from Pod 2</h1>
    <h3>pod 2 is using react</h3>
</body>
</html>
EOF

# Create a macvlan network with the same subnet as cluster_network
echo "Creating macvlan Docker network..."
docker network create \
  -d macvlan \
  --subnet=172.20.0.0/16 \
  --ip-range=172.20.0.0/24 \
  -o parent=eth0 \
  node2-macvlan

# Start 2 pod with nginx
echo "Starting 2 pod.."
docker run -d \
  --name node2-pod-1 \
  --network node2-macvlan \
  --ip 172.20.1.3 \
  -v /tmp/webroot:/usr/share/nginx/html \
  nginx:alpine

docker run -d \
  --name node2-pod-2 \
  --network node2-macvlan \
  --ip 172.20.1.4 \
  -v /tmp/webroot2:/usr/share/nginx/html \
  nginx:alpine

echo "Pod 1 started with IP 172.20.1.3 and port 80"
echo "Pod 2 started with IP 172.20.1.4 and port 80"

# Create service with some ip 
ip link add dummy0 type dummy
ip addr add 172.20.1.5/32 dev dummy0
ip link set dummy0 up

# This for node access because when host call ip that is ip itself the OUTPUT mode will route
iptables -t nat -A OUTPUT -d 172.20.1.5 -p tcp --dport 80 -m statistic --mode random --probability 0.5 -j DNAT --to-destination 172.20.1.3:80
iptables -t nat -A OUTPUT -d 172.20.1.5 -p tcp --dport 80 -j DNAT --to-destination 172.20.1.4:80

# This for pod access because pod using macvlan network that's not a host network so the request from pod is an external request in this case prerouting will working on it
iptables -t nat -A PREROUTING -d 172.20.1.5 -p tcp --dport 80 -m statistic --mode random --probability 0.5 -j DNAT --to-destination 172.20.1.3:80
iptables -t nat -A PREROUTING -d 172.20.1.5 -p tcp --dport 80 -j DNAT --to-destination 172.20.1.4:80

# Create a nodeport
echo "Setting up NodePort on 30080 -> $SERVICEIP:$SERVICEPORT"
apk add --no-cache socat
socat TCP-LISTEN:30080,fork TCP:172.20.1.5:80 &

# Enable IP forwarding
sysctl -w net.ipv4.ip_forward=1

# Masquerade (if needed)
iptables  -t nat -A POSTROUTING -o eth0 -j MASQUERADE

# Keep container running
tail -f /dev/null