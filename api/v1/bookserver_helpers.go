package v1

import (
	"strings"
)

func (b *BookServer) DeploymentName() string {
	return strings.Join([]string{b.Name, "deployment"}, "-")
}

func (b *BookServer) ServiceName() string {
	return strings.Join([]string{b.Name, "deployment"}, "-")
}
