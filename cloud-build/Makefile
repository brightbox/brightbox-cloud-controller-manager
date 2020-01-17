#!/bin/env make
first_release = 14
last_release = 17
environment = default
versions = environments/$(environment)/versions

$(versions): FORCE
	seq $(first_release) $(last_release) | xargs -I{} git describe --always release-1.{} | sed 's/^v\([0-9]*\.[0-9]*\.[0-9]*\).*$$/\1/' > $@

FORCE:

.PHONY: show
show: $(versions)
	tk show environments/$(environment)

apply: $(versions)
	tk apply environments/$(environment)
