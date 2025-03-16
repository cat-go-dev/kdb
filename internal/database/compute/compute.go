package compute

import (
	"context"
	"log/slog"
	"strings"
)

type Compute struct {
	logger *slog.Logger
}

func NewCompute(logger *slog.Logger) (*Compute, error) {
	if logger == nil {
		return nil, errInvalidLogger
	}

	return &Compute{
		logger: logger,
	}, nil
}

const minTokensNum = 2

func (c Compute) Parse(ctx context.Context, query string) (*Command, error) {
	logAttrs := []any{
		slog.String("component", "compute"),
		slog.String("method", "Parse"),
	}

	tokens := strings.Split(query, " ")

	if len(tokens) < minTokensNum {
		err := errNotEnoughArguments
		c.logger.InfoContext(ctx, err.Error(), logAttrs...)
		return nil, err
	}

	commandType, err := c.getCommandType(tokens[0])
	if err != nil {
		c.logger.InfoContext(ctx, err.Error(), logAttrs...)
		return nil, err
	}

	arguments, err := c.getArguments(tokens)
	if err != nil {
		c.logger.InfoContext(ctx, err.Error(), logAttrs...)
		return nil, err
	}

	return &Command{
		Type:      commandType,
		Arguments: arguments,
	}, nil
}

func (c Compute) getCommandType(command string) (CommandType, error) {
	switch command {
	case "GET":
		return Get, nil
	case "SET":
		return Set, nil
	case "DEL":
		return Del, nil
	default:
		return Unknown, errUnknownCommandType
	}
}

func (c Compute) getArguments(tokens []string) (Arguments, error) {
	if len(tokens) < minTokensNum {
		return Arguments{}, errNotEnoughArguments
	}

	keyArg := tokens[1]
	arguments := Arguments{
		Key: Argument(keyArg),
	}
	if len(tokens) > minTokensNum {
		valArg := tokens[2]
		arguments.Value = Argument(valArg)
	}

	return arguments, nil
}
