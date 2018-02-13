#! /bin/sh
env_pid="$(pgrep -u "$USER" mate-session)"
DBUS_SESSION_BUS_ADDRESS="$(grep -z ^DBUS_SESSION_BUS_ADDRESS= "/proc/$env_pid/environ" | sed s/DBUS_SESSION_BUS_ADDRESS=//)"
DISPLAY="$(grep -z ^DISPLAY= "/proc/$env_pid/environ" | sed s/DISPLAY=//)"
echo "DBUS_SESSION_BUS_ADDRESS=$DBUS_SESSION_BUS_ADDRESS"
echo "DISPLAY=$DISPLAY"
