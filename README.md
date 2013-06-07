# Overload

## Overview

Overload is a very simple tool for load testing and benchmarking HTTP servers and applications.

## Usage

    overload [options] URL

Application Options:

    -r, --num-requests= Number of requests to make (1)
    -c, --concurrent=   Number of concurrent connections to make (1)
    -k, --keep-alive    Use keep alive connection

Sample:

    jack@hk-47~$ overload -r 500 -c 4 http://localhost:8080/
    # Requests: 500
    # Successes: 500
    # Failures: 0
    # Unavailable: 0
    Duration: 1.719238256s
    Average Request Duration: 13.575435ms

## Why another HTTP load tester / benchmark?

ab is broken on the Mac and siege does not support keep alive. Plus, it was fun.
