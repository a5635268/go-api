package redis

const (
	Prefix = "goapi:"
)

func getDefaultKay(name string) string{
	return Prefix + name
}
