// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file break-refs.colf for package static.
// The compiler used schema file break.colf for package void.

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


export namespace static_ {

	// The upper limit for serial byte sizes.
	var colferSizeMax = 16 * 1024 * 1024;
	// The upper limit for the number of elements in a list.
	var colferListMax = 64 * 1024;

	// Constructor.
	// Int is a cross-package reference for void.class.
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class Int {

		try_:  Array<text> = [];

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		
		// Serializes the object into an Uint8Array.
		// All null entries in property try_ will be replaced with an empty String.
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);


			if (this.try_ && this.try_.length) {
				var at = this.try_;
				if (at.length > colferListMax)
					throw new Error('colfer: static.int.try exceeds colferListMax');
				buf[i++] = 0;
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

	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: static.int serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
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
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: static.int.try length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				this.try_ = new Array(l);
				for (var n = 0; n < l; ++n) {
					var size = readVarint();
					if (size < 0 || size > colferSizeMax)
						throw new Error('colfer: static.int.try[' + this.try_.length + '] size ' + size + ' exceeds ' + colferSizeMax + ' bytes');
	
					var start = i;
					i += size;
					if (i > data.length) throw new Error(EOF);
					this.try_[n] = decodeUTF8(data.subarray(start, i));
				}
				readHeader();
			}

			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: static.int serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
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


// Package void tries to break the generated code.
// Note that void is a reserved keyword in all supported languages except for Go.
export namespace void_ {

	// The upper limit for serial byte sizes.
	var colferSizeMax = 16 * 1024 * 1024;
	// The upper limit for the number of elements in a list.
	var colferListMax = 64 * 1024;

	// Constructor.
	// Class has local and cross-package refereces.
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class Class {

		extends_: void_.Int | undefined = undefined;

		public_:  Array<static_.Int> = [];

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		
		// Serializes the object into an Uint8Array.
		// All null entries in property public_ will be replaced with a new static_.Int.
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);


			if (this.extends_) {
				buf[i++] = 0;
				var ba = this.extends_.marshal();
				buf.set(ba, i);
				i += ba.length;
			}

			if (this.public_ && this.public_.length) {
				var al = this.public_;
				if (al.length > colferListMax)
					throw new Error('colfer: void.class.public exceeds colferListMax');
				buf[i++] = 1;
				i = encodeVarint(buf, i, al.length);
				al.forEach(function(v, vi) {
					if (v == null) {
						v = new static_.Int();
						al[vi] = v;
					}
					var bi = v.marshal();
					buf.set(bi, i);
					i += bi.length;
				});
			}

	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: void.class serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
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
				var oh = new void_.Int();
				i += oh.unmarshal(data.subarray(i));
				this.extends_ = oh;
				readHeader();
			}

			if (header == 1) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: void.class.public length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				for (var n = 0; n < l; ++n) {
					var on = new static_.Int();
					i += on.unmarshal(data.subarray(i));
					this.public_[n] = on;
				}
				readHeader();
			}

			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: void.class serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
			return i;
		}
	}

	// Constructor.
	// Int is a circular dependency.
	// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
	export class Int {

		throw_:  Array<void_.Class> = [];

		finally_:  Array<void_.Class> = [];

		// When init is provided all enumerable properties are merged into the new object a.k.a. shallow cloning.
		constructor(init: Record<string, any> = {}) {
			// @ts-ignore
			for (let p in init) this[p] = init[p];
		}

		
		// Serializes the object into an Uint8Array.
		// All null entries in property throw_ will be replaced with a new void_.Class.
		// All null entries in property finally_ will be replaced with a new void_.Class.
		public marshal(b?: Uint8Array): Uint8Array {
			const buf: Uint8Array = !b || !b.length ? new Uint8Array(colferSizeMax) : b;
			var i = 0;
			var view = new DataView(buf.buffer);


			if (this.throw_ && this.throw_.length) {
				var al = this.throw_;
				if (al.length > colferListMax)
					throw new Error('colfer: void.int.throw exceeds colferListMax');
				buf[i++] = 0;
				i = encodeVarint(buf, i, al.length);
				al.forEach(function(v, vi) {
					if (v == null) {
						v = new void_.Class();
						al[vi] = v;
					}
					var bi = v.marshal();
					buf.set(bi, i);
					i += bi.length;
				});
			}

			if (this.finally_ && this.finally_.length) {
				var al = this.finally_;
				if (al.length > colferListMax)
					throw new Error('colfer: void.int.finally exceeds colferListMax');
				buf[i++] = 1;
				i = encodeVarint(buf, i, al.length);
				al.forEach(function(v, vi) {
					if (v == null) {
						v = new void_.Class();
						al[vi] = v;
					}
					var bi = v.marshal();
					buf.set(bi, i);
					i += bi.length;
				});
			}

	
			buf[i++] = 127;
			if (i >= colferSizeMax)
				throw new Error('colfer: void.int serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
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
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: void.int.throw length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				for (var n = 0; n < l; ++n) {
					var on = new void_.Class();
					i += on.unmarshal(data.subarray(i));
					this.throw_[n] = on;
				}
				readHeader();
			}

			if (header == 1) {
				var l = readVarint();
				if (l < 0 || l > colferListMax)
					throw new Error('colfer: void.int.finally length ' + l + ' exceeds ' + colferListMax + ' elements');
	
				for (var n = 0; n < l; ++n) {
					var on = new void_.Class();
					i += on.unmarshal(data.subarray(i));
					this.finally_[n] = on;
				}
				readHeader();
			}

			if (header != 127) throw new Error('colfer: unknown header at byte ' + (i - 1));
			if (i > colferSizeMax)
				throw new Error('colfer: void.int serial size ' + i + ' exceeds ' + colferSizeMax + ' bytes');
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


