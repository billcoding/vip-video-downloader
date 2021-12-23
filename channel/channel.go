package channel

type Channel interface {
	Parse(URL string) (string, string, bool)
}

var channels = map[string]Channel{
	"c1": C1(),
}

func GetChannel(name string) Channel {
	return channels[name]
}
