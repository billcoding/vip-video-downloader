package channel

type Channel interface {
	Parse(URL string) string
}

var channels = map[string]Channel{
	"lqiyi": newLQiYi(),
}

func GetChannel(name string) Channel {
	return channels[name]
}
