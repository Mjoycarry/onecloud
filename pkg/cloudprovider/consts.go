package cloudprovider

import (
	"errors"
)

const (
	CloudVMStatusRunning      = "running"
	CloudVMStatusSuspend      = "suspend"
	CloudVMStatusStopped      = "stopped"
	CloudVMStatusChangeFlavor = "change_flavor"
	CloudVMStatusDeploying    = "deploying"
	CloudVMStatusOther        = "other"
)

var ErrNotFound = errors.New("id not found")
var ErrDuplicateId = errors.New("duplicate id")
var ErrInvalidStatus = errors.New("invalid status")
var ErrTimeout = errors.New("timeout")
var ErrNotImplemented = errors.New("Not implemented")
var ErrNotSupported = errors.New("Not supported")
var ErrInvalidProvider = errors.New("Invalid provider")

const (
	VM_AZURE_DEFAULT_LOGIN_USER = "toor"
)
