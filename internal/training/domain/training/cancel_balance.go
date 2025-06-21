package training

func CancelBalanceDelta(tr Training) int {
	if tr.CanBeCanceledForFree() {
		return 1
	}

	return 0
}
