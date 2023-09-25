package seq

func Map[T any, U any](ts []T, f func(T) U) (res []U) {
	if ts == nil {
		return nil
	}
	for _, v := range ts {
		res = append(res, f(v))
	}
	return
}
