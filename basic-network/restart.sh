#! /bin/sh
#
# restart.sh
# Copyright (C) 2018 jack <jack@fabric>
#
# Distributed under terms of the MIT license.
#

./stop.sh
./teardown.sh
./start.sh
./installcc.sh
