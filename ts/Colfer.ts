// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf for package gen.

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

// Package gen tests all field mapping options.
export namespace gen {

	// The upper limit for serial byte sizes.
	var colferSizeMax = 16 * 1024 * 1024;
	// The upper limit for the number of elements in a list.
	var colferListMax = 64 * 1024;

	// Constructor.
	// O contains all supported data types.
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class O {
		// B tests booleans.
		b: bool = false;
		// U32 tests unsigned 32-bit integers.
		u32: uint32 = 0;
		// U64 tests unsigned 64-bit integers.
		u64: uint64 = 0;
		// I32 tests signed 32-bit integers.
		i32: int32 = 0;
		// I64 tests signed 64-bit integers.
		i64: int64 = 0;
		// F32 tests 32-bit floating points.
		f32: float32 = 0;
		// F64 tests 64-bit floating points.
		f64: float64 = 0;
		// T tests timestamps.
		t: timestamp = undefined;
		t_ns: number = 0;
		// S tests text.
		s: text = undefined;;
		// A tests binaries.
		a: binary = new Uint8Array(0);
		// O tests nested data structures.
		o: gen.O | undefined = undefined;
		// Os tests data structure lists.
		os:  Array<gen.O> = [];
		// Ss tests text lists.
		ss:  Array<text> = [];
		// As tests binary lists.
		as:  Array<binary> = [];
		// U8 tests unsigned 8-bit integers.
		u8: uint8 = 0;
		// U16 tests unsigned 16-bit integers.
		u16: uint16 = 0;
		// F32s tests 32-bit floating point lists.
		f32s:  Float32Array = new Float32Array(0);
		// F64s tests 64-bit floating point lists.
		f64s:  Float64Array = new Float64Array(0);

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		
		// Serializes the object into an Uint8Array.
		// All null entries in property os will be replaced with a new gen.O.
		// All null entries in property ss will be replaced with an empty String.
		// All null entries in property as will be replaced with an empty Array.
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);


			if (this.b)
				buf[i++] = 0;

			if (this.u32) {
				if (this.u32 > 4294967295 || this.u32 < 0)
					throw new Error('colfer: gen.o.u32 out of reach: ' + this.u32);
				if (this.u32 < 0x200000) {
					buf[i++] = 1;
					i = encodeVarint(buf, i, this.u32);
				} else {
					buf[i++] = 1 | 128;
					view.setUint32(i, this.u32);
					i += 4;
				}
			}

			if (this.u64) {
				if (this.u64 < 0)
					throw new Error('colfer: gen.o.u64 out of reach: ' + this.u64);
				if (this.u64 > Number.MAX_SAFE_INTEGER)
					throw new Error('colfer: gen.o.u64 exceeds Number.MAX_SAFE_INTEGER');
				if (this.u64 < 0x2000000000000) {
					buf[i++] = 2;
					i = encodeVarint(buf, i, this.u64);
				} else {
					buf[i++] = 2 | 128;
					view.setUint32(i, this.u64 / 0x100000000);
					i += 4;
					view.setUint32(i, this.u64 % 0x100000000);
					i += 4;
				}
			}

			if (this.i32) {
				if (this.i32 < 0) {
					buf[i++] = 3 | 128;
					if (this.i32 < -2147483648)
						throw new Error('colfer: gen.o.i32 exceeds 32-bit range');
					i = encodeVarint(buf, i, -this.i32);
				} else {
					buf[i++] = 3; 
					if (this.i32 > 2147483647)
						throw new Error('colfer: gen.o.i32 exceeds 32-bit range');
					i = encodeVarint(buf, i, this.i32);
				}
			}

			if (this.i64) {
				if (this.i64 < 0) {
					buf[i++] = 4 | 128;
					if (this.i64 < Number.MIN_SAFE_INTEGER)
						throw new Error('colfer: gen.o.i64 exceeds Number.MIN_SAFE_INTEGER');
					i = encodeVarint(buf, i, -this.i64);
				} else {
					buf[i++] = 4; 
					if (this.i64 > Number.MAX_SAFE_INTEGER)
						throw new Error('colfer: gen.o.i64 exceeds Number.MAX_SAFE_INTEGER');
					i = encodeVarint(buf, i, this.i64);
				}
			}

			if (this.f32) {
				if (this.f32 > 3.4028234663852886E38 || this.f32 < -3.4028234663852886E38)
					throw new Error('colfer: gen.o.f32 exceeds 32-bit range');
				buf[i++] = 5;
				view.setFloat32(i, this.f32);
				i += 4;
			} else if (Number.isNaN(this.f32)) {
				buf.set([5, 0x7f, 0xc0, 0, 0], i);
				i += 5;
			}

			if (this.f64) {
				buf[i++] = 6;
				view.setFloat64(i, this.f64);
				i += 8;
			} else if (Number.isNaN(this.f64)) {
				buf.set([6, 0x7f, 0xf8, 0, 0, 0, 0, 0, 0], i);
				i += 9;
			}

			if ((this.t && this.t.getTime()) || this.t_ns) {
				var ms = this.t ? this.t.getTime() : 0;
				var s = ms / 1E3;
	
				var ns = this.t_ns || 0;
				if (ns < 0 || ns >= 1E6)
					throw new Error('colfer: gen.o.t ns not in range (0, 1ms>');
				var msf = ms % 1E3;
				if (ms < 0 && msf) {
					s--
					msf = 1E3 + msf;
				}
				ns += msf * 1E6;
	
				if (s > 0xffffffff || s < 0) {
					buf[i++] = 7 | 128;
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
					buf[i++] = 7;
					view.setUint32(i, s);
					i += 4;
					view.setUint32(i, ns);
					i += 4;
				}
			}

			if (this.s) {
				buf[i++] = 8;
				var utf8 = encodeUTF8(this.s);
				i = encodeVarint(buf, i, utf8.length);
				buf.set(utf8, i);
				i += utf8.length;
			}

			if (this.a && this.a.length) {
				buf[i++] = 9;
				var bn = this.a;
				i = encodeVarint(buf, i, bn.length);
				buf.set(bn, i);
				i += bn.length;
			}

			if (this.o) {
				buf[i++] = 10;
				var ba = this.o.marshal();
				buf.set(ba, i);
				i += ba.length;
			}

			if (this.os && this.os.length) {
				var al = this.os;
				if (al.length > colferListMax)
					throw new Error('colfer: gen.o.os exceeds colferListMax');
				buf[i++] = 11;
				i = encodeVarint(buf, i, al.length);
				al.forEach(function(v, vi) {
					if (v == null) {
						v = new gen.O();
						al[vi] = v;
					}
					var bi = v.marshal();
					buf.set(bi, i);
					i += bi.length;
				});
			}

			if (this.ss && this.ss.length) {
				var at = this.ss;
				if (at.length > colferListMax)
					throw new Error('colfer: gen.o.ss exceeds colferListMax');
				buf[i++] = 12;
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

			if (this.as && this.as.length) {
				var ab = this.as;
				if (ab.length > colferListMax)
					throw new Error('colfer: gen.o.as exceeds colferListMax');
				buf[i++] = 13;
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

			if (this.u8) {
				if (this.u8 > 255 || this.u8 < 0)
					throw new Error('colfer: gen.o.u8 out of reach: ' + this.u8);
				buf[i++] = 14;
				buf[i++] = this.u8;
			}

			if (this.u16) {
				if (this.u16 > 65535 || this.u16 < 0)
					throw new Error('colfer: gen.o.u16 out of reach: ' + this.u16);
				if (this.u16 < 256) {
					buf[i++] = 15 | 128;
					buf[i++] = this.u16;
				} else {
					buf[i++] = 15;
					buf[i++] = this.u16 >>> 0;
					buf[i++] = this.u16 & 255;
				}
			}

			if (this.f32s && this.f32s.length) {
				var a32 = this.f32s;
				if (a32.length > colferListMax)
					throw new Error('colfer: gen.o.f32s exceeds colferListMax');
				buf[i++] = 16;
				i = encodeVarint(buf, i, a32.length);
				a32.forEach(function(f, fi) {
					if (f > 3.4028234663852886E38 || f < -3.4028234663852886E38)
						throw new Error('colfer: gen.o.f32s[' + fi + '] exceeds 32-bit range');
					view.setFloat32(i, f);
					i += 4;
				});
			}

			if (this.f64s && this.f64s.length) {
				var a64 = this.f64s;
				if (a64.length > colferListMax)
					throw new Error('colfer: gen.o.f64s exceeds colferListMax');
				buf[i++] = 17;
				i = encodeVarint(buf, i, a64.length);
				a64.forEach(function(f) {
					view.setFloat64(i, f);
					i += 8;
				});
			}

	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: gen.o serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return buf.subarray(0, i);
		}
		
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

			if (header == 0) {
				this.b = true;
				readHeader();
			}

			if (header == 1) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: gen.o.u32 exceeds Number.MAX_SAFE_INTEGER');
				this.u32 = x;
				readHeader();
			} else if (header == (1 | 128)) {
				if (i + 4 > data.length) throw new Error(EOF);
				this.u32 = view.getUint32(i);
				i += 4;
				readHeader();
			}

			if (header == 2) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: gen.o.u64 exceeds Number.MAX_SAFE_INTEGER');
				this.u64 = x;
				readHeader();
			} else if (header == (2 | 128)) {
				if (i + 8 > data.length) throw new Error(EOF);
				var x = view.getUint32(i) * 0x100000000;
				x += view.getUint32(i + 4);
				if (x > Number.MAX_SAFE_INTEGER)
					throw new Error('colfer: gen.o.u64 exceeds Number.MAX_SAFE_INTEGER');
				this.u64 = x;
				i += 8;
				readHeader();
			}

			if (header == 3) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: gen.o.i32 exceeds Number.MAX_SAFE_INTEGER');
				this.i32 = x;
				readHeader();
			} else if (header == (3 | 128)) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: gen.o.i32 exceeds Number.MAX_SAFE_INTEGER');
				this.i32 = -1 * x;
				readHeader();
			}

			if (header == 4) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: gen.o.i64 exceeds Number.MAX_SAFE_INTEGER');
				this.i64 = x;
				readHeader();
			} else if (header == (4 | 128)) {
				var x = readVarint();
				if (x < 0) throw new Error('colfer: gen.o.i64 exceeds Number.MAX_SAFE_INTEGER');
				this.i64 = -1 * x;
				readHeader();
			}

			if (header == 5) {
				if (i + 4 > data.length) throw new Error(EOF);
				this.f32 = view.getFloat32(i);
				i += 4;
				readHeader();
				}

			if (header == 6) {
				if (i + 8 > data.length) throw new Error(EOF);
				this.f64 = view.getFloat64(i);
				i += 8;
				readHeader();
			}

			if (header == 7) {
				if (i + 8 > data.length) throw new Error(EOF);
	
				var ms = view.getUint32(i) * 1E3;
				var ns = view.getUint32(i + 4);
				ms += Math.floor(ns / 1E6);
				this.t = new Date(ms);
				this.t_ns = ns % 1E6;
	
				i += 8;
				readHeader();
			} else if (header == (7 | 128)) {
				if (i + 12 > data.length) throw new Error(EOF);
	
				var ms = decodeInt64(data, i) * 1E3;
				var ns = view.getUint32(i + 8);
				ms += Math.floor(ns / 1E6);
				if (ms < -864E13 || ms > 864E13)
					throw new Error('colfer: gen.o.t exceeds ECMA Date range');
				this.t = new Date(ms);
				this.t_ns = ns % 1E6;
	
				i += 12;
				readHeader();
			}

			if (header == 8) {
				var size = readVarint();
				if (size < 0 || size > colferSizeMax)
					throw new Error('colfer: gen.o.s size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
				var start = i;
				i += size;
				if (i > data.length) throw new Error(EOF);
				this.s = decodeUTF8(data.subarray(start, i));
				readHeader();
			}

			if (header == 9) {
				var size = readVarint();
				if (size < 0 || size > colferSizeMax)
					throw new Error('colfer: gen.o.a size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
				var start = i;
				i += size;
				if (i > data.length) throw new Error(EOF);
				this.a = data.slice(start, i);
				readHeader();
			}

			if (header == 10) {
				var oh = new gen.O();
				i += oh.unmarshal(data.subarray(i));
				this.o = oh;
				readHeader();
			}

			if (header == 11) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: gen.o.os length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				for (var n = 0; n < l; ++n) {
					var on = new gen.O();
					i += on.unmarshal(data.subarray(i));
					this.os[n] = on;
				}
				readHeader();
			}

			if (header == 12) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: gen.o.ss length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				this.ss = new Array(l);
				for (var n = 0; n < l; ++n) {
					var size = readVarint();
					if (size < 0 || size > colferSizeMax)
						throw new Error('colfer: gen.o.ss[' + this.ss.length + '] size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
					var start = i;
					i += size;
					if (i > data.length) throw new Error(EOF);
					this.ss[n] = decodeUTF8(data.subarray(start, i));
				}
				readHeader();
			}

			if (header == 13) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: gen.o.as length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				this.as = new Array(l);
				for (var n = 0; n < l; ++n) {
					var size = readVarint();
					if (size < 0 || size > colferSizeMax)
						throw new Error('colfer: gen.o.as[' + this.as.length + '] size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
					var start = i;
					i += size;
					if (i > data.length) throw new Error(EOF);
					this.as[n] = data.slice(start, i);
				}
				readHeader();
			}

			if (header == 14) {
				if (i + 1 >= data.length) throw new Error(EOF);
				this.u8 = data[i++];
				header = data[i++];
			}

			if (header == 15) {
				if (i + 2 >= data.length) throw new Error(EOF);
				this.u16 = (data[i++] << 8) | data[i++];
				header = data[i++];
			} else if (header == (15 | 128)) {
				if (i + 1 >= data.length) throw new Error(EOF);
				this.u16 = data[i++];
				header = data[i++];
			}

			if (header == 16) {
				var l = readVarint();
				if (l < 0) throw new Error('colfer: gen.o.f32s length exceeds Number.MAX_SAFE_INTEGER');
				if (l > colferListMax)
					throw new Error('colfer: gen.o.f32s length ' + l + ' exceeds ' + colferListMax + ' elements');
				if (i + l * 4 > data.length) throw new Error(EOF);
	
				this.f32s = new Float32Array(l);
				for (var n = 0; n < l; ++n) {
					this.f32s[n] = view.getFloat32(i);
					i += 4;
				}
				readHeader();
				}

			if (header == 17) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: gen.o.f64s length ' + l + ' exceeds ' + colferListMax + ' elements');
				if (i + l * 8 > data.length) throw new Error(EOF);
	
				this.f64s = new Float64Array(l);
				for (var n = 0; n < l; ++n) {
					this.f64s[n] = view.getFloat64(i);
					i += 8;
				}
				readHeader();
			}

			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: gen.o serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return i;
		}
	}

	// Constructor.
	// DromedaryCase oposes name casings.
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class DromedaryCase {

		pascalCase: text = undefined;;

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		
		// Serializes the object into an Uint8Array.
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);


			if (this.pascalCase) {
				buf[i++] = 0;
				var utf8 = encodeUTF8(this.pascalCase);
				i = encodeVarint(buf, i, utf8.length);
				buf.set(utf8, i);
				i += utf8.length;
			}

	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: gen.dromedaryCase serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return buf.subarray(0, i);
		}
		
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

			if (header == 0) {
				var size = readVarint();
				if (size < 0 || size > colferSizeMax)
					throw new Error('colfer: gen.dromedaryCase.PascalCase size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
				var start = i;
				i += size;
				if (i > data.length) throw new Error(EOF);
				this.pascalCase = decodeUTF8(data.subarray(start, i));
				readHeader();
			}

			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: gen.dromedaryCase serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return i;
		}
	}

	// Constructor.
	// EmbedO has an inner object only.
	// Covers regression of issue #66.
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class EmbedO {

		inner: gen.O | undefined = undefined;

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		
		// Serializes the object into an Uint8Array.
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);


			if (this.inner) {
				buf[i++] = 0;
				var ba = this.inner.marshal();
				buf.set(ba, i);
				i += ba.length;
			}

	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: gen.EmbedO serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return buf.subarray(0, i);
		}
		
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

			if (header == 0) {
				var oh = new gen.O();
				i += oh.unmarshal(data.subarray(i));
				this.inner = oh;
				readHeader();
			}

			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: gen.EmbedO serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return i;
		}
	}


	// private section

	function encodeVarint(bytes: Uint8Array, i: number, x: number) {
		while (x > 127) {
			bytes[i++] = (x & 127) | 128;
			x /= 128;
		}
		bytes[i++] = x & 127;
		return i;
	}

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


