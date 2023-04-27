package ramda

func findPath(path []any, obj any, goForward bool) any {
	if len(path) == 0 && goForward {
		return nil
	}

	if len(path) == 0 && !goForward {
		return obj
	}

	var entry any

	switch cast := obj.(type) {
	case map[string][]any:
		idx, ok := path[0].(string)
		if !ok {
			return nil
		}

		entry, ok = cast[idx]
		if !ok {
			return nil
		}

	case map[string]any:
		idx, ok := path[0].(string)
		if !ok {
			return nil
		}

		entry, ok = cast[idx]
		if !ok {
			return nil
		}

	case []any:
		idx, ok := path[0].(int)
		if !ok {
			return nil
		}

		if idx >= len(cast) {
			return nil
		}

		entry = cast[idx]

	default:
		return nil
	}

	newPath := append([]any{}, path[1:]...)
	goForward = len(newPath) > 0

	return findPath(newPath, entry, goForward)
}
