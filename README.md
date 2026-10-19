# vz

Go binding for Apple's [Virtualization.framework], forked from
[Code-Hex/vz] for the hosts [ossein] and [zerogravity] run on: Apple silicon,
the current macOS, darwin only.

What the fork carries that upstream does not, in short:

- **One vmnet network per virtual machine.** `VZVmnetNetworkDeviceAttachment`
  (macOS 26) and the `vmnet` package, from upstream's unmerged pull request
  [#205], so a VM owns its network instead of joining the host-wide shared one
  that any other VM product can renumber. `vmnet.NewPerVM` creates one and
  derives the gateway and guest addresses from the subnet the framework
  allocated.
- **A virtio socket you can rely on.** `VirtioSocketListener` closes without
  deadlocking or panicking, `VirtioSocketConnection` can half-close, and
  `VsockConnect`, `VsockConnectRetry` and `VsockListen` are the three wrappers
  every consumer wrote around the first socket device.
- **A build that fails on the wrong SDK** instead of a VM that boots with no
  network at run time, and every Go file constrained to darwin so a non-darwin
  toolchain can load the module.
- **macOS guests, unattended.** The restore image, installer and Mac
  platform bindings zerogravity needs, and `WithMacOSGuestProvisioning`, the
  macOS 27 API that answers Setup Assistant on the first boot after a restore
  (account, auto-login, Remote Login), which upstream does not bind yet.
- **Only what its consumers call.** Devices and notification channels nothing
  here uses were removed, and come back one file at a time when a consumer
  needs them.

Each change, its reason and its evidence: [FORK.md](./FORK.md). Where the code
stands against upstream: [UPSTREAM.md](./UPSTREAM.md).

```go
import (
	"github.com/forkcloser/vz"
	"github.com/forkcloser/vz/vmnet"
)
```

The binding links Virtualization.framework through cgo; a consumer's binary
needs the `com.apple.security.virtualization` entitlement, and the vmnet
attachment needs the macOS 26 SDK to compile.

## Developing

Every workflow runs through `just`: `just lint`, `just test`, `just fix`;
`just --list` shows the rest. The contributor rules are in `AGENTS.md`.

[Virtualization.framework]: https://developer.apple.com/documentation/virtualization
[Code-Hex/vz]: https://github.com/Code-Hex/vz
[ossein]: https://github.com/farcloser/ossein
[zerogravity]: https://github.com/farcloser/zerogravity
[#205]: https://github.com/Code-Hex/vz/pull/205
