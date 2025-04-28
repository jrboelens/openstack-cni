# Changelog

## 0.0.28 (2025-03-21)
 - Added portbindings port options
 -  now retrying mac lookups
   -added a test for mac address lookup failures

## 0.0.27 (2025-02-21)
 - Added `enable_port_security` as an optional cni configuration

## 0.0.26 (2024-10-30)

- Interfaces are now referred to by index rather than name in order to avoid udev race conditions
- Added a check to ensure that `eth0` is never used as the destination interface name
- Added WaitForUdev feature
    - If enabled, before netlink configuration is applied, the interface name is compared to the `CNI_WAIT_FOR_UDEV_PREFIX`.  A matching prefixes causes a delay of `CNI_WAIT_FOR_UDEV_DELAY_MS` up to a total of `CNI_WAIT_FOR_UDEV_TIMEOUT_MS`.
    This logic is intended to avoid race condition that can be created when udev is manipulating interfaces.
- Added WaitForUdev related configuration:
    - `CNI_WAIT_FOR_UDEV (default: true)`
    - `CNI_WAIT_FOR_UDEV_PREFIX (default: 'eth')`
    - `CNI_WAIT_FOR_UDEV_DELAY_MS (default: '100')`
    - `CNI_WAIT_FOR_UDEV_TIMEOUT_MS (default: '5000')`


## 0.0.25 (2024-10-23)

- Added logfile support (`CNI_LOG_FILENAME`)
- Improved reaper behavior
- Added contextual log message to CNI plugin errors
- Added extra CNI plugin logging
- Updateid default CNI plugin configuration to include `CNI_LOG_FILENAME=/opt/cni/bin/openstack-cni.log`
- Added `host` tag to all ports
- Added `CNI_SKIP_REAPING` configuration in order to disable the reaping of ports
- Only `DOWN` ports will now be reaped
- Now replacing IP Addresses rather than adding them (`AddrReplace` vs `AddrAdd`)