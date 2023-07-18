package logger

import (
	"FuAuth/app/model"

	"github.com/gogf/gf/v2/frame/g"
)

func LogInfo(baseStruct model.BaseStruct, logString string, ctx g.Ctx) {
	baseStruct.FuLogger.Info(ctx, logString)
}
func LogWarn(baseStruct model.BaseStruct, logString string, ctx g.Ctx) {
	baseStruct.FuLogger.Warning(ctx, logString)
}
func LogError(baseStruct model.BaseStruct, logString string, ctx g.Ctx) {
	baseStruct.FuLogger.Error(ctx, logString)
}
