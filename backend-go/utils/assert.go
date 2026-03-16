package utils

func Assert(cond bool, mes string) {
	if !cond {
		panic(mes)
	}
}
