package jutzo

// ConfigurationProvider provides an interface for getting
// configuration information from the mechanism used by the
// client environment. The client of the Jutzo engine needs
// to provide an object that satisfies this interface
type ConfigurationProvider interface {

	// GetConfigurationString returns the configuration item
	// value. The second return determines if the item was
	// actually provided. If the item was not provided the
	// routine should return ("", false); otherwise it should
	// provide (value, true)
	GetConfigurationString(name string) (string, bool)

	// GetConfigurationInt is similar to GetConfigurationString
	// but validates that the item is both provided and an integer
	// If the item is not provided or not an integer, this should
	// return (0, false); otherwise it should return (value, true)
	GetConfigurationInt(name string) (int, bool)
}
