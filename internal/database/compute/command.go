package compute

type Command struct {
	Type      CommandType
	Arguments Arguments
}

type CommandType string

const (
	Get     CommandType = "GET"
	Set     CommandType = "SET"
	Del     CommandType = "DEL"
	Unknown CommandType = "unknown"
)

func (c CommandType) IsGet() bool {
	return c == Get
}

func (c CommandType) IsSet() bool {
	return c == Set
}

func (c CommandType) IsDel() bool {
	return c == Del
}

type Arguments struct {
	Key   Argument
	Value Argument
}

type Argument string
