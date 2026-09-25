# Key-Value in-memory store

Simple in-memory KV store in pure golang. Simple, dumb and little bit stupid. For personal use and learning purposes.

## Implementations

- **MutexStore**: single `map[string][]byte` guarded by one `sync.RWMutex`. Simplest option, but all goroutines contend on the same lock.
- **ShardStore**: splits keys across N shards (N must be a power of two), each with its own map and `sync.RWMutex`. Key is hashed with `maphash` to pick a shard, reducing lock contention under concurrent access. Each shard is cache-line padded to avoid false sharing.

All implementations satisfy the `Store` interface (`Get`, `Set`, `Has`, `Delete`), are safe for concurrent use, and copy values in and out to avoid aliasing internal state.

## Testing

Run tests with the race detector:

```sh
go test ./... -race
```

Run benchmarks:

```sh
go test ./kv/... -bench=. -run=^$
```
