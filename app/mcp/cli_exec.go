package mcp

import (
	"bytes"
	"context"
	"fmt"
	"gin/pkg/cli"
	"gin/pkg/serviceprovider/mcp"
	"io"
	"os"
	"strings"
)

func init() { mcp.Register(&CliExec{}) }

// CliExec CLI命令执行工具
type CliExec struct{}

func (t *CliExec) Name() string { return "execute_cli" }
func (t *CliExec) Description() string {
	return `执行CLI命令,可用命令: 
				make:controller(创建控制器/file=路径,desc=描述), 
				make:model(创建模型/name=名称,table=表名), 
				make:service(创建服务/name=名称), 
				make:request(创建验证器/name=名称), 
				make:middleware(创建中间件/name=名称), 
				make:router(创建路由/file=路径), 
				make:docs(生成Swagger文档), 
				permission:sync(同步用户权限到Redis),
				route:list(查看路由列表), 
				job:list(查看Job列表)`
}

func (t *CliExec) InputSchema() mcp.InputSchema {
	return mcp.InputSchema{
		Type: "object",
		Properties: map[string]mcp.Property{
			"command": {Type: "string", Description: "CLI命令名称,如make:controller"},
			"args":    {Type: "string", Description: "命令参数,格式: name=User&file=v1/user&desc=用户管理,多个参数用&分隔"},
		},
		Required: []string{"command"},
	}
}

func (t *CliExec) Call(ctx context.Context, args map[string]any) (any, error) {
	cmdName, _ := args["command"].(string)
	if cmdName == "" {
		return map[string]any{"error": "请指定命令名称,如make:controller"}, nil
	}

	// 查找命令
	cmd, exists := cli.Get(cmdName)
	if !exists {
		available := t.listAvailable()
		return map[string]any{"error": fmt.Sprintf("命令 %s 不存在。可用命令: %s", cmdName, available)}, nil
	}

	// 解析参数
	values := make(map[string]string)
	if argsStr, ok := args["args"].(string); ok && argsStr != "" {
		for _, pair := range strings.Split(argsStr, "&") {
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) == 2 {
				values[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			} else if len(kv) == 1 {
				values[strings.TrimSpace(kv[0])] = ""
			}
		}
	}

	// 如果提供了name但没提供file,自动填充
	if name, ok := values["name"]; ok && values["file"] == "" {
		values["file"] = name
	}

	// 捕获stdout执行命令
	output := t.captureOutput(func() {
		cmd.Execute(values)
	})

	// 构建结果文件路径
	resultFile := t.buildResultPath(cmdName, values)
	result := map[string]any{
		"command": cmdName,
		"args":    values,
		"output":  output,
	}
	if resultFile != "" {
		result["created"] = resultFile
	}
	result["message"] = fmt.Sprintf("命令 %s 执行完成", cmdName)

	return result, nil
}

// captureOutput 捕获stdout输出
func (t *CliExec) captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		done <- buf.String()
	}()

	fn()

	_ = w.Close()
	os.Stdout = old

	return <-done
}

func (t *CliExec) listAvailable() string {
	names := []string{
		"make:controller", "make:model", "make:service", "make:request",
		"make:middleware", "make:router", "make:docs",
		"permission:sync", "permission:grant-admin",
		"route:list", "job:list",
	}
	return strings.Join(names, ", ")
}

func (t *CliExec) buildResultPath(cmdName string, values map[string]string) string {
	file, _ := values["file"]
	name, _ := values["name"]
	if file == "" {
		file = name
	}
	if file == "" {
		return ""
	}

	switch cmdName {
	case "make:controller":
		return "app/controller/" + file + ".go"
	case "make:model":
		return "app/model/" + file + ".go"
	case "make:service":
		return "app/service/" + file + ".go"
	case "make:request":
		return "app/request/" + file + ".go"
	case "make:middleware":
		return "app/middleware/" + file + ".go"
	case "make:router":
		return "router/" + file + ".go"
	default:
		return ""
	}
}
