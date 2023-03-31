package scenario

import (
	"encoding/json"
	"fmt"
	"github.com/finishy1995/effibot/http/internal/gpt"
	"github.com/finishy1995/effibot/http/internal/types"
	"io"
	"os"
	"strings"
)

const (
	jsonPath = "../data/rules/default.json"
)

var (
	rules      = map[string]*rule{}
	gptHandler *gpt.Handler
)

type rule struct {
	Id        string              `json:"id"`
	Guide     string              `json:"guide"`
	NextRule  []string            `json:"next_rule"`
	Params    []map[string]string `json:"params,optional"`
	Root      bool                `json:"root,default=false"`
	Clear     bool                `json:"clear,default=false"`
	Scripts   []string            `json:"scripts,optional"`
	paramsMap map[string]string
}

// Init rules 初始化规则
func Init() error {
	return loadJsonConfig(jsonPath)
}

// SetGPTHandler set gpt handler 设置 GPT 处理器
func SetGPTHandler(h *gpt.Handler) {
	gptHandler = h
}

// GetAllRules get all rules 获取所有规则
func GetAllRules() []types.Types {
	t := make([]types.Types, 0, len(rules))

	for id, item := range rules {
		t = append(t, types.Types{
			ID:     id,
			NextID: item.NextRule,
			IsRoot: item.Root,
			Params: item.GetParams(),
		})
	}
	return t
}

// GetParams get params 获取参数
func (r *rule) GetParams() []string {
	result := make([]string, 0, len(r.paramsMap))
	for key := range r.paramsMap {
		result = append(result, key)
	}
	return result
}

// Generate generate content 生成内容
func (r *rule) Generate(content string) string {
	return content // TODO: using params and scripts to generate content
}

func (r *rule) getParamsMap(params string) map[string]string {
	result := map[string]string{}
	for k, v := range r.paramsMap {
		result[k] = v
	}
	if params == "" {
		return result
	}

	arr := strings.Split(params, ",")
	for _, item := range arr {
		kvPair := strings.Split(item, ":")
		if len(kvPair) != 2 {
			continue
		}
		result[kvPair[0]] = kvPair[1]
	}
	return result
}

// GenerateGuide params: key1:value1,key2:value2 生成指南
func (r *rule) GenerateGuide(params string) string {
	m := r.getParamsMap(params)

	result := r.Guide
	for k, v := range m {
		result = strings.Replace(result, fmt.Sprintf("{%s}", k), v, -1)
	}

	return result
}

// Reset reset params 重置参数
func (r *rule) Reset() {
	r.paramsMap = map[string]string{}
	if r.Params == nil || len(r.Params) == 0 {
		return
	}
	for _, item := range r.Params {
		if name, ok := item["name"]; ok {
			value, ok := item["default"]
			if ok {
				r.paramsMap[name] = value
			} else {
				r.paramsMap[name] = ""
			}
		}
	}
}

// Allow new rule 是否允许新规则
func (r *rule) Allow(newRule string) bool {
	for _, item := range r.NextRule {
		if item == newRule {
			return true
		}
	}
	return false
}

func loadJsonConfig(filePath string) error {
	// open json file 打开本地 JSON 文件
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// read json data 读取 JSON 数据
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// parse json data 解析 JSON 数据
	var jsonData []*rule
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return err
	}

	for _, ruleItem := range jsonData {
		ruleItem.Reset()
		if ruleItem.Id != "" {
			rules[ruleItem.Id] = ruleItem
		}
	}

	return nil
}
