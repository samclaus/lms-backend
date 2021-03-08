package main

func scrub(values []byte) {
	for i := 0; i < len(values); i++ {
		values[i] = 0
	}
}
