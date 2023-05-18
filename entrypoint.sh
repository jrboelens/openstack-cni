#!/bin/sh
set -x
set -e

if [[ "$SKIP_ENTRY" = "1" ]]; then
  /usr/bin/openstack-cni-daemon
  exit 0
fi

# Set known directories.
HOST_CNI_BIN_DIR="/host/opt/cni/bin"
HOST_CNI_ETC_DIR="/host/etc/cni/net.d"
CNI_BIN_FILE="/usr/bin/openstack-cni"

# Loop through and verify each location each.
for i in $HOST_CNI_BIN_DIR $HOST_CNI_ETC_DIR $CNI_BIN_FILE
do
  if [ ! -e "$i" ]; then
    echo "Location $i does not exist"
    exit 1
  fi
done

# Copy the CNI binary into place
if cp -f "$CNI_BIN_FILE" "$HOST_CNI_BIN_DIR"; then
    echo "Openstack CNI installed Success!"
else
    echo "Could not copy file $CNI_BIN_FILE"
    exit 1
fi

# Write out config that the CNI needs to run in kubelet's context
if [ "$CNI_API_URL" = "" ]; then
  CNI_API_URL="http://127.0.0.1:4242"
fi
echo "Using Api URL $CNI_API_URL"
echo "CNI_API_URL=$CNI_API_URL" > "$HOST_CNI_ETC_DIR/openstack-cni.conf"

# allow the binary to be injected from the host's filesystem
# this allows for testing without shipping new images
OVERRIDE_BINARY=$HOST_CNI_BIN_DIR/openstack-cni-daemon
if [ -f "$OVERRIDE_BINARY" ]; then
  echo "Found override.  Using $OVERRIDE_BINARY"
  cp $OVERRIDE_BINARY /usr/bin/openstack-cni-daemon
fi

/usr/bin/openstack-cni-daemon