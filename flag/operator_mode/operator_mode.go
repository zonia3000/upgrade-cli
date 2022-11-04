package operatormode

import "upgrade-cli/flag"

type OperatorMode string

const (
	OLM   OperatorMode = "OLM"
	Plain OperatorMode = "Plain"
	Auto  OperatorMode = "Auto"
)

func GetOperatorModeFlag() *flag.EnumFlag {
	return flag.NewEnumFlag(GetOperatorModeValues(), string(Auto))
}

func GetOperatorModeValues() []string {
	return []string{string(OLM), string(Plain), string(Auto)}
}
