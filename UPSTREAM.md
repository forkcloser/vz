# Upstream

Fork of <https://github.com/Code-Hex/vz> (`main`).

- Forked from: [`0d35cf3`](https://github.com/Code-Hex/vz/commit/0d35cf3a3a8b834ee3b5bf61e4946971b2c0d61a), the base of pull request [#205] as of 2026-09-28, with that pull request's head [`e27a5fb`](https://github.com/Code-Hex/vz/pull/205/commits/e27a5fb55e5936b69a62590f4ce326d8772641a9) applied
- Reviewed through: [`0d35cf3`](https://github.com/Code-Hex/vz/commit/0d35cf3a3a8b834ee3b5bf61e4946971b2c0d61a), reviewed 2026-10-19

The code arrived here from [ossein], where it lived in-tree as `third_party/vz`
from 2026-09-06 to 2026-10-19 with the changes below; this repository's first
commits move it, unchanged beyond the import path.

Every upstream commit after the fork point, through the reviewed one, is in
exactly one of the tables below. A port is a commit of this repository,
adapted to the fork; nothing is merged from upstream.

## Incorporated

| Upstream commit | Here | What |
| --- | --- | --- |

## Not incorporated

| Upstream commit | Why |
| --- | --- |

## Changes upstream does not have

The reasoning and evidence for each is in [FORK.md](./FORK.md).

| Change | Here |
| --- | --- |
| Pull request [#205]: `VZVmnetNetworkDeviceAttachment` and the `vmnet` package (macOS 26), without its `vmnet/fileadapter` and `xpc` packages | initial import |
| The SDK gate in `internal/osversion/virtualization_helper.h` is `#error`: a pre-26 SDK fails the build instead of raising at run time | initial import |
| `IPv4Subnet` and `IPv6Prefix` carry the `macOSAvailable(26)` guard | initial import |
| `VirtioSocketListener.Close` deregisters, releases every waiter, joins in-flight deliveries and closes buffered connections; `Accept` after `Close` returns `ErrListenerClosed` | initial import |
| `VirtioSocketConnection.CloseWrite` | initial import |
| `VirtualMachine.StateChangedNotify`, `NetworkDeviceAttachmentWasDisconnected` and the `NetworkBlockDeviceStorageDeviceAttachment` family removed, with the `go-infinity-channel` dependency | initial import |
| Every Go file constrained to `darwin` | initial import |
| Devices no consumer uses removed: audio, clipboard, console, graphics, keyboard, pointing device, USB, disk images, the Mac platform and installer, Rosetta's directory share, the view; and upstream's tests and examples | initial import |
| `vmnet.NewPerVM`: one vmnet network per VM with derived gateway and guest addressing (from ossein `pkg/vm/network.go`) | initial import |
| `VsockConnect`, `VsockConnectRetry`, `VsockListen` over the VM's first socket device (from ossein `pkg/vm/vm.go`) | initial import |

[#205]: https://github.com/Code-Hex/vz/pull/205
[ossein]: https://github.com/farcloser/ossein
