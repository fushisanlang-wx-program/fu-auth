package model

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
)

type BaseStruct struct {
	FuLogger *glog.Logger
	FuDb     gdb.DB
	FuConf   *gcfg.Config
}
