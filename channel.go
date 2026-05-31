package svcroot

import (
	"errors"
	"fmt"
)

//
// ────────────────────────────────────────
// logical ipc channels.
//

// Channel names a logical local IPC endpoint under a service root.
type Channel string

const (
	// ChannelRPC is the primary control RPC socket.
	ChannelRPC Channel = "rpc"
	// ChannelObserve is the read-only observe socket when configured.
	ChannelObserve Channel = "observe"
)

// ErrChannelDisabled is returned when a channel is not configured in layout.
var ErrChannelDisabled = errors.New("svcroot: channel disabled")

// ChannelUnix returns the Unix socket path for channel under root.
func ChannelUnix(root string, layout *Layout, channel Channel) (string, error) {
	switch channel {
	case ChannelRPC:
		return Socket(root, layout), nil
	case ChannelObserve:
		path, ok := layout.ObserveSocketPath(root)
		if !ok {
			return "", ErrChannelDisabled
		}
		return path, nil
	default:
		return "", fmt.Errorf("svcroot: unknown channel %q", channel)
	}
}
