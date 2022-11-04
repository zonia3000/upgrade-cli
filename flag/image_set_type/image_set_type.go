package imagesettype

import "upgrade-cli/flag"

type ImageSetType string

const (
	Community       ImageSetType = "Community"
	RedhatCertified ImageSetType = "RedhatCertified"
	Auto            ImageSetType = "Auto"
)

func GetImageSetTypeFlag() *flag.EnumFlag {
	return flag.NewEnumFlag(GetImageSetTypeValues(), string(Auto))
}

func GetImageSetTypeValues() []string {
	return []string{string(RedhatCertified), string(Community), string(Auto)}
}
