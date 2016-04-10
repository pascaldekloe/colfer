var testdata = new function() {
	const EOF = 'colfer: EOF';
	const StructMis = 'colfer: struct header mismatch';
	const Overflow = 'colfer: varint overflow';
	const UnknownField = 'colfer: unknown field header';

	this.unmarshalO = function(data) {
		if (!data || ! data.length) return null;
		var i = 0;
		if (data[i++] != 0x80) throw StructMis;

		var readVarint = function() {
			var pos = 0, result = 0;
			while (true) {
				var c = data[i+pos];
				result += (c&127) * Math.pow(128, pos);
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
			if (header&0x80) o.i32 *= -1;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 4) {
			o.i64 = readVarint();
			if (header&0x80) o.i64 *= -1;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 5) {
			if (data.length < (i+4)) throw EOF;
			o.f32 = new DataView(data.buffer).getFloat32(i);
			i += 4;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 6) {
			if (data.length < (i+8)) throw EOF;
			o.f64 = new DataView(data.buffer).getFloat64(i);
			i += 8;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 7) {
			if (data.length < (i+8)) throw EOF;
			var view = new DataView(data.buffer);
			// BUG(pascaldekloe): negative time offset not supported
			var ms = view.getUint32(i) * Math.pow(2, 32);
			ms += view.getUint32(i+4);
			ms *= 1000;
			i += 8;
			if (header&0x80) {
				if (data.length < (i+4)) throw EOF;
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
			o.s = decodeUTF8(data.slice(i, to));
			i = to;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		if (field == 9) {
			var to = readVarint() + i;
			if (to > data.length) throw EOF;
			o.a = data.slice(i, to);
			i = to;

			if (i == data.length) return o;
			header = data[i++];
			field = header & 128;
		}

		return UnknownField;
	}

	// Unmarshals Uint8Array to string.
	var decodeUTF8 = function(bytes) {
		var s = '';
		var i = 0;
		while (i < bytes.length) {
			var c = bytes[i++];
			if (c > 127) {
				if (c > 191 && c < 224) {
					if (i >= bytes.length) throw 'UTF-8 decode: incomplete 2-byte sequence';
					c = ((c&31)<<6) | (bytes[i]&63);
				} else if (c > 223 && c < 240) {
					if (i+1 >= bytes.length) throw 'UTF-8 decode: incomplete 3-byte sequence';
					c = ((c&15)<<12) | ((bytes[i]&63)<<6) | (bytes[++i]&63);
				} else if (c > 239 && c < 248) {
					if (i+2 >= bytes.length) throw 'UTF-8 decode: incomplete 4-byte sequence';
					c = ((c&7)<<18) | ((bytes[i]&63)<<12) | ((bytes[++i]&63)<<6) | (bytes[++i]&63);
				} else throw 'UTF-8 decode: unknown multibyte start 0x' + c.toString(16) + ' at index ' + (i-1);
				++i;
			}

			if (c <= 0xFFFF) s += String.fromCharCode(c);
			else {
				if (!String.fromCodePoint) throw 'UTF-8 decode: need code point support conform ECMAScript version 6';
				s += String.fromCodePoint(c);
			}
		}
		return s;
	}
}
