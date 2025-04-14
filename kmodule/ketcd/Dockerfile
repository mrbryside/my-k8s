FROM quay.io/coreos/etcd:v3.5.0

# Set up environment variables
ENV ETCD_DATA_DIR=/etcd-data \
    ETCD_NAME=default \
    ETCD_INITIAL_CLUSTER_STATE=new

# Expose ports
EXPOSE 2379 2380

# Command to run etcd
CMD ["/usr/local/bin/etcd"]