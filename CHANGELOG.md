# Changelog

## 0.0.26 (2024-10-30)

- Interfaces are now referred to by index rather than name in order to avoid udev race conditions

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