package wikipedia

type LimitTooHigh struct {
}

func (e LimitTooHigh) Error() string {
	return "the limit provided can't be more than 500"
}
