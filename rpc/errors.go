package rpc

import "errors"

var E_NO_PROBE = errors.New("No such probe exists")
var E_NO_SCHEDULE = errors.New("No such schedule exists")
var E_PROBE_ID = errors.New("Probe id must be specified")
var E_SCHEDULE_ID = errors.New("Schedule id must be specified")
var E_EMPTY_SCHEDULE = errors.New("Empty schedule specified")
