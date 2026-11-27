package error

import "errors"

var (
	ErrRepairJobNotFound = errors.New("Repair job not found")
)

var RepairJobErrors = []error{
	ErrRepairJobNotFound,
}
