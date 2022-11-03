package flag

type BoolOrAutoType string

const (
	True  BoolOrAutoType = "true"
	False BoolOrAutoType = "false"
	Auto  BoolOrAutoType = "auto"
)

func GetBoolOrAutoFlag() *enumFlag {
	return NewEnumFlag(GetBoolOrAutoValues(), string(Auto))
}

func GetBoolOrAutoValues() []string {
	return []string{string(True), string(False), string(Auto)}
}
