#
# Makefile
# jack, 2018-04-17 19:30
#

all: cc

CCC := reg point map

cc:
	for cc in ${CCC} ; do \
		(cd $$cc; go build .); \
	done

.phony: cc

# vim:ft=make
#
