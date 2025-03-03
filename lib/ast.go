package lib

import (
	"fmt"
	"strings"
)

type UniqueID struct {
	Val string // is a number, but we don't bother to parse it
}

func (u UniqueID) IsSet() bool  { return u.Val != "" }
func (u UniqueID) IsZero() bool { return u.Val == "" }

type Root struct {
	Id    UniqueID
	Stmts []Statement
}

func (r *Root) DoInventory(i *Inventory) {
	for _, s := range r.Stmts {
		s.DoInventory(i)
	}
}

type Statement interface {
	DoInventory(i *Inventory)
	emit(e Emitter)
}

type BaseStatement struct {
	Dec Decorators
}

type BaseTypedef struct {
	BaseStatement
	Ident    Identifier
	UniqueID UniqueID
}

type Import struct {
	Path string
	Name string
	Lang Language
}

type Typedef struct {
	BaseTypedef
	Type Type
}

func (b BaseTypedef) DoInventory(i *Inventory) { b.UniqueID.DoInventory(i) }
func (t Typedef) DoInventory(i *Inventory)     { i.Typedef = true; t.BaseTypedef.DoInventory(i) }
func (s Struct) DoInventory(i *Inventory)      { i.Struct = true; s.BaseTypedef.DoInventory(i) }
func (v Variant) DoInventory(i *Inventory)     { i.Variant = true; v.BaseTypedef.DoInventory(i) }
func (b Import) DoInventory(i *Inventory)      { i.Import = true }
func (u UniqueID) DoInventory(i *Inventory) {
	if u.IsSet() {
		i.Unique = true
	}
}

type Decorators struct {
	Doc Docstring
}

type Docstring struct {
	Raw string
}

type Identifier struct {
	Name string
}

type Type interface {
	Emit(e Emitter)
	EmitInternal(e Emitter)
	EmitExport(e Emitter, nm string)
	EmitImport(e Emitter, nm string)
	EmitBytes(e Emitter, nm string)
	EmitFutureLink(e Emitter, child string)
	MakeOptional() Type
	IsPrimitiveType() bool
	IsList() bool
	DerivedPrefix() string
	EnumPrefix() string
	IsVoid() bool
}

type BaseType struct {
}

type Future struct {
	BaseType
	Type Type
}

type List struct {
	BaseType
	Type Type
}

type Blob struct {
	BaseType
	Count int
}

type DerivedType struct {
	BaseType
	Name         Identifier
	ImportedFrom Identifier
}

type Field struct {
	Ident Identifier
	Pos   int
	Type  Type
}

type Struct struct {
	BaseTypedef
	Fields []Field
}

type CaseLabel interface {
	CaseLabelToString(e Emitter, switchType Type) string
	GetterMethodName(e Emitter) string
	ConstructorName(e Emitter, swtch string) string
}

type CaseLabelIdentifier struct {
	Ident Identifier
}

func (c CaseLabelIdentifier) CaseLabelToString(e Emitter, switchType Type) string {
	return e.ToEnumConstant(switchType, c.Ident.Name)
}

func (c CaseLabelIdentifier) GetterMethodName(e Emitter) string {
	return e.GetterMethodNameForConstant(c.Ident.Name)
}

func (c CaseLabelIdentifier) ConstructorName(e Emitter, swtch string) string {
	return e.ConstructorNameForConstant(swtch, c.Ident.Name)
}

type CaseLabelNumber struct {
	Num int
}

func (c CaseLabelNumber) CaseLabelToString(e Emitter, switchType Type) string {
	return fmt.Sprintf("%d", c.Num)
}

func (c CaseLabelNumber) GetterMethodName(e Emitter) string {
	return e.GetterMethodNameForInt(c.Num)
}

func (c CaseLabelNumber) ConstructorName(e Emitter, swtch string) string {
	return e.ConstructorNameForInt(swtch, c.Num)
}

type CaseLabelBool struct {
	Bool bool
}

func (c CaseLabelBool) CaseLabelToString(e Emitter, switchType Type) string {
	if c.Bool {
		return "true"
	}
	return "false"
}

func (c CaseLabelBool) GetterMethodName(e Emitter) string {
	return e.GetterMethodNameForBool(c.Bool)
}

func (c CaseLabelBool) ConstructorName(e Emitter, swtch string) string {
	return e.ConstructorNameForBool(swtch, c.Bool)
}

type Case struct {
	Labels   []CaseLabel // nil for default case
	Position *int        // will be nil for void data; 0 is a valid position
	Type     Type
}

func (c Case) HasData() bool { return !c.Type.IsVoid() }

type Variant struct {
	BaseTypedef
	SwitchVar  Identifier
	SwitchType Type
	Cases      []Case
}

type Option struct {
	BaseType
	Type Type
}

type Void struct {
	BaseType
}

type EnumValue struct {
	Ident Identifier
	Num   int
}

type Enum struct {
	BaseTypedef
	Values []EnumValue
}

type Protocol struct {
	BaseTypedef
	Modifiers ProtocolModifiers
	Methods   []Method
}

type ProtocolModifier interface {
}

func (p Protocol) DoInventory(i *Inventory) {
	i.Rpc = true
	p.BaseTypedef.DoInventory(i)
}

func NewProtocolModifiers(pms []ProtocolModifier) (*ProtocolModifiers, error) {
	var errors *Errors
	var argHeader *ArgHeader
	var resHeader *ResHeader

	for _, pm := range pms {
		switch pm := pm.(type) {
		case Errors:
			if errors != nil {
				return nil, fmt.Errorf("multiple errors protocol modifiers found")
			}
			errors = &pm
		case ArgHeader:
			if argHeader != nil {
				return nil, fmt.Errorf("multiple arg_header protocol modifiers found")
			}
			argHeader = &pm
		case ResHeader:
			if resHeader != nil {
				return nil, fmt.Errorf("multiple res_header protocol modifiers found")
			}
			resHeader = &pm
		}
	}

	if errors == nil {
		return nil, fmt.Errorf("missing errors protocol modifier")
	}
	return &ProtocolModifiers{
		Errors:    *errors,
		ArgHeader: argHeader,
		ResHeader: resHeader,
	}, nil
}

type ProtocolModifiers struct {
	Errors    Errors
	ArgHeader *ArgHeader
	ResHeader *ResHeader
}

type Errors struct {
	Type Type
}

type ArgHeader struct {
	Type Type
}

type ResHeader struct {
	Type Type
}

type Param struct {
	Ident Identifier
	Type  Type
	Pos   int
}

func (p Param) ToField() Field {
	return Field{
		Ident: p.Ident,
		Pos:   p.Pos,
		Type:  p.Type,
	}
}

type Method struct {
	BaseTypedef
	Pos     int
	Params  []Param
	ArgType Identifier
	ResType Type
}

func (m Method) ParamsToStruct(n string) Struct {
	var fields []Field
	for _, p := range m.Params {
		fields = append(fields, p.ToField())
	}
	return Struct{
		BaseTypedef: BaseTypedef{
			Ident: Identifier{Name: n},
		},
		Fields: fields,
	}
}

type Text struct {
	BaseType
}
type Uint struct {
	BaseType
}
type Int struct {
	BaseType
}
type Bool struct {
	BaseType
}

var _ CaseLabel = CaseLabelIdentifier{}
var _ CaseLabel = CaseLabelNumber{}
var _ CaseLabel = CaseLabelBool{}
var _ Statement = Typedef{}
var _ Statement = Enum{}
var _ Statement = Variant{}
var _ Statement = Protocol{}
var _ Type = List{}
var _ Type = Future{}
var _ Type = Blob{}
var _ Type = Text{}
var _ Type = Uint{}
var _ Type = Int{}
var _ Type = Bool{}
var _ Type = Option{}
var _ Type = Void{}
var _ Type = DerivedType{}
var _ ProtocolModifier = Errors{}
var _ ProtocolModifier = ArgHeader{}
var _ ProtocolModifier = ResHeader{}

func (v Void) Emit(e Emitter)        { e.EmitVoid(v) }
func (l List) Emit(e Emitter)        { e.EmitList(l) }
func (f Future) Emit(e Emitter)      { e.EmitFuture(f) }
func (b Blob) Emit(e Emitter)        { e.EmitBlob(b) }
func (t Text) Emit(e Emitter)        { e.EmitText(t) }
func (u Uint) Emit(e Emitter)        { e.EmitUint(u) }
func (i Int) Emit(e Emitter)         { e.EmitInt(i) }
func (b Bool) Emit(e Emitter)        { e.EmitBool(b) }
func (o Option) Emit(e Emitter)      { e.EmitOption(o) }
func (d DerivedType) Emit(e Emitter) { e.EmitDerivedType(d) }

func (v Void) EmitInternal(e Emitter)        { e.EmitVoid(v) }
func (l List) EmitInternal(e Emitter)        { e.EmitListInternal(l) }
func (f Future) EmitInternal(e Emitter)      { e.EmitFuture(f) }
func (b Blob) EmitInternal(e Emitter)        { e.EmitBlob(b) }
func (t Text) EmitInternal(e Emitter)        { e.EmitText(t) }
func (u Uint) EmitInternal(e Emitter)        { e.EmitUint(u) }
func (i Int) EmitInternal(e Emitter)         { e.EmitInt(i) }
func (b Bool) EmitInternal(e Emitter)        { e.EmitBool(b) }
func (o Option) EmitInternal(e Emitter)      { e.EmitOptionInternal(o) }
func (d DerivedType) EmitInternal(e Emitter) { e.EmitDerivedTypeInternal(d) }

func (v Void) IsPrimitiveType() bool        { return false }
func (l List) IsPrimitiveType() bool        { return false }
func (f Future) IsPrimitiveType() bool      { return true }
func (b Blob) IsPrimitiveType() bool        { return true }
func (t Text) IsPrimitiveType() bool        { return true }
func (u Uint) IsPrimitiveType() bool        { return true }
func (i Int) IsPrimitiveType() bool         { return true }
func (b Bool) IsPrimitiveType() bool        { return true }
func (o Option) IsPrimitiveType() bool      { return false }
func (d DerivedType) IsPrimitiveType() bool { return false }

func (v Void) EmitExport(e Emitter, nm string)        { e.EmitVoid(v) }
func (l List) EmitExport(e Emitter, nm string)        { e.EmitExportList(l, nm) }
func (b Blob) EmitExport(e Emitter, nm string)        { e.EmitExportBlob(b, nm) }
func (t Text) EmitExport(e Emitter, nm string)        { e.EmitExportText(t, nm) }
func (u Uint) EmitExport(e Emitter, nm string)        { e.EmitExportUint(u, nm) }
func (i Int) EmitExport(e Emitter, nm string)         { e.EmitExportInt(i, nm) }
func (b Bool) EmitExport(e Emitter, nm string)        { e.EmitExportBool(b, nm) }
func (f Future) EmitExport(e Emitter, nm string)      { e.EmitExportFuture(f, nm) }
func (o Option) EmitExport(e Emitter, nm string)      { e.EmitExportOption(o, nm) }
func (d DerivedType) EmitExport(e Emitter, nm string) { e.EmitExportDerivedType(d, nm) }

func (v Void) EmitImport(e Emitter, nm string)        {}
func (l List) EmitImport(e Emitter, nm string)        { e.EmitImportList(l, nm) }
func (b Blob) EmitImport(e Emitter, nm string)        { e.EmitImportBlob(b, nm) }
func (t Text) EmitImport(e Emitter, nm string)        { e.EmitImportText(t, nm) }
func (u Uint) EmitImport(e Emitter, nm string)        { e.EmitImportUint(u, nm) }
func (i Int) EmitImport(e Emitter, nm string)         { e.EmitImportInt(i, nm) }
func (b Bool) EmitImport(e Emitter, nm string)        { e.EmitImportBool(b, nm) }
func (f Future) EmitImport(e Emitter, nm string)      { e.EmitImportFuture(f, nm) }
func (o Option) EmitImport(e Emitter, nm string)      { e.EmitImportOption(o, nm) }
func (d DerivedType) EmitImport(e Emitter, nm string) { e.EmitImportDerivedType(d, nm) }

func (v Void) EmitBytes(e Emitter, nm string)        { e.EmitNil() }
func (l List) EmitBytes(e Emitter, nm string)        { e.EmitNil() }
func (f Future) EmitBytes(e Emitter, nm string)      { e.EmitBlobToBytes(nm) }
func (b Blob) EmitBytes(e Emitter, nm string)        { e.EmitBlobToBytes(nm) }
func (t Text) EmitBytes(e Emitter, nm string)        { e.EmitNil() }
func (u Uint) EmitBytes(e Emitter, nm string)        { e.EmitNil() }
func (i Int) EmitBytes(e Emitter, nm string)         { e.EmitNil() }
func (b Bool) EmitBytes(e Emitter, nm string)        { e.EmitNil() }
func (o Option) EmitBytes(e Emitter, nm string)      { e.EmitNil() }
func (d DerivedType) EmitBytes(e Emitter, nm string) { e.EmitBytesDowncast(d.FullTypeName(), nm) }

func (v Void) MakeOptional() Type        { return Option{Type: v} }
func (l List) MakeOptional() Type        { return Option{Type: l} }
func (f Future) MakeOptional() Type      { return Option{Type: f} }
func (b Blob) MakeOptional() Type        { return Option{Type: b} }
func (t Text) MakeOptional() Type        { return Option{Type: t} }
func (u Uint) MakeOptional() Type        { return Option{Type: u} }
func (i Int) MakeOptional() Type         { return Option{Type: i} }
func (b Bool) MakeOptional() Type        { return Option{Type: b} }
func (o Option) MakeOptional() Type      { return o }
func (d DerivedType) MakeOptional() Type { return Option{Type: d} }

func (b BaseType) EmitFutureLink(e Emitter, child string) {}
func (b BaseType) DerivedPrefix() string                  { return "" }
func (b BaseType) EnumPrefix() string                     { return "" }
func (b BaseType) IsVoid() bool                           { return false }
func (b BaseType) IsList() bool                           { return false }

func (v Void) IsVoid() bool { return true }
func (l List) IsList() bool { return true }

func (d DerivedType) DerivedPrefix() string {
	if d.ImportedFrom.Name != "" {
		return d.ImportedFrom.Name + "."
	}
	return ""
}

func (d DerivedType) EnumPrefix() string {
	return d.Name.Name
}

func (f Future) EmitFutureLink(e Emitter, child string) {
	e.EmitFutureLink(f.Type, child)
}

func (d DerivedType) FullTypeName() string {
	var parts []string
	if d.ImportedFrom.Name != "" {
		parts = append(parts, d.ImportedFrom.Name)
	}
	parts = append(parts, d.Name.Name)
	return strings.Join(parts, ".")
}

func (m Method) makeArgName(g Emitter) string {
	return g.MethodArgName(m.Ident.Name, m.ArgType.Name)
}

func (m Method) singleArg() bool {
	return len(m.Params) == 1 && m.Params[0].Pos == 0
}
