#! /bin/bash
#
# upgrade.sh
# Copyright (C) 2018 jack <jack@HP-WorkStation>
#
# Distributed under terms of the MIT license.
#

[[ -z $1 ]] && echo "No chaincode specified." && exit
[[ -z $2 ]] && echo "No version specified." && exit

. ./fabric.conf
for i in $1
do
    upgrade $i $2
done
