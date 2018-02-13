#! /bin/sh
environ_grep() {
    awk -v "re=$1" 'BEGIN { RS="\0"; FS="=" } $0 ~ re { print $2 }' "$2"
}

env_pid="$(pgrep -u "$USER" mate-session)"
DBUS_SESSION_BUS_ADDRESS="$(environ_grep DBUS_SESSION_BUS_ADDRESS "/proc/$env_pid/environ")"
DISPLAY="$(environ_grep DISPLAY "/proc/$env_pid/environ")"
echo "DBUS_SESSION_BUS_ADDRESS=$DBUS_SESSION_BUS_ADDRESS"
echo "DISPLAY=$DISPLAY"
