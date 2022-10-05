package main

func or(channels ...<-chan interface{}) <-chan interface{} {
	res := make(chan interface{})

	defer close(res)

	// если не переданы аргументы
	length := len(channels)
	if length == 0 {
		return res
	}

	// итерируюсь по циклу пока не удастся прочитать из какого-нибудь канала
	for i := 0; ; i++ {
		if i == length {
			i = 0
		}
		select {
		case <-channels[i]:
			return res
		default:
		}
	}
}
