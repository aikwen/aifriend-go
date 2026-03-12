package tools


import (
	"context"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type getTimeParams struct{}

func NewGetTimeTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "get_time",
		Desc: "返回当前精确时间。用户询问现在几点、当前时间、今天日期、当前日期时间时，应调用此工具。返回中国时区时间，格式为 2006-01-02 15:04:05",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	}

	return utils.NewTool(info, func(ctx context.Context, input *getTimeParams) (string, error) {
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			return time.Now().Format("2006-01-02 15:04:05"), nil
		}
		return time.Now().In(loc).Format("2006-01-02 15:04:05"), nil
	})
}