#! /bin/sh
#
# restart.sh
# Copyright (C) 2018 jack <jack@fabric>
#
# Distributed under terms of the MIT license.
#


./stop.sh
docker rm `docker ps -qa`
./teardown.sh
./installcc.sh
