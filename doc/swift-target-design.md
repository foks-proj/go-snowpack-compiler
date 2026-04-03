# Swift Target Design for Snowpack Compiler

This document describes the proposed Swift code generation target for the Snowpack
compiler (`snowpc`). It maps every Snowpack language construct to idiomatic Swift,
with comparison to the existing Go output.

The examples below use real `.snowp` definitions from the `go-foks` project
(`proto-src/lib/device.snowp`, `proto-src/lib/kv.snowp`, `proto-src/infra/quota.snowp`).

---

## Runtime Protocol

Every generated type conforms to `SnowpackCodable`, the Swift equivalent of Go's
`Encode`/`Decode`/`Bytes` methods:

```swift
protocol SnowpackCodable: Codable {
    func encode(to encoder: SnowpackEncoder) throws
    static func decode(from decoder: SnowpackDecoder) throws -> Self
    func bytes() -> Data?
}
```

Swift's `Codable` system replaces Go's dual-representation pattern entirely. There
is no need for `Internal__` types — a single struct handles both the public API and
wire serialization via `CodingKeys`.

---

## 1. Typedefs

### Input

```snowp
typedef DeviceNameNormalized = Text;
typedef DeviceName = Text;
typedef DeviceSerial = Uint;
typedef FileID = Blob(16);
```

### Go output (~30 lines each)

```go
type DeviceNameNormalized string
type DeviceNameNormalizedInternal__ string
func (d DeviceNameNormalized) Export() *DeviceNameNormalizedInternal__ { ... }
func (d DeviceNameNormalizedInternal__) Import() DeviceNameNormalized { ... }
func (d *DeviceNameNormalized) Encode(enc rpc.Encoder) error { ... }
func (d *DeviceNameNormalized) Decode(dec rpc.Decoder) error { ... }
func (d DeviceNameNormalized) Bytes() []byte { return nil }
```

### Swift output

```swift
// typealias creates a name alias — fully interchangeable with the
// original type, no conversion needed.
typealias DeviceNameNormalized = String
typealias DeviceName = String
typealias DeviceSerial = UInt64

// Fixed-size blobs get a wrapper struct with a size invariant.
struct FileID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "FileID must be exactly 16 bytes")
        self.data = data
    }
    init() {
        self.data = Data(count: 16)
    }
}
```

For simple type aliases, Swift needs 1 line where Go needs ~30. Fixed-size blobs
need a small wrapper struct.

---

## 2. Enums

### Input

```snowp
enum DeviceStatus {
    ACTIVE @0;
    REVOKED @1;
}
```

### Go output (~25 lines)

```go
type DeviceStatus int
const (
    DeviceStatus_ACTIVE  DeviceStatus = 0
    DeviceStatus_REVOKED DeviceStatus = 1
)
var DeviceStatusMap = map[string]DeviceStatus{"ACTIVE": 0, "REVOKED": 1}
var DeviceStatusRevMap = map[DeviceStatus]string{0: "ACTIVE", 1: "REVOKED"}
type DeviceStatusInternal__ DeviceStatus
func (d DeviceStatusInternal__) Import() DeviceStatus { ... }
func (d DeviceStatus) Export() *DeviceStatusInternal__ { ... }
```

### Swift output

```swift
// Swift enums are first-class types. `: Int` gives each case an integer
// raw value. Forward/reverse maps come for free via RawRepresentable.
// CaseIterable auto-generates an allCases array.
// Codable gives encode/decode for free for raw-value enums.
enum DeviceStatus: Int, Codable, CaseIterable, SnowpackEnum {
    case active = 0
    case revoked = 1

    // Computed property — called without parens: myStatus.label
    var label: String {
        // switch must be exhaustive — the compiler forces you to
        // handle every case. No fallthrough by default.
        switch self {
        case .active:  return "ACTIVE"
        case .revoked: return "REVOKED"
        }
    }

    // Failable initializer — returns nil if string doesn't match.
    init?(label: String) {
        switch label {
        case "ACTIVE":  self = .active
        case "REVOKED": self = .revoked
        default: return nil
        }
    }
}

// All generated enums conform to this protocol.
protocol SnowpackEnum: RawRepresentable where RawValue == Int {
    var label: String { get }
    init?(label: String)
}
```

Key advantage: Swift compiler guarantees exhaustive switch matching. Adding a new
case to the `.snowp` file causes a compile error at every unhandled switch site.

---

## 3. Structs

### Input

```snowp
struct DeviceLabel @0x9650227205486122 {
    deviceType @0 : DeviceType;
    name @1 : DeviceNameNormalized;
    serial @2 : DeviceSerial;
}
```

### Go output (~60 lines)

Public struct + `Internal__` struct with codec tags + Import (nil-checking closures)
+ Export + Encode + Decode + TypeUniqueID + Bytes.

### Swift output

```swift
struct DeviceLabel: SnowpackCodable, Equatable, HasTypeUniqueID {
    var deviceType: DeviceType
    var name: DeviceNameNormalized       // String
    var serial: DeviceSerial             // UInt64

    static let typeUniqueID: TypeUniqueID = 0x9650227205486122

    // CodingKeys maps fields to msgpack array positions, replacing
    // Go's `codec:",toarray"` struct tag.
    private enum CodingKeys: Int, CodingKey {
        case deviceType = 0
        case name = 1
        case serial = 2
    }

    // `throws` replaces Go's (result, error) return pattern.
    // Callers must use `try` — you can't ignore the error.
    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(deviceType)
        try container.encode(name)
        try container.encode(serial)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.deviceType = try container.decode(DeviceType.self)
        self.name = try container.decode(DeviceNameNormalized.self)
        self.serial = try container.decode(DeviceSerial.self)
    }

    // Memberwise init (Swift auto-generates this, but we include it
    // explicitly since the Codable init above suppresses auto-generation).
    init(deviceType: DeviceType, name: DeviceNameNormalized, serial: DeviceSerial) {
        self.deviceType = deviceType
        self.name = name
        self.serial = serial
    }

    func bytes() -> Data? { return nil }
}
```

No `Internal__` dual representation. Swift's `Codable` handles public-to-wire
conversion in one type.

---

## 4. Structs with Optional Fields

### Input

```snowp
struct DeviceInfo {
    status @0 : DeviceStatus;
    dn @1 : Option(DeviceLabelAndName);
    key @2 : MemberRole;
    ctime @3 : Time;
    provisioned @4 : ProvisionInfo;
    revoked @5 : Option(RevokeInfo);
}
```

### Go output

`Option(T)` becomes `*T`. The Import function has ~50 lines of nested nil-checking
closures per optional field.

### Swift output

```swift
struct DeviceInfo: SnowpackCodable, Equatable {
    var status: DeviceStatus
    var dn: DeviceLabelAndName?         // T? is Optional<T>
    var key: MemberRole
    var ctime: Time
    var provisioned: ProvisionInfo
    var revoked: RevokeInfo?

    // Optional<T> is a built-in enum: .none (nil) or .some(T).
    // The compiler forces you to unwrap before using:
    //   if let label = info.dn {
    //       print(label.name)      // safe — label is non-optional here
    //   }

    private enum CodingKeys: Int, CodingKey {
        case status = 0, dn = 1, key = 2, ctime = 3, provisioned = 4, revoked = 5
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(status)
        try container.encodeIfPresent(dn)       // writes nil for .none
        try container.encode(key)
        try container.encode(ctime)
        try container.encode(provisioned)
        try container.encodeIfPresent(revoked)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.status = try container.decode(DeviceStatus.self)
        self.dn = try container.decodeIfPresent(DeviceLabelAndName.self)
        self.key = try container.decode(MemberRole.self)
        self.ctime = try container.decode(Time.self)
        self.provisioned = try container.decode(ProvisionInfo.self)
        self.revoked = try container.decodeIfPresent(RevokeInfo.self)
    }

    init(status: DeviceStatus, dn: DeviceLabelAndName? = nil, key: MemberRole,
         ctime: Time, provisioned: ProvisionInfo, revoked: RevokeInfo? = nil) {
        self.status = status
        self.dn = dn
        self.key = key
        self.ctime = ctime
        self.provisioned = provisioned
        self.revoked = revoked
    }

    func bytes() -> Data? { return nil }
}
```

Go needs ~80 lines of nil-checking closures for Import. Swift handles it with `T?`
and `decodeIfPresent`.

---

## 5. Variants (Tagged Unions)

This is the biggest difference between Go and Swift output.

### Input

```snowp
enum KVListPaginationType {
    None @0;
    MAC @1;
    Time @2;
}

variant KVListPagination switch (t : KVListPaginationType) {
    case None: void;
    case MAC @1: HMAC;
    case Time @2: TimeMicro;
}
```

### Go output (~130 lines)

A struct with tag field + optional pointer fields, getters that panic on wrong case,
constructors, plus Internal__/Switch__ machinery:

```go
type KVListPagination struct {
    T     KVListPaginationType
    F_1__ *HMAC      `json:"f1,omitempty"`
    F_2__ *TimeMicro `json:"f2,omitempty"`
}

func (k KVListPagination) Mac() HMAC {
    if k.F_1__ == nil {
        panic("unexpected nil case; should have been checked")
    }
    if k.T != KVListPaginationType_MAC {
        panic(fmt.Sprintf("unexpected switch value (%v) when Mac is called", k.T))
    }
    return *k.F_1__
}

func NewKVListPaginationWithMac(v HMAC) KVListPagination { ... }
// + Internal__, Switch__, Import, Export, Encode, Decode
```

### Swift output

```swift
// Swift enums with associated values are a native tagged union. Each
// case carries different data. The compiler enforces correct access.
enum KVListPagination: SnowpackCodable, Equatable {
    case none                        // void
    case mac(HMAC)                   // carries an HMAC
    case time(TimeMicro)             // carries a TimeMicro

    // Pattern matching to access data:
    //   switch pagination {
    //   case .none:
    //       print("no pagination")
    //   case .mac(let hmac):            // binds the value
    //       print("paginate by MAC: \(hmac)")
    //   case .time(let t):
    //       print("paginate by time: \(t)")
    //   }
    //
    // Accessing .mac's data when the value is .time is a compile
    // error, not a runtime panic.

    var paginationType: KVListPaginationType {
        switch self {
        case .none:    return .none
        case .mac:     return .mac
        case .time:    return .time
        }
    }

    // Wire format: [tag, {field: value}] matching Go's format.
    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        switch self {
        case .none:
            try container.encode(KVListPaginationType.none)
        case .mac(let hmac):
            try container.encode(KVListPaginationType.mac)
            try container.encode(hmac)
        case .time(let t):
            try container.encode(KVListPaginationType.time)
            try container.encode(t)
        }
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        let tag = try container.decode(KVListPaginationType.self)
        switch tag {
        case .none:
            self = .none
        case .mac:
            self = .mac(try container.decode(HMAC.self))
        case .time:
            self = .time(try container.decode(TimeMicro.self))
        }
    }

    func bytes() -> Data? { return nil }
}

// The standalone tag enum (needed for wire format).
enum KVListPaginationType: Int, Codable, CaseIterable, SnowpackEnum {
    case none = 0
    case mac = 1
    case time = 2

    var label: String {
        switch self {
        case .none: return "None"
        case .mac:  return "MAC"
        case .time: return "Time"
        }
    }

    init?(label: String) {
        switch label {
        case "None": self = .none
        case "MAC":  self = .mac
        case "Time": self = .time
        default: return nil
        }
    }
}
```

| Aspect | Go | Swift |
|--------|-----|-------|
| Lines | ~130 | ~55 |
| Wrong-case access | Runtime panic | Compile error |
| Exhaustive matching | No (silent default) | Yes (compiler enforced) |
| Types needed | 3 (main, Internal__, Switch__) | 1 |
| Constructors | Manual (`NewKVListPaginationWithMac`) | Built-in (`.mac(value)`) |

---

## 6. Protocols (RPC Services)

### Input

```snowp
protocol Quota errors lib.Status @0xfa82ee9c {
    poke @0 () : QuotaPokeArg;
    testSetConfig @1 (config @0 : QuotaConfig);
    makePlan @4 (plan @0 : Plan, opts @1 : MakePlanOpts) -> Plan;
}
```

### Swift output

```swift
// --- Server Interface ---
// `async throws` replaces Go's (context.Context, Arg) -> (Result, error).

protocol QuotaInterface {
    func poke() async throws
    func testSetConfig(config: QuotaConfig) async throws
    func makePlan(plan: Plan, opts: MakePlanOpts) async throws -> Plan
    func wrapError(_ error: Error) -> Status
}

// --- Client Stub ---

struct QuotaClient {
    let cli: GenericClient
    let errorUnwrapper: (Status) -> Error?

    static let protocolID: ProtocolUniqueID = 0xfa82ee9c

    func poke() async throws {
        let arg = QuotaPokeArg()
        try await cli.call(
            method: Method(QuotaClient.protocolID, position: 0, name: "Quota.poke"),
            arg: arg,
            res: Void.self
        )
    }

    func testSetConfig(config: QuotaConfig) async throws {
        let arg = TestSetConfigArg(config: config)
        try await cli.call(
            method: Method(QuotaClient.protocolID, position: 1, name: "Quota.testSetConfig"),
            arg: arg,
            res: Void.self
        )
    }

    func makePlan(plan: Plan, opts: MakePlanOpts) async throws -> Plan {
        let arg = MakePlanArg(plan: plan, opts: opts)
        return try await cli.call(
            method: Method(QuotaClient.protocolID, position: 4, name: "Quota.makePlan"),
            arg: arg,
            res: Plan.self
        )
    }
}

// --- Arg Structs ---

struct QuotaPokeArg: SnowpackCodable {
    func encode(to encoder: Encoder) throws {}
    init(from decoder: Decoder) throws {}
    init() {}
    func bytes() -> Data? { return nil }
}

struct TestSetConfigArg: SnowpackCodable {
    var config: QuotaConfig

    private enum CodingKeys: Int, CodingKey { case config = 0 }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(config)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.config = try container.decode(QuotaConfig.self)
    }
    init(config: QuotaConfig) { self.config = config }
    func bytes() -> Data? { return nil }
}

struct MakePlanArg: SnowpackCodable {
    var plan: Plan
    var opts: MakePlanOpts

    private enum CodingKeys: Int, CodingKey { case plan = 0, opts = 1 }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(plan)
        try container.encode(opts)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.plan = try container.decode(Plan.self)
        self.opts = try container.decode(MakePlanOpts.self)
    }
    init(plan: Plan, opts: MakePlanOpts) {
        self.plan = plan
        self.opts = opts
    }
    func bytes() -> Data? { return nil }
}

// --- Server-side Dispatch ---

struct QuotaProtocolRegistration {
    let impl: QuotaInterface

    func register() -> ProtocolV2 {
        return ProtocolV2(
            name: "Quota",
            id: QuotaClient.protocolID,
            methods: [
                0: makeHandler { _ in try await impl.poke() },
                1: makeHandler { (arg: TestSetConfigArg) in
                    try await impl.testSetConfig(config: arg.config)
                },
                4: makeHandler { (arg: MakePlanArg) -> Plan in
                    try await impl.makePlan(plan: arg.plan, opts: arg.opts)
                },
            ]
        )
    }
}
```

No error unwrapper adapter needed — Swift's `async throws` and typed error
handling replaces Go's `(result, error)` convention and `ErrorUnwrapper` interface.

---

## 7. Type Mapping Reference

| Snowpack | Go | Swift |
|----------|-----|-------|
| `Text` | `string` | `String` |
| `Uint` | `uint64` | `UInt64` |
| `Int` | `int64` | `Int64` |
| `Bool` | `bool` | `Bool` |
| `Blob` | `[]byte` | `Data` |
| `Blob(N)` | `[N]byte` | `Data` (size-checked wrapper) |
| `List(T)` | `[]T` | `[T]` |
| `Option(T)` | `*T` | `T?` |
| `Future(T)` | `[]byte` | `Data` (decoded lazily) |
| `void` | (omitted) | `Void` |

---

## 8. TypeUniqueID / Domain Separation

```swift
typealias TypeUniqueID = UInt64
typealias ProtocolUniqueID = UInt64

// Singleton registry (class = reference type, shared instance).
final class UniqueIDRegistry {
    static let shared = UniqueIDRegistry()
    private init() {}

    private var registered = Set<UInt64>()

    func register(_ id: UInt64) {
        let (inserted, _) = registered.insert(id)
        if !inserted {
            preconditionFailure("Duplicate TypeUniqueID: \(String(id, radix: 16))")
        }
    }
}

// Types with a unique ID conform to this protocol.
protocol HasTypeUniqueID {
    static var typeUniqueID: TypeUniqueID { get }
}

// Registration at load time (replaces Go's init()):
private let _register: Void = {
    UniqueIDRegistry.shared.register(DeviceLabel.typeUniqueID)
}()
```

---

## 9. Architectural Differences Summary

### No Internal__ dual representation

Go generates two types per definition (public + `Internal__` with msgpack tags).
Swift's `Codable` with `CodingKeys` handles both in one type, roughly halving
generated code volume.

### No nil-checking closure trees

Go's Import functions use inline closures for nil-safe conversion of every field.
Swift's `Optional<T>` (`T?`) and `decodeIfPresent` handle this in the type system.

### Native tagged unions

Go fakes variants with a struct + pointer fields + panic-on-wrong-case getters.
Swift enums with associated values are a direct language-level match, with
compile-time safety.

### async/throws replaces context + error tuples

Go's `func(context.Context, Arg) (Result, error)` becomes Swift's
`func(arg: Arg) async throws -> Result`. Error handling is compiler-enforced.

### Estimated code reduction

Based on the `go-foks/proto` output (~79K lines of Go), a Swift target would
produce roughly 35-45K lines for the same definitions.
