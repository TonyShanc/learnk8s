#!/bin/bash

memtester 200M 1
# mkdir /tmp/memory
# mount -t tmpfs -o size=1024M tmpfs /tmp/memory
# dd if=/dev/zero of=/tmp/memory/block
# sleep 30
# rm /tmp/memory/block
# umount /tmp/memory
# rmdir /tmp/memory