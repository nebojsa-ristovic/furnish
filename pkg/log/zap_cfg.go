package log

// Config is used to initialize a logger.
type Config struct {
	// LogVerbosity is is used to declare a level of logging.
	// Valid values are log.Debug, log.Info, log.Warn and log.Error
	// Values go from the lowest to highest. Choosing a lower value will print all higher values.
	// Ignoring this field or setting it to an incompatible verbosity will default the logger to log.Info.
	LogVerbosity string

	// Name of the executable using the logger.
	ServiceName string
}
