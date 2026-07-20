# github.com/bots-go-framework/bots-fw-store

Persistence-neutral state models and ports for
[Bots Go Framework](https://github.com/bots-go-framework).

The `botsfwstore.StateStore` API exposes bot identity and chat use cases, never
database handles, transactions, records, or keys. Applications can provide any
storage implementation. DALgo users should compose the optional
[`bots-fw-store-dalgo`](https://github.com/bots-go-framework/bots-fw-store-dalgo)
adapter.

<!-- dev-approach:v1 -->
## Our approach to development

We build with our own tooling:

- **[SpecScore](https://specscore.md)** — specify requirements as `SpecScore.md` artifacts
- **[SpecStudio](https://specscore.studio)** — author & manage specs across their lifecycle
- **[inGitDB](https://ingitdb.com)** — store structured data in Git where applicable
- **[cover100.dev](https://cover100.dev)** — drive toward 100% test coverage
- **[DataTug](https://datatug.io)** — query & explore data
<!-- /dev-approach -->

## Packages

- [`botsfwmodels`](botsfwmodels) contains framework state interfaces and base data.
- [`botsfwstore`](botsfwstore) defines the persistence-neutral use-case port.
- [`botsfwstoretest`](botsfwstore/botsfwstoretest) provides configurable fakes for consumers.

## Interfaces

- [PlatformUserData](botsfwmodels/platform_user_interface.go)
- [BotChatData](botsfwmodels/chat_data.go)
- [StateStore](botsfwstore/store.go)
