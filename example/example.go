package example

func Bla(input int, bla int) int {
	if input >= 10 {
		return 10
	}
	if input >= 5 && bla == 10 {
		return 5
	}
	if bla == 17 {
		return 17
	}
	return 0
}
