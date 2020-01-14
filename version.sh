#!/bin/bash
sed 's/.*"\(.*\)".*/\1/' <<< "`grep "	Version" util/declarations.go`"