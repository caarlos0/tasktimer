#!/bin/sh
set -e
rm -rf manpages
mkdir manpages
go run . man | gzip -c >manpages/tt.1.gz
