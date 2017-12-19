package middleware

type Counter struct {
	totalCount uint32;
	todayCount uint32;
	currentHourCount uint32;
	currentMinuteCount uint32;
}

func NewCounter() *Counter {
	return &Counter{
		0, 0, 0, 0,
	}
}

func (this *Counter) handle() bool {
	return true
}