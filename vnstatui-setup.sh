#!/bin/bash

vnstat_iface=$(ip route show to match 8.8.8.8 | grep default -m 1 | awk '{print $5}')

# uninstall eventually already installed vnstatui
if [ -f /usr/bin/vnstatui ]; then
  rm -rf /usr/bin/vnstatui
  systemctl disable vnstatui@${vnstat_iface[0]}.service 2>&1 >> /dev/null
  systemctl stop vnstatui@${vnstat_iface[0]}.service 2>&1 >> /dev/null
  systemctl daemon-reload
  rm -rf /usr/lib/systemd/system/vnstatui@.service
  rm -rf /tmp/vnstatui
fi
echo

# install program
cp -rf vnstatui /usr/bin/
chmod +x /usr/bin/vnstatui

# install systemd service
cp -rf vnstatui@.service /usr/lib/systemd/system/

# start and enable service
systemctl daemon-reload
systemctl enable vnstatui@${vnstat_iface[0]}.service
systemctl start vnstatui@${vnstat_iface[0]}.service
