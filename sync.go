package horloge

const PUBSUB_CHANNEL = "horloge"

type Sync interface {
	Write([]*Job) error
	Read() []Job
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
