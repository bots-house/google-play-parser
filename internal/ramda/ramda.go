package ramda

func Path(path []any, obj any) any {
	return findPath(path, obj, true)
}
