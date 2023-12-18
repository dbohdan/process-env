#! /usr/bin/env python3

from __future__ import annotations

import argparse
import os
import shlex
import subprocess as sp
from typing import Literal, assert_never

import psutil

SESSION = "mate-session"
VARS = ["DBUS_SESSION_BUS_ADDRESS", "DISPLAY", "SSH_AUTH_SOCK"]

Shell = Literal["fish", "posix"]


def pgrep(user: str, proc_name: str) -> list[psutil.Process]:
    return [
        p for p in psutil.process_iter() if p.username() == user and proc_name == p.name()
    ]


def fish_quote(s: str) -> str:
    return sp.run(
        ["fish", "-c", 'string escape "$argv[1]"', s],
        capture_output=True,
        check=True,
        text=True,
    ).stdout.rstrip()


def quote_var(name: str, value: str, *, shell: Shell) -> str:
    match shell:
        case "fish":
            return f"set -x {fish_quote(name)} {fish_quote(value)}"
        case "posix":
            return f"export {shlex.quote(name)}={shlex.quote(value)}"
        case _ as unreachable:
            assert_never(unreachable)


def cli() -> Shell:
    parser = argparse.ArgumentParser(
        description="Print shell commands to set environment variables like "
        "`DISPLAY` to those of the current user's desktop session.",
    )

    parser.add_argument(
        "shell",
        choices=["fish", "posix"],
        help="what shell to print commands for",
    )

    args = parser.parse_args()
    return args.shell


def main() -> None:
    shell = cli()

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
        print(quote_var(var, env[var], shell=shell))


if __name__ == "__main__":
    main()
