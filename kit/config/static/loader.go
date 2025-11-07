package static

type ConfigurationLoader interface {
	LoadConfiguration(context string) error
}
