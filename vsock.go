//go:build darwin

package vz

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"
)

// FORK: host-side helpers over the VM's first virtio socket device. Upstream
// exposes the device and its Connect and Listen; what every consumer then
// writes is the same three wrappers, so they live here once.

// ErrVsock indicates a failure on the host<->guest vsock channel.
var ErrVsock = errors.New("vsock failure")

// halfCloser is the shutdown(SHUT_WR) half of a duplex conn.
//
// The assertion below is the contract VsockConnect and VsockListen advertise:
// the conns they hand out can be shut down one direction at a time, so a relay
// between a vsock hop and anything else can propagate a half-close instead of
// tearing the whole connection down. This fork gave VirtioSocketConnection
// CloseWrite; a refresh that drops it fails to compile here rather than
// silently degrading every proxied connection into "closes when the client
// stops writing".
type halfCloser interface {
	net.Conn
	CloseWrite() error
}

var _ halfCloser = (*VirtioSocketConnection)(nil)

// VsockConnect dials a guest vsock port on the VM's first socket device. The
// guest side must be listening. The returned conn supports CloseWrite.
func VsockConnect(machine *VirtualMachine, port uint32) (net.Conn, error) {
	dev, err := vsockDevice(machine)
	if err != nil {
		return nil, err
	}

	conn, err := dev.Connect(port)
	if err != nil {
		return nil, fmt.Errorf("%w: connect port %d: %w", ErrVsock, port, err)
	}

	return conn, nil
}

// VsockConnectRetry dials a guest vsock port, retrying until timeout — for a
// control channel right after boot. The guest is usually up within a few
// milliseconds, so retries start at 2ms and back off mildly (×1.5, capped at
// 8ms) instead of a fixed coarse poll that would tax every boot.
//
// A failed vsock connect to a not-yet-listening guest costs microseconds,
// while every unit of backoff is added BOOT LATENCY once the guest turns ready
// inside the gap: with a 50ms cap, readiness at ~150ms was detected up to 50ms
// late, and run-to-run boot times wobbled by exactly that. The cap is low on
// purpose: a handful of extra no-op dials in exchange for bounded, small
// detection latency.
func VsockConnectRetry(
	ctx context.Context, machine *VirtualMachine, port uint32, timeout time.Duration,
) (net.Conn, error) {
	const (
		initialDelay = 2 * time.Millisecond
		maxDelay     = 8 * time.Millisecond
	)

	deadline := time.Now().Add(timeout)
	delay := initialDelay

	var last error

	for time.Now().Before(deadline) {
		conn, err := VsockConnect(machine, port)
		if err == nil {
			return conn, nil
		}

		last = err

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: port %d: %w", ErrVsock, port, ctx.Err())
		case <-time.After(delay):
		}

		delay = min(delay*3/2, maxDelay)
	}

	return nil, fmt.Errorf("%w: port %d not answering after %s: %w", ErrVsock, port, timeout, last)
}

// VsockListen opens a host-side listener on a guest vsock port (the guest
// connects out) on the VM's first socket device.
func VsockListen(machine *VirtualMachine, port uint32) (net.Listener, error) {
	dev, err := vsockDevice(machine)
	if err != nil {
		return nil, err
	}

	listener, err := dev.Listen(port)
	if err != nil {
		return nil, fmt.Errorf("%w: listen port %d: %w", ErrVsock, port, err)
	}

	return listener, nil
}

func vsockDevice(machine *VirtualMachine) (*VirtioSocketDevice, error) {
	devs := machine.SocketDevices()
	if len(devs) == 0 {
		return nil, fmt.Errorf("%w: no socket device on VM", ErrVsock)
	}

	return devs[0], nil
}
