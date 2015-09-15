/*
Package llog is a logger which controls the output by level.

It is essentially a frontend to the standard lib log.
Messages are written depending on the set level.
Several helpers for consistent formating are provided.

Only 1 main log at one time is possible.

A simple logrotate is also provided.
*/
package llog
