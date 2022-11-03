package flag

type InstallationType string

const (
	Community       InstallationType = "Community"
	RedhatCertified InstallationType = "RedhatCertified"

	DefaultInstallationType = Community
)

func GetInstallationTypeFlag() *enumFlag {
	return NewEnumFlag(GetInstallationTypeValues(), string(DefaultInstallationType))
}

func GetInstallationTypeValues() []string {
	return []string{string(RedhatCertified), string(Community)}
}
