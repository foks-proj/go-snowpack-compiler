package lib

import (
	"fmt"
	"io"
	"strings"
)

const version = "0.0.3"
const name = "snowpc"
const url = "https://github.com/foks-proj/go-snowpack-compiler"

type Inventory struct {
	Rpc     bool
	Variant bool
	Struct  bool
	Typedef bool
	Unique  bool
	Import  bool
}

func (i *Inventory) imports() []string {
	var ret []string
	if i.Rpc || i.Variant {
		ret = append(ret, "errors")
	}
	if i.Variant {
		ret = append(ret, "fmt")
	}
	if i.Rpc {
		ret = append(ret, "context")
		ret = append(ret, "time")
	}

	if i.Rpc || i.Unique || i.Struct || i.Variant || i.Typedef {
		ret = append(ret, "github.com/foks-proj/go-snowpack-rpc/rpc")
	}
	return ret
}

type Emitter interface {
	EmitEnum(e Enum)
	EmitTypedef(t Typedef)
	EmitVoid(v Void)
	EmitList(l List)
	EmitOption(o Option)
	EmitFuture(f Future)
	EmitBlob(b Blob)
	EmitText(t Text)
	EmitUint(u Uint)
	EmitInt(i Int)
	EmitBool(b Bool)
	EmitProtocol(p Protocol)
	EmitListInternal(l List)
	EmitOptionInternal(o Option)
	EmitStruct(s Struct)
	EmitVariant(v Variant)
	EmitImport(i Import)
	EmitDerivedType(d DerivedType)
	EmitExportList(l List, param string)
	EmitExportBlob(b Blob, param string)
	EmitExportText(t Text, param string)
	EmitExportUint(u Uint, param string)
	EmitExportInt(i Int, param string)
	EmitExportBool(b Bool, param string)
	EmitExportFuture(f Future, param string)
	EmitExportOption(o Option, param string)
	EmitDerivedTypeInternal(d DerivedType)
	EmitExportDerivedType(d DerivedType, param string)
	EmitImportList(l List, param string)
	EmitImportBlob(b Blob, param string)
	EmitImportText(t Text, param string)
	EmitImportUint(u Uint, param string)
	EmitImportInt(i Int, param string)
	EmitImportBool(b Bool, param string)
	EmitImportFuture(f Future, param string)
	EmitImportOption(o Option, param string)
	EmitImportDerivedType(d DerivedType, param string)
	EmitBytesDowncast(nm string, param string)
	EmitNil()
	EmitBlobToBytes(nm string)
	EmitFutureLink(t Type, child string)
	ToEnumConstant(t Type, nm string) string
	GetterMethodNameForBool(b bool) string
	GetterMethodNameForInt(i int) string
	GetterMethodNameForConstant(s string) string
	ConstructorNameForConstant(vrnt string, cnst string) string
	ConstructorNameForInt(vrnt string, i int) string
	ConstructorNameForBool(vrnt string, b bool) string
	MethodArgName(mthd string, arnName string) string
}

type ImportFlavors struct {
	m map[Language]Import
}

func NewImportFlavors() *ImportFlavors {
	return &ImportFlavors{
		m: make(map[Language]Import),
	}
}

type BaseEmitter struct {
	uniques   []string
	md        *Metadata
	dst       io.Writer
	nTabs     int
	isNewline bool
	imports   map[string]*ImportFlavors
}

func (b *BaseEmitter) storeImport(i Import) {
	flav := b.imports[i.Name]
	if flav == nil {
		flav = NewImportFlavors()
		b.imports[i.Name] = flav
	}
	flav.m[i.Lang] = i
}

func NewBaseEmitter(m *Metadata, dst io.Writer) *BaseEmitter {
	return &BaseEmitter{
		imports:   make(map[string]*ImportFlavors),
		md:        m,
		dst:       dst,
		isNewline: true,
	}
}

func (g *BaseEmitter) addUnique(s string) {
	g.uniques = append(g.uniques, s)
}

func (g *BaseEmitter) outputString(s string) {
	n, err := g.dst.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	if n != len(s) {
		panic(fmt.Sprintf("short io.Writer write: wrote %d bytes, expected %d", n, len(s)))
	}
}

func (g *BaseEmitter) output(
	s string,
	isFrag bool,
) {
	if g.isNewline {
		g.outputString(g.tabs())
		g.isNewline = false
	}
	g.outputString(s)
	if !isFrag {
		g.outputString("\n")
		g.isNewline = true
	} else if len(s) > 0 {
		g.isNewline = false
	}
}

func (g *BaseEmitter) foutputLine(s string, args ...any) {
	g.output(fmt.Sprintf(s, args...), false)
}

func (g *BaseEmitter) foutputFrag(s string, args ...any) {
	g.output(fmt.Sprintf(s, args...), true)
}

func (g *BaseEmitter) outputFrag(s string) { g.output(s, true) }
func (g *BaseEmitter) emptyLine()          { g.outputLine("") }
func (g *BaseEmitter) outputLine(s string) { g.output(s, false) }
func (g *BaseEmitter) tabs() string        { return strings.Repeat("\t", g.nTabs) }
func (g *BaseEmitter) tab()                { g.nTabs++ }

func (g *BaseEmitter) untab() {
	g.nTabs--
	if g.nTabs < 0 {
		panic("untab() called too many times")
	}
}

type GoEmitter struct {
	*BaseEmitter
}

func NewGoEmitter(m *Metadata, dst io.Writer) *GoEmitter {
	return &GoEmitter{
		BaseEmitter: NewBaseEmitter(m, dst),
	}
}

func (g *GoEmitter) emitPreamble(r *Root) {
	g.outputLine(
		fmt.Sprintf("%s %s %s (%s)",
			`// Auto-generated to Go types and interfaces using`,
			name, version, url,
		),
	)
	inv := &Inventory{}
	r.DoInventory(inv)

	g.outputLine(`//  Input file:` + g.md.infile.Name())
	g.emptyLine()
	g.outputLine(`package ` + g.md.pkg)
	g.emptyLine()
	g.outputLine(`import (`)
	g.tab()
	imps := inv.imports()
	for _, imp := range imps {
		g.outputLine("\"" + imp + "\"")
	}
	g.untab()
	g.outputLine(`)`)
	g.emptyLine()
}

func (g *GoEmitter) emitPostamble(r *Root) {
	if len(g.uniques) == 0 {
		return
	}
	g.emptyLine()
	g.outputLine("func init() {")
	g.tab()
	for _, u := range g.uniques {
		g.foutputLine("rpc.AddUnique(%s)", u)
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) Emit(r *Root) {
	g.emitPreamble(r)
	r.emit(g)
	g.emitPostamble(r)
}

// emitDoc will work for any target language that has //-style comments
// This isn't exact but it's close enough for now. Feel free to revisit
// when it's wrong.
func (g *BaseEmitter) emitDoc(d Docstring) {
	s := d.Raw
	if len(s) == 0 {
		return
	}
	isEmpty := func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	}
	lines := strings.Split(s, "\n")

	// Remove leading and trailing empty lines
	for len(lines) > 0 && isEmpty(lines[0]) {
		lines = lines[1:]
	}
	for len(lines) > 0 && isEmpty(lines[len(lines)-1]) {
		lines = lines[:len(lines)-1]
	}
	for _, line := range lines {
		g.outputLine("// " + line)
	}

}

func (g *BaseEmitter) emitDecorators(d Decorators) {
	g.emitDoc(d.Doc)
}

func (g *BaseEmitter) emitStatementPremable(s BaseStatement) {
	g.emitDecorators(s.Dec)
}

func (g *GoEmitter) emitEnumConstants(e Enum) {
	exsym := g.exportSymbol(e.Ident.Name)
	g.outputLine("const (")
	g.tab()
	for _, v := range e.Values {
		g.foutputLine("%s_%s %s = %d", exsym, v.Ident.Name, exsym, v.Num)
	}
	g.untab()
	g.outputLine(")")
}

func (g *GoEmitter) emitEnumMap(e Enum) {
	exsym := g.exportSymbol(e.Ident.Name)
	g.outputLine("var " + exsym + "Map = map[string]" + exsym + "{")
	g.tab()
	for _, v := range e.Values {
		g.foutputLine("\"%s\" : %d,", v.Ident.Name, v.Num)
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitEnumRevMap(e Enum) {
	exsym := g.exportSymbol(e.Ident.Name)
	g.outputLine("var " + exsym + "RevMap = map[" + exsym + "]string{")
	g.tab()
	for _, v := range e.Values {
		g.foutputLine("%d : \"%s\",", v.Num, v.Ident.Name)
	}
	g.untab()
	g.outputLine("}")
}

// Fixme for generic- or typescript-style imports
func (g *GoEmitter) EmitDerivedPrefix(d DerivedType) {
	if d.ImportedFrom.Name != "" {
		g.outputFrag(d.ImportedFrom.Name + ".")
	}
}

func (g *GoEmitter) EmitDerivedType(d DerivedType) {
	g.EmitDerivedPrefix(d)
	g.outputFrag(d.Name.Name)
}

func (g *GoEmitter) EmitDerivedTypeInternal(d DerivedType) {
	g.EmitDerivedPrefix(d)
	g.outputFrag(g.internalStructName(d.Name.Name))
}

func (g *GoEmitter) exportSymbol(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func (g *GoEmitter) internalStructName(s string) string {
	return s + "Internal__"
}

func (g *GoEmitter) thisVariableName(s string) string {
	return strings.ToLower(s[:1])
}

func (g *GoEmitter) baseTypeNames(b BaseTypedef) (string, string, string) {
	nm := b.Ident.Name
	tv := g.thisVariableName(nm)
	isn := g.internalStructName(nm)
	exsym := g.exportSymbol(nm)
	return tv, isn, exsym
}

func (g *GoEmitter) emitEnumImport(e Enum) {
	tv := g.thisVariableName(e.Ident.Name)
	es := g.exportSymbol(e.Ident.Name)
	isn := g.internalStructName(e.Ident.Name)
	g.foutputLine("func (%s %s) Import() %s {", tv, isn, es)
	g.tab()
	g.foutputLine("return %s(%s)", es, tv)
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitEnumExport(e Enum) {
	tv := g.thisVariableName(e.Ident.Name)
	es := g.exportSymbol(e.Ident.Name)
	isn := g.internalStructName(e.Ident.Name)
	g.foutputLine("func (%s %s) Export() *%s {", tv, es, isn)
	g.tab()
	g.foutputLine("return ((*%s)(&%s))", isn, tv)
	g.untab()
	g.outputLine("}")

}

func (g *GoEmitter) EmitEnum(e Enum) {
	g.emitStatementPremable(e.BaseStatement)
	exsym := g.exportSymbol(e.Ident.Name)
	g.outputLine("type " + exsym + " int")
	g.emptyLine()
	g.emitEnumConstants(e)
	g.emitEnumMap(e)
	g.emitEnumRevMap(e)
	g.outputLine("type " + g.internalStructName(e.Ident.Name) + " " + exsym)
	g.emitEnumImport(e)
	g.emitEnumExport(e)
}

func (g *GoEmitter) emitTypedefInternal(t Typedef) {
	g.foutputFrag("type %s ", g.internalStructName(t.Ident.Name))
	t.Type.EmitInternal(g)
	g.emptyLine()
}

func (g *GoEmitter) emitTypedefExport(t Typedef) {
	tv, isn, exsym := g.baseTypeNames(t.BaseTypedef)
	g.foutputLine("func (%s %s) Export() *%s {", tv, exsym, isn)
	g.tab()
	tmp := "tmp"
	g.foutputFrag("%s := ((", tmp)
	t.Type.Emit(g)
	g.foutputLine(")(%s))", tv)
	g.foutputFrag("return ((*%s)(", isn)
	t.Type.EmitExport(g, tmp)
	g.outputLine("))")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitTypedefImport(t Typedef) {
	tv, isn, exsym := g.baseTypeNames(t.BaseTypedef)
	g.foutputLine("func (%s %s) Import() %s {", tv, isn, exsym)
	g.tab()
	g.outputFrag("tmp := (")
	t.Type.EmitInternal(g)
	g.foutputLine(")(%s)", tv)
	g.foutputFrag("return %s(", exsym)
	t.Type.EmitImport(g, "&tmp")
	g.outputLine(")")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitCodec(b BaseTypedef) {
	tv, isn, exsym := g.baseTypeNames(b)

	g.foutputLine("func (%s *%s) Encode(enc rpc.Encoder) error {", tv, exsym)
	g.tab()
	g.foutputLine("return enc.Encode(%s.Export())", tv)
	g.untab()
	g.outputLine("}")
	g.emptyLine()

	g.foutputLine("func (%s *%s) Decode(dec rpc.Decoder) error {", tv, exsym)
	g.tab()
	g.foutputLine("var tmp %s", isn)
	g.foutputLine("err := dec.Decode(&tmp)")
	g.outputLine("if err != nil {")
	g.tab()
	g.foutputLine("return err")
	g.untab()
	g.outputLine("}")
	g.foutputLine("*%s = tmp.Import()", tv)
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.emptyLine()
}

func (g *GoEmitter) emitID(b BaseTypedef) {
	if !b.UniqueID.IsSet() {
		return
	}
	tv, _, exsym := g.baseTypeNames(b)
	tuid := "TypeUniqueID"
	nm := exsym + tuid
	g.foutputLine("var %s = rpc.%s(%s)", nm, tuid, b.UniqueID.Val)
	g.foutputLine("func (%s *%s) Get%s() rpc.%s {", tv, exsym, tuid, tuid)
	g.tab()
	g.foutputLine("return %s", nm)
	g.untab()
	g.outputLine("}")
	g.addUnique(nm)
}

func (g *GoEmitter) emitBytesTypedef(t Typedef) {
	tv, _, exsym := g.baseTypeNames(t.BaseTypedef)
	g.foutputLine("func (%s %s) Bytes() []byte {", tv, exsym)
	g.tab()
	g.outputFrag("return ")
	t.Type.EmitBytes(g, tv)
	g.emptyLine()
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitType(t Type)         { t.Emit(g) }
func (g *GoEmitter) emitTypeInternal(t Type) { t.EmitInternal(g) }

func (g *GoEmitter) EmitTypedef(t Typedef) {
	g.emitStatementPremable(t.BaseStatement)
	exsym := g.exportSymbol(t.Ident.Name)
	g.foutputFrag("type %s ", exsym)
	g.emitType(t.Type)
	g.emptyLine()
	g.emitTypedefInternal(t)
	g.emitTypedefExport(t)
	g.emitTypedefImport(t)
	g.emptyLine()
	g.emitCodec(t.BaseTypedef)
	g.emitID(t.BaseTypedef)
	g.emitBytesTypedef(t)

	// If we've typedef'ed to a Future(Foo) type, then we need to link
	// the Unique IDs of this child object to the parent's.
	t.Type.EmitFutureLink(g, t.Ident.Name)
}

func (g *GoEmitter) EmitVoid(v Void)     {}
func (g *GoEmitter) EmitBlob(b Blob)     { g.emitBlob(b.Count) }
func (g *GoEmitter) EmitFuture(f Future) { g.emitBlob(0) }
func (g *GoEmitter) EmitText(t Text)     { g.outputFrag("string") }
func (g *GoEmitter) EmitUint(u Uint)     { g.outputFrag("uint64") }
func (g *GoEmitter) EmitInt(i Int)       { g.outputFrag("int64") }
func (g *GoEmitter) EmitBool(b Bool)     { g.outputFrag("bool") }

func (g *GoEmitter) EmitList(l List) {
	g.outputFrag("[]")
	g.emitType(l.Type)
}

func (g *GoEmitter) EmitOptionInternal(o Option) {
	g.outputFrag("*")
	g.emitTypeInternal(o.Type)
}

func (g *GoEmitter) EmitListInternal(l List) {
	g.outputFrag("[](")
	if !l.Type.IsPrimitiveType() {
		g.outputFrag("*")
	}
	g.emitTypeInternal(l.Type)
	g.outputFrag(")")
}

func (g *GoEmitter) EmitOption(o Option) {
	g.outputFrag("*")
	g.emitType(o.Type)
}

func (g *GoEmitter) emitBlob(n int) {
	g.foutputFrag("[")
	if n > 0 {
		g.foutputFrag("%d", n)
	}
	g.foutputFrag("]byte")
}

func (g *GoEmitter) outputParamsMaybe(p string) {
	if len(p) > 0 {
		g.foutputFrag("(%s)", p)
	}
}

func (g *GoEmitter) EmitExportList(l List, param string) {
	g.outputFrag("(func (x ")
	l.Emit(g)
	g.outputFrag(") * ")
	l.EmitInternal(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("if len(x) == 0 {")
	g.tab()
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.outputFrag("ret := make(")
	l.EmitInternal(g)
	g.outputLine(", len(x))")

	// Satisfy the golinter
	if l.Type.IsPrimitiveType() {
		g.outputLine("copy(ret, x)")
	} else {
		g.outputLine("for k,v := range x {")
		g.tab()
		g.outputFrag("ret[k] = ")
		l.Type.EmitExport(g, "v")
		g.emptyLine()
		g.untab()
		g.outputLine("}")
	}
	g.outputLine("return &ret")
	g.untab()
	g.outputFrag("})")
	g.outputParamsMaybe(param)
}

func (g *GoEmitter) emitImportSignature(t Type) {
	g.outputFrag("(func (x *")
	t.EmitInternal(g)
	g.outputFrag(") (ret ")
	t.Emit(g)
	g.outputLine(") {")
	g.tab()
}

func (g *GoEmitter) emitImportPreamble(t Type) {
	g.emitImportSignature(t)
	g.outputLine("if x == nil {")
	g.tab()
	g.outputLine("return ret")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) EmitImportList(l List, param string) {
	g.emitImportSignature(l)
	g.outputLine("if x == nil || len(*x) == 0 {")
	g.tab()
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.outputFrag("ret = make(")
	l.Emit(g)
	g.outputLine(", len(*x))")
	g.outputLine("for k,v := range *x {")
	g.tab()
	if !l.Type.IsPrimitiveType() {
		g.outputLine("if v == nil {")
		g.tab()
		g.outputLine("continue")
		g.untab()
		g.outputLine("}")
	}
	g.outputFrag("ret[k] = ")
	var mkref string
	if l.Type.IsPrimitiveType() {
		mkref = "&"
	}
	l.Type.EmitImport(g, mkref+"v")
	g.emptyLine()
	g.untab()
	g.outputLine("}")
	g.outputLine("return ret")
	g.untab()
	g.outputFrag("})")
	g.outputParamsMaybe(param)
}

func (g *GoEmitter) emitImportPrimitiveType(t Type, param string) {
	g.emitImportPreamble(t)
	g.outputLine("return *x")
	g.untab()
	g.outputFrag("})")
	g.outputParamsMaybe(param)
}

func (g *GoEmitter) EmitImportOption(o Option, param string) {
	g.outputFrag("(func (x *")
	o.Type.EmitInternal(g)
	g.outputFrag(") *")
	o.Type.Emit(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("if x == nil {")
	g.tab()
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.outputFrag("tmp := ")
	o.Type.EmitImport(g, "x")
	g.emptyLine()
	g.outputLine("return &tmp")
	g.untab()
	g.outputFrag("})")
	g.outputParamsMaybe(param)
}

func (g *GoEmitter) EmitImportDerivedType(d DerivedType, param string) {
	g.emitImportPreamble(d)
	g.outputLine("return x.Import()")
	g.untab()
	g.outputFrag("})")
	g.outputParamsMaybe(param)
}

func (g *GoEmitter) EmitImportBlob(b Blob, param string) { g.emitImportPrimitiveType(b, param) }
func (g *GoEmitter) EmitImportText(t Text, param string) { g.emitImportPrimitiveType(t, param) }
func (g *GoEmitter) EmitImportUint(u Uint, param string) { g.emitImportPrimitiveType(u, param) }
func (g *GoEmitter) EmitImportInt(i Int, param string)   { g.emitImportPrimitiveType(i, param) }
func (g *GoEmitter) EmitImportBool(b Bool, param string) { g.emitImportPrimitiveType(b, param) }
func (g *GoEmitter) EmitImportFuture(f Future, param string) {
	g.emitImportPrimitiveType(Blob{}, param)
}

func (g *GoEmitter) emitExportPrimitiveType(t Type, param string) {
	// optimization
	if param != "" {
		g.foutputFrag("&%s", param)
		return
	}
	g.outputFrag("(func (x ")
	t.EmitInternal(g)
	g.outputFrag(") * ")
	t.Emit(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("return &x")
	g.untab()
	g.outputFrag("})")
}

func (g *GoEmitter) EmitExportOption(o Option, param string) {
	if o.Type.IsPrimitiveType() && param != "" {
		g.outputFrag(param)
		return
	}
	g.outputFrag("(func (x *")
	o.Type.Emit(g)
	g.outputFrag(") * ")
	o.Type.EmitInternal(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("if x == nil {")
	g.tab()
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.outputFrag("return ")
	o.Type.EmitExport(g, "(*x)")
	g.emptyLine()
	g.untab()
	g.outputFrag("})")
	g.outputParamsMaybe(param)
}

func (g *GoEmitter) EmitExportDerivedType(d DerivedType, param string) {
	if param != "" {
		g.foutputFrag("%s.Export()", param)
		return
	}
	g.outputFrag("(func (x ")
	d.Emit(g)
	g.outputFrag(") * ")
	d.EmitInternal(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("return x.Export()")
	g.untab()
	g.outputFrag("})")
}

func (g *GoEmitter) EmitExportBlob(b Blob, param string) { g.emitExportPrimitiveType(b, param) }
func (g *GoEmitter) EmitExportText(t Text, param string) { g.emitExportPrimitiveType(t, param) }
func (g *GoEmitter) EmitExportUint(u Uint, param string) { g.emitExportPrimitiveType(u, param) }
func (g *GoEmitter) EmitExportInt(i Int, param string)   { g.emitExportPrimitiveType(i, param) }
func (g *GoEmitter) EmitExportBool(b Bool, param string) { g.emitExportPrimitiveType(b, param) }
func (g *GoEmitter) EmitExportFuture(f Future, param string) {
	g.emitExportPrimitiveType(Blob{}, param)
}

func (g *GoEmitter) EmitNil() {
	g.outputFrag("nil")
}

func (g *GoEmitter) EmitBlobToBytes(nm string) {
	g.foutputFrag("(%s)[:]", nm)
}

func (g *GoEmitter) EmitBytesDowncast(klass string, varName string) {
	g.foutputFrag("((%s)(%s)).Bytes()", klass, varName)
}

func (g *GoEmitter) EmitFutureLink(parent Type, child string) {
	nm := g.exportSymbol(child)
	tv := g.thisVariableName(child)

	g.foutputFrag("func (%s *%s) AllocAndDecode(f rpc.DecoderFactory) (*", tv, nm)
	parent.Emit(g)
	g.outputLine(", error) {")
	g.tab()
	g.outputFrag("var ret ")
	parent.Emit(g)
	g.emptyLine()
	g.foutputLine("src := f.NewDecoderBytes(&ret, %s.Bytes())", tv)
	g.outputLine("if err := ret.Decode(src); err != nil {")
	g.tab()
	g.outputLine("return nil, err")
	g.untab()
	g.outputLine("}")
	g.outputLine("return &ret, nil")
	g.untab()
	g.foutputLine("}")

	g.foutputFrag("func (%s *%s) AssertNormalized() error { return nil }", tv, nm)
	g.emptyLine()

	g.foutputFrag("func (%s *", tv)
	parent.Emit(g)
	g.foutputLine(") EncodeTyped(f rpc.EncoderFactory) (*%s, error) {", nm)
	g.tab()
	g.outputLine("var tmp []byte")
	g.outputLine("enc := f.NewEncoderBytes(&tmp)")
	g.outputLine("if err := enc.Encode(enc); err != nil {")
	g.tab()
	g.outputLine("return nil, err")
	g.untab()
	g.outputLine("}")
	g.foutputLine("ret := %s(tmp)", nm)
	g.outputLine("return &ret, nil")
	g.untab()
	g.outputLine("}")

	g.foutputFrag("func (%s *", tv)
	parent.Emit(g)
	b := "__b"
	g.foutputLine(")  ChildBlob(%s []byte) %s {", b, nm)
	g.tab()
	g.foutputLine("return %s(%s)", nm, b)
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitStructVisibleField(f Field) {
	nm := g.exportSymbol(f.Ident.Name)
	g.outputFrag(nm + " ")
	f.Type.Emit(g)
	g.emptyLine()
}

func (g *GoEmitter) emitStructVisible(s Struct) {
	t := g.exportSymbol(s.Ident.Name)
	g.foutputLine("type %s struct {", t)
	g.tab()
	for _, f := range s.Fields {
		g.emitStructVisibleField(f)
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitStructInternalField(f Field) {
	nm := g.exportSymbol(f.Ident.Name)
	g.outputFrag(nm + " ")
	f.Type.MakeOptional().EmitInternal(g)
	g.emptyLine()
}

func (g *GoEmitter) emitMsgpackStructOpts() {
	g.outputLine(
		"_struct struct{} `codec:\",toarray\"` //lint:ignore U1000 msgpack internal field",
	)
}

func (g *GoEmitter) emitStructInternal(s Struct) {
	isn := g.internalStructName(s.Ident.Name)
	g.foutputLine("type %s struct {", isn)
	g.tab()
	g.emitMsgpackStructOpts()
	i := 0
	for _, f := range s.Fields {
		for i < f.Pos {
			g.foutputLine("Deprecated%d *struct{}", i)
			i++
		}
		g.emitStructInternalField(f)
		i++
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitStructImport(s Struct) {
	tv, isn, exsym := g.baseTypeNames(s.BaseTypedef)
	g.foutputLine("func (%s %s) Import() %s {", tv, isn, exsym)
	g.tab()
	g.foutputLine("return %s {", exsym)
	g.tab()
	for _, f := range s.Fields {
		fn := g.exportSymbol(f.Ident.Name)
		g.foutputFrag("%s: ", fn)
		f.Type.EmitImport(g, tv+"."+fn)
		g.outputLine(",")
	}
	g.untab()
	g.outputLine("}")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitStructExport(s Struct) {
	tv, isn, exsym := g.baseTypeNames(s.BaseTypedef)
	g.foutputLine("func (%s %s) Export() *%s {", tv, exsym, isn)
	g.tab()
	g.foutputLine("return &%s {", isn)
	g.tab()
	for _, f := range s.Fields {
		fn := g.exportSymbol(f.Ident.Name)
		g.foutputFrag("%s: ", fn)
		f.Type.EmitExport(g, tv+"."+fn)
		g.outputLine(",")
	}
	g.untab()
	g.outputLine("}")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitBytesNil(b BaseTypedef) {
	tv, _, exsym := g.baseTypeNames(b)
	g.foutputLine("func (%s *%s) Bytes() []byte { return nil }", tv, exsym)
}

func (g *GoEmitter) EmitStruct(s Struct) {
	g.emitStatementPremable(s.BaseStatement)
	g.emitStructVisible(s)
	g.emitStructInternal(s)
	g.emitStructImport(s)
	g.emitStructExport(s)
	g.emitCodec(s.BaseTypedef)
	g.emitID(s.BaseTypedef)
	g.emitBytesNil(s.BaseTypedef)
}

func (g *GoEmitter) variantCasePositionToVariable(i int) string {
	return fmt.Sprintf("F_%d__", i)
}

func (g *GoEmitter) emitVariantStructCase(v Variant, c Case, isInternal bool) {
	if c.Position == nil {
		return
	}
	g.foutputFrag("%s *", g.variantCasePositionToVariable(*c.Position))
	if isInternal {
		p := b64encode(*c.Position)
		c.Type.EmitInternal(g)
		g.foutputLine(" `codec:\"%s\"`", p)
	} else {
		c.Type.Emit(g)
		g.foutputLine(" `json:\"f%d,omitempty\"`", *c.Position)
	}
}

func (g *GoEmitter) emitVariantTopStruct(v Variant) {
	snm := g.exportSymbol(v.Ident.Name)
	sv := g.exportSymbol(v.SwitchVar.Name)
	g.foutputLine("type %s struct {", snm)
	g.tab()
	g.outputFrag(sv + " ")
	v.SwitchType.Emit(g)
	g.emptyLine()
	for _, c := range v.Cases {
		g.emitVariantStructCase(v, c, false)
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) switchStructName() string {
	return "Switch__"
}
func (g *GoEmitter) switchInternalStructType(s string) string {
	return s + "InternalSwitch__"
}

func (g *GoEmitter) emitVariantInternalStruct(v Variant) {
	ism := g.internalStructName(v.Ident.Name)
	sv := g.exportSymbol(v.SwitchVar.Name)
	g.foutputLine("type %s struct {", ism)
	g.tab()
	g.emitMsgpackStructOpts()
	g.outputFrag(sv + " ")
	v.SwitchType.Emit(g)
	g.emptyLine()
	g.foutputLine("%s %s",
		g.switchStructName(),
		g.switchInternalStructType(v.Ident.Name),
	)
	g.emptyLine()
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitVariantInternalSwitchStruct(v Variant) {
	g.foutputLine("type %s struct {", g.switchInternalStructType(v.Ident.Name))
	g.tab()
	g.outputLine(
		"_struct struct{} `codec:\",omitempty\"` //lint:ignore U1000 msgpack internal field",
	)
	for _, c := range v.Cases {
		g.emitVariantStructCase(v, c, true)
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) caseDataAccess(v Variant, c Case) string {
	p := c.Position
	if p == nil {
		return ""
	}
	tv := g.thisVariableName(v.Ident.Name)
	return strings.Join([]string{
		tv,
		g.variantCasePositionToVariable(*p),
	}, ".")
}

func (g *GoEmitter) emitVariantSwitchAccessorCase(v Variant, c Case, tv string) {
	if len(c.Labels) == 0 {
		g.outputLine("default:")
	} else {
		var labels []string
		for _, l := range c.Labels {
			labels = append(labels, l.CaseLabelToString(g, v.SwitchType))
		}
		g.foutputLine("case %s:", strings.Join(labels, ", "))
	}
	g.tab()
	p := c.Position
	if p == nil {
		g.outputLine("break")
	} else {
		cda := g.caseDataAccess(v, c)
		if cda == "" {
			panic("case data access is nil, should never be")
		}
		g.foutputLine("if %s == nil {", cda)
		g.tab()
		g.foutputLine("return ret, errors.New(\"unexpected nil case for %s\")",
			g.variantCasePositionToVariable(*p),
		)
		g.untab()
		g.outputLine("}")
	}
	g.untab()
}

func (g *GoEmitter) emitVariantSwitchAccessor(v Variant) {
	tv := g.thisVariableName(v.Ident.Name)
	exsym := g.exportSymbol(v.Ident.Name)
	lclsv := g.exportSymbol(v.SwitchVar.Name)
	sv := tv + "." + lclsv

	g.foutputFrag("func (%s %s) Get%s() (ret ", tv, exsym, lclsv)
	v.SwitchType.Emit(g)
	g.outputLine(", err error) {")
	g.tab()
	g.foutputLine("switch %s {", sv)
	g.tab()
	for _, c := range v.Cases {
		g.emitVariantSwitchAccessorCase(v, c, tv)
	}
	g.untab()
	g.outputLine("}")
	g.foutputLine("return %s, nil", sv)
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) switchValue(v Variant) string {
	tv := g.thisVariableName(v.Ident.Name)
	sv := g.exportSymbol(v.SwitchVar.Name)
	return fmt.Sprintf("%s.%s", tv, sv)
}

func (g *GoEmitter) emitVariantDataAcceessorsCase(v Variant, c Case) {
	cda := g.caseDataAccess(v, c)
	if cda == "" {
		return
	}

	type pair struct {
		getterMethodName string
		caseLabel        string
	}

	var pairs []pair
	if len(c.Labels) == 0 {
		pairs = []pair{{getterMethodName: "Default"}}
	} else {
		for _, l := range c.Labels {
			pairs = append(pairs, pair{
				getterMethodName: l.GetterMethodName(g),
				caseLabel:        l.CaseLabelToString(g, v.SwitchType),
			})
		}
	}

	sv := g.switchValue(v)
	exsym := g.exportSymbol(v.Ident.Name)
	tv := g.thisVariableName(v.Ident.Name)

	for _, p := range pairs {
		g.foutputFrag("func (%s %s) %s() ", tv, exsym, p.getterMethodName)
		c.Type.Emit(g)
		g.outputLine(" {")
		g.tab()
		g.foutputLine("if %s == nil {", cda)
		g.tab()
		g.foutputLine("panic(\"unexpected nil case; should have been checked\")")
		g.untab()
		g.outputLine("}")
		if len(p.caseLabel) > 0 {
			switch p.caseLabel {
			case "true":
				g.foutputLine("if !%s {", sv)
			case "false":
				g.foutputLine("if %s {", sv)
			default:
				g.foutputLine("if %s != %s {", sv, p.caseLabel)
			}
			g.tab()
			g.outputLine(
				`panic(fmt.Sprintf("unexpected switch value (%v) when ` + p.getterMethodName +
					` is called", ` + sv + `))`,
			)
			g.untab()
			g.outputLine("}")
		}
		g.foutputLine("return *%s", cda)
		g.untab()
		g.outputLine("}")
	}
}

func (g *GoEmitter) emitVariantConstructorCase(v Variant, c Case) {
	tt := g.exportSymbol(v.Ident.Name)
	defConstructor := "New" + tt + "Default"

	type pair struct {
		constructorName string
		caseLabel       string
	}
	var pairs []pair
	if len(c.Labels) == 0 {
		pairs = []pair{{constructorName: defConstructor, caseLabel: "s"}}
	} else {
		for _, l := range c.Labels {
			pairs = append(pairs, pair{
				constructorName: l.ConstructorName(g, v.Ident.Name),
				caseLabel:       l.CaseLabelToString(g, v.SwitchType),
			})
		}
	}
	for _, p := range pairs {
		g.foutputFrag("func %s(", p.constructorName)
		var didOutput bool
		if len(c.Labels) == 0 {
			g.outputFrag("s ")
			v.SwitchType.Emit(g)
			didOutput = true
		}
		if c.HasData() {
			if didOutput {
				g.outputFrag(", ")
			}
			g.outputFrag("v ")
			c.Type.Emit(g)
		}

		g.foutputLine(") %s {", tt)
		g.tab()
		g.foutputLine("return %s{", tt)
		g.tab()
		g.foutputLine("%s: %s,",
			g.exportSymbol(v.SwitchVar.Name),
			p.caseLabel,
		)
		if c.HasData() && c.Position != nil {
			cda := g.variantCasePositionToVariable(*c.Position)
			g.foutputLine("%s: &v,", cda)
		}
		g.untab()
		g.outputLine("}")
		g.untab()
		g.outputLine("}")
	}
}

func (g *GoEmitter) emitVariantDataAccessors(v Variant) {
	for _, c := range v.Cases {
		g.emitVariantDataAcceessorsCase(v, c)
	}
}

func (g *GoEmitter) emitVariantConstructors(v Variant) {
	for _, c := range v.Cases {
		g.emitVariantConstructorCase(v, c)
	}
}

func (g *GoEmitter) emitVariantImportCase(v Variant, c Case) {
	p := c.Position
	if p == nil {
		return
	}
	tv := g.thisVariableName(v.Ident.Name)
	field := g.variantCasePositionToVariable(*p)
	g.foutputFrag("%s: ", field)
	source := strings.Join([]string{
		tv,
		g.switchStructName(),
		field,
	}, ".")

	// optimization
	if c.Type.IsPrimitiveType() {
		g.foutputLine("%s,", source)
		return
	}

	g.outputFrag("(func (x *")
	c.Type.EmitInternal(g)
	g.outputFrag(") *")
	c.Type.Emit(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("if x == nil {")
	g.tab()
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.outputFrag("tmp := ")
	c.Type.EmitImport(g, "x")
	g.emptyLine()
	g.foutputLine("return &tmp")
	g.untab()
	g.foutputLine("})(%s),", source)
}

func (g *GoEmitter) emitVariantImport(v Variant) {
	tv, isn, exsym := g.baseTypeNames(v.BaseTypedef)
	g.foutputLine("func (%s %s) Import() %s {", tv, isn, exsym)
	g.tab()
	g.foutputLine("return %s{", exsym)
	g.tab()
	sv := g.exportSymbol(v.SwitchVar.Name)
	g.foutputLine("%s: %s.%s,", sv, tv, sv)
	for _, c := range v.Cases {
		g.emitVariantImportCase(v, c)
	}
	g.untab()
	g.outputLine("}")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitVariantExportCase(v Variant, c Case) {
	if c.Position == nil {
		return
	}
	tv := g.thisVariableName(v.Ident.Name)
	field := g.variantCasePositionToVariable(*c.Position)
	g.foutputFrag("%s: ", field)

	// optimization!
	if c.Type.IsPrimitiveType() {
		g.foutputLine("%s.%s,", tv, field)
		return
	}

	g.outputFrag("(func (x *")
	c.Type.Emit(g)
	g.outputFrag(") *")
	c.Type.EmitInternal(g)
	g.outputLine(" {")
	g.tab()
	g.outputLine("if x == nil {")
	g.tab()
	g.outputLine("return nil")
	g.untab()
	g.outputLine("}")
	g.outputFrag("return ")
	c.Type.EmitExport(g, "(*x)")
	g.emptyLine()
	g.untab()
	g.foutputLine("})(%s.%s),", tv, field)
}

func (g *GoEmitter) emitVariantExport(v Variant) {
	tv, isn, exsym := g.baseTypeNames(v.BaseTypedef)
	g.foutputLine("func (%s %s) Export() *%s {", tv, exsym, isn)
	g.tab()
	g.foutputLine("return &%s{", isn)
	g.tab()
	sv := g.exportSymbol(v.SwitchVar.Name)
	g.foutputLine("%s: %s.%s,", sv, tv, sv)
	g.foutputLine("%s: %s{",
		g.switchStructName(),
		g.switchInternalStructType(v.Ident.Name),
	)
	g.tab()
	for _, c := range v.Cases {
		g.emitVariantExportCase(v, c)
	}
	g.untab()
	g.outputLine("},")
	g.untab()
	g.outputLine("}")
	g.untab()
	g.outputLine("}")

}

func (g *GoEmitter) EmitVariant(v Variant) {
	g.emitStatementPremable(v.BaseStatement)
	g.emitVariantTopStruct(v)
	g.emitVariantInternalStruct(v)
	g.emitVariantInternalSwitchStruct(v)
	g.emitVariantSwitchAccessor(v)
	g.emitVariantDataAccessors(v)
	g.emitVariantConstructors(v)
	g.emitVariantImport(v)
	g.emitVariantExport(v)
	g.emitCodec(v.BaseTypedef)
	g.emitID(v.BaseTypedef)
	g.emitBytesNil(v.BaseTypedef)
}

func (g *GoEmitter) protocolID(p Protocol) string {
	return g.exportSymbol(p.Ident.Name) + "ProtocolID"
}

func (g *GoEmitter) emitProtocolID(p Protocol) {
	nm := g.protocolID(p)
	g.foutputLine(
		`var %s rpc.ProtocolUniqueID = rpc.ProtocolUniqueID(%s)`,
		nm,
		p.UniqueID.Val,
	)
	g.addUnique(nm)
}

func (g *GoEmitter) emitMethodArgs(p Protocol, m Method) {
	argName := m.makeArgName(g)
	s := m.ParamsToStruct(argName)
	s.emit(g)
}

func (g *GoEmitter) emitMethodsArgs(p Protocol) {
	for _, m := range p.Methods {
		g.emitMethodArgs(p, m)
	}
}

func (g *GoEmitter) emitServerHookSignature(p Protocol, m Method) {
	g.emitDecorators(m.Dec)
	exsym := g.exportSymbol(m.Ident.Name)
	g.foutputFrag("%s(context.Context", exsym)
	if len(m.Params) > 0 {
		g.outputFrag(", ")
		if len(m.Params) == 1 && m.Params[0].Pos == 0 {
			m.Params[0].Type.Emit(g)
		} else {
			g.outputFrag(m.makeArgName(g))
		}
	}
	g.outputFrag(") (")
	if !m.ResType.IsVoid() {
		m.ResType.Emit(g)
		g.outputFrag(", ")
	}
	g.outputLine("error)")
}

func (g *GoEmitter) emitServerInterface(p Protocol) {
	g.emitStatementPremable(p.BaseStatement)
	nm := g.exportSymbol(p.Ident.Name)
	g.foutputLine("type %sInterface interface {", nm)
	g.tab()
	for _, m := range p.Methods {
		g.emitServerHookSignature(p, m)
	}
	g.outputFrag("ErrorWrapper() func(error) ")
	p.Modifiers.Errors.Type.Emit(g)
	g.emptyLine()
	if p.Modifiers.ArgHeader != nil {
		g.outputFrag("CheckArgHeader(ctx context.Context, h ")
		p.Modifiers.ArgHeader.Type.Emit(g)
		g.outputFrag(") error")
		g.emptyLine()
	}
	if p.Modifiers.ResHeader != nil {
		g.outputFrag("MakeResHeader() ")
		p.Modifiers.ResHeader.Type.Emit(g)
		g.emptyLine()
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitServerWrapError(p Protocol) {
	exsym := g.exportSymbol(p.Ident.Name)

	g.foutputLine(
		`func %sMakeGenericErrorWrapper(f %sErrorWrapper) rpc.WrapErrorFunc {`,
		exsym, exsym,
	)
	g.tab()
	g.outputLine("return func(err error) interface{} {")
	g.tab()
	g.outputLine("if err == nil {")
	g.tab()
	g.outputLine("return err")
	g.untab()
	g.outputLine("}")
	var exp string
	if !p.Modifiers.Errors.Type.IsPrimitiveType() {
		exp = ".Export()"
	}
	g.foutputLine("return f(err)%s", exp)
	g.untab()
	g.outputLine("}")
	g.untab()
	g.outputLine("}")

}
func (g *GoEmitter) emitClientErrorUnwrapperType(p Protocol) {
	exsym := g.exportSymbol(p.Ident.Name)
	g.foutputFrag("type %sErrorUnwrapper func(", exsym)
	p.Modifiers.Errors.Type.Emit(g)
	g.outputLine(") error")
}

func (g *GoEmitter) emitClientErrorWrapperType(p Protocol) {
	exsym := g.exportSymbol(p.Ident.Name)
	g.foutputFrag("type %sErrorWrapper func(error) ", exsym)
	p.Modifiers.Errors.Type.Emit(g)
	g.emptyLine()
	g.emptyLine()
}
func (g *GoEmitter) privateSymbol(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func (g *GoEmitter) emitClientErrorUnwrapperAdapterStruct(p Protocol) {
	nm := g.privateSymbol(p.Ident.Name) + "ErrorUnwrapperAdapter"
	hook := g.exportSymbol(p.Ident.Name) + "ErrorUnwrapper"
	tv := g.thisVariableName(p.Ident.Name)

	g.foutputLine("type %s struct {", nm)
	g.tab()
	g.foutputLine("h %s", hook)
	g.untab()
	g.outputLine("}")
	g.emptyLine()

	g.foutputLine("func (%s %s) MakeArg() interface{} {", tv, nm)
	g.tab()
	g.foutputFrag("return &")
	p.Modifiers.Errors.Type.EmitInternal(g)
	g.outputLine("{}")
	g.untab()
	g.outputLine("}")
	g.emptyLine()

	g.foutputLine("func (%s %s) UnwrapError(raw interface{}) (appError error, dispatchError error) {",
		tv, nm)
	g.tab()
	convVar := "sTmp"
	g.foutputFrag("%s, ok := raw.(*", convVar)
	p.Modifiers.Errors.Type.EmitInternal(g)
	g.outputLine(")")
	g.outputLine("if !ok {")
	g.tab()
	g.outputLine("return nil, errors.New(\"error converting to internal type in UnwrapError\")")
	g.untab()
	g.outputLine("}")
	g.foutputLine("if %s == nil {", convVar)
	g.tab()
	g.outputLine("return nil, nil")
	g.untab()
	g.outputLine("}")
	g.foutputLine("return %s.h(%s.Import()), nil", tv, convVar)
	g.untab()
	g.outputLine("}")
	g.emptyLine()
	g.foutputLine("var _ rpc.ErrorUnwrapper = %s{}", nm)

}

func (g *GoEmitter) emitClientErrorUnwrapper(p Protocol) {
	g.emitClientErrorUnwrapperType(p)
	g.emitClientErrorWrapperType(p)
	g.emitClientErrorUnwrapperAdapterStruct(p)
}

func (g *GoEmitter) emitClientStub(p Protocol) {
	exsym := g.exportSymbol(p.Ident.Name)
	g.foutputLine("type %sClient struct {", exsym)
	g.tab()
	g.outputLine("Cli rpc.GenericClient")
	g.foutputLine("ErrorUnwrapper %sErrorUnwrapper", exsym)
	if p.Modifiers.ArgHeader != nil {
		g.outputFrag(`MakeArgHeader func() `)
		p.Modifiers.ArgHeader.Type.Emit(g)
		g.emptyLine()
	}
	if p.Modifiers.ResHeader != nil {
		g.outputFrag(`CheckResHeader func(context.Context, `)
		p.Modifiers.ResHeader.Type.Emit(g)
		g.outputLine(") error")
	}
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitClientMethod(p Protocol, m Method) {
	pn := g.exportSymbol(p.Ident.Name)
	mn := g.exportSymbol(m.Ident.Name)

	g.foutputFrag("func (c %sClient) %s (ctx context.Context", pn, mn)
	argStructName := m.makeArgName(g) // mn + "Arg"
	if len(m.Params) > 0 {
		g.outputFrag(", ")
		if m.singleArg() {
			g.foutputFrag("%s ", m.Params[0].Ident.Name)
			m.Params[0].Type.Emit(g)
		} else {
			g.foutputFrag("arg %s", argStructName)
		}
	}
	g.outputFrag(") (")
	if !m.ResType.IsVoid() {
		g.foutputFrag("res ")
		m.ResType.Emit(g)
		g.outputFrag(", ")
	}
	g.outputLine("err error) {")
	g.tab()

	if m.singleArg() {
		g.foutputLine("arg := %s{", argStructName)
		g.tab()
		n := m.Params[0].Ident.Name
		g.foutputLine("%s: %s,", g.exportSymbol(n), n)
		g.untab()
		g.outputLine("}")
	} else if len(m.Params) == 0 {
		g.foutputLine("var arg %s", argStructName)
	}

	if p.Modifiers.ArgHeader != nil {
		g.foutputFrag("warg := &rpc.DataWrap[")
		p.Modifiers.ArgHeader.Type.Emit(g)
		g.outputFrag(", *")
		argType := g.internalStructName(m.makeArgName(g))
		g.outputFrag(argType)
		g.outputLine("] {")
		g.tab()
		g.outputLine("Data: arg.Export(),")
		g.untab()
		g.outputLine("}")
		g.outputLine("if c.MakeArgHeader != nil {")
		g.tab()
		g.outputLine("warg.Header = c.MakeArgHeader()")
		g.untab()
		g.outputLine("}")
	} else {
		g.outputLine("warg := arg.Export()")
	}

	var nilRes bool

	if p.Modifiers.ResHeader != nil {
		g.outputFrag("var tmp rpc.DataWrap[")
		p.Modifiers.ResHeader.Type.Emit(g)
		g.outputFrag(", ")
		if m.ResType.IsVoid() {
			g.outputFrag("interface{}")
		} else {
			m.ResType.EmitInternal(g)
		}
		g.outputLine("]")
	} else if !m.ResType.IsVoid() {
		g.outputFrag("var tmp ")
		m.ResType.Emit(g)
		g.emptyLine()
	} else {
		nilRes = true
	}

	res := "&tmp"
	if nilRes {
		res = "nil"
	}
	method := fmt.Sprintf("rpc.NewMethodV2(%s, %d, \"%s.%s\")",
		g.protocolID(p), m.Pos, p.Ident.Name, m.Ident.Name)

	adapter := g.privateSymbol(p.Ident.Name) + "ErrorUnwrapperAdapter{" +
		"h: c.ErrorUnwrapper}"

	g.foutputLine("err = c.Cli.Call2(ctx, %s, warg, %s, 0 * time.Millisecond, %s)",
		method, res, adapter)

	g.outputLine("if err != nil {")
	g.tab()
	g.outputLine("return")
	g.untab()
	g.outputLine("}")

	if p.Modifiers.ResHeader != nil {
		g.outputLine("if c.CheckResHeader != nil {")
		g.tab()
		g.outputLine("err = c.CheckResHeader(ctx, tmp.Header)")
		g.outputLine("if err != nil {")
		g.tab()
		g.outputLine("return")
		g.untab()
		g.outputLine("}")
		g.untab()
		g.outputLine("}")
	}

	if !m.ResType.IsVoid() {
		tmp := "tmp"
		if p.Modifiers.ResHeader != nil {
			tmp = "tmp.Data"
		}
		g.outputFrag("res = ")
		if m.ResType.IsPrimitiveType() {
			g.outputLine(tmp)
		} else if m.ResType.IsList() {
			m.ResType.EmitImport(g, "&"+tmp)
			g.emptyLine()
		} else {
			g.foutputLine("%s.Import()", tmp)
		}
	}
	g.outputLine("return")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) emitClientMethods(p Protocol) {
	for _, m := range p.Methods {
		g.emitClientMethod(p, m)
	}
}

func (g *GoEmitter) emitServerProtocolHandler(p Protocol, m Method) {
	argType := g.internalStructName(m.makeArgName(g))
	g.foutputLine("%d: {", m.Pos)
	g.tab()

	g.outputLine("ServeHandlerDescription: rpc.ServeHandlerDescription{")
	g.tab()
	g.outputLine("MakeArg : func() interface{} {")
	g.tab()
	if p.Modifiers.ArgHeader != nil {
		g.outputFrag("var ret rpc.DataWrap[")
		p.Modifiers.ArgHeader.Type.Emit(g)
		g.outputFrag(", *")
		g.outputFrag(argType)
		g.outputLine("]")
	} else {
		g.foutputLine("var ret %s", argType)
	}
	g.outputLine("return &ret")
	g.untab()
	g.outputLine("},")

	g.outputLine("Handler: func(ctx context.Context, args interface{}) (interface{}, error) {")
	g.tab()

	if p.Modifiers.ArgHeader != nil {
		g.foutputFrag("typedWrappedArg, ok := args.(*rpc.DataWrap[")
		p.Modifiers.ArgHeader.Type.Emit(g)
		g.outputFrag(", *")
		g.outputFrag(argType)
		g.outputLine("])")
		g.outputLine("if !ok {")
		g.tab()
		g.outputFrag("err := rpc.NewTypeError((*rpc.DataWrap[")
		p.Modifiers.ArgHeader.Type.Emit(g)
		g.outputFrag(", *")
		g.outputFrag(argType)
		g.outputLine("])(nil), args)")
		g.outputLine("return nil, err")
		g.untab()
		g.outputLine("}")
		g.outputLine("if err := i.CheckArgHeader(ctx, typedWrappedArg.Header); err != nil {")
		g.tab()
		g.outputLine("return nil, err")
		g.untab()
		g.outputLine("}")
		if len(m.Params) > 0 {
			g.outputLine("typedArg := typedWrappedArg.Data")
		}
	} else {
		typedArgs := "typedArg"
		if len(m.Params) == 0 {
			typedArgs = "_"
		}
		g.foutputLine(("%s, ok := args.(*%s)"), typedArgs, argType)
		g.outputLine("if !ok {")
		g.tab()
		g.foutputLine("err := rpc.NewTypeError((*%s)(nil), args)", argType)
		g.outputLine("return nil, err")
		g.untab()
		g.outputLine("}")
	}

	ret := "tmp, "
	if m.ResType.IsVoid() {
		ret = ""
	}
	arg := ", (typedArg.Import())"
	if m.singleArg() {
		arg += "." + g.exportSymbol(m.Params[0].Ident.Name)
	} else if len(m.Params) == 0 {
		arg = ""
	}
	g.foutputLine("%serr := i.%s(ctx%s)", ret, g.exportSymbol(m.Ident.Name), arg)
	g.outputLine("if err != nil {")
	g.tab()
	g.outputLine("return nil, err")
	g.untab()
	g.outputLine("}")

	if !m.ResType.IsVoid() && m.ResType.IsList() {
		g.outputFrag("lst := ")
		m.ResType.EmitExport(g, "tmp")
		g.emptyLine()
	}

	if p.Modifiers.ResHeader != nil {
		g.outputFrag("ret := rpc.DataWrap[")
		p.Modifiers.ResHeader.Type.Emit(g)
		g.outputFrag(", ")
		if m.ResType.IsVoid() {
			g.outputFrag("interface{}")
		} else {
			if !m.ResType.IsPrimitiveType() && !m.ResType.IsList() {
				g.outputFrag("*")
			}
			m.ResType.EmitInternal(g)
		}
		g.outputLine("]{")
		g.tab()
		if m.ResType.IsVoid() || m.ResType.IsList() {
			// noop
		} else if m.ResType.IsPrimitiveType() {
			g.foutputLine("Data: tmp,")
		} else {
			g.foutputLine("Data: tmp.Export(),")
		}
		g.outputLine("Header : i.MakeResHeader(),")
		g.untab()
		g.outputLine("}")
		if m.ResType.IsList() {
			g.outputLine("if lst != nil {")
			g.tab()
			g.foutputLine("ret.Data = *lst")
			g.untab()
			g.outputLine("}")
		}
		g.outputLine("return &ret, nil")
	} else {
		if m.ResType.IsVoid() {
			g.outputLine("return nil, nil")
		} else if m.ResType.IsPrimitiveType() {
			g.foutputLine("return tmp, nil")
		} else if m.ResType.IsList() {
			g.outputLine("return lst, nil")
		} else {
			g.foutputLine("return tmp.Export(), nil")
		}
	}

	g.untab()
	g.outputLine("},")
	g.untab()

	g.outputLine("},")
	g.foutputLine("Name: \"%s\",", m.Ident.Name)
	g.untab()
	g.outputLine("},")
}

func (g *GoEmitter) emitServerProtocol(p Protocol) {
	exsym := g.exportSymbol(p.Ident.Name)
	g.foutputLine("func %sProtocol(i %sInterface) rpc.ProtocolV2 {", exsym, exsym)
	g.tab()
	g.foutputLine("return rpc.ProtocolV2{")
	g.tab()
	g.foutputLine("Name: \"%s\",", p.Ident.Name)
	g.foutputLine("ID: %s,", g.protocolID(p))
	g.foutputLine("Methods: map[rpc.Position]rpc.ServeHandlerDescriptionV2{")
	g.tab()
	for _, m := range p.Methods {
		g.emitServerProtocolHandler(p, m)
	}
	g.untab()
	g.outputLine("},")
	g.foutputLine("WrapError: %sMakeGenericErrorWrapper(i.ErrorWrapper()),", exsym)
	g.untab()
	g.outputLine("}")
	g.untab()
	g.outputLine("}")
}

func (g *GoEmitter) EmitProtocol(p Protocol) {
	g.emitProtocolID(p)
	g.emitMethodsArgs(p)
	g.emitServerInterface(p)
	g.emitServerWrapError(p)
	g.emitClientErrorUnwrapper(p)
	g.emitClientStub(p)
	g.emitClientMethods(p)
	g.emitServerProtocol(p)
}

func (g *GoEmitter) ToEnumConstant(t Type, nm string) string {
	var parts []string
	prfx := t.DerivedPrefix()
	if len(prfx) > 0 {
		parts = append(parts, prfx)
	}
	prfx = t.EnumPrefix()
	if len(prfx) > 0 {
		parts = append(parts, g.exportSymbol(prfx)+"_"+nm)
	} else {
		parts = append(parts, nm)
	}
	return strings.Join(parts, "")
}

func (g *GoEmitter) GetterMethodNameForBool(b bool) string {
	if b {
		return "True"
	}
	return "False"

}

func (g *GoEmitter) GetterMethodNameForInt(i int) string {
	if i >= 0 {
		return fmt.Sprintf("P%d", i)
	}
	return fmt.Sprintf("N%d", -i)
}

func (g *GoEmitter) ConstructorNameForConstant(vrnt string, cnst string) string {
	parts := []string{
		"New",
		g.exportSymbol(vrnt),
		"With",
		g.snakeToCamelCase(cnst),
	}
	return strings.Join(parts, "")
}

func (g *GoEmitter) ConstructorNameForInt(vrnt string, i int) string {
	return g.ConstructorNameForConstant(vrnt, g.GetterMethodNameForInt(i))
}

func (g *GoEmitter) ConstructorNameForBool(vrnt string, b bool) string {
	return g.ConstructorNameForConstant(vrnt, g.GetterMethodNameForBool(b))
}

func (g *GoEmitter) snakeToCamelCase(s string) string {
	s = strings.ToLower(s)
	parts := strings.Split(s, "_")
	for i, p := range parts {
		parts[i] = g.exportSymbol(p)
	}
	return strings.Join(parts, "")
}

func (g *GoEmitter) GetterMethodNameForConstant(s string) string {
	return g.exportSymbol(g.snakeToCamelCase(s))
}

func (g *GoEmitter) MethodArgName(mthd string, argName string) string {
	var ret string
	if len(argName) > 0 {
		ret = argName
	} else {
		ret = mthd + "Arg"
	}
	return g.exportSymbol(ret)
}

func (g *GoEmitter) EmitImport(i Import) {
	g.storeImport(i)
	if i.Lang == LangGo {
		g.foutputLine("import %s \"%s\"", i.Name, i.Path)
	}
}

var _ Emitter = (*GoEmitter)(nil)

func (r Root) emit(g Emitter) {
	for _, s := range r.Stmts {
		s.emit(g)
	}
}

func (t Typedef) emit(g Emitter)  { g.EmitTypedef(t) }
func (e Enum) emit(g Emitter)     { g.EmitEnum(e) }
func (v Variant) emit(g Emitter)  { g.EmitVariant(v) }
func (s Struct) emit(g Emitter)   { g.EmitStruct(s) }
func (p Protocol) emit(g Emitter) { g.EmitProtocol(p) }
func (i Import) emit(g Emitter)   { g.EmitImport(i) }
