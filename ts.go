package colfer

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pascaldekloe/name"
)

// TSKeywords are the reserved tokens for Typescript.
// Some entries are redundant due to the use of a Go parser.
var tsKeywords = map[string]struct{}{
	"break": {}, "case": {}, "catch": {}, "class": {}, "const": {},
	"continue": {}, "debugger": {}, "default": {}, "delete": {}, "do": {},
	"else": {}, "enum": {}, "export": {}, "extends": {}, "false": {}, "finally": {},
	"for": {}, "function": {}, "if": {}, "import": {}, "in": {}, "instanceof": {},
	"new": {}, "null": {}, "return": {}, "super": {}, "switch": {}, "this": {},
	"throw": {}, "true": {}, "try": {}, "typeof": {}, "var": {}, "void": {},
	"while": {}, "with": {}, "implements": {}, "interface": {}, "let": {},
	"package": {}, "private": {}, "protected": {}, "public": {}, "static": {}, "yield": {},
	"any": {}, "boolean": {}, "constructor": {}, "declare": {}, "get": {}, "module": {},
	"require": {}, "number": {}, "set": {}, "string": {}, "symbol": {}, "type": {}, "from": {}, "of": {},
}

// GenerateTS writes the code into file "Colfer.ts".
func GenerateTS(basedir string, packages Packages) error {
	for _, p := range packages {
		p.NameNative = strings.Replace(p.Name, "/", "_", -1)
		if _, ok := tsKeywords[p.NameNative]; ok {
			p.NameNative += "_"
		}

		for _, t := range p.Structs {
			t.NameNative = name.CamelCase(t.Name, true)
			for _, f := range t.Fields {
				f.NameNative = name.CamelCase(f.Name, false)
				if _, ok := tsKeywords[f.NameNative]; ok {
					f.NameNative += "_"
				}
			}
		}
	}

	t := template.New("tsCode")
	template.Must(t.Parse(tsCode))
	template.Must(t.New("marshal").Parse(tsMarshal))
	template.Must(t.New("unmarshal").Parse(tsUnmarshal))

	if err := os.MkdirAll(basedir, os.ModeDir|os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(basedir, "Colfer.ts"))
	if err != nil {
		return err
	}
	defer f.Close()
	return t.Execute(f, packages)
}

const tsCode = `// Code generated by colf(1); DO NOT EDIT.
{{- range .}}
// The compiler used schema file {{.SchemaFileList}} for package {{.Name}}.
{{- end}}

const EOF = 'colfer: EOF';

// Type Aliases
type bool = Boolean;
type uint8 = number;
type uint16 = number;
type uint32 = number;
type uint64 = number;
type int32 = number;
type int64 = number;
type float32 = number;
type float64 = number;
type timestamp = Date | undefined;
type text = string | undefined;
type binary = Uint8Array;
{{range .}}
{{.DocText "// "}}
export namespace {{.NameNative}} {

	// The upper limit for serial byte sizes.
	var colferSizeMax = {{.SizeMax}};
{{- if .HasList}}
	// The upper limit for the number of elements in a list.
	var colferListMax = {{.ListMax}};
{{- end}}
{{range .Structs}}
	// Constructor.
{{.DocText "\t// "}}
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class {{.NameNative}} {
{{- range .Fields}}
{{.DocText "\t\t// "}}
		{{.NameNative}}:
{{- if .TypeList}} {{if eq .Type "float32"}} Float32Array = new Float32Array(0){{else if eq .Type "float64"}} Float64Array = new Float64Array(0){{- else if .TypeRef}} Array<{{.TypeRef.Pkg.NameNative}}.{{.TypeRef.NameNative}}> = []{{else}} Array<{{.Type}}> = []{{end}}
{{- else if eq .Type "bool"}} {{.Type}} = false
{{- else if eq .Type "timestamp"}} {{.Type}} = undefined;
		{{.NameNative}}_ns: number = 0
{{- else if eq .Type "text"}} {{.Type}} = undefined;
{{- else if eq .Type "binary"}} {{.Type}} = new Uint8Array(0)
{{- else if .TypeRef}} {{.TypeRef.Pkg.NameNative}}.{{.TypeRef.NameNative}} | undefined = undefined
{{- else}} {{.Type}} = 0
{{- end}};{{end}}

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		{{template "marshal" .}}
		{{template "unmarshal" .}}
	}
{{end}}

	// private section

	function encodeVarint(bytes: Uint8Array, i: number, x: number) {
		while (x > 127) {
			bytes[i++] = (x & 127) | 128;
			x /= 128;
		}
		bytes[i++] = x & 127;
		return i;
	}
{{if .HasTimestamp}}
	function decodeInt64(data: Uint8Array, i: number) {
		var v = 0, j = i + 7, m = 1;
		if (data[i] & 128) {
			// two's complement
			for (var carry = 1; j >= i; --j, m *= 256) {
				const by = (data[j] ^ 255) + carry;
				carry = by >> 8;
				v += (by & 255) * m;
			}
			v = -v;
		} else {
			for (; j >= i; --j, m *= 256)
				v += data[j] * m;
		}
		return v;
	}
{{end}}
	
	function encodeUTF8(s: string) {
		var i = 0, bytes = new Uint8Array(s.length * 4);
		for (var ci = 0; ci != s.length; ci++) {
			var c = s.charCodeAt(ci);
			if (c < 128) {
				bytes[i++] = c;
				continue;
			}
			if (c < 2048) {
				bytes[i++] = c >> 6 | 192;
			} else {
				if (c > 0xd7ff && c < 0xdc00) {
					if (++ci >= s.length) {
						bytes[i++] = 63;
						continue;
					}
					var c2 = s.charCodeAt(ci);
					if (c2 < 0xdc00 || c2 > 0xdfff) {
						bytes[i++] = 63;
						--ci;
						continue;
					}
					c = 0x10000 + ((c & 0x03ff) << 10) + (c2 & 0x03ff);
					bytes[i++] = c >> 18 | 240;
					bytes[i++] = c >> 12 & 63 | 128;
				} else bytes[i++] = c >> 12 | 224;
				bytes[i++] = c >> 6 & 63 | 128;
			}
			bytes[i++] = c & 63 | 128;
		}
		return bytes.subarray(0, i);
	}
	
	function decodeUTF8(bytes: Uint8Array) {
		var i = 0, s = '';
		while (i < bytes.length) {
			var c = bytes[i++];
			if (c > 127) {
				if (c > 191 && c < 224) {
					c = (i >= bytes.length) ? 63 : (c & 31) << 6 | bytes[i++] & 63;
				} else if (c > 223 && c < 240) {
					c = (i + 1 >= bytes.length) ? 63 : (c & 15) << 12 | (bytes[i++] & 63) << 6 | bytes[i++] & 63;
				} else if (c > 239 && c < 248) {
					c = (i + 2 >= bytes.length) ? 63 : (c & 7) << 18 | (bytes[i++] & 63) << 12 | (bytes[i++] & 63) << 6 | bytes[i++] & 63;
				} else c = 63
			}
			if (c <= 0xffff) s += String.fromCharCode(c);
			else if (c > 0x10ffff) s += '?';
			else {
				c -= 0x10000;
				s += String.fromCharCode(c >> 10 | 0xd800)
				s += String.fromCharCode(c & 0x3FF | 0xdc00)
			}
		}
		return s;
	}

}

{{end}}
`

const tsMarshal = `
		// Serializes the object into an Uint8Array.
{{- range .Fields}}{{if .TypeList}}{{if eq .Type "float32" "float64"}}{{else}}
		// All null entries in property {{.NameNative}} will be replaced with {{if eq .Type "text"}}an empty String{{else if eq .Type "binary"}}an empty Array{{else}}a new {{.TypeRef.Pkg.NameNative}}.{{.TypeRef.NameNative}}{{end}}.
{{- end}}{{end}}{{end}}
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);

{{range .Fields}}{{if eq .Type "bool"}}
			if (this.{{.NameNative}})
				buf[i++] = {{.Index}};
{{else if eq .Type "uint8"}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} > 255 || this.{{.NameNative}} < 0)
					throw new Error('colfer: {{.String}} out of reach: ' + this.{{.NameNative}});
				buf[i++] = {{.Index}};
				buf[i++] = this.{{.NameNative}};
			}
{{else if eq .Type "uint16"}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} > 65535 || this.{{.NameNative}} < 0)
					throw new Error('colfer: {{.String}} out of reach: ' + this.{{.NameNative}});
				if (this.{{.NameNative}} < 256) {
					buf[i++] = {{.Index}} | 128;
					buf[i++] = this.{{.NameNative}};
				} else {
					buf[i++] = {{.Index}};
					buf[i++] = this.{{.NameNative}} >>> 0;
					buf[i++] = this.{{.NameNative}} & 255;
				}
			}
{{else if eq .Type "uint32"}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} > 4294967295 || this.{{.NameNative}} < 0)
					throw new Error('colfer: {{.String}} out of reach: ' + this.{{.NameNative}});
				if (this.{{.NameNative}} < 0x200000) {
					buf[i++] = {{.Index}};
					i = encodeVarint(buf, i, this.{{.NameNative}});
				} else {
					buf[i++] = {{.Index}} | 128;
					view.setUint32(i, this.{{.NameNative}});
					i += 4;
				}
			}
{{else if eq .Type "uint64"}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} < 0)
					throw new Error('colfer: {{.String}} out of reach: ' + this.{{.NameNative}});
				if (this.{{.NameNative}} > Number.MAX_SAFE_INTEGER)
					throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				if (this.{{.NameNative}} < 0x2000000000000) {
					buf[i++] = {{.Index}};
					i = encodeVarint(buf, i, this.{{.NameNative}});
				} else {
					buf[i++] = {{.Index}} | 128;
					view.setUint32(i, this.{{.NameNative}} / 0x100000000);
					i += 4;
					view.setUint32(i, this.{{.NameNative}} % 0x100000000);
					i += 4;
				}
			}
{{else if eq .Type "int32"}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} < 0) {
					buf[i++] = {{.Index}} | 128;
					if (this.{{.NameNative}} < -2147483648)
						throw new Error('colfer: {{.String}} exceeds 32-bit range');
					i = encodeVarint(buf, i, -this.{{.NameNative}});
				} else {
					buf[i++] = {{.Index}}; 
					if (this.{{.NameNative}} > 2147483647)
						throw new Error('colfer: {{.String}} exceeds 32-bit range');
					i = encodeVarint(buf, i, this.{{.NameNative}});
				}
			}
{{else if eq .Type "int64"}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} < 0) {
					buf[i++] = {{.Index}} | 128;
					if (this.{{.NameNative}} < Number.MIN_SAFE_INTEGER)
						throw new Error('colfer: {{.String}} exceeds Number.MIN_SAFE_INTEGER');
					i = encodeVarint(buf, i, -this.{{.NameNative}});
				} else {
					buf[i++] = {{.Index}}; 
					if (this.{{.NameNative}} > Number.MAX_SAFE_INTEGER)
						throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
					i = encodeVarint(buf, i, this.{{.NameNative}});
				}
			}
{{else if eq .Type "float32"}}
 {{- if .TypeList}}
			if (this.{{.NameNative}} && this.{{.NameNative}}.length) {
				var a32 = this.{{.NameNative}};
				if (a32.length > colferListMax)
					throw new Error('colfer: {{.String}} exceeds colferListMax');
				buf[i++] = {{.Index}};
				i = encodeVarint(buf, i, a32.length);
				a32.forEach(function(f, fi) {
					if (f > 3.4028234663852886E38 || f < -3.4028234663852886E38)
						throw new Error('colfer: {{.String}}[' + fi + '] exceeds 32-bit range');
					view.setFloat32(i, f);
					i += 4;
				});
			}
 {{- else}}
			if (this.{{.NameNative}}) {
				if (this.{{.NameNative}} > 3.4028234663852886E38 || this.{{.NameNative}} < -3.4028234663852886E38)
					throw new Error('colfer: {{.String}} exceeds 32-bit range');
				buf[i++] = {{.Index}};
				view.setFloat32(i, this.{{.NameNative}});
				i += 4;
			} else if (Number.isNaN(this.{{.NameNative}})) {
				buf.set([{{.Index}}, 0x7f, 0xc0, 0, 0], i);
				i += 5;
			}
 {{- end}}
{{else if eq .Type "float64"}}
 {{- if .TypeList}}
			if (this.{{.NameNative}} && this.{{.NameNative}}.length) {
				var a64 = this.{{.NameNative}};
				if (a64.length > colferListMax)
					throw new Error('colfer: {{.String}} exceeds colferListMax');
				buf[i++] = {{.Index}};
				i = encodeVarint(buf, i, a64.length);
				a64.forEach(function(f) {
					view.setFloat64(i, f);
					i += 8;
				});
			}
 {{- else}}
			if (this.{{.NameNative}}) {
				buf[i++] = {{.Index}};
				view.setFloat64(i, this.{{.NameNative}});
				i += 8;
			} else if (Number.isNaN(this.{{.NameNative}})) {
				buf.set([{{.Index}}, 0x7f, 0xf8, 0, 0, 0, 0, 0, 0], i);
				i += 9;
			}
 {{- end}}
{{else if eq .Type "timestamp"}}
			if ((this.{{.NameNative}} && this.{{.NameNative}}.getTime()) || this.{{.NameNative}}_ns) {
				var ms = this.{{.NameNative}} ? this.{{.NameNative}}.getTime() : 0;
				var s = ms / 1E3;
	
				var ns = this.{{.NameNative}}_ns || 0;
				if (ns < 0 || ns >= 1E6)
					throw new Error('colfer: {{.String}} ns not in range (0, 1ms>');
				var msf = ms % 1E3;
				if (ms < 0 && msf) {
					s--
					msf = 1E3 + msf;
				}
				ns += msf * 1E6;
	
				if (s > 0xffffffff || s < 0) {
					buf[i++] = {{.Index}} | 128;
					if (s > 0) {
						view.setUint32(i, s / 0x100000000);
						view.setUint32(i + 4, s);
					} else {
						s = -s;
						view.setUint32(i, s / 0x100000000);
						view.setUint32(i + 4, s);
						let carry = 1;
						for (let j = i + 7; j >= i; j--) {
							const by = (buf[j] ^ 255) + carry;
							buf[j] = by & 255;
							carry = by >> 8;
						}
					}
					view.setUint32(i + 8, ns);
					i += 12;
				} else {
					buf[i++] = {{.Index}};
					view.setUint32(i, s);
					i += 4;
					view.setUint32(i, ns);
					i += 4;
				}
			}
{{else if eq .Type "text"}}
 {{- if .TypeList}}
			if (this.{{.NameNative}} && this.{{.NameNative}}.length) {
				var at = this.{{.NameNative}};
				if (at.length > colferListMax)
					throw new Error('colfer: {{.String}} exceeds colferListMax');
				buf[i++] = {{.Index}};
				i = encodeVarint(buf, i, at.length);
	
				at.forEach(function(s, si) {
					if (s == null) {
						s = "";
						at[si] = s;
					}
					var utf8 = encodeUTF8(s);
					i = encodeVarint(buf, i, utf8.length);
					buf.set(utf8, i);
					i += utf8.length;
				});
			}
 {{- else}}
			if (this.{{.NameNative}}) {
				buf[i++] = {{.Index}};
				var utf8 = encodeUTF8(this.{{.NameNative}});
				i = encodeVarint(buf, i, utf8.length);
				buf.set(utf8, i);
				i += utf8.length;
			}
 {{- end}}
{{else if eq .Type "binary"}}
 {{- if .TypeList}}
			if (this.{{.NameNative}} && this.{{.NameNative}}.length) {
				var ab = this.{{.NameNative}};
				if (ab.length > colferListMax)
					throw new Error('colfer: {{.String}} exceeds colferListMax');
				buf[i++] = {{.Index}};
				i = encodeVarint(buf, i, ab.length);
				ab.forEach(function(ba, bi) {
					if (ba == null) {
						(ba as any) = "";
						ab[bi] = ba;
					}
					i = encodeVarint(buf, i, ba.length);
					buf.set(ba, i);
					i += ba.length;
				});
			}
 {{- else}}
			if (this.{{.NameNative}} && this.{{.NameNative}}.length) {
				buf[i++] = {{.Index}};
				var bn = this.{{.NameNative}};
				i = encodeVarint(buf, i, bn.length);
				buf.set(bn, i);
				i += bn.length;
			}
 {{- end}}
{{else if .TypeList}}
			if (this.{{.NameNative}} && this.{{.NameNative}}.length) {
				var al = this.{{.NameNative}};
				if (al.length > colferListMax)
					throw new Error('colfer: {{.String}} exceeds colferListMax');
				buf[i++] = {{.Index}};
				i = encodeVarint(buf, i, al.length);
				al.forEach(function(v, vi) {
					if (v == null) {
						v = new {{.TypeRef.Pkg.NameNative}}.{{.TypeRef.NameNative}}();
						al[vi] = v;
					}
					var bi = v.marshal();
					buf.set(bi, i);
					i += bi.length;
				});
			}
{{else}}
			if (this.{{.NameNative}}) {
				buf[i++] = {{.Index}};
				var ba = this.{{.NameNative}}.marshal();
				buf.set(ba, i);
				i += ba.length;
			}
{{end}}{{end}}
	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: {{.String}} serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return buf.subarray(0, i);
		}`

const tsUnmarshal = `
		// Deserializes the object from an Uint8Array and returns the number of bytes read.
		public unmarshal(data: Uint8Array): number {
			if (!data || ! data.length) throw new Error(EOF);
			var header = data[0];
			var i = 1;
			var readHeader = function() {
				if (i >= data.length) throw new Error(EOF);
				header = data[i++];
			}
	
			var view = new DataView(data.buffer, data.byteOffset, data.byteLength);
	
			var readVarint = function() {
				var pos = 0, result = 0;
				while (pos != 8) {
					var c = data[i+pos];
					result += (c & 127) * Math.pow(128, pos);
					++pos;
					if (c < 128) {
						i += pos;
						if (result > Number.MAX_SAFE_INTEGER) break;
						return result;
					}
					if (pos == data.length) throw new Error(EOF);
				}
				return -1;
			}
{{range .Fields}}{{if eq .Type "bool"}}
			if (header == {{.Index}}) {
				this.{{.NameNative}} = true;
				readHeader();
			}
{{else if eq .Type "uint8"}}
			if (header == {{.Index}}) {
				if (i + 1 >= data.length) throw new Error(EOF);
				this.{{.NameNative}} = data[i++];
				header = data[i++];
			}
{{else if eq .Type "uint16"}}
			if (header == {{.Index}}) {
				if (i + 2 >= data.length) throw new Error(EOF);
				this.{{.NameNative}} = (data[i++] << 8) | data[i++];
				header = data[i++];
			} else if (header == ({{.Index}} | 128)) {
				if (i + 1 >= data.length) throw new Error(EOF);
				this.{{.NameNative}} = data[i++];
				header = data[i++];
			}
{{else if eq .Type "uint32"}}
			if (header == {{.Index}}) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = x;
				readHeader();
			} else if (header == ({{.Index}} | 128)) {
				if (i + 4 > data.length) throw new Error(EOF);
				this.{{.NameNative}} = view.getUint32(i);
				i += 4;
				readHeader();
			}
{{else if eq .Type "uint64"}}
			if (header == {{.Index}}) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = x;
				readHeader();
			} else if (header == ({{.Index}} | 128)) {
				if (i + 8 > data.length) throw new Error(EOF);
				var x = view.getUint32(i) * 0x100000000;
				x += view.getUint32(i + 4);
				if (x > Number.MAX_SAFE_INTEGER)
					throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = x;
				i += 8;
				readHeader();
			}
{{else if eq .Type "int32"}}
			if (header == {{.Index}}) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = x;
				readHeader();
			} else if (header == ({{.Index}} | 128)) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = -1 * x;
				readHeader();
			}
{{else if eq .Type "int64"}}
			if (header == {{.Index}}) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = x;
				readHeader();
			} else if (header == ({{.Index}} | 128)) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: {{.String}} exceeds Number.MAX_SAFE_INTEGER');
				this.{{.NameNative}} = -1 * x;
				readHeader();
			}
{{else if eq .Type "float32"}}
			if (header == {{.Index}}) {
 {{- if .TypeList}}
				var l = readVarint();
				if (l < 0) throw new Error('colfer: {{.String}} length exceeds Number.MAX_SAFE_INTEGER');
				if (l > colferListMax)
					throw new Error('colfer: {{.String}} length ' + l + ' exceeds ' + colferListMax + ' elements');
				if (i + l * 4 > data.length) throw new Error(EOF);
	
				this.{{.NameNative}} = new Float32Array(l);
				for (var n = 0; n < l; ++n) {
					this.{{.NameNative}}[n] = view.getFloat32(i);
					i += 4;
				}
 {{- else}}
				if (i + 4 > data.length) throw new Error(EOF);
				this.{{.NameNative}} = view.getFloat32(i);
				i += 4;
 {{- end}}
				readHeader();
				}
{{else if eq .Type "float64"}}
			if (header == {{.Index}}) {
 {{- if .TypeList}}
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: {{.String}} length ' + l + ' exceeds ' + colferListMax + ' elements');
				if (i + l * 8 > data.length) throw new Error(EOF);
	
				this.{{.NameNative}} = new Float64Array(l);
				for (var n = 0; n < l; ++n) {
					this.{{.NameNative}}[n] = view.getFloat64(i);
					i += 8;
				}
 {{- else}}
				if (i + 8 > data.length) throw new Error(EOF);
				this.{{.NameNative}} = view.getFloat64(i);
				i += 8;
 {{- end}}
				readHeader();
			}
{{else if eq .Type "timestamp"}}
			if (header == {{.Index}}) {
				if (i + 8 > data.length) throw new Error(EOF);
	
				var ms = view.getUint32(i) * 1E3;
				var ns = view.getUint32(i + 4);
				ms += Math.floor(ns / 1E6);
				this.{{.NameNative}} = new Date(ms);
				this.{{.NameNative}}_ns = ns % 1E6;
	
				i += 8;
				readHeader();
			} else if (header == ({{.Index}} | 128)) {
				if (i + 12 > data.length) throw new Error(EOF);
	
				var ms = decodeInt64(data, i) * 1E3;
				var ns = view.getUint32(i + 8);
				ms += Math.floor(ns / 1E6);
				if (ms < -864E13 || ms > 864E13)
					throw new Error('colfer: {{.String}} exceeds ECMA Date range');
				this.{{.NameNative}} = new Date(ms);
				this.{{.NameNative}}_ns = ns % 1E6;
	
				i += 12;
				readHeader();
			}
{{else if eq .Type "text"}}
			if (header == {{.Index}}) {
 {{- if .TypeList}}
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: {{.String}} length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				this.{{.NameNative}} = new Array(l);
				for (var n = 0; n < l; ++n) {
					var size = readVarint();
					if (size < 0 || size > colferSizeMax)
						throw new Error('colfer: {{.String}}[' + this.{{.NameNative}}.length + '] size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
					var start = i;
					i += size;
					if (i > data.length) throw new Error(EOF);
					this.{{.NameNative}}[n] = decodeUTF8(data.subarray(start, i));
				}
 {{- else}}
				var size = readVarint();
				if (size < 0 || size > colferSizeMax)
					throw new Error('colfer: {{.String}} size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
				var start = i;
				i += size;
				if (i > data.length) throw new Error(EOF);
				this.{{.NameNative}} = decodeUTF8(data.subarray(start, i));
 {{- end}}
				readHeader();
			}
{{else if eq .Type "binary"}}
			if (header == {{.Index}}) {
 {{- if .TypeList}}
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: {{.String}} length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				this.{{.NameNative}} = new Array(l);
				for (var n = 0; n < l; ++n) {
					var size = readVarint();
					if (size < 0 || size > colferSizeMax)
						throw new Error('colfer: {{.String}}[' + this.{{.NameNative}}.length + '] size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
					var start = i;
					i += size;
					if (i > data.length) throw new Error(EOF);
					this.{{.NameNative}}[n] = data.slice(start, i);
				}
 {{- else}}
				var size = readVarint();
				if (size < 0 || size > colferSizeMax)
					throw new Error('colfer: {{.String}} size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
				var start = i;
				i += size;
				if (i > data.length) throw new Error(EOF);
				this.{{.NameNative}} = data.slice(start, i);
 {{- end}}
				readHeader();
			}
{{else if .TypeList}}
			if (header == {{.Index}}) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: {{.String}} length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				for (var n = 0; n < l; ++n) {
					var on = new {{.TypeRef.Pkg.NameNative}}.{{.TypeRef.NameNative}}();
					i += on.unmarshal(data.subarray(i));
					this.{{.NameNative}}[n] = on;
				}
				readHeader();
			}
{{else}}
			if (header == {{.Index}}) {
				var oh = new {{.TypeRef.Pkg.NameNative}}.{{.TypeRef.NameNative}}();
				i += oh.unmarshal(data.subarray(i));
				this.{{.NameNative}} = oh;
				readHeader();
			}
{{end}}{{end}}
			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: {{.String}} serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return i;
		}`
