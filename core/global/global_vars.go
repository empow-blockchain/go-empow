package global

import (
	"github.com/empow-blockchain/go-empow/common"
)

// BuildTime build time
var BuildTime string

// GitHash git hash
var GitHash string

// CodeVersion is the version string of code
var CodeVersion string

var globalConf *common.Config

// SetGlobalConf ...
func SetGlobalConf(conf *common.Config) {
	globalConf = conf
}

// GetGlobalConf ...
func GetGlobalConf() *common.Config {
	return globalConf
}
