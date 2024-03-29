package parser

import (
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/bots-house/google-play-parser/internal/ramda"
	"github.com/bots-house/google-play-parser/internal/shared"
)

func ExtractDataWithServiceRequestID(parsed shared.ParsedObject, spec shared.ParsedSpec) any {
	serviceMapping := shared.Keys(parsed.ServiceData)

	filteredMapping := shared.Filter(serviceMapping, func(key string) bool {
		service := parsed.ServiceData[key]
		return service.ID == spec.Clusters.UserServiceID
	})

	path := spec.Clusters.Path

	if len(filteredMapping) > 0 {
		path = append([]any{filteredMapping[0]}, path...)
	}

	return ramda.Path(path, parsed.Data)
}

func Extract[T any, M shared.Mapping](rawApp any, mapping M) (outValue T, ok bool) {
	out := reflect.Indirect(reflect.ValueOf(&outValue))

	reflectMapping := reflect.Indirect(reflect.ValueOf(mapping))
	mappingFields := reflect.VisibleFields(reflectMapping.Type())

	for _, field := range mappingFields {
		appField := out.FieldByName(field.Name)

		if !appField.IsValid() {
			continue
		}

		var (
			path     reflect.Value
			fun      reflect.Value
			withFunc bool
		)

		switch field.Type.Kind() {
		case reflect.Array, reflect.Slice:
			path = reflectMapping.FieldByName(field.Name)
			withFunc = false
		case reflect.Struct:
			innerField := reflectMapping.FieldByName(field.Name)
			path = innerField.FieldByName("Path")
			withFunc = true
			fun = innerField.FieldByName("Fun")
		default:
			continue
		}

		if path.IsZero() {
			continue
		}

		pathValue, ok := path.Interface().([]any)
		if !ok {
			log.Debug().Msg("broken mapping path")
			continue
		}

		result := ramda.Path(pathValue, rawApp)
		if result == nil {
			continue
		}

		reflectResult := reflect.ValueOf(result)

		if withFunc {
			if fnType := fun.Type().In(0); fnType != reflectResult.Type() && fnType.Kind() != reflect.Interface {
				log.Debug().
					Str("func_parameter", fnType.String()).
					Str("argument", reflectResult.Type().String()).
					Msg("can call fun with ")
				continue
			}

			funResult := fun.Call([]reflect.Value{reflectResult})
			reflectResult = funResult[0]
		}

		if reflectResult.Type().Kind() != appField.Type().Kind() && !withFunc {
			log.Debug().
				Str("field", field.Name).
				Any("value", reflectResult.Interface()).
				Str("type", reflectResult.Type().String()).
				Msg("founded value has different app field type")

			continue
		}

		appField.Set(reflectResult)
	}

	result, ok := out.Interface().(T)

	return result, ok
}
