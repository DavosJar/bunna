package decorator

import (
    "reflect"
    "github.com/davosjar/bunna/services/identidad/internal/infrastructure/telemetry"
)

// validationError marks validation errors.
type validationError interface{ ValidationError() }
// businessError marks business errors.
type businessError interface{ BusinessError() }
// infrastructureError marks infra errors.
type infrastructureError interface{ InfrastructureError() }

func classifyResult(err error) string {
    if err == nil {
        return "success"
    }
    switch {
    case implements(err, (*validationError)(nil)):
        return "validation_error"
    case implements(err, (*businessError)(nil)):
        return "business_error"
    case implements(err, (*infrastructureError)(nil)):
        return "infrastructure_error"
    default:
        return "unexpected_error"
    }
}

func implements(err error, target interface{}) bool {
    if err == nil {
        return false
    }
    targetType := reflect.TypeOf(target).Elem()
    errType := reflect.TypeOf(err)
    return errType.Implements(targetType)
}

func determineLevel(result string) string {
    switch result {
    case "success":
        return "INFO"
    case "validation_error":
        return "WARN"
    case "business_error", "infrastructure_error", "unexpected_error":
        return "ERROR"
    default:
        return "ERROR"
    }
}

func safeCommand(cmd interface{}) map[string]any {
    if cmd == nil {
        return map[string]any{}
    }
    v := reflect.ValueOf(cmd)
    method := v.MethodByName("ToLog")
    if !method.IsValid() {
        return map[string]any{}
    }
    // call without arguments
    results := method.Call(nil)
    if len(results) == 0 {
        return map[string]any{}
    }
    if m, ok := results[0].Interface().(map[string]any); ok {
        return m
    }
    return map[string]any{}
}

// LogPayload is an alias for telemetry.LogPayload used in decorators.
type LogPayload = telemetry.LogPayload
