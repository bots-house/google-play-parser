package ramda

import "reflect"

func Path(path []any, obj any) any {
	return findPath(path, obj, true)
}

func Map[R, T, U any](obj R, fn func(T) U) (out R) {
	v := reflect.ValueOf(obj)

	outValue := reflect.ValueOf(out)

	switch v.Type().Kind() {
	case reflect.Array, reflect.Slice:
		len := v.Len()

		for i := 0; i < len; i++ {
			value, ok := v.Index(i).Interface().(T)
			if !ok {
				continue
			}

			outValue = reflect.Append(outValue, reflect.ValueOf(fn(value)))
		}

	case reflect.Map:
		outValue = reflect.MakeMap(outValue.Type())
		vIter := v.MapRange()

		for vIter.Next() {
			key := vIter.Key()
			value := vIter.Value()

			rawValue, ok := value.Interface().(T)
			if !ok {
				continue
			}

			newValue := reflect.ValueOf(fn(rawValue))

			outValue.SetMapIndex(key, newValue)
		}

	default:
		return obj
	}

	return outValue.Interface().(R)
}

