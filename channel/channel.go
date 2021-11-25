package channel

type Channel interface {
	Parse(URL string) (string, string, bool)
}

var channels = map[string]Channel{
	"lqiyi": newLQiYi(),
}

func GetChannel(name string) Channel {
	return channels[name]
}
