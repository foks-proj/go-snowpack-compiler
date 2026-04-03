// Auto-generated to Swift types and interfaces using snowpc (https://github.com/foks-proj/go-snowpack-compiler)
//  Input file: ../../proto-src/lib/kv.snowp

import Foundation
import SnowpackRPC

// ---------------------------------------------------------------
//  typedefs (simple aliases)
// ---------------------------------------------------------------

typealias KVShardID = UInt64
typealias KVVersion = UInt64
typealias LocalFSPath = String
typealias Size = UInt64
typealias KVPathComponent = String
typealias KVNameNonce = NaclNonce
typealias Offset = UInt64
typealias Chunk = Data
typealias KVMacKey = HMACKey
typealias KVBoxKey = SecretBoxKey

// ---------------------------------------------------------------
//  typedefs with TypeUniqueID
// ---------------------------------------------------------------

struct FileKeySeed: SnowpackCodable, Equatable, HasTypeUniqueID {
    let raw: SecretSeed32

    static let typeUniqueID: TypeUniqueID = 0xb6487f843b10cfce

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(raw)
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        self.raw = try container.decode(SecretSeed32.self)
    }

    init(_ raw: SecretSeed32) { self.raw = raw }

    func bytes() -> Data? { return raw.bytes() }
}

struct DirKeySeed: SnowpackCodable, Equatable, HasTypeUniqueID {
    let raw: SecretSeed32

    static let typeUniqueID: TypeUniqueID = 0x8aece6566b244356

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(raw)
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        self.raw = try container.decode(SecretSeed32.self)
    }

    init(_ raw: SecretSeed32) { self.raw = raw }

    func bytes() -> Data? { return raw.bytes() }
}

struct KVNamePlaintext: SnowpackCodable, Equatable, HasTypeUniqueID {
    let raw: String

    static let typeUniqueID: TypeUniqueID = 0x9f9fe83fe9f0d475

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(raw)
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        self.raw = try container.decode(String.self)
    }

    init(_ raw: String) { self.raw = raw }

    func bytes() -> Data? { return nil }
}

struct KVPath: SnowpackCodable, Equatable, HasTypeUniqueID {
    let raw: String

    static let typeUniqueID: TypeUniqueID = 0xc05a71f8cf1dfa5c

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(raw)
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        self.raw = try container.decode(String.self)
    }

    init(_ raw: String) { self.raw = raw }

    func bytes() -> Data? { return nil }
}

struct ChunkPlaintext: SnowpackCodable, Equatable, HasTypeUniqueID {
    let raw: Data

    static let typeUniqueID: TypeUniqueID = 0x8d2e290071d39250

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(raw)
    }

    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        self.raw = try container.decode(Data.self)
    }

    init(_ raw: Data) { self.raw = raw }

    func bytes() -> Data? { return raw }
}

// ---------------------------------------------------------------
//  typedefs (fixed-size blobs)
// ---------------------------------------------------------------

struct ShortPartyID: SnowpackCodable, Equatable {
    let data: Data  // 17 bytes

    init(data: Data) {
        precondition(data.count == 17, "ShortPartyID must be exactly 17 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 17) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 17, "ShortPartyID must be exactly 17 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct FileID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "FileID must be exactly 16 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 16) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 16, "FileID must be exactly 16 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct SmallFileID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "SmallFileID must be exactly 16 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 16) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 16, "SmallFileID must be exactly 16 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct SymlinkID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "SymlinkID must be exactly 16 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 16) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 16, "SymlinkID must be exactly 16 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct DirID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "DirID must be exactly 16 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 16) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 16, "DirID must be exactly 16 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct DirentID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "DirentID must be exactly 16 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 16) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 16, "DirentID must be exactly 16 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct KVNodeID: SnowpackCodable, Equatable {
    let data: Data  // 17 bytes

    init(data: Data) {
        precondition(data.count == 17, "KVNodeID must be exactly 17 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 17) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 17, "KVNodeID must be exactly 17 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

struct KVUploadID: SnowpackCodable, Equatable {
    let data: Data  // 16 bytes

    init(data: Data) {
        precondition(data.count == 16, "KVUploadID must be exactly 16 bytes")
        self.data = data
    }
    init() { self.data = Data(count: 16) }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        try container.encode(data)
    }
    init(from decoder: Decoder) throws {
        let container = try decoder.singleValueContainer()
        let d = try container.decode(Data.self)
        precondition(d.count == 16, "KVUploadID must be exactly 16 bytes")
        self.data = d
    }
    func bytes() -> Data? { return data }
}

// ---------------------------------------------------------------
//  enum KVNodeType
// ---------------------------------------------------------------

enum KVNodeType: Int, Codable, CaseIterable, SnowpackEnum {
    case none = 0
    case dir = 1
    case file = 2
    case smallFile = 3
    case symlink = 4

    var label: String {
        switch self {
        case .none:      return "None"
        case .dir:       return "Dir"
        case .file:      return "File"
        case .smallFile: return "SmallFile"
        case .symlink:   return "Symlink"
        }
    }

    init?(label: String) {
        switch label {
        case "None":      self = .none
        case "Dir":       self = .dir
        case "File":      self = .file
        case "SmallFile": self = .smallFile
        case "Symlink":   self = .symlink
        default: return nil
        }
    }
}

// ---------------------------------------------------------------
//  enum KVDirStatus
// ---------------------------------------------------------------

enum KVDirStatus: Int, Codable, CaseIterable, SnowpackEnum {
    case active = 0
    case encrypting = 1
    case dead = 2

    var label: String {
        switch self {
        case .active:     return "Active"
        case .encrypting: return "Encrypting"
        case .dead:       return "Dead"
        }
    }

    init?(label: String) {
        switch label {
        case "Active":     self = .active
        case "Encrypting": self = .encrypting
        case "Dead":       self = .dead
        default: return nil
        }
    }
}

// ---------------------------------------------------------------
//  enum KVListPaginationType
// ---------------------------------------------------------------

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

// ---------------------------------------------------------------
//  struct KVNameNonceInput
// ---------------------------------------------------------------

struct KVNameNonceInput: SnowpackCodable, Equatable, HasTypeUniqueID {
    var parentDir: DirID
    var name: KVNamePlaintext

    static let typeUniqueID: TypeUniqueID = 0x8f02a15745471874

    private enum CodingKeys: Int, CodingKey {
        case parentDir = 0
        case name = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(parentDir)
        try container.encode(name)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.parentDir = try container.decode(DirID.self)
        self.name = try container.decode(KVNamePlaintext.self)
    }

    init(parentDir: DirID, name: KVNamePlaintext) {
        self.parentDir = parentDir
        self.name = name
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct PathVersionVector
// ---------------------------------------------------------------

struct PathVersionVector: SnowpackCodable, Equatable {
    var root: KVVersion
    var path: [DirVersion]

    private enum CodingKeys: Int, CodingKey {
        case root = 0
        case path = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(root)
        try container.encode(path)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.root = try container.decode(KVVersion.self)
        self.path = try container.decode([DirVersion].self)
    }

    init(root: KVVersion, path: [DirVersion]) {
        self.root = root
        self.path = path
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct DirVersion
// ---------------------------------------------------------------

struct DirVersion: SnowpackCodable, Equatable {
    var id: DirID
    var vers: KVVersion
    var de: [DirentVersion]

    private enum CodingKeys: Int, CodingKey {
        case id = 0
        case vers = 1
        case de = 2
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(id)
        try container.encode(vers)
        try container.encode(de)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.id = try container.decode(DirID.self)
        self.vers = try container.decode(KVVersion.self)
        self.de = try container.decode([DirentVersion].self)
    }

    init(id: DirID, vers: KVVersion, de: [DirentVersion]) {
        self.id = id
        self.vers = vers
        self.de = de
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct DirentVersion
// ---------------------------------------------------------------

struct DirentVersion: SnowpackCodable, Equatable {
    var id: DirentID
    var vers: KVVersion

    private enum CodingKeys: Int, CodingKey {
        case id = 0
        case vers = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(id)
        try container.encode(vers)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.id = try container.decode(DirentID.self)
        self.vers = try container.decode(KVVersion.self)
    }

    init(id: DirentID, vers: KVVersion) {
        self.id = id
        self.vers = vers
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct RolePair
// ---------------------------------------------------------------

struct RolePair: SnowpackCodable, Equatable {
    var read: Role
    var write: Role

    private enum CodingKeys: Int, CodingKey {
        case read = 0
        case write = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(read)
        try container.encode(write)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.read = try container.decode(Role.self)
        self.write = try container.decode(Role.self)
    }

    init(read: Role, write: Role) {
        self.read = read
        self.write = write
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct RolePairOpt
// ---------------------------------------------------------------

struct RolePairOpt: SnowpackCodable, Equatable {
    var read: Role?
    var write: Role?

    private enum CodingKeys: Int, CodingKey {
        case read = 0
        case write = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encodeIfPresent(read)
        try container.encodeIfPresent(write)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.read = try container.decodeIfPresent(Role.self)
        self.write = try container.decodeIfPresent(Role.self)
    }

    init(read: Role? = nil, write: Role? = nil) {
        self.read = read
        self.write = write
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVDirent
// ---------------------------------------------------------------

struct KVDirent: SnowpackCodable, Equatable {
    var parentDir: DirID
    var id: DirentID
    var value: KVNodeID
    var version: KVVersion
    var dirVersion: KVVersion
    var writeRole: Role
    var nameMac: HMAC
    var nameBox: SecretBox
    var dirStatus: KVDirStatus
    var bindingMac: HMAC
    var ctime: TimeMicro

    private enum CodingKeys: Int, CodingKey {
        case parentDir = 0
        case id = 1
        case value = 2
        case version = 3
        case dirVersion = 4
        case writeRole = 5
        case nameMac = 6
        case nameBox = 7
        case dirStatus = 8
        case bindingMac = 9
        case ctime = 10
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(parentDir)
        try container.encode(id)
        try container.encode(value)
        try container.encode(version)
        try container.encode(dirVersion)
        try container.encode(writeRole)
        try container.encode(nameMac)
        try container.encode(nameBox)
        try container.encode(dirStatus)
        try container.encode(bindingMac)
        try container.encode(ctime)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.parentDir = try container.decode(DirID.self)
        self.id = try container.decode(DirentID.self)
        self.value = try container.decode(KVNodeID.self)
        self.version = try container.decode(KVVersion.self)
        self.dirVersion = try container.decode(KVVersion.self)
        self.writeRole = try container.decode(Role.self)
        self.nameMac = try container.decode(HMAC.self)
        self.nameBox = try container.decode(SecretBox.self)
        self.dirStatus = try container.decode(KVDirStatus.self)
        self.bindingMac = try container.decode(HMAC.self)
        self.ctime = try container.decode(TimeMicro.self)
    }

    init(parentDir: DirID, id: DirentID, value: KVNodeID, version: KVVersion,
         dirVersion: KVVersion, writeRole: Role, nameMac: HMAC, nameBox: SecretBox,
         dirStatus: KVDirStatus, bindingMac: HMAC, ctime: TimeMicro) {
        self.parentDir = parentDir
        self.id = id
        self.value = value
        self.version = version
        self.dirVersion = dirVersion
        self.writeRole = writeRole
        self.nameMac = nameMac
        self.nameBox = nameBox
        self.dirStatus = dirStatus
        self.bindingMac = bindingMac
        self.ctime = ctime
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct SmallFileBox
// ---------------------------------------------------------------

struct SmallFileBox: SnowpackCodable, Equatable {
    var rg: RoleAndGen
    var dataBox: NaclCiphertext

    private enum CodingKeys: Int, CodingKey {
        case rg = 0
        case dataBox = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(rg)
        try container.encode(dataBox)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.rg = try container.decode(RoleAndGen.self)
        self.dataBox = try container.decode(NaclCiphertext.self)
    }

    init(rg: RoleAndGen, dataBox: NaclCiphertext) {
        self.rg = rg
        self.dataBox = dataBox
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct SeedBox
// ---------------------------------------------------------------

struct SeedBox: SnowpackCodable, Equatable {
    var rg: RoleAndGen
    var box: SecretBox

    private enum CodingKeys: Int, CodingKey {
        case rg = 0
        case box = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(rg)
        try container.encode(box)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.rg = try container.decode(RoleAndGen.self)
        self.box = try container.decode(SecretBox.self)
    }

    init(rg: RoleAndGen, box: SecretBox) {
        self.rg = rg
        self.box = box
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  variant KVListPagination (tagged union)
// ---------------------------------------------------------------

enum KVListPagination: SnowpackCodable, Equatable {
    case none
    case mac(HMAC)
    case time(TimeMicro)

    var paginationType: KVListPaginationType {
        switch self {
        case .none:    return .none
        case .mac:     return .mac
        case .time:    return .time
        }
    }

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

// ---------------------------------------------------------------
//  struct KVUsageStats
// ---------------------------------------------------------------

struct KVUsageStats: SnowpackCodable, Equatable {
    var num: UInt64
    var sum: Size

    private enum CodingKeys: Int, CodingKey {
        case num = 0
        case sum = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(num)
        try container.encode(sum)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.num = try container.decode(UInt64.self)
        self.sum = try container.decode(Size.self)
    }

    init(num: UInt64, sum: Size) {
        self.num = num
        self.sum = sum
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVUsageStatsChunked
// ---------------------------------------------------------------

struct KVUsageStatsChunked: SnowpackCodable, Equatable {
    var base: KVUsageStats
    var numChunks: UInt64

    private enum CodingKeys: Int, CodingKey {
        case base = 0
        case numChunks = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(base)
        try container.encode(numChunks)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.base = try container.decode(KVUsageStats.self)
        self.numChunks = try container.decode(UInt64.self)
    }

    init(base: KVUsageStats, numChunks: UInt64) {
        self.base = base
        self.numChunks = numChunks
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVUsage
// ---------------------------------------------------------------

struct KVUsage: SnowpackCodable, Equatable {
    var small: KVUsageStats
    var large: KVUsageStatsChunked

    private enum CodingKeys: Int, CodingKey {
        case small = 0
        case large = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(small)
        try container.encode(large)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.small = try container.decode(KVUsageStats.self)
        self.large = try container.decode(KVUsageStatsChunked.self)
    }

    init(small: KVUsageStats, large: KVUsageStatsChunked) {
        self.small = small
        self.large = large
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct SeedBoxExternalNonce
// ---------------------------------------------------------------

struct SeedBoxExternalNonce: SnowpackCodable, Equatable {
    var rg: RoleAndGen
    var ctext: NaclCiphertext

    private enum CodingKeys: Int, CodingKey {
        case rg = 0
        case ctext = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(rg)
        try container.encode(ctext)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.rg = try container.decode(RoleAndGen.self)
        self.ctext = try container.decode(NaclCiphertext.self)
    }

    init(rg: RoleAndGen, ctext: NaclCiphertext) {
        self.rg = rg
        self.ctext = ctext
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVDir
// ---------------------------------------------------------------

struct KVDir: SnowpackCodable, Equatable {
    var id: DirID
    var version: KVVersion
    var box: SeedBoxExternalNonce
    var writeRole: Role
    var status: KVDirStatus

    private enum CodingKeys: Int, CodingKey {
        case id = 0
        case version = 1
        case box = 2
        case writeRole = 3
        case status = 4
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(id)
        try container.encode(version)
        try container.encode(box)
        try container.encode(writeRole)
        try container.encode(status)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.id = try container.decode(DirID.self)
        self.version = try container.decode(KVVersion.self)
        self.box = try container.decode(SeedBoxExternalNonce.self)
        self.writeRole = try container.decode(Role.self)
        self.status = try container.decode(KVDirStatus.self)
    }

    init(id: DirID, version: KVVersion, box: SeedBoxExternalNonce,
         writeRole: Role, status: KVDirStatus) {
        self.id = id
        self.version = version
        self.box = box
        self.writeRole = writeRole
        self.status = status
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVRoot
// ---------------------------------------------------------------

struct KVRoot: SnowpackCodable, Equatable {
    var root: DirID
    var vers: KVVersion
    var rg: RoleAndGen
    var bindingMac: HMAC

    private enum CodingKeys: Int, CodingKey {
        case root = 0
        case vers = 1
        case rg = 2
        case bindingMac = 3
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(root)
        try container.encode(vers)
        try container.encode(rg)
        try container.encode(bindingMac)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.root = try container.decode(DirID.self)
        self.vers = try container.decode(KVVersion.self)
        self.rg = try container.decode(RoleAndGen.self)
        self.bindingMac = try container.decode(HMAC.self)
    }

    init(root: DirID, vers: KVVersion, rg: RoleAndGen, bindingMac: HMAC) {
        self.root = root
        self.vers = vers
        self.rg = rg
        self.bindingMac = bindingMac
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct LargeFileMetadata
// ---------------------------------------------------------------

struct LargeFileMetadata: SnowpackCodable, Equatable {
    var rg: RoleAndGen
    var keySeed: SecretBox
    var vers: KVVersion
    var customMetadata: SecretBox?

    private enum CodingKeys: Int, CodingKey {
        case rg = 0
        case keySeed = 1
        case vers = 2
        case customMetadata = 3
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(rg)
        try container.encode(keySeed)
        try container.encode(vers)
        try container.encodeIfPresent(customMetadata)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.rg = try container.decode(RoleAndGen.self)
        self.keySeed = try container.decode(SecretBox.self)
        self.vers = try container.decode(KVVersion.self)
        self.customMetadata = try container.decodeIfPresent(SecretBox.self)
    }

    init(rg: RoleAndGen, keySeed: SecretBox, vers: KVVersion,
         customMetadata: SecretBox? = nil) {
        self.rg = rg
        self.keySeed = keySeed
        self.vers = vers
        self.customMetadata = customMetadata
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct UploadFinal
// ---------------------------------------------------------------

struct UploadFinal: SnowpackCodable, Equatable {
    var sz: Size
    var chunkSum: StdHash

    private enum CodingKeys: Int, CodingKey {
        case sz = 0
        case chunkSum = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(sz)
        try container.encode(chunkSum)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.sz = try container.decode(Size.self)
        self.chunkSum = try container.decode(StdHash.self)
    }

    init(sz: Size, chunkSum: StdHash) {
        self.sz = sz
        self.chunkSum = chunkSum
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct UploadChunk
// ---------------------------------------------------------------

struct UploadChunk: SnowpackCodable, Equatable {
    var data: NaclCiphertext
    var offset: Offset
    var final: UploadFinal?

    private enum CodingKeys: Int, CodingKey {
        case data = 0
        case offset = 1
        case final_ = 2
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(data)
        try container.encode(offset)
        try container.encodeIfPresent(self.final)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.data = try container.decode(NaclCiphertext.self)
        self.offset = try container.decode(Offset.self)
        self.final = try container.decodeIfPresent(UploadFinal.self)
    }

    init(data: NaclCiphertext, offset: Offset, final: UploadFinal? = nil) {
        self.data = data
        self.offset = offset
        self.final = final
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVExtendedDirent
// ---------------------------------------------------------------

struct KVExtendedDirent: SnowpackCodable, Equatable {
    var pos: UInt64
    var sfb: SmallFileBox

    private enum CodingKeys: Int, CodingKey {
        case pos = 0
        case sfb = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(pos)
        try container.encode(sfb)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.pos = try container.decode(UInt64.self)
        self.sfb = try container.decode(SmallFileBox.self)
    }

    init(pos: UInt64, sfb: SmallFileBox) {
        self.pos = pos
        self.sfb = sfb
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVDirentIDPair
// ---------------------------------------------------------------

struct KVDirentIDPair: SnowpackCodable, Equatable {
    var parentDirID: DirID
    var direntID: DirentID

    private enum CodingKeys: Int, CodingKey {
        case parentDirID = 0
        case direntID = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(parentDirID)
        try container.encode(direntID)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.parentDirID = try container.decode(DirID.self)
        self.direntID = try container.decode(DirentID.self)
    }

    init(parentDirID: DirID, direntID: DirentID) {
        self.parentDirID = parentDirID
        self.direntID = direntID
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVDirPair
// ---------------------------------------------------------------

struct KVDirPair: SnowpackCodable, Equatable {
    var active: KVDir
    var encrypting: KVDir?

    private enum CodingKeys: Int, CodingKey {
        case active = 0
        case encrypting = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(active)
        try container.encodeIfPresent(encrypting)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.active = try container.decode(KVDir.self)
        self.encrypting = try container.decodeIfPresent(KVDir.self)
    }

    init(active: KVDir, encrypting: KVDir? = nil) {
        self.active = active
        self.encrypting = encrypting
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVDirentBindingPayload
// ---------------------------------------------------------------

struct KVDirentBindingPayload: SnowpackCodable, Equatable, HasTypeUniqueID {
    var parentDir: DirID
    var id: DirentID
    var value: KVNodeID
    var version: KVVersion
    var dirVersion: KVVersion
    var writeRole: Role
    var nameMac: HMAC
    var nameBox: SecretBox

    static let typeUniqueID: TypeUniqueID = 0x9cc37c8363dc39fa

    private enum CodingKeys: Int, CodingKey {
        case parentDir = 0
        case id = 1
        case value = 2
        case version = 3
        case dirVersion = 4
        case writeRole = 5
        case nameMac = 6
        case nameBox = 7
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(parentDir)
        try container.encode(id)
        try container.encode(value)
        try container.encode(version)
        try container.encode(dirVersion)
        try container.encode(writeRole)
        try container.encode(nameMac)
        try container.encode(nameBox)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.parentDir = try container.decode(DirID.self)
        self.id = try container.decode(DirentID.self)
        self.value = try container.decode(KVNodeID.self)
        self.version = try container.decode(KVVersion.self)
        self.dirVersion = try container.decode(KVVersion.self)
        self.writeRole = try container.decode(Role.self)
        self.nameMac = try container.decode(HMAC.self)
        self.nameBox = try container.decode(SecretBox.self)
    }

    init(parentDir: DirID, id: DirentID, value: KVNodeID, version: KVVersion,
         dirVersion: KVVersion, writeRole: Role, nameMac: HMAC, nameBox: SecretBox) {
        self.parentDir = parentDir
        self.id = id
        self.value = value
        self.version = version
        self.dirVersion = dirVersion
        self.writeRole = writeRole
        self.nameMac = nameMac
        self.nameBox = nameBox
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct KVRootBindingPayload
// ---------------------------------------------------------------

struct KVRootBindingPayload: SnowpackCodable, Equatable, HasTypeUniqueID {
    var party: FQParty
    var rg: RoleAndGen
    var root: DirID
    var vers: KVVersion

    static let typeUniqueID: TypeUniqueID = 0xcfacdd4eab213a36

    private enum CodingKeys: Int, CodingKey {
        case party = 0
        case rg = 1
        case root = 2
        case vers = 3
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(party)
        try container.encode(rg)
        try container.encode(root)
        try container.encode(vers)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.party = try container.decode(FQParty.self)
        self.rg = try container.decode(RoleAndGen.self)
        self.root = try container.decode(DirID.self)
        self.vers = try container.decode(KVVersion.self)
    }

    init(party: FQParty, rg: RoleAndGen, root: DirID, vers: KVVersion) {
        self.party = party
        self.rg = rg
        self.root = root
        self.vers = vers
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  Unique ID registration
// ---------------------------------------------------------------

private let _registerKVTypes: Void = {
    UniqueIDRegistry.shared.register(FileKeySeed.typeUniqueID)
    UniqueIDRegistry.shared.register(DirKeySeed.typeUniqueID)
    UniqueIDRegistry.shared.register(KVNamePlaintext.typeUniqueID)
    UniqueIDRegistry.shared.register(KVPath.typeUniqueID)
    UniqueIDRegistry.shared.register(ChunkPlaintext.typeUniqueID)
    UniqueIDRegistry.shared.register(KVNameNonceInput.typeUniqueID)
    UniqueIDRegistry.shared.register(KVDirentBindingPayload.typeUniqueID)
    UniqueIDRegistry.shared.register(KVRootBindingPayload.typeUniqueID)
}()
