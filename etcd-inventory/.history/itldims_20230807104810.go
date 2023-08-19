package command

import (
	"fmt"

	"github.com/spf13/cobra"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// NewGetCommand returns the cobra command for "get".
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <key> [range_end]",
		Short: "Gets the key or a range of keys",
		Run:   getCommandFunc,
	}
	return cmd
}

// getCommandFunc executes the "get" command.
func getCommandFunc(cmd *cobra.Command, args []string) {
	key, opts := getGetOp(args)
	ctx, cancel := commandCtx(cmd)
	resp, err := mustClientFromCmd(cmd).Get(ctx, key, opts...)
	cancel()
	if err != nil {
		cobrautl.ExitWithError(cobrautl.ExitError, err)
	}
	display.Get(*resp)
}

func getGetOp(args []string) (string, []clientv3.OpOption) {
	if len(args) == 0 {
		cobrautl.ExitWithError(cobrautl.ExitBadArgs, fmt.Errorf("get command needs one argument as key and an optional argument as range_end"))
	}

	var opts []clientv3.OpOption
	key := args[0]
	if len(args) > 1 {
		opts = append(opts, clientv3.WithRange(args[1]))
	}

	return key, opts
}
