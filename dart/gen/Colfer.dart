// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf for package gen.

import 'dart:convert';
import 'dart:typed_data';

/// Package gen tests all field mapping options.

/// The upper limit for serial byte sizes.
const colferSizeMax = 16 * 1024 * 1024;

/// The upper limit for the number of elements in a list.
const colferListMax = 64 * 1024;

/// O contains all supported data types.
class O {
  O({
    this.b = false,
    this.u32 = 0,
    this.u64 = 0,
    this.i32 = 0,
    this.i64 = 0,
    this.f32 = 0.0,
    this.f64 = 0.0,
    this.t,
    this.s = '',
    Uint8List? a,
    this.o,
    List<O>? os,
    List<String>? ss,
    List<Uint8List>? as_0,
    this.u8 = 0,
    this.u16 = 0,
    Float32List? f32s,
    Float64List? f64s,
  })  : a = a ?? Uint8List(0),
        os = os ?? [],
        ss = ss ?? [],
        as_0 = as_0 ?? [],
        f32s = f32s ?? Float32List(0),
        f64s = f64s ?? Float64List(0);

  /// B tests booleans.
  bool b;

  /// U32 tests unsigned 32-bit integers.
  int u32;

  /// U64 tests unsigned 64-bit integers.
  int u64;

  /// I32 tests signed 32-bit integers.
  int i32;

  /// I64 tests signed 64-bit integers.
  int i64;

  /// F32 tests 32-bit floating points.
  double f32;

  /// F64 tests 64-bit floating points.
  double f64;

  /// T tests timestamps.
  DateTime? t;

  /// S tests text.
  String s;

  /// A tests binaries.
  Uint8List a;

  /// O tests nested data structures.
  O? o;

  /// Os tests data structure lists.
  List<O?> os;

  /// Ss tests text lists.
  List<String> ss;

  /// As tests binary lists.
  List<Uint8List> as_0;

  /// U8 tests unsigned 8-bit integers.
  int u8;

  /// U16 tests unsigned 16-bit integers.
  int u16;

  /// F32s tests 32-bit floating point lists.
  Float32List f32s;

  /// F64s tests 64-bit floating point lists.
  Float64List f64s;

  /// Returns an over estimatation of marshal length.
  ///
  /// Throws [RangeError] if the size of a list exceeds [colferListMax].
  /// Returns an over estimated length for the required buffer. String
  /// characters are counted for 4 bytes, everything has its exact size.
  int marshalLen() {
    int _l = 1;
    {
      if (b) {
        _l++;
      }
    }
    {
      int _x = u32;
      if (_x >= 1 << 21) {
        _l += 5;
      } else if (_x != 0) {
        for (_l += 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    {
      int _x = u64;
      if (_x < 0 || 0x2000000000000 <= _x) {
        _l += 9;
      } else if (_x != 0) {
        for (_l += 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    {
      int _x = i32;
      if (_x != 0) {
        if (_x < 0) {
          _x = -_x;
        }
        for (_l += 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    {
      int _x = i64;
      if (_x != 0) {
        if (-_x == _x) {
          _l += 10;
        } else {
          _l += 2;
          if (_x < 0) {
            _x = -_x;
          }
          for (; _x >= 0x80; _l++) {
            _x >>= 7;
          }
        }
      }
    }
    {
      if (f32 != 0) {
        _l += 5;
      }
    }
    {
      if (f64 != 0) {
        _l += 9;
      }
    }
    {
      DateTime? _v = t;
      if (_v != null) {
        int _us = _v.microsecondsSinceEpoch;
        int _s = _us ~/ 1E6;
        if (_s >= 1 << 33 || _us < 0) {
          _l += 13;
        } else {
          _l += 9;
        }
      }
    }
    {
      int _x = s.length;
      if (_x != 0) {
        _x *= 4;
        for (_l += _x + 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    {
      int _x = a.length;
      if (_x != 0) {
        for (_l += _x + 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    {
      if (o != null) {
        _l += o!.marshalLen() + 1;
      }
    }
    {
      int _x = os.length;
      if (_x != 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.os', 'colfer');
        }
        for (_l += 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
        for (final _v in os) {
          if (_v == null) {
            _l++;
            continue;
          }
          _l += _v.marshalLen();
        }
      }
    }
    {
      int _x = ss.length;
      if (_x != 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.ss', 'colfer');
        }
        for (_l += 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
        for (final _a in ss) {
          _x = _a.length;
          _x *= 4;
          for (_l += _x + 1; _x >= 0x80; _l++) {
            _x >>= 7;
          }
        }
      }
    }
    {
      int _x = as_0.length;
      if (_x != 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.as', 'colfer');
        }
        for (_l += 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
        for (final _a in as_0) {
          _x = _a.length;
          for (_l += _x + 1; _x >= 0x80; _l++) {
            _x >>= 7;
          }
        }
      }
    }
    {
      if (u8 != 0) {
        _l += 2;
      }
    }
    {
      if (u16 >= 1 << 8) {
        _l += 3;
      } else if (u16 != 0) {
        _l += 2;
      }
    }
    {
      int _x = f32s.length;
      if (_x != 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.f32s', 'colfer');
        }
        for (_l += 2 + _x * 4; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    {
      int _x = f64s.length;
      if (_x != 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.f64s', 'colfer');
        }
        for (_l += 2 + _x * 8; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    if (_l > colferSizeMax) {
      return colferSizeMax;
    }
    return _l;
  }

  /// Encodes as Colfer into [_buf].
  ///
  /// Throws [RangeError] if uint8, uint16, uint32 or int32 value overflows or
  /// underflows, or when the size of a list exceeds [colferListMax], or if a
  /// text, binary, or [_buf] exceeds [colferSizeMax]. Returns the number of
  /// bytes written.
  int marshalTo(Uint8List _buf) {
    var _view = _buf.buffer.asByteData(_buf.offsetInBytes);
    int _i = 0;
    {
      if (b) {
        _buf[_i] = 0;
        _i++;
      }
    }
    {
      int _x = u32;
      if (_x != 0) {
        if (_x > 4294967295 || _x < 0) {
          throw RangeError.range(_x, 0, 4294967295, 'gen.o.u32', 'colfer');
        }
        if (_x < 0x200000) {
          _buf[_i] = 1;
          _i++;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _i++;
        } else {
          _buf[_i] = 1 | 128;
          _view.setInt32(_i + 1, _x);
          _i += 5;
        }
      }
    }
    {
      int _x = u64;
      if (_x != 0) {
        if (0 < _x && _x < 0x2000000000000) {
          _buf[_i] = 2;
          _i++;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _i++;
        } else {
          _buf[_i] = 2 | 128;
          _view.setInt64(_i + 1, _x);
          _i += 9;
        }
      }
    }
    {
      int _x = i32;
      if (_x != 0) {
        if (_x < 0) {
          if (_x < -2147483648) {
            throw RangeError.range(_x, -2147483648, null, 'gen.o.i32', 'colfer');
          }
          _buf[_i] = 3 | 128;
          _i++;
          _x = -_x;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _i++;
        } else {
          if (_x > 2147483647) {
            throw RangeError.range(_x, null, 2147483647, 'gen.o.i32', 'colfer');
          }
          _buf[_i] = 3;
          _i++;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _i++;
        }
      }
    }
    {
      int _x = i64;
      if (_x != 0) {
        if (_x < 0) {
          _buf[_i] = 4 | 128;
          _i++;
          _x = -_x;
          if (i64 == _x) {
            _buf.fillRange(_i, 10, 128);
            _i += 9;
          } else {
            while (_x > 127) {
              _buf[_i] = (_x & 127) | 128;
              _i++;
              _x >>= 7;
            }
            _buf[_i] = _x;
            _i++;
          }
        } else {
          _buf[_i] = 4;
          _i++;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _i++;
        }
      }
    }
    {
      if (f32 != 0) {
        _buf[_i] = 5;
        if (f32.isNaN) {
          _buf[_i + 1] = 0x7f;
          _buf[_i + 2] = 0xc0;
        } else {
          _view.setFloat32(_i + 1, f32);
        }
        _i += 5;
      }
    }
    {
      if (f64 != 0) {
        _buf[_i] = 6;
        if (f64.isNaN) {
          _buf[_i + 1] = 0x7f;
          _buf[_i + 2] = 0xf8;
          _buf[_i + 8] = 1;
        } else {
          _view.setFloat64(_i + 1, f64);
        }
        _i += 9;
      }
    }
    {
      if (t != null) {
        int _us = t!.microsecondsSinceEpoch;
        int _res = _us % 1000000;
        _us -= _res;
        int _s = _us ~/ 1E6;
        int _ns = _res * 1000;

        if (_s >= 1 << 33 || _us < 0) {
          _buf[_i] = 7 | 128;
          _view.setInt64(_i + 1, _s);
          _i += 9;
        } else {
          _buf[_i] = 7;
          _view.setInt32(_i + 1, _s);
          _i += 5;
        }
        _view.setInt32(_i, _ns);
        _i += 4;
      }
    }
    {
      int _x = s.length;
      if (_x > 0) {
        _buf[_i] = 8;
        _i++;
        var _v = utf8.encode(s);
        _x = _v.length;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _buf.setAll(_i + 1, _v);
        _i += 1 + _v.length;
      }
    }
    {
      int _x = a.length;
      if (_x > 0) {
        _buf[_i] = 9;
        _i++;
        var _v = a;
        _x = _v.length;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _buf.setAll(_i + 1, _v);
        _i += 1 + _v.length;
      }
    }
    {
      if (o != null) {
        _buf[_i] = 10;
        _i++;
        _i += o!.marshalTo(Uint8List.view(_buf.buffer, _buf.offsetInBytes + _i));
      }
    }
    {
      int _x = os.length;
      if (_x > 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.os', 'colfer');
        }
        _buf[_i] = 11;
        _i++;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _i++;
        for (var _vi in os) {
          _vi ??= O();
          _i += _vi.marshalTo(Uint8List.view(_buf.buffer, _buf.offsetInBytes + _i));
        }
      }
    }
    {
      int _x = ss.length;
      if (_x > 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.ss', 'colfer');
        }
        _buf[_i] = 12;
        _i++;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _i++;
        for (final _vi in ss) {
          var _v = utf8.encode(_vi);
          _x = _v.length;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _buf.setAll(_i + 1, _v);
          _i += 1 + _v.length;
        }
      }
    }
    {
      int _x = as_0.length;
      if (_x > 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.as', 'colfer');
        }
        _buf[_i] = 13;
        _i++;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _i++;
        for (final _vi in as_0) {
          var _v = _vi;
          _x = _v.length;
          while (_x > 127) {
            _buf[_i] = (_x & 127) | 128;
            _i++;
            _x >>= 7;
          }
          _buf[_i] = _x;
          _buf.setAll(_i + 1, _v);
          _i += 1 + _v.length;
        }
      }
    }
    {
      if (u8 != 0) {
        if (u8 > 255 || u8 < 0) {
          throw RangeError.range(u8, 0, 255, 'gen.o.u8', 'colfer');
        }
        _buf[_i] = 14;
        _buf[_i + 1] = u8;
        _i += 2;
      }
    }
    {
      if (u16 != 0) {
        if (u16 > 65535 || u16 < 0) {
          throw RangeError.range(u16, 0, 65535, 'gen.o.u16', 'colfer');
        }
        if (u16 < 256) {
          _buf[_i] = 15 | 128;
          _buf[_i + 1] = u16;
          _i += 2;
        } else {
          _buf[_i] = 15;
          _buf[_i + 1] = u16 >> 8;
          _buf[_i + 2] = u16;
          _i += 3;
        }
      }
    }
    {
      int _x = f32s.length;
      if (_x > 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.f32s', 'colfer');
        }
        _buf[_i] = 16;
        _i++;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _i++;
        for (final _vi in f32s) {
          if (_vi.isNaN) {
            _buf[_i] = 0x7f;
            _buf[_i + 1] = 0xc0;
          } else {
            _view.setFloat32(_i, _vi);
          }
          _i += 4;
        }
      }
    }
    {
      int _x = f64s.length;
      if (_x > 0) {
        if (_x > colferListMax) {
          throw RangeError.range(_x, null, colferListMax, 'gen.o.f64s', 'colfer');
        }
        _buf[_i] = 17;
        _i++;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _i++;
        for (final _vi in f64s) {
          if (_vi.isNaN) {
            _buf[_i] = 0x7f;
            _buf[_i + 1] = 0xf8;
            _buf[_i + 7] = 1;
          } else {
            _view.setFloat64(_i, _vi);
          }
          _i += 8;
        }
      }
    }

    _buf[_i] = 127;
    _i++;
    if (_i > colferSizeMax) {
      throw RangeError.range(_i, null, colferSizeMax, 'gen.o', 'colfer');
    }
    return _i;
  }

  /// Decodes [_data] as Colfer.
  ///
  /// Throws [RangeError] if there is an unexpexted end of data, if a list
  /// exceeds [colferListMax], or if a text, binary or [_data] exceeds
  /// [colferSizeMax]. Throws [StateError] if ending header mismatches.
  /// Returns the number of bytes read.
  int unmarshal(Uint8List _data) {
    int _header = 0;
    int _i = 0;
    var _view = ByteData.view(_data.buffer);
    _header = _data[_i];
    _i++;

    if (_header == 0) {
      b = true;
      _header = _data[_i];
      _i++;
    }

    if (_header == 1) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      u32 = _c;
      _header = _data[_i];
      _i++;
    } else if (_header == (1 | 128)) {
      u32 = _view.getUint32(_i);
      _header = _data[_i + 4];
      _i += 5;
    }

    if (_header == 2) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      u64 = _c;
      _header = _data[_i];
      _i++;
    } else if (_header == (2 | 128)) {
      int _v = _view.getUint32(_i) * 0x100000000;
      _v += _view.getUint32(_i + 4);
      u64 = _v;
      _header = _data[_i + 8];
      _i += 9;
    }

    if (_header == 3) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      i32 = _c;
      _header = _data[_i];
      _i++;
    } else if (_header == (3 | 128)) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      i32 = -1 * _c;
      _header = _data[_i];
      _i++;
    }

    if (_header == 4) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      i64 = _c;
      _header = _data[_i];
      _i++;
    } else if (_header == (4 | 128)) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      i64 = -1 * _c;
      _header = _data[_i];
      _i++;
    }

    if (_header == 5) {
      f32 = _view.getFloat32(_i);
      _i += 4;
      _header = _data[_i];
      _i++;
    }

    if (_header == 6) {
      f64 = _view.getFloat64(_i);
      _i += 8;
      _header = _data[_i];
      _i++;
    }

    if (_header == 7) {
      int _s = _view.getUint32(_i);
      int _us = _view.getUint32(_i + 4) ~/ 1000;
      t = DateTime.fromMicrosecondsSinceEpoch(_s * 1000000 + _us);
      _i += 8;
      _header = _data[_i];
      _i++;
    } else if (_header == (7 | 128)) {
      int _s = _view.getInt64(_i);
      int _us = _view.getUint32(_i + 8) ~/ 1000;
      t = DateTime.fromMicrosecondsSinceEpoch(_s * 1000000 + _us);
      _i += 12;
      _header = _data[_i];
      _i++;
    }

    if (_header == 8) {
      int _size = _data[_i];
      _i++;
      if (_size >= 0x80) {
        _size &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _size |= _b << _shift;
            break;
          }
          _size |= (_b & 0x7f) << _shift;
        }
      }
      if (_size < 0 || _size > colferSizeMax) {
        throw RangeError.range(_size, 0, colferSizeMax, 'gen.o.s', 'colfer');
      }

      int _s = _i;
      _i += _size;
      s = utf8.decode(_data.sublist(_s, _i));
      _header = _data[_i];
      _i++;
    }

    if (_header == 9) {
      int _size = _data[_i];
      _i++;
      if (_size >= 0x80) {
        _size &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _size |= _b << _shift;
            break;
          }
          _size |= (_b & 0x7f) << _shift;
        }
      }
      if (_size < 0 || _size > colferSizeMax) {
        throw RangeError.range(_size, 0, colferSizeMax, 'gen.o.a', 'colfer');
      }

      int _start = _i;
      _i += _size;
      a = _data.sublist(_start, _i);
      _header = _data[_i];
      _i++;
    }

    if (_header == 10) {
      var _s = O();
      _i += _s.unmarshal(_data.sublist(_i));
      o = _s;
      _header = _data[_i];
      _i++;
    }

    if (_header == 11) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      if (_c < 0 || _c > colferListMax) {
        throw RangeError.range(_c, 0, colferListMax, 'gen.o.os', 'colfer');
      }

      if (os.length != _c) {
        os = List<O>.filled(_c, O());
      }
      for (var _ci in os) {
        _ci ??= O();
        _i += _ci.unmarshal(_data.sublist(_i));
      }
      _header = _data[_i];
      _i++;
    }

    if (_header == 12) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      if (_c < 0 || _c > colferListMax) {
        throw RangeError.range(_c, 0, colferListMax, 'gen.o.ss', 'colfer');
      }

      if (ss.length != _c) {
        ss = List<String>.filled(_c, '');
      }
      for (int _ci = 0; _ci < _c; _ci++) {
        int _size = _data[_i];
        _i++;
        if (_size >= 0x80) {
          _size &= 0x7f;
          for (int _shift = 7;; _shift += 7) {
            int _b = _data[_i];
            _i++;
            if (_b < 0x80 || _shift == 56) {
              _size |= _b << _shift;
              break;
            }
            _size |= (_b & 0x7f) << _shift;
          }
        }
        if (_size < 0 || _size > colferSizeMax) {
          throw RangeError.range(_size, 0, colferSizeMax, 'gen.o.ss', 'colfer');
        }

        int _s = _i;
        _i += _size;
        ss[_ci] = utf8.decode(_data.sublist(_s, _i));
      }
      _header = _data[_i];
      _i++;
    }

    if (_header == 13) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      if (_c < 0 || _c > colferListMax) {
        throw RangeError.range(_c, 0, colferListMax, 'gen.o.as', 'colfer');
      }

      if (as_0.length != _c) {
        as_0 = List<Uint8List>.filled(_c, Uint8List(0));
      }
      for (int _ci = 0; _ci < _c; _ci++) {
        int _size = _data[_i];
        _i++;
        if (_size >= 0x80) {
          _size &= 0x7f;
          for (int _shift = 7;; _shift += 7) {
            int _b = _data[_i];
            _i++;
            if (_b < 0x80 || _shift == 56) {
              _size |= _b << _shift;
              break;
            }
            _size |= (_b & 0x7f) << _shift;
          }
        }
        if (_size < 0 || _size > colferSizeMax) {
          throw RangeError.range(_size, 0, colferSizeMax, 'gen.o.as', 'colfer');
        }

        int _s = _i;
        _i += _size;
        as_0[_ci] = _data.sublist(_s, _i);
      }
      _header = _data[_i];
      _i++;
    }

    if (_header == 14) {
      u8 = _data[_i];
      _header = _data[_i + 1];
      _i += 2;
    }

    if (_header == 15) {
      u16 = (_data[_i] << 8) | _data[_i + 1];
      _header = _data[_i + 2];
      _i += 3;
    } else if (_header == (15 | 128)) {
      u16 = _data[_i];
      _header = _data[_i + 1];
      _i += 2;
    }

    if (_header == 16) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      if (_c < 0 || _c > colferListMax) {
        throw RangeError.range(_c, 0, colferListMax, 'gen.o.f32s', 'colfer');
      }

      if (f32s.length != _c) {
        f32s = Float32List(_c);
      }
      for (int _ci = 0; _ci < _c; _ci++) {
        f32s[_ci] = _view.getFloat32(_i);
        _i += 4;
      }
      _header = _data[_i];
      _i++;
    }

    if (_header == 17) {
      int _c = _data[_i];
      _i++;
      if (_c >= 0x80) {
        _c &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _c |= _b << _shift;
            break;
          }
          _c |= (_b & 0x7f) << _shift;
        }
      }
      if (_c < 0 || _c > colferListMax) {
        throw RangeError.range(_c, 0, colferListMax, 'gen.o.f64s', 'colfer');
      }

      if (f64s.length != _c) {
        f64s = Float64List(_c);
      }
      for (int _ci = 0; _ci < _c; _ci++) {
        f64s[_ci] = _view.getFloat64(_i);
        _i += 8;
      }
      _header = _data[_i];
      _i++;
    }

    if (_header != 127) {
      throw StateError('colfer: unknown header $_header at byte ${_i - 1}');
    }
    if (_i > colferSizeMax) {
      throw RangeError.range(_i, null, colferSizeMax, 'gen.o', 'colfer');
    }
    return _i;
  }

  @override
  bool operator ==(_other) {
    if (_other is! O ||
        b != _other.b ||
        u32 != _other.u32 ||
        u64 != _other.u64 ||
        i32 != _other.i32 ||
        i64 != _other.i64 ||
        f32 != _other.f32 ||
        f64 != _other.f64 ||
        t != _other.t ||
        s != _other.s ||
        a.length != _other.a.length ||
        o != _other.o ||
        os.length != _other.os.length ||
        ss.length != _other.ss.length ||
        as_0.length != _other.as_0.length ||
        u8 != _other.u8 ||
        u16 != _other.u16 ||
        f32s.length != _other.f32s.length ||
        f64s.length != _other.f64s.length) return false;
    for (int _i = 0; _i < a.length; _i++) if (a[_i] != _other.a[_i]) return false;
    for (int _i = 0; _i < os.length; _i++) if (os[_i] != _other.os[_i]) return false;
    for (int _i = 0; _i < ss.length; _i++) if (ss[_i] != _other.ss[_i]) return false;
    for (int _i = 0; _i < as_0.length; _i++) {
      var _l1 = as_0[_i];
      var _l2 = _other.as_0[_i];
      if (_l1.length != _l2.length) return false;
      for (int _i = 0; _i < _l1.length; _i++) if (_l1[_i] != _l2[_i]) return false;
    }
    for (int _i = 0; _i < f32s.length; _i++) if (f32s[_i] != _other.f32s[_i]) return false;
    for (int _i = 0; _i < f64s.length; _i++) if (f64s[_i] != _other.f64s[_i]) return false;
    return true;
  }

  @override
  int get hashCode {
    int _h = 0;
    _h = 31 * _h + b.hashCode;
    _h = 31 * _h + u32.hashCode;
    _h = 31 * _h + u64.hashCode;
    _h = 31 * _h + i32.hashCode;
    _h = 31 * _h + i64.hashCode;
    _h = 31 * _h + f32.hashCode;
    _h = 31 * _h + f64.hashCode;
    _h = 31 * _h + t.hashCode;
    _h = 31 * _h + s.hashCode;
    _h = 31 * _h + a.hashCode;
    _h = 31 * _h + o.hashCode;
    _h = 31 * _h + os.length;
    for (var _e in os) _h = 31 * _h + _e.hashCode;
    _h = 31 * _h + ss.length;
    for (var _e in ss) _h = 31 * _h + _e.hashCode;
    _h = 31 * _h + as_0.length;
    for (var _e in as_0) _h = 31 * _h + _e.hashCode;
    _h = 31 * _h + u8.hashCode;
    _h = 31 * _h + u16.hashCode;
    _h = 31 * _h + f32s.length;
    for (var _e in f32s) _h = 31 * _h + _e.hashCode;
    _h = 31 * _h + f64s.length;
    for (var _e in f64s) _h = 31 * _h + _e.hashCode;
    return _h;
  }

  @override
  String toString() => 'class O {'
      'b: ${b.toString()}'
      ', u32: ${u32.toString()}'
      ', u64: ${u64.toString()}'
      ', i32: ${i32.toString()}'
      ', i64: ${i64.toString()}'
      ', f32: ${f32.toString()}'
      ', f64: ${f64.toString()}'
      ', t: ${t.toString()}'
      ', s: "$s"'
      ', a: ${a.toString()}'
      ', o: ${o.toString()}'
      ', os: List<O>${os.toString()}'
      ', ss: [${ss.isNotEmpty ? "\"ss.join('\", \"')}" : ""}]'
      ', as_0: List<Uint8List>${as_0.toString()}'
      ', u8: ${u8.toString()}'
      ', u16: ${u16.toString()}'
      ', f32s: List<Float32List>${f32s.toString()}'
      ', f64s: List<Float64List>${f64s.toString()}'
      '}';
}

/// DromedaryCase oposes name casings.
class DromedaryCase {
  DromedaryCase({
    this.pascalCase = '',
  });

  String pascalCase;

  /// Returns an over estimatation of marshal length.
  ///
  /// Throws [RangeError] if the size of a list exceeds [colferListMax].
  /// Returns an over estimated length for the required buffer. String
  /// characters are counted for 4 bytes, everything has its exact size.
  int marshalLen() {
    int _l = 1;
    {
      int _x = pascalCase.length;
      if (_x != 0) {
        _x *= 4;
        for (_l += _x + 2; _x >= 0x80; _l++) {
          _x >>= 7;
        }
      }
    }
    if (_l > colferSizeMax) {
      return colferSizeMax;
    }
    return _l;
  }

  /// Encodes as Colfer into [_buf].
  ///
  /// Throws [RangeError] if uint8, uint16, uint32 or int32 value overflows or
  /// underflows, or when the size of a list exceeds [colferListMax], or if a
  /// text, binary, or [_buf] exceeds [colferSizeMax]. Returns the number of
  /// bytes written.
  int marshalTo(Uint8List _buf) {
    int _i = 0;
    {
      int _x = pascalCase.length;
      if (_x > 0) {
        _buf[_i] = 0;
        _i++;
        var _v = utf8.encode(pascalCase);
        _x = _v.length;
        while (_x > 127) {
          _buf[_i] = (_x & 127) | 128;
          _i++;
          _x >>= 7;
        }
        _buf[_i] = _x;
        _buf.setAll(_i + 1, _v);
        _i += 1 + _v.length;
      }
    }

    _buf[_i] = 127;
    _i++;
    if (_i > colferSizeMax) {
      throw RangeError.range(_i, null, colferSizeMax, 'gen.dromedaryCase', 'colfer');
    }
    return _i;
  }

  /// Decodes [_data] as Colfer.
  ///
  /// Throws [RangeError] if there is an unexpexted end of data, if a list
  /// exceeds [colferListMax], or if a text, binary or [_data] exceeds
  /// [colferSizeMax]. Throws [StateError] if ending header mismatches.
  /// Returns the number of bytes read.
  int unmarshal(Uint8List _data) {
    int _header = 0;
    int _i = 0;
    _header = _data[_i];
    _i++;

    if (_header == 0) {
      int _size = _data[_i];
      _i++;
      if (_size >= 0x80) {
        _size &= 0x7f;
        for (int _shift = 7;; _shift += 7) {
          int _b = _data[_i];
          _i++;
          if (_b < 0x80 || _shift == 56) {
            _size |= _b << _shift;
            break;
          }
          _size |= (_b & 0x7f) << _shift;
        }
      }
      if (_size < 0 || _size > colferSizeMax) {
        throw RangeError.range(_size, 0, colferSizeMax, 'gen.dromedaryCase.PascalCase', 'colfer');
      }

      int _s = _i;
      _i += _size;
      pascalCase = utf8.decode(_data.sublist(_s, _i));
      _header = _data[_i];
      _i++;
    }

    if (_header != 127) {
      throw StateError('colfer: unknown header $_header at byte ${_i - 1}');
    }
    if (_i > colferSizeMax) {
      throw RangeError.range(_i, null, colferSizeMax, 'gen.dromedaryCase', 'colfer');
    }
    return _i;
  }

  @override
  bool operator ==(_other) {
    if (_other is! DromedaryCase || pascalCase != _other.pascalCase) return false;
    return true;
  }

  @override
  int get hashCode {
    int _h = 0;
    _h = 31 * _h + pascalCase.hashCode;
    return _h;
  }

  @override
  String toString() => 'class DromedaryCase {'
      'pascalCase: "$pascalCase"'
      '}';
}

/// EmbedO has an inner object only.
/// Covers regression of issue #66.
class EmbedO {
  EmbedO({
    this.inner,
  });

  O? inner;

  /// Returns an over estimatation of marshal length.
  ///
  /// Throws [RangeError] if the size of a list exceeds [colferListMax].
  /// Returns an over estimated length for the required buffer. String
  /// characters are counted for 4 bytes, everything has its exact size.
  int marshalLen() {
    int _l = 1;
    {
      if (inner != null) {
        _l += inner!.marshalLen() + 1;
      }
    }
    if (_l > colferSizeMax) {
      return colferSizeMax;
    }
    return _l;
  }

  /// Encodes as Colfer into [_buf].
  ///
  /// Throws [RangeError] if uint8, uint16, uint32 or int32 value overflows or
  /// underflows, or when the size of a list exceeds [colferListMax], or if a
  /// text, binary, or [_buf] exceeds [colferSizeMax]. Returns the number of
  /// bytes written.
  int marshalTo(Uint8List _buf) {
    int _i = 0;
    {
      if (inner != null) {
        _buf[_i] = 0;
        _i++;
        _i += inner!.marshalTo(Uint8List.view(_buf.buffer, _buf.offsetInBytes + _i));
      }
    }

    _buf[_i] = 127;
    _i++;
    if (_i > colferSizeMax) {
      throw RangeError.range(_i, null, colferSizeMax, 'gen.EmbedO', 'colfer');
    }
    return _i;
  }

  /// Decodes [_data] as Colfer.
  ///
  /// Throws [RangeError] if there is an unexpexted end of data, if a list
  /// exceeds [colferListMax], or if a text, binary or [_data] exceeds
  /// [colferSizeMax]. Throws [StateError] if ending header mismatches.
  /// Returns the number of bytes read.
  int unmarshal(Uint8List _data) {
    int _header = 0;
    int _i = 0;
    _header = _data[_i];
    _i++;

    if (_header == 0) {
      var _s = O();
      _i += _s.unmarshal(_data.sublist(_i));
      inner = _s;
      _header = _data[_i];
      _i++;
    }

    if (_header != 127) {
      throw StateError('colfer: unknown header $_header at byte ${_i - 1}');
    }
    if (_i > colferSizeMax) {
      throw RangeError.range(_i, null, colferSizeMax, 'gen.EmbedO', 'colfer');
    }
    return _i;
  }

  @override
  bool operator ==(_other) {
    if (_other is! EmbedO || inner != _other.inner) return false;
    return true;
  }

  @override
  int get hashCode {
    int _h = 0;
    _h = 31 * _h + inner.hashCode;
    return _h;
  }

  @override
  String toString() => 'class EmbedO {'
      'inner: ${inner.toString()}'
      '}';
}
