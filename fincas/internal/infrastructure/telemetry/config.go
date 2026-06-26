package telemetry

// ServiceInfo identifica el microservicio en todos los eventos de telemetría.
type ServiceInfo struct {
	Name        string
	Environment string
}
