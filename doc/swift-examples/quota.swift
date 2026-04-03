// Auto-generated to Swift types and interfaces using snowpc (https://github.com/foks-proj/go-snowpack-compiler)
//  Input file: ../../proto-src/infra/quota.snowp

import Foundation
import SnowpackRPC
import Lib

// ---------------------------------------------------------------
//  struct QuotaConfig
// ---------------------------------------------------------------

struct QuotaConfig: SnowpackCodable, Equatable {
    var slacks: Slacks
    var delay: Lib.DurationMilli
    var noPlanMaxTeams: Int64
    var noResurrection: Bool

    private enum CodingKeys: Int, CodingKey {
        case slacks = 0
        case delay = 1
        case noPlanMaxTeams = 2
        case noResurrection = 3
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(slacks)
        try container.encode(delay)
        try container.encode(noPlanMaxTeams)
        try container.encode(noResurrection)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.slacks = try container.decode(Slacks.self)
        self.delay = try container.decode(Lib.DurationMilli.self)
        self.noPlanMaxTeams = try container.decode(Int64.self)
        self.noResurrection = try container.decode(Bool.self)
    }

    init(slacks: Slacks, delay: Lib.DurationMilli, noPlanMaxTeams: Int64,
         noResurrection: Bool) {
        self.slacks = slacks
        self.delay = delay
        self.noPlanMaxTeams = noPlanMaxTeams
        self.noResurrection = noResurrection
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  typedefs
// ---------------------------------------------------------------

typealias Cents = UInt64
typealias SignedCents = Int64
typealias StripeProdID = String
typealias StripePriceID = String
typealias StripeSessionID = String
typealias StripeCustomerID = String
typealias StripeSubscriptionID = String
typealias StripeInvoiceID = String
typealias StripeChargeID = String
typealias StripeEventID = String

// ---------------------------------------------------------------
//  enum QuotaScope
// ---------------------------------------------------------------

enum QuotaScope: Int, Codable, CaseIterable, SnowpackEnum {
    case none = 0
    case teams = 1
    case vHost = 2

    var label: String {
        switch self {
        case .none:  return "None"
        case .teams: return "Teams"
        case .vHost: return "VHost"
        }
    }

    init?(label: String) {
        switch label {
        case "None":  self = .none
        case "Teams": self = .teams
        case "VHost": self = .vHost
        default: return nil
        }
    }
}

// ---------------------------------------------------------------
//  struct Plan
// ---------------------------------------------------------------

struct Plan: SnowpackCodable, Equatable {
    var id: Lib.PlanID
    var name: String
    var maxSeats: UInt64
    var quota: Lib.Size
    var displayName: String
    var stripeProdId: StripeProdID
    var points: [String]
    var prices: [PlanPrice]
    var promoted: Bool
    var scope: QuotaScope
    var maxVhosts: UInt64
    var sso: Bool

    private enum CodingKeys: Int, CodingKey {
        case id = 0
        case name = 1
        case maxSeats = 2
        case quota = 3
        case displayName = 4
        case stripeProdId = 5
        case points = 6
        case prices = 7
        case promoted = 8
        case scope = 9
        case maxVhosts = 10
        case sso = 11
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(id)
        try container.encode(name)
        try container.encode(maxSeats)
        try container.encode(quota)
        try container.encode(displayName)
        try container.encode(stripeProdId)
        try container.encode(points)
        try container.encode(prices)
        try container.encode(promoted)
        try container.encode(scope)
        try container.encode(maxVhosts)
        try container.encode(sso)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.id = try container.decode(Lib.PlanID.self)
        self.name = try container.decode(String.self)
        self.maxSeats = try container.decode(UInt64.self)
        self.quota = try container.decode(Lib.Size.self)
        self.displayName = try container.decode(String.self)
        self.stripeProdId = try container.decode(StripeProdID.self)
        self.points = try container.decode([String].self)
        self.prices = try container.decode([PlanPrice].self)
        self.promoted = try container.decode(Bool.self)
        self.scope = try container.decode(QuotaScope.self)
        self.maxVhosts = try container.decode(UInt64.self)
        self.sso = try container.decode(Bool.self)
    }

    init(id: Lib.PlanID, name: String, maxSeats: UInt64, quota: Lib.Size,
         displayName: String, stripeProdId: StripeProdID, points: [String],
         prices: [PlanPrice], promoted: Bool, scope: QuotaScope,
         maxVhosts: UInt64, sso: Bool) {
        self.id = id
        self.name = name
        self.maxSeats = maxSeats
        self.quota = quota
        self.displayName = displayName
        self.stripeProdId = stripeProdId
        self.points = points
        self.prices = prices
        self.promoted = promoted
        self.scope = scope
        self.maxVhosts = maxVhosts
        self.sso = sso
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct PaymentInterval
// ---------------------------------------------------------------

struct PaymentInterval: SnowpackCodable, Equatable {
    var interval: Interval
    var count: UInt64

    private enum CodingKeys: Int, CodingKey {
        case interval = 0
        case count = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(interval)
        try container.encode(count)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.interval = try container.decode(Interval.self)
        self.count = try container.decode(UInt64.self)
    }

    init(interval: Interval, count: UInt64) {
        self.interval = interval
        self.count = count
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  struct PlanPrice
// ---------------------------------------------------------------

struct PlanPrice: SnowpackCodable, Equatable {
    var id: Lib.PriceID
    var stripePriceId: StripePriceID
    var cents: Cents
    var pi: PaymentInterval
    var promoted: Bool
    var pri: Int64

    private enum CodingKeys: Int, CodingKey {
        case id = 0
        case stripePriceId = 1
        case cents = 2
        case pi = 3
        case promoted = 4
        case pri = 5
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(id)
        try container.encode(stripePriceId)
        try container.encode(cents)
        try container.encode(pi)
        try container.encode(promoted)
        try container.encode(pri)
    }

    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.id = try container.decode(Lib.PriceID.self)
        self.stripePriceId = try container.decode(StripePriceID.self)
        self.cents = try container.decode(Cents.self)
        self.pi = try container.decode(PaymentInterval.self)
        self.promoted = try container.decode(Bool.self)
        self.pri = try container.decode(Int64.self)
    }

    init(id: Lib.PriceID, stripePriceId: StripePriceID, cents: Cents,
         pi: PaymentInterval, promoted: Bool, pri: Int64) {
        self.id = id
        self.stripePriceId = stripePriceId
        self.cents = cents
        self.pi = pi
        self.promoted = promoted
        self.pri = pri
    }

    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  enum MakePlanOpts (placeholder — defined elsewhere in real code)
// ---------------------------------------------------------------

// MakePlanOpts, Slacks, Interval, etc. would be generated from
// their respective .snowp definitions. Omitted here for brevity.

// ===============================================================
//  Protocol: Quota
// ===============================================================

let QuotaProtocolID: ProtocolUniqueID = 0xfa82ee9c

// ---------------------------------------------------------------
//  Server interface
// ---------------------------------------------------------------

protocol QuotaInterface {
    func poke() async throws
    func testSetConfig(config: QuotaConfig) async throws
    func testUnsetConfig() async throws
    func testBumpUsage(arg: TestBumpUsageArg) async throws
    func makePlan(arg: MakePlanArg) async throws -> Plan
    func setPlan(arg: SetPlanArg) async throws -> Lib.CancelID
    func cancelPlan(fqu: Lib.FQUser) async throws -> Lib.CancelID
    func assignQuotaMaster(arg: AssignQuotaMasterArg) async throws
    func unassignQuotaMaster(arg: UnassignQuotaMasterArg) async throws
    func errorWrapper() -> (Error) -> Lib.Status
}

// ---------------------------------------------------------------
//  Arg structs
// ---------------------------------------------------------------

struct QuotaPokeArg: SnowpackCodable, Equatable {
    func encode(to encoder: Encoder) throws {
        _ = encoder.unkeyedContainer()
    }
    init(from decoder: Decoder) throws {
        _ = try decoder.unkeyedContainer()
    }
    init() {}
    func bytes() -> Data? { return nil }
}

struct TestSetConfigArg: SnowpackCodable, Equatable {
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

struct TestUnsetConfigArg: SnowpackCodable, Equatable {
    func encode(to encoder: Encoder) throws {
        _ = encoder.unkeyedContainer()
    }
    init(from decoder: Decoder) throws {
        _ = try decoder.unkeyedContainer()
    }
    init() {}
    func bytes() -> Data? { return nil }
}

struct TestBumpUsageArg: SnowpackCodable, Equatable {
    var hid: Lib.HostID
    var pid: Lib.PartyID
    var amt: Lib.Size

    private enum CodingKeys: Int, CodingKey {
        case hid = 0
        case pid = 1
        case amt = 2
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(hid)
        try container.encode(pid)
        try container.encode(amt)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.hid = try container.decode(Lib.HostID.self)
        self.pid = try container.decode(Lib.PartyID.self)
        self.amt = try container.decode(Lib.Size.self)
    }
    init(hid: Lib.HostID, pid: Lib.PartyID, amt: Lib.Size) {
        self.hid = hid
        self.pid = pid
        self.amt = amt
    }
    func bytes() -> Data? { return nil }
}

struct MakePlanArg: SnowpackCodable, Equatable {
    var plan: Plan
    var opts: MakePlanOpts

    private enum CodingKeys: Int, CodingKey {
        case plan = 0
        case opts = 1
    }

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

struct SetPlanArg: SnowpackCodable, Equatable {
    var fqu: Lib.FQUser
    var plan: Lib.PlanID
    var price: Lib.PriceID
    var replace: Bool
    var validFor: Lib.DurationSecs
    var stripeSubId: StripeSubscriptionID

    private enum CodingKeys: Int, CodingKey {
        case fqu = 0
        case plan = 1
        case price = 2
        case replace = 3
        case validFor = 4
        case stripeSubId = 5
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(fqu)
        try container.encode(plan)
        try container.encode(price)
        try container.encode(replace)
        try container.encode(validFor)
        try container.encode(stripeSubId)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.fqu = try container.decode(Lib.FQUser.self)
        self.plan = try container.decode(Lib.PlanID.self)
        self.price = try container.decode(Lib.PriceID.self)
        self.replace = try container.decode(Bool.self)
        self.validFor = try container.decode(Lib.DurationSecs.self)
        self.stripeSubId = try container.decode(StripeSubscriptionID.self)
    }
    init(fqu: Lib.FQUser, plan: Lib.PlanID, price: Lib.PriceID,
         replace: Bool, validFor: Lib.DurationSecs, stripeSubId: StripeSubscriptionID) {
        self.fqu = fqu
        self.plan = plan
        self.price = price
        self.replace = replace
        self.validFor = validFor
        self.stripeSubId = stripeSubId
    }
    func bytes() -> Data? { return nil }
}

struct CancelPlanArg: SnowpackCodable, Equatable {
    var fqu: Lib.FQUser

    private enum CodingKeys: Int, CodingKey { case fqu = 0 }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(fqu)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.fqu = try container.decode(Lib.FQUser.self)
    }
    init(fqu: Lib.FQUser) { self.fqu = fqu }
    func bytes() -> Data? { return nil }
}

struct AssignQuotaMasterArg: SnowpackCodable, Equatable {
    var fqu: Lib.FQUser
    var team: Lib.TeamID

    private enum CodingKeys: Int, CodingKey {
        case fqu = 0
        case team = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(fqu)
        try container.encode(team)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.fqu = try container.decode(Lib.FQUser.self)
        self.team = try container.decode(Lib.TeamID.self)
    }
    init(fqu: Lib.FQUser, team: Lib.TeamID) {
        self.fqu = fqu
        self.team = team
    }
    func bytes() -> Data? { return nil }
}

struct UnassignQuotaMasterArg: SnowpackCodable, Equatable {
    var fqu: Lib.FQUser
    var team: Lib.TeamID

    private enum CodingKeys: Int, CodingKey {
        case fqu = 0
        case team = 1
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.unkeyedContainer()
        try container.encode(fqu)
        try container.encode(team)
    }
    init(from decoder: Decoder) throws {
        var container = try decoder.unkeyedContainer()
        self.fqu = try container.decode(Lib.FQUser.self)
        self.team = try container.decode(Lib.TeamID.self)
    }
    init(fqu: Lib.FQUser, team: Lib.TeamID) {
        self.fqu = fqu
        self.team = team
    }
    func bytes() -> Data? { return nil }
}

// ---------------------------------------------------------------
//  Client stub
// ---------------------------------------------------------------

struct QuotaClient {
    let cli: GenericClient
    let errorUnwrapper: (Lib.Status) -> Error?

    func poke() async throws {
        let arg = QuotaPokeArg()
        try await cli.call(
            method: Method(QuotaProtocolID, position: 0, name: "Quota.poke"),
            arg: arg, res: Void.self)
    }

    func testSetConfig(config: QuotaConfig) async throws {
        let arg = TestSetConfigArg(config: config)
        try await cli.call(
            method: Method(QuotaProtocolID, position: 1, name: "Quota.testSetConfig"),
            arg: arg, res: Void.self)
    }

    func testUnsetConfig() async throws {
        let arg = TestUnsetConfigArg()
        try await cli.call(
            method: Method(QuotaProtocolID, position: 2, name: "Quota.testUnsetConfig"),
            arg: arg, res: Void.self)
    }

    func testBumpUsage(hid: Lib.HostID, pid: Lib.PartyID, amt: Lib.Size) async throws {
        let arg = TestBumpUsageArg(hid: hid, pid: pid, amt: amt)
        try await cli.call(
            method: Method(QuotaProtocolID, position: 3, name: "Quota.testBumpUsage"),
            arg: arg, res: Void.self)
    }

    func makePlan(plan: Plan, opts: MakePlanOpts) async throws -> Plan {
        let arg = MakePlanArg(plan: plan, opts: opts)
        return try await cli.call(
            method: Method(QuotaProtocolID, position: 4, name: "Quota.makePlan"),
            arg: arg, res: Plan.self)
    }

    func setPlan(fqu: Lib.FQUser, plan: Lib.PlanID, price: Lib.PriceID,
                 replace: Bool, validFor: Lib.DurationSecs,
                 stripeSubId: StripeSubscriptionID) async throws -> Lib.CancelID {
        let arg = SetPlanArg(fqu: fqu, plan: plan, price: price,
                             replace: replace, validFor: validFor, stripeSubId: stripeSubId)
        return try await cli.call(
            method: Method(QuotaProtocolID, position: 5, name: "Quota.setPlan"),
            arg: arg, res: Lib.CancelID.self)
    }

    func cancelPlan(fqu: Lib.FQUser) async throws -> Lib.CancelID {
        let arg = CancelPlanArg(fqu: fqu)
        return try await cli.call(
            method: Method(QuotaProtocolID, position: 6, name: "Quota.cancelPlan"),
            arg: arg, res: Lib.CancelID.self)
    }

    func assignQuotaMaster(fqu: Lib.FQUser, team: Lib.TeamID) async throws {
        let arg = AssignQuotaMasterArg(fqu: fqu, team: team)
        try await cli.call(
            method: Method(QuotaProtocolID, position: 7, name: "Quota.assignQuotaMaster"),
            arg: arg, res: Void.self)
    }

    func unassignQuotaMaster(fqu: Lib.FQUser, team: Lib.TeamID) async throws {
        let arg = UnassignQuotaMasterArg(fqu: fqu, team: team)
        try await cli.call(
            method: Method(QuotaProtocolID, position: 8, name: "Quota.unassignQuotaMaster"),
            arg: arg, res: Void.self)
    }
}

// ---------------------------------------------------------------
//  Server-side dispatch
// ---------------------------------------------------------------

func quotaProtocol(_ impl: QuotaInterface) -> ProtocolV2 {
    let wrapError = impl.errorWrapper()
    return ProtocolV2(
        name: "Quota",
        id: QuotaProtocolID,
        methods: [
            0: ServeHandlerDescriptionV2(
                name: "poke",
                handler: { decoder in
                    _ = try QuotaPokeArg(from: decoder)
                    try await impl.poke()
                    return nil
                }
            ),
            1: ServeHandlerDescriptionV2(
                name: "testSetConfig",
                handler: { decoder in
                    let arg = try TestSetConfigArg(from: decoder)
                    try await impl.testSetConfig(config: arg.config)
                    return nil
                }
            ),
            2: ServeHandlerDescriptionV2(
                name: "testUnsetConfig",
                handler: { decoder in
                    _ = try TestUnsetConfigArg(from: decoder)
                    try await impl.testUnsetConfig()
                    return nil
                }
            ),
            3: ServeHandlerDescriptionV2(
                name: "testBumpUsage",
                handler: { decoder in
                    let arg = try TestBumpUsageArg(from: decoder)
                    try await impl.testBumpUsage(arg: arg)
                    return nil
                }
            ),
            4: ServeHandlerDescriptionV2(
                name: "makePlan",
                handler: { decoder in
                    let arg = try MakePlanArg(from: decoder)
                    let res = try await impl.makePlan(arg: arg)
                    return res
                }
            ),
            5: ServeHandlerDescriptionV2(
                name: "setPlan",
                handler: { decoder in
                    let arg = try SetPlanArg(from: decoder)
                    let res = try await impl.setPlan(arg: arg)
                    return res
                }
            ),
            6: ServeHandlerDescriptionV2(
                name: "cancelPlan",
                handler: { decoder in
                    let arg = try CancelPlanArg(from: decoder)
                    let res = try await impl.cancelPlan(fqu: arg.fqu)
                    return res
                }
            ),
            7: ServeHandlerDescriptionV2(
                name: "assignQuotaMaster",
                handler: { decoder in
                    let arg = try AssignQuotaMasterArg(from: decoder)
                    try await impl.assignQuotaMaster(arg: arg)
                    return nil
                }
            ),
            8: ServeHandlerDescriptionV2(
                name: "unassignQuotaMaster",
                handler: { decoder in
                    let arg = try UnassignQuotaMasterArg(from: decoder)
                    try await impl.unassignQuotaMaster(arg: arg)
                    return nil
                }
            ),
        ],
        wrapError: { error in
            return wrapError(error)
        }
    )
}

// ---------------------------------------------------------------
//  Unique ID registration
// ---------------------------------------------------------------

private let _registerQuotaTypes: Void = {
    UniqueIDRegistry.shared.register(QuotaProtocolID)
}()
