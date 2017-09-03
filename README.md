# goxel

-----------------

## Description

	goxel like axel tool

## Installation

	$ go get github.com/wayne666/goxel

## Usage

	goxel [options...] <url>
	
	Options:
		-n  Numbers of blocks to run(must).
		-H  Add header string.
		-v  More status information.
		-o  Specify local output file(must).
		-h  Help information.
		-V  Version information.
		-cpus Number of used cpu cores(Default is current machine cores).

#### Example:

	goxel -n 10 -o outfile http://xxx.com

## AUTHOR

	Written by WayneZhou, cumtxhzyy[at]gmail.com

## COPYRIGHT

	Copyright (c) 2015 WayneZhou. This library is free software; you can redistribute it and/or modify it.
