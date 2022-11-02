package spawn

import (
	"fmt"
)

const SwitchPrefixShort = "-"

type Res struct {
	Stdout string
	Stderr string
}

type Environ map[string]string

type SubOptionArg struct {
	Tag   string
	Name  string
	Value string
}

func MkRawSpawnArgs(args []interface{}) []string {
	var ret []string
	recursiveAdd(&ret, args)
	return ret
}

func recursiveAdd(ret *[]string, args []interface{}) {
	for _, ga := range args {
		switch tp := ga.(type) {
		case string:
			*ret = append(*ret, ga.(string))
		case WordArg:
			*ret = append(*ret, ga.(WordArg).Name)
		case SubOptionArg:
			soa := ga.(SubOptionArg)
			*ret = append(*ret, SwitchPrefixShort+soa.Tag, soa.Name+"="+soa.Value)
		case []interface{}:
			recursiveAdd(ret, ga.([]interface{}))
		default:
			panic(fmt.Errorf("illegal type %T detected", tp))
		}
	}
}

type WordArg struct {
	Name string
}

func AppendToArgs(toThis []interface{}, this interface{}) []interface{} {
	return append(toThis, this)
}

func PrependToArgs(toThis []interface{}, this interface{}) []interface{} {
	return append([]interface{}{this}, toThis...)
}

func ComposeArgs(argsArr ...[]interface{}) []interface{} {
	var toThis []interface{}
	for _, this := range argsArr {
		for _, in := range this {
			toThis = append(toThis, in)
		}
	}
	return toThis
}

func addEnv(toThisEnv []string, thisEnv Environ) []string {
	for k, v := range thisEnv {
		return append(toThisEnv, fmt.Sprintf("%s=%s", k, v))
	}
	return toThisEnv
}
