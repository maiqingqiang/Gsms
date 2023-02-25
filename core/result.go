package core

import "fmt"

const StatusSuccess = "success"

const StatusFailure = "failure"

type Result struct {
	Gateway  string
	Status   string
	Template string
	Result   string
	Error    error
}

func (r *Result) String() string {
	return fmt.Sprintf("gateway: %s, status: %s, template: %s, result: %s, error: %v", r.Gateway, r.Status, r.Template, r.Result, r.Error)
}
