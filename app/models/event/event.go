package event

type Event struct {
	ID        int64
	Name      string `binding:"required"`
	Date      string
	CreatedAt string
	UserId    int64
}

func NewEventModel() *Event {
	return &Event{}
}
