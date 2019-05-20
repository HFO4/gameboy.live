#!/bin/sh

cd `dirname $0`
fyne bundle -package fyne frame.svg > bunded.go

