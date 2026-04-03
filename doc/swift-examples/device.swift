// Auto-generated to Swift types and interfaces using snowpc (https://github.com/foks-proj/go-snowpack-compiler)
//  Input file: ../../proto-src/lib/device.snowp

import Foundation
import SnowpackRPC

// ---------------------------------------------------------------
//  enum DeviceStatus
// ---------------------------------------------------------------

enum DeviceStatus: Int, Codable, CaseIterable, SnowpackEnum {
    case active = 0
    case revoked = 1

    var label: String {
        switch self {
        case .active:  return "ACTIVE"
        case .revoked: return "REVOKED"
        }
    }

    init?(label: String) {
        switch label {
        case "ACTIVE":  self = .active
        case "REVOKED": self = .revoked
        default: return nil
        }
    }
}

// ---------------------------------------------------------------
//  enum NormalizationVersion
// ---------------------------------------------------------------

enum NormalizationVersion: Int, Codable, CaseIterable, SnowpackEnum {
    case v0 = 0

    var label: String {
        switch self {
        case .v0: return "V0"
        }
    }

    init?(label: String) {
        switch label {
        case "V0": self = .v0
        default: return nil
        }
    }
}

// ---------------------------------------------------------------
//  typedefs
// ---------------------------------------------------------------

typealias DeviceNameNormalized = String
typealias DeviceName = String
typealias DeviceSerial = UInt64

// ---------------------------------------------------------------
//  struct DeviceLabel
// ---------------------------------------------------------------

struct DeviceLabel: SnowpackCodable, Equatable, HasTypeUniqueID {
    var deviceType: DeviceType
    var name: DeviceNameNormalized
    var serial: DeviceSerial

    static let typeUniqueID: TypeUniqueID = 0x9650227205486122

    private enum CodingKeys: Int, CodingKey {
        case deviceType = 0
        case name = 1
        case serial = 2
    }

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

    init(deviceType: DeviceType, name: DeviceNameNormalized, serial: DeviceSerial) {
        self.deviceType = deviceType
        self.name = name
        self.serial = serial
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct DeviceNameNormalizationPreimage
// ---------------------------------------------------------------

struct DeviceNameNormalizationPreimage: SnowpackCodable, Equatable {
    var nv: NormalizationVersion
    var name: DeviceName

    private enum CodingKeys: Int, CodingKey {
        case nv = 0
        case name = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(nv)
        try container.encode(name)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.nv = try container.decode(NormalizationVersion.self)
        self.name = try container.decode(DeviceName.self)
    }

    init(nv: NormalizationVersion, name: DeviceName) {
        self.nv = nv
        self.name = name
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct DeviceLabelAndName
// ---------------------------------------------------------------

struct DeviceLabelAndName: SnowpackCodable, Equatable {
    var label: DeviceLabel
    var nv: NormalizationVersion
    var name: DeviceName

    private enum CodingKeys: Int, CodingKey {
        case label = 0
        case nv = 1
        case name = 2
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(label)
        try container.encode(nv)
        try container.encode(name)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.label = try container.decode(DeviceLabel.self)
        self.nv = try container.decode(NormalizationVersion.self)
        self.name = try container.decode(DeviceName.self)
    }

    init(label: DeviceLabel, nv: NormalizationVersion, name: DeviceName) {
        self.label = label
        self.nv = nv
        self.name = name
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct ProvisionInfo
// ---------------------------------------------------------------

struct ProvisionInfo: SnowpackCodable, Equatable {
    var signer: EntityID
    var chain: BaseChainer
    var leaf: MerkleLeaf

    private enum CodingKeys: Int, CodingKey {
        case signer = 0
        case chain = 1
        case leaf = 2
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(signer)
        try container.encode(chain)
        try container.encode(leaf)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.signer = try container.decode(EntityID.self)
        self.chain = try container.decode(BaseChainer.self)
        self.leaf = try container.decode(MerkleLeaf.self)
    }

    init(signer: EntityID, chain: BaseChainer, leaf: MerkleLeaf) {
        self.signer = signer
        self.chain = chain
        self.leaf = leaf
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct DeviceInfo
// ---------------------------------------------------------------

struct DeviceInfo: SnowpackCodable, Equatable {
    var status: DeviceStatus
    var dn: DeviceLabelAndName?
    var key: MemberRole
    var ctime: Time
    var provisioned: ProvisionInfo
    var revoked: RevokeInfo?

    private enum CodingKeys: Int, CodingKey {
        case status = 0
        case dn = 1
        case key = 2
        case ctime = 3
        case provisioned = 4
        case revoked = 5
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(status)
        try container.encodeIfPresent(dn)
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

// ---------------------------------------------------------------
//  struct DeviceNagInfo
// ---------------------------------------------------------------

struct DeviceNagInfo: SnowpackCodable, Equatable {
    var numDevices: UInt64
    var cleared: Bool

    private enum CodingKeys: Int, CodingKey {
        case numDevices = 0
        case cleared = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(numDevices)
        try container.encode(cleared)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.numDevices = try container.decode(UInt64.self)
        self.cleared = try container.decode(Bool.self)
    }

    init(numDevices: UInt64, cleared: Bool) {
        self.numDevices = numDevices
        self.cleared = cleared
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  Unique ID registration
// ---------------------------------------------------------------

private let _registerDeviceTypes: Void = {
    UniqueIDRegistry.shared.register(DeviceLabel.typeUniqueID)
}()
