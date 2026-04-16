# Podman Auto-Update ZFS Snapshotter

This script calls `podman auto-update --dry-run` to check for pending container updates then makes a ZFS snapshot of the corresponding volume. The directory containing the volume must have the same name as the container.

This script is intended to be called before a proper auto-update. This can be done using systemd:

```
/etc/systemd/system/pmau-zfs-snap.service

[Unit]
Description=Podman auto-update ZFS snapshotter
Wants=network-online.target
After=network-online.target
Before=podman-auto-update.service

[Service]
Type=oneshot
ExecStart=/usr/local/bin/pmau-zfs-snap --root-dataset rpool/data/volumes

[Install]
WantedBy=default.target
```

Also copy `podman-auto-update.service` from /usr/lib/systemd/system to /etc/systemd/system/podman-auto-update.service with the following addition:

```
[Unit]
...
After=pmau-zfs-snap.service
```
