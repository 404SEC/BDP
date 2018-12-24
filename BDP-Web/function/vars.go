package function

import (
	"regexp"

	"github.com/micro/go-plugins/registry/zookeeper"
)

var (
	Cors      = map[string]bool{"*": true}
	VersionRe = regexp.MustCompilePOSIX("^v[0-9]+$")
	Namespace = "BDP-micro"
	Reg       = zookeeper.NewRegistry()
)
