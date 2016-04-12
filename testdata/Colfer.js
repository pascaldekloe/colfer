var testdata = new function() {
	const EOF = 'colfer: EOF';
	const StructMis = 'colfer: struct header mismatch';
	const UnknownField = 'colfer: unknown field header';

	this.marshalO = function(o) {
		var segs = [[128]];
		if (o.b) {
			segs.push([0]);
		}

		if (o.u32) {
			if (o.u32 > 4294967295)
				throw 'colfer: field "u32" overflow: ' + o.u32;
			var seg = [1];
			encodeVarint(seg, o.u32);
			segs.push(seg);
		}

		if (o.u64) {
			if (o.u64 > 18446744073709551616)
				throw 'colfer: field "u64" overflow: ' + o.u64;
			var seg = [2];
			encodeVarint(seg, o.u64);
			segs.push(seg);
		}

		if (o.i32) {
			if (o.i32 > 2147483647 || o.i32 < -2147483648)
				throw 'colfer: field "i32" overflow: ' + o.i32;
			var seg = [3];
			if (o.i32 < 0) {
				seg[0] |= 128;
				encodeVarint(seg, -o.i32);
			} else	encodeVarint(seg, o.i32);
			segs.push(seg);
		}

		if (o.i64) {
			if (o.i64 > 9223372036854775807 || o.i64 < -9223372036854775808)
				throw 'colfer: field "i64" overflow: ' + o.i64;
			var seg = [4];
			if (o.i64 < 0) {
				seg[0] |= 128;
				encodeVarint(seg, -o.i64);
			} else	encodeVarint(seg, o.i64);
			segs.push(seg);
		}

		if (o.f32 || Number.isNaN(o.f32)) {
			if (o.f32 > 3.4028234663852886E38 || o.f32 < -3.4028234663852886E38)
				throw 'colfer: field "f32" overflow: ' + o.f32;
			var bytes = new Uint8Array(5);
			bytes[0] = 5;
			new DataView(bytes.buffer).setFloat32(1, o.f32);
			segs.push(bytes);
		}

		if (o.f64 || Number.isNaN(o.f64)) {
			var bytes = new Uint8Array(9);
			bytes[0] = 6;
			new DataView(bytes.buffer).setFloat64(1, o.f64);
			segs.push(bytes);
		}

		if (o.t) {
			var ms = o.t.getTime()
			var s = ms / 1000;
			var ns = (ms % 1000) * 1E6;
			if (o.t_ns) ns += o.t_ns % 1E6;

			var bytes = new Uint8Array((ns) ? 13 : 9);
			bytes[0] = 7;
			var view = new DataView(bytes.buffer);
			view.setUint32(1, s / Math.pow(2, 32));
			view.setUint32(5, s % Math.pow(2, 32));
			if (ns) {
				bytes[0] |= 128;
				view.setUint32(9, ns);
			}
			segs.push(bytes);
		}

		if (o.s) {
			var utf = encodeUTF8(o.s);
			var seg = [8];
			encodeVarint(seg, utf.length);
			segs.push(seg);
			segs.push(utf)
		}

		if (o.a) {
			var seg = [9];
			encodeVarint(seg, o.a.length);
			segs.push(seg);
			segs.push(o.a);
		}

		return joinSegs(segs);
	}

	this.unmarshalO = function(data) {
		if (!data || ! data.length) return null;
		var i = 0;
		if (data[i++] != 0x80) throw StructMis;

		var readVarint = function() {
			var pos = 0, result = 0;
			while (true) {
				var c = data[i+pos];
				result += (c & 127) * Math.pow(128, pos);
				++pos;
				if (c < 128) {
					i += pos;
					return result;
				}
				if (pos == data.length) throw EOF;
			}
		}

		var o = {};
		if (i == data.length) return o;

		var header = data[i++];
		var field = header & 127;

		if (field == 0) {
			o.b = true;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 127;
		}

		if (field == 1) {
			o.u32 = readVarint();

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 2) {
			o.u64 = readVarint();

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 3) {
			o.i32 = readVarint();
			if (header & 0x80) o.i32 *= -1;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 4) {
			o.i64 = readVarint();
			if (header & 0x80) o.i64 *= -1;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 5) {
			if (data.length < i + 4) throw EOF;
			o.f32 = new DataView(data.buffer).getFloat32(i);
			i += 4;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 6) {
			if (data.length < i + 8) throw EOF;
			o.f64 = new DataView(data.buffer).getFloat64(i);
			i += 8;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 7) {
			if (data.length < i + 8) throw EOF;
			var view = new DataView(data.buffer);
			// BUG(pascaldekloe): negative time offset not supported
			var ms = view.getUint32(i) * Math.pow(2, 32);
			ms += view.getUint32(i + 4);
			ms *= 1000;
			i += 8;
			if (header&0x80) {
				if (data.length < i + 4) throw EOF;
				var ns = view.getUint32(i);
				i += 4;
				ms += ns / 1E6;
				o.t_ns = ns % 1E6;
			}
			o.t = new Date();
			o.t.setTime(ms);

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 8) {
			var to = readVarint() + i;
			if (to > data.length) throw EOF;
			o.s = decodeUTF8(data.subarray(i, to));
			i = to;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 9) {
			var to = readVarint() + i;
			if (to > data.length) throw EOF;
			o.a = data.subarray(i, to);
			i = to;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		return UnknownField;
	}

	var joinSegs = function(segs) {
		var size = 0;
		segs.forEach(function(seg) {
			size += seg.length;
		});

		var data = new Uint8Array(size);
		var i = 0;
		segs.forEach(function(seg) {
			data.set(seg, i);
			i += seg.length;
		});
		return data;
	}

	var encodeVarint = function(bytes, x) {
		while (x > 127) {
			bytes.push(x|128);
			x /= 128;
		}
		bytes.push(x&127);
		return bytes;
	}

	// Marshals a string to Uint8Array.
	var encodeUTF8 = function(s) {
		var i = 0;
		var bytes = new Uint8Array(s.length * 4);
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
					if (++ci == s.length) throw 'UTF-8 encode: incomplete surrogate pair';
					var c2 = s.charCodeAt(ci);
					if (c2 < 0xdc00 || c2 > 0xdfff) throw 'UTF-8 encode: second char code 0x' + c2.toString(16) + ' at index ' + ci + ' in surrogate pair out of range';
					c = 0x10000 + ((c & 0x03ff) << 10) + (c2 & 0x03ff);
					bytes[i++] = c >> 18 | 240;
					bytes[i++] = c>> 12 & 63 | 128;
				} else { // c <= 0xffff
					bytes[i++] = c >> 12 | 224;
				}
				bytes[i++] = c >> 6 & 63 | 128;
			}
			bytes[i++] = c & 63 | 128;
		}
		return bytes.subarray(0, i);
	}

	// Unmarshals an Uint8Array to string.
	var decodeUTF8 = function(bytes) {
		var s = '';
		var i = 0;
		while (i < bytes.length) {
			var c = bytes[i++];
			if (c > 127) {
				if (c > 191 && c < 224) {
					if (i >= bytes.length) throw 'UTF-8 decode: incomplete 2-byte sequence';
					c = (c & 31) << 6 | bytes[i] & 63;
				} else if (c > 223 && c < 240) {
					if (i + 1 >= bytes.length) throw 'UTF-8 decode: incomplete 3-byte sequence';
					c = (c & 15) << 12 | (bytes[i] & 63) << 6 | bytes[++i] & 63;
				} else if (c > 239 && c < 248) {
					if (i+2 >= bytes.length) throw 'UTF-8 decode: incomplete 4-byte sequence';
					c = (c & 7) << 18 | (bytes[i] & 63) << 12 | (bytes[++i] & 63) << 6 | bytes[++i] & 63;
				} else throw 'UTF-8 decode: unknown multibyte start 0x' + c.toString(16) + ' at index ' + (i - 1);
				++i;
			}

			if (c <= 0xffff) s += String.fromCharCode(c);
			else if (c <= 0x10ffff) {
				c -= 0x10000;
				s += String.fromCharCode(c >> 10 | 0xd800)
				s += String.fromCharCode(c & 0x3FF | 0xdc00)
			} else throw 'UTF-8 decode: code point 0x' + c.toString(16) + ' exceeds UTF-16 reach';
		}
		return s;
	}
}
