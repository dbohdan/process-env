#! /usr/bin/env python3
from __future__ import annotations

import os
import shlex

import psutil

SESSION = "mate-session"
VARS = ["DBUS_SESSION_BUS_ADDRESS", "DISPLAY"]


def pgrep(user: str, program: str) -> list[psutil.Process]:
    return [
        p for p in psutil.process_iter() if p.username() == user and program in p.name()
    ]


def main() -> None:
    username = os.environ["USER"]
    session_procs = pgrep(username, SESSION)

    if not session_procs:
        msg = "no session found"
        raise ValueError(msg)
    if len(session_procs) > 1:
        msg = "more than one session found"
        raise ValueError(msg)

    session = session_procs[0]

    env = session.environ()
    for var in VARS:
        print(f"{var}={shlex.quote(env[var])}")


if __name__ == "__main__":
    main()
