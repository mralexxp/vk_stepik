package main

import (
	"encoding/json"
	"strings"
)

type ACL struct {
	Directory map[string][]string
}

func NewACL(ACLData string) (*ACL, error) {
	acl := ACL{make(map[string][]string)}

	err := json.Unmarshal([]byte(ACLData), &acl.Directory)
	if err != nil {
		return nil, err
	}

	return &acl, err
}

func (acl *ACL) CheckAccess(consumer, RequestedMethod string) bool {
	if methods, ok := acl.Directory[consumer]; ok {
		for _, AvailableMethod := range methods {
			if idx := strings.Index(AvailableMethod, "*"); idx != -1 {
				AvailableMethod, RequestedMethod = AvailableMethod[:idx], RequestedMethod[:idx]
			}

			if RequestedMethod == AvailableMethod {
				return true
			}
		}
	}

	return false
}
