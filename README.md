# go-snowpack-compiler

A compiler for the Snowpack IDL (Interface Definition Language) that generates type-safe Go code for RPC services and cryptographic operations. Define your protocols once in `.snowp` files and get structs, enums, tagged unions, serialization, full client/server RPC stubs, and cryptographic domain separation stubs generated automatically.

Ported from the original [TypeScript implementation](https://github.com/foks-proj/node-snowpack-compiler).

## Motivation

Building distributed systems in Go typically means writing a lot of boilerplate: struct definitions, serialization/deserialization, RPC dispatch, client stubs, error wrapping. When protocols span multiple services (client, server, internal infrastructure), keeping types in sync becomes a real problem.

Cryptographic systems add another dimension: every hash, signature, encryption, and commitment must be domain-separated to prevent cross-type attacks. If the same byte sequence can be valid as both a "user profile" and a "payment authorization", an attacker can trick a system into signing one when it thinks it's signing the other. Manually managing domain separators is error-prone and audit-unfriendly.

Snowpack solves both problems by letting you define your types and protocols in a single IDL, then generating all the Go code you need. The generated code uses [go-snowpack-rpc](https://github.com/foks-proj/go-snowpack-rpc) for the wire format (msgpack) and RPC transport, giving you:

- Strongly typed RPC methods with compile-time checking
- Automatic msgpack encode/decode via internal wire types
- Discriminated unions (variants) with type-safe accessors and constructors
- Optional types, lists, fixed-size byte arrays, and futures
- Cryptographic domain separation enforced through the type system
- Clean separation between public API types and internal serialization types

The [go-foks](https://github.com/foks-proj/go-foks) project uses Snowpack extensively, with 60+ `.snowp` files defining protocols across its client-server architecture (local agent protocols, remote server protocols, and internal infrastructure protocols).

## Installation

```bash
go install github.com/foks-proj/go-snowpack-compiler/snowpc@latest
```

Or add it as a tool dependency in your `go.mod`:

```
tool github.com/foks-proj/go-snowpack-compiler/snowpc
```

## Usage

### Single file

```bash
snowpc -l go -p mypackage -i input.snowp -o output.go
```

### Batch (entire directory)

```bash
snowpc -l go -p mypackage -I ./proto-src -O ./proto -v
```

### Stdin/stdout

```bash
cat input.snowp | snowpc -l go -p mypackage
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--lang` | `-l` | Output language (`go`) | `go` |
| `--package` | `-p` | Go package name (required) | |
| `--infile` | `-i` | Input file | stdin |
| `--outfile` | `-o` | Output file | stdout |
| `--input-dir` | `-I` | Input directory (batch mode) | |
| `--output-dir` | `-O` | Output directory (batch mode) | |
| `--ext` | `-e` | File extension to match | `.snowp` |
| `--verbose` | `-v` | Print progress to stderr | `false` |

## IDL Reference

Every `.snowp` file starts with a randomly generated unique file ID (64-bit):

```
@0xdcb1c7e83fa16a34;
```

Generate this randomly when creating a new file. It serves as a namespace identifier for the file.

### Imports

Import types from other packages using language-specific imports:

```
go:import "github.com/myorg/myproject/proto/lib" as lib;
```

Generic imports (language-agnostic) are also supported:

```
import "foo.bar/pickle" as bizzle;
```

### Primitive Types

| Type | Go equivalent | Description |
|------|--------------|-------------|
| `Text` | `string` | UTF-8 string |
| `Uint` | `uint64` | Unsigned 64-bit integer |
| `Int` | `int64` | Signed 64-bit integer |
| `Bool` | `bool` | Boolean |
| `Blob` | `[]byte` | Variable-length byte array |
| `Blob(n)` | `[n]byte` | Fixed-length byte array |
| `Void` | — | No data |

### Container Types

| Type | Go equivalent | Description |
|------|--------------|-------------|
| `List(T)` | `[]T` | Slice of T |
| `Option(T)` | `*T` | Nullable/optional value |
| `Future(T)` | `[]byte` | Deferred data (stored as blob, decoded later) |

### Typedef

Type aliases, optionally with a unique ID:

```
typedef EntityID = Blob;
typedef UID = Blob(33);
typedef HostID @0xa2f3bf71638e062d = Blob(33);
typedef TeamList = List(TeamEntry);
```

### Struct

Product types with explicitly positioned fields:

```
struct FQUser @0x98c90258bf748a2f {
    uid @0 : UID;
    hostID @1 : HostID;
}
```

Each field has a position (`@0`, `@1`, ...) that determines its order in the serialized msgpack array. Fields can use any type including `Option(T)` for optional values:

```
struct LoadUserChainArg {
    uid @0 : lib.UID;
    start @1 : lib.Seqno;
    username @2 : Option(NameSeqnoPair);
    auth @3 : LoadUserChainAuth;
}
```

### Unique IDs and Cryptographic Domain Separation

The `@0x...` tag on structs, variants, and typedefs is a randomly generated 64-bit hex identifier (16 hex chars) used for **cryptographic domain separation**. Protocols also take a unique ID, but use a shorter 32-bit identifier (8 hex chars) to save bandwidth since protocol IDs are transmitted with every RPC message.

Every tag must be globally unique across the entire project. This uniqueness is enforced at three levels:

1. **Compile-time**: A checker tool (see [go-foks's snowp-checker](https://github.com/foks-proj/go-foks/blob/main/tools/snowp-checker/main.go)) scans all `.snowp` source files and rejects duplicate IDs before any code is generated.
2. **Init-time**: The generated code calls `rpc.AddUnique()` for every tagged type in an `init()` function. The RPC runtime tracks all registered IDs and can detect collisions at startup.
3. **Runtime**: Applications can call `rpc.AllUniques()` at startup and panic on duplicates, providing a belt-and-suspenders check.

When a type has a unique ID, the compiler generates:

```go
// A module-level constant for the ID
var FQUserTypeUniqueID = rpc.TypeUniqueID(0x98c90258bf748a2f)

// A getter method on the type
func (f *FQUser) GetTypeUniqueID() rpc.TypeUniqueID {
    return FQUserTypeUniqueID
}
```

This `GetTypeUniqueID()` method is the key to domain separation. In the runtime, `rpc.TypeUniqueID` is a `uint64` that serializes as 8 big-endian bytes. These 8 bytes serve as a cryptographic prefix, ensuring that data hashed, signed, or encrypted under one type can never be confused with another type.

#### How domain separation works in practice

The generated `GetTypeUniqueID()` method makes tagged types satisfy the `CryptoPayloader` interface:

```go
type CryptoPayloader interface {
    Codecable                                    // has Encode/Decode
    GetTypeUniqueID() rpc.TypeUniqueID           // has domain separator
}
```

Cryptographic operations in the application then use this interface to enforce domain separation structurally through the type system. Any function that accepts a `CryptoPayloader` is guaranteed to get a domain separator — you can't accidentally pass in an untagged type:

- **Hashing**: `SHA512/256(TypeUniqueID || msgpack(object))` — the 8-byte type ID is prepended before the serialized object, so the same data under two different types produces different hashes.
- **Signing**: The message signed is `TypeUniqueID || msgpack(object)`, so a signature over a `FQUser` can never validate as a signature over a `TempDHKeySigTemplate`, even if they happened to serialize identically.
- **Symmetric encryption**: The type ID occupies the first 8 bytes of the 24-byte NaCl nonce, so ciphertext produced for one type cannot be decrypted as another type.
- **Public-key encryption**: Same nonce scheme — the type ID is embedded in the nonce.
- **HMAC commitments**: `HMAC-SHA512/256(key, TypeUniqueID || msgpack(object))`.

This means the compiler output directly enables a pattern where the Go type system prevents you from hashing a `GrantLocalViewPermissionPayload` with the domain separator for a `GrantRemoteViewPermissionPayload`. The function signatures enforce it:

```go
// This function ONLY accepts types with a unique ID tag.
// Untagged types won't compile here.
func PrefixedHash(o CryptoPayloader) (StdHash, error) {
    h := sha512.New512_256()
    var buf [8]byte
    o.GetTypeUniqueID().EncodeToBytes(buf[:])  // 8-byte domain prefix
    h.Write(buf[:])
    o.Encode(codec.NewEncoder(h, mh))          // msgpack body
    // ...
}
```

Not every type needs a tag. Only types that participate in cryptographic operations (types that will be hashed, signed, encrypted, or committed to) require one. Types used purely for RPC arguments or local data structures can omit the tag.

### Enum

Enumerated constants with explicit ordinals:

```
enum EntityType {
    User @1;
    Host @2;
    Team @3;
    Device @4;
}
```

### Variant

Discriminated unions (tagged unions / sum types). Each variant switches on a typed discriminator:

```
variant LoadUserChainAuth switch (t : LoadUserChainAuthType) {
    case AsLocalUser: void;
    case Token @1 : lib.PermissionToken;
    case SelfToken @2 : lib.PermissionToken;
    case AsLocalTeam @3 : TeamVOBearerToken;
    case OpenVHost : void;
}
```

- Cases with `void` carry no data
- Cases with data have a position (`@1`, `@2`, ...) for serialization
- A `default` case is supported
- Variants can have a unique ID: `variant Box switch (t : BoxType) @0xb83ffa8f8cfd2ee3 { ... }`

Multiple enum values can map to a single case by comma-separating the labels:

```
variant ParsedUser switch (s : bool) {
    case false @0 : UID;
    case true @1 : NameUtf8;
}

enum ViewType {
    Public @0;
    Basic @1;
    Full @2;
}

variant ViewData switch (t : ViewType) {
    case Public, Basic @0 : BasicView;   // two enum values, same data type
    case Full @1 : FullView;
}
```

When multiple labels share a case, the compiler generates a getter and constructor for each label individually. The switch accessor handles all labels in a single `case` branch.

The generated Go code includes type-safe accessors (`GetToken()`, `GetSwitch()`) and constructors (`NewLoadUserChainAuthWithToken(v)`).

### Protocol

RPC service definitions. Protocols require an error type and optionally specify request/response header types. The protocol's unique ID is 32-bit (8 hex chars), shorter than type IDs to save bandwidth since it is sent with every RPC message:

```
protocol User
    errors lib.Status
    argHeader lib.Header
    resHeader lib.Header @0x823f0899 {

    ping @0 () -> lib.UID;

    setPassphrase @1 (
        key @0 : lib.EntityID,
        salt @1 : lib.PassphraseSalt,
        skwkBox @2 : lib.SecretBox
    );

    loadUserChain @9 (
        a @0 : LoadUserChainArg
    ) : UserLoadUserChainArg -> UserChain;
}
```

**Protocol modifiers:**
- `errors <Type>` (required) — the error type returned by all methods
- `argHeader <Type>` — request header type (sent with every call)
- `resHeader <Type>` — response header type (returned with every call)

**Method syntax:**
- `methodName @position (params) -> ReturnType;` — standard method
- `methodName @position (params);` — method with no return value
- `methodName @position (params) : ArgTypeName -> ReturnType;` — method with a custom arg wrapper type name

Each method has a unique position for dispatch. Parameters follow the same `name @position : Type` syntax as struct fields.

### Doc Comments

Block comments with `/** ... */` are preserved as Go doc comments on the generated types:

```
/** FQUser is a fully-qualified user identifier. */
struct FQUser {
    uid @0 : UID;
    hostID @1 : HostID;
}
```

## Generated Code

For each type, the compiler generates:

1. **Public type** — the API surface you use in your code
2. **Internal type** (suffixed `Internal__`) — msgpack-compatible wire representation
3. **Export/Import methods** — bidirectional conversion between public and internal types
4. **Encode/Decode methods** — serialization via `rpc.Encoder`/`rpc.Decoder`
5. **Unique ID constant and getter** (when tagged with `@0x...`) — for cryptographic domain separation

### Struct example

Given:
```
struct UserChain {
    links @0 : List(lib.LinkOuter);
    merkle @1 : lib.MerklePathsCompressed;
}
```

Generated:
```go
type UserChain struct {
    Links  []lib.LinkOuter
    Merkle lib.MerklePathsCompressed
}

type UserChainInternal__ struct {
    _struct struct{}                           `codec:",toarray"`
    Links   *[](*lib.LinkOuterInternal__)
    Merkle  *lib.MerklePathsCompressedInternal__
}

func (u UserChain) Export() *UserChainInternal__     { /* ... */ }
func (u UserChainInternal__) Import() UserChain      { /* ... */ }
func (u *UserChain) Encode(enc rpc.Encoder) error    { /* ... */ }
func (u *UserChain) Decode(dec rpc.Decoder) error    { /* ... */ }
```

### Tagged struct example

Given:
```
struct GrantLocalViewPermissionPayload @0xf620e4a9845fa063 {
    viewee @0 : lib.PartyID;
    viewer @1 : lib.PartyID;
    tm @2 : lib.Time;
}
```

Generated (in addition to the standard struct output):
```go
var GrantLocalViewPermissionPayloadTypeUniqueID = rpc.TypeUniqueID(0xf620e4a9845fa063)

func (g *GrantLocalViewPermissionPayload) GetTypeUniqueID() rpc.TypeUniqueID {
    return GrantLocalViewPermissionPayloadTypeUniqueID
}
```

This type now satisfies `CryptoPayloader` and can be passed to cryptographic functions that enforce domain separation at the type level.

### Variant example

Given:
```
variant Box switch (t : BoxType) @0xb83ffa8f8cfd2ee3 {
    case NACL @0 : NaclBox;
    case YUBI @1 : YubiBox;
}
```

Generated:
```go
type Box struct {
    T     BoxType
    F_0__ *NaclBox   // NACL
    F_1__ *YubiBox   // YUBI
}

func NewBoxWithNACL(v NaclBox) Box        { /* ... */ }
func NewBoxWithYUBI(v YubiBox) Box        { /* ... */ }
func (b Box) GetNACL() (NaclBox, error)   { /* ... */ }
func (b Box) GetYUBI() (YubiBox, error)   { /* ... */ }
func (b Box) GetSwitch() (BoxType, error) { /* ... */ }

var BoxTypeUniqueID = rpc.TypeUniqueID(0xb83ffa8f8cfd2ee3)
func (b *Box) GetTypeUniqueID() rpc.TypeUniqueID { return BoxTypeUniqueID }
```

### Protocol example

Given:
```
protocol User
    errors lib.Status
    argHeader lib.Header
    resHeader lib.Header @0x823f0899 {

    ping @0 () -> lib.UID;
    loadUserChain @9 (a @0 : LoadUserChainArg) : UserLoadUserChainArg -> UserChain;
}
```

Generated:
```go
// Server interface — implement this
type UserInterface interface {
    Ping(context.Context) (lib.UID, error)
    LoadUserChain(context.Context, LoadUserChainArg) (UserChain, error)
    ErrorWrapper() func(error) lib.Status
    CheckArgHeader(ctx context.Context, h lib.Header) error
    MakeResHeader() lib.Header
}

// Protocol registration
var UserProtocolID rpc.ProtocolUniqueID = rpc.ProtocolUniqueID(0x823f0899)
func UserProtocol(i UserInterface) rpc.ProtocolV2 { /* ... */ }

// Client stub — use this to make RPC calls
type UserClient struct {
    Cli            rpc.GenericClient
    ErrorUnwrapper UserErrorUnwrapper
    MakeArgHeader  func() lib.Header
    CheckResHeader func(context.Context, lib.Header) error
}

func (c UserClient) Ping(ctx context.Context) (lib.UID, error)                              { /* ... */ }
func (c UserClient) LoadUserChain(ctx context.Context, a LoadUserChainArg) (UserChain, error) { /* ... */ }
```

## Integrating Into Your Project

### 1. Define your protocols

Create a directory for your `.snowp` source files:

```
proto-src/
  lib/
    types.snowp
  rem/
    service.snowp
```

`proto-src/lib/types.snowp`:
```
@0xaabbccdd11223344;

typedef UserID = Blob(33);
typedef Name = Text;

struct UserInfo {
    id @0 : UserID;
    name @1 : Name;
    active @2 : Bool;
}

enum Role {
    Admin @0;
    Member @1;
    Guest @2;
}
```

`proto-src/rem/service.snowp`:
```
@0x1122334455667788;

go:import "github.com/myorg/myproject/proto/lib" as lib;

struct CreateUserArg {
    name @0 : lib.Name;
    role @1 : lib.Role;
}

protocol UserService
    errors lib.Status @0xdeadbeef {

    createUser @0 (
        arg @0 : CreateUserArg
    ) -> lib.UserInfo;

    getUser @1 (
        id @0 : lib.UserID
    ) -> lib.UserInfo;

    listUsers @2 () -> List(lib.UserInfo);
}
```

### 2. Set up code generation

Create output directories and `gen.go` files with `go:generate` directives:

```
proto/
  lib/
    gen.go
  rem/
    gen.go
```

`proto/lib/gen.go`:
```go
package lib

//go:generate go tool github.com/foks-proj/go-snowpack-compiler/snowpc -l go -p lib -I ../../proto-src/lib -O . -v
//go:generate go fmt .
```

`proto/rem/gen.go`:
```go
package rem

//go:generate go tool github.com/foks-proj/go-snowpack-compiler/snowpc -l go -p rem -I ../../proto-src/rem -O . -v
//go:generate go fmt .
```

### 3. Generate code

```bash
go generate ./...
```

Or add a Makefile target:

```makefile
.PHONY: proto
proto:
	go generate ./...
```

### 4. Implement the server

```go
package server

import (
    "context"
    "github.com/myorg/myproject/proto/lib"
    "github.com/myorg/myproject/proto/rem"
    "github.com/foks-proj/go-snowpack-rpc/rpc"
)

type UserServer struct{}

func (s *UserServer) CreateUser(ctx context.Context, arg rem.CreateUserArg) (lib.UserInfo, error) {
    // your implementation here
}

func (s *UserServer) GetUser(ctx context.Context, id lib.UserID) (lib.UserInfo, error) {
    // your implementation here
}

func (s *UserServer) ListUsers(ctx context.Context) ([]lib.UserInfo, error) {
    // your implementation here
}

func (s *UserServer) ErrorWrapper() func(error) lib.Status {
    return func(err error) lib.Status { /* wrap error */ }
}

// Register with the RPC server:
func register(srv *rpc.Server) {
    srv.RegisterV2(rem.UserServiceProtocol(&UserServer{}))
}
```

### 5. Use the client

```go
client := rem.UserServiceClient{
    Cli: genericClient, // from go-snowpack-rpc
}

user, err := client.GetUser(ctx, someUserID)
```

## How the Compiler Works

The compiler follows a classic multi-stage pipeline:

```
.snowp source → Lexer → Parser → AST → Emitter → .go output
```

### Lexer (`lib/lexer.go`)

A hand-written state-machine lexer that tokenizes `.snowp` input. It uses a push/pop state stack for handling nested contexts (braces, parentheses) and produces tokens via a Go channel. The lexer recognizes 40+ token types including keywords (`struct`, `enum`, `variant`, `protocol`, `typedef`), primitive type names (`Text`, `Uint`, `Blob`, etc.), structural tokens (`{`, `}`, `@`, `->`, `:`, `;`), and literals (identifiers, integers, hex values, quoted strings). Doc comments (`/** ... */`) are preserved as tokens for attachment to generated types.

### Parser (`lib/parser.y` → `lib/parser.go`)

An LALR parser generated by [goyacc](https://pkg.go.dev/golang.org/x/tools/cmd/goyacc) from a yacc grammar file. The grammar defines the full Snowpack language structure: file IDs, imports, typedefs, structs with positioned fields, enums, variants with switch/case, and protocols with methods and modifiers. The parser produces a typed AST.

To regenerate the parser after modifying the grammar:

```bash
go generate ./lib/...
```

### AST (`lib/ast.go`)

The abstract syntax tree uses Go interfaces for polymorphism:

- `Statement` interface — `Typedef`, `Struct`, `Enum`, `Variant`, `Protocol`, `Import`
- `Type` interface — `Text`, `Uint`, `Int`, `Bool`, `Blob`, `Void`, `List`, `Option`, `Future`, `DerivedType`
- `CaseLabel` interface — `CaseLabelIdentifier`, `CaseLabelNumber`, `CaseLabelBool`

Each AST node carries enough information for code generation: field positions, unique IDs, decorators (doc comments), and type references.

### Emitter (`lib/emit.go`)

The `Emitter` interface abstracts language-specific code generation. Currently only `GoEmitter` is implemented. The emitter walks the AST and generates:

1. **Preamble** — package declaration, imports (computed from an inventory pass over the AST)
2. **Types** — for each typedef, struct, enum, variant:
   - Public Go types (exported, PascalCase)
   - Internal wire types (suffixed `Internal__`, tagged for msgpack with `codec:",toarray"`)
   - Import/Export conversion methods
   - Encode/Decode serialization methods
   - Unique ID constant and `GetTypeUniqueID()` getter (when tagged with `@0x...`)
3. **Protocols** — for each protocol:
   - Server interface with handler methods
   - Protocol factory function for registration
   - Client stub struct with typed RPC call methods
   - Arg/response types for methods with multiple parameters
   - Error wrapper/unwrapper adapters
   - Header check/make hooks
4. **Postamble** — an `init()` function that calls `rpc.AddUnique()` for every tagged type and protocol, registering all unique IDs with the RPC runtime for collision detection

## Runtime Dependency

Generated code depends on the [go-snowpack-rpc](https://github.com/foks-proj/go-snowpack-rpc) runtime library:

```bash
go get github.com/foks-proj/go-snowpack-rpc
```

## License

MIT
