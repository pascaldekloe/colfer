// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf for package gen.
// Package gen tests all field mapping options.

import 'dart:typed_data';
import 'package:collection/collection.dart';
import 'package:quiver/core.dart';
import 'dart:convert';

/// The upper limit for serial byte sizes.
final colferSizeMax = 16 * 1024 * 1024;

/// The upper limit for the number of elements in a list.
final colferListMax = 64 * 1024;

void encodeVarint(BytesBuilder buf, int x) {
  while (x > 127) {
    buf.addByte((x & 127) | 128);
    x >>= 7;
  }
  buf.addByte(x);
}

int decodeInt64(Uint8List data, int i) {
  int v = 0, j = i + 7, m = 1;
  if (data[i] & 128 != 0) {
    // two's complement
    for (int carry = 1; j >= i; j--, m *= 256) {
      int b = (data[j] ^ 255) + carry;
      carry = b >> 8;
      v += (b & 255) * m;
    }
    v = -v;
  } else {
    for (; j >= i; j--, m *= 256) v += (data[j] * m).toInt();
  }
  return v;
}

// O contains all supported data types.
class O {
  // B tests booleans.
  bool b;
  // U32 tests unsigned 32-bit integers.
  int u32;
  // U64 tests unsigned 64-bit integers.
  int u64;
  // I32 tests signed 32-bit integers.
  int i32;
  // I64 tests signed 64-bit integers.
  int i64;
  // F32 tests 32-bit floating points.
  double f32;
  // F64 tests 64-bit floating points.
  double f64;
  // T tests timestamps.
  DateTime? t;
  // S tests text.
  String s_;
  // A tests binaries.
  Uint8List a_;
  // O tests nested data structures.
  O? o;
  // Os tests data structure lists.
  List<O> os;
  // Ss tests text lists.
  List<String> ss;
  // As tests binary lists.
  List<Uint8List> as_;
  // U8 tests unsigned 8-bit integers.
  int u8;
  // U16 tests unsigned 16-bit integers.
  int u16;
  // F32s tests 32-bit floating point lists.
  Float32List f32s;
  // F64s tests 64-bit floating point lists.
  Float64List f64s;

  @override
  bool operator ==(other) {
    if (!(other is O)) {
      print('compare type doesn\'t match $this != $other');
      return false;
    }
    if (other.b != b) {
      print('compare prop b doesn\'t match $b != ${other.b}');
      return false;
    }
    if (other.u32 != u32) {
      print('compare prop u32 doesn\'t match $u32 != ${other.u32}');
      return false;
    }
    if (other.u64 != u64) {
      print('compare prop u64 doesn\'t match $u64 != ${other.u64}');
      return false;
    }
    if (other.i32 != i32) {
      print('compare prop i32 doesn\'t match $i32 != ${other.i32}');
      return false;
    }
    if (other.i64 != i64) {
      print('compare prop i64 doesn\'t match $i64 != ${other.i64}');
      return false;
    }
    if (other.f32 != f32) {
      print('compare prop f32 doesn\'t match $f32 != ${other.f32}');
      return false;
    }
    if (other.f64 != f64) {
      print('compare prop f64 doesn\'t match $f64 != ${other.f64}');
      return false;
    }
    if (other.t != t) {
      print('compare prop t doesn\'t match $t != ${other.t}');
      return false;
    }
    if (other.s_ != s_) {
      print('compare prop s_ doesn\'t match $s_ != ${other.s_}');
      return false;
    }
    if (!(IterableEquality().equals(other.a_, a_))) {
      print('compare prop a_ doesn\'t match $a_ != ${other.a_}');
      return false;
    }
    if (other.o != o) {
      print('compare prop o doesn\'t match $o != ${other.o}');
      return false;
    }
    if (!(IterableEquality().equals(other.os, os))) {
      print('compare prop os doesn\'t match $os != ${other.os}');
      return false;
    }
    if (!(IterableEquality().equals(other.ss, ss))) {
      print('compare prop ss doesn\'t match $ss != ${other.ss}');
      return false;
    }
    if (other.as_.length != as_.length) {
      print('compare prop length of as_ doesn\'t match ${as_.length} != ${other.as_.length}');
      return false;
    }
    for (int vi = 0; vi < as_.length; vi++) {
      if (!(IterableEquality().equals(other.as_[vi], as_[vi]))) {
        print('compare prop as_[$vi] doesn\'t match ${as_[vi]} != ${other.as_[vi]}');
        return false;
      }
    }
    if (other.u8 != u8) {
      print('compare prop u8 doesn\'t match $u8 != ${other.u8}');
      return false;
    }
    if (other.u16 != u16) {
      print('compare prop u16 doesn\'t match $u16 != ${other.u16}');
      return false;
    }
    if (!(IterableEquality().equals(other.f32s, f32s))) {
      print('compare prop f32s doesn\'t match $f32s != ${other.f32s}');
      return false;
    }
    if (!(IterableEquality().equals(other.f64s, f64s))) {
      print('compare prop f64s doesn\'t match $f64s != ${other.f64s}');
      return false;
    }
    return true;
  }

  @override
  int get hashCode => hashObjects(
      [b, u32, u64, i32, i64, f32, f64, t, s_, a_, o, os, ss, as_, u8, u16, f32s, f64s]);

  O({
    this.b = false,
    this.u32 = 0,
    this.u64 = 0,
    this.i32 = 0,
    this.i64 = 0,
    this.f32 = 0.0,
    this.f64 = 0.0,
    this.t,
    this.s_ = '',
    Uint8List? a_,
    this.o,
    List<O>? os,
    List<String>? ss,
    List<Uint8List>? as_,
    this.u8 = 0,
    this.u16 = 0,
    Float32List? f32s,
    Float64List? f64s,
  })  : a_ = a_ ?? Uint8List(0),
        os = os ?? [],
        ss = ss ?? [],
        as_ = as_ ?? [],
        f32s = f32s ?? Float32List(0),
        f64s = f64s ?? Float64List(0);

  /// Encodes as Colfer into buf and returns the number of bytes written.
  /// It May throw RangeError or colfer related Exception.
  int marshal(BytesBuilder buf) {
    int i = buf.length;

    if (b) {
      buf.addByte(0);
    }

    if (u32 != 0) {
      if (u32 > 4294967295 || u32 < 0) {
        throw Exception('colfer: $u32 out of reach: u32');
      }
      if (u32 < 0x200000) {
        buf.addByte(1);
        encodeVarint(buf, u32);
      } else {
        buf.addByte(1 | 128);
        buf.add(Uint8List(4)..buffer.asByteData().setInt32(0, u32));
      }
    }

    if (u64 != 0) {
      if (0 < u64 && u64 < 0x2000000000000) {
        buf.addByte(2);
        encodeVarint(buf, u64);
      } else {
        buf.addByte(2 | 128);
        buf.add(Uint8List(8)..buffer.asByteData().setInt64(0, u64));
      }
    }

    if (i32 != 0) {
      if (i32 < 0) {
        buf.addByte(3 | 128);
        if (i32 < -2147483648) {
          throw Exception('colfer: $i32 out of reach: i32');
        }
        encodeVarint(buf, -i32);
      } else {
        buf.addByte(3);
        if (i32 > 2147483647) {
          throw Exception('colfer: $i32 out of reach: i32');
        }
        encodeVarint(buf, i32);
      }
    }

    if (i64 != 0) {
      int a = i64;
      if (a < 0) {
        buf.addByte(4 | 128);
        a = -a;
        if (i64 == a) {
          buf.add(List.filled(9, 128));
        } else {
          encodeVarint(buf, a);
        }
      } else {
        buf.addByte(4);
        encodeVarint(buf, a);
      }
    }

    if (f32 != 0) {
      buf.addByte(5);
      if (f32.isNaN) {
        buf.add([0x7f, 0xc0, 0, 0]);
      } else {
        buf.add(Uint8List(4)..buffer.asByteData().setFloat32(0, f32));
      }
    }

    if (f64 != 0) {
      buf.addByte(6);
      if (f64.isNaN) {
        buf.add([0x7f, 0xf8, 0, 0, 0, 0, 0, 1]);
      } else {
        buf.add(Uint8List(8)..buffer.asByteData().setFloat64(0, f64));
      }
    }

    if (t != null) {
      int us = t!.microsecondsSinceEpoch;
      int res = us % 1000000;
      us -= res;
      int s = us ~/ 1E6;
      int ns = res * 1000;

      if (s >= 1 << 33 || us < 0) {
        buf.addByte(7 | 128);
        buf.add(Uint8List(8)..buffer.asByteData().setInt64(0, s));
      } else {
        buf.addByte(7);
        buf.add(Uint8List(4)..buffer.asByteData().setInt32(0, s));
      }
      buf.add(Uint8List(4)..buffer.asByteData().setInt32(0, ns));
    }

    if (s_.isNotEmpty) {
      buf.addByte(8);
      var v = utf8.encode(s_);
      encodeVarint(buf, v.length);
      buf.add(v);
    }

    if (a_.isNotEmpty) {
      buf.addByte(9);
      encodeVarint(buf, a_.length);
      buf.add(a_);
    }

    if (o != null) {
      buf.addByte(10);
      o?.marshal(buf);
    }

    if (os.isNotEmpty) {
      if (os.length > colferListMax) {
        throw Exception('colfer: gen.o.os size ${os.length} exceeds $colferListMax bytes');
      }
      buf.addByte(11);
      encodeVarint(buf, os.length);
      for (var vi in os) {
        vi.marshal(buf);
      }
    }

    if (ss.isNotEmpty) {
      if (ss.length > colferListMax) {
        throw Exception('colfer: gen.o.ss size ${ss.length} exceeds $colferListMax bytes');
      }
      buf.addByte(12);
      encodeVarint(buf, ss.length);

      for (final vi in ss) {
        var v = utf8.encode(vi);
        encodeVarint(buf, v.length);
        buf.add(v);
      }
    }

    if (as_.isNotEmpty) {
      if (as_.length > colferListMax) {
        throw Exception('colfer: gen.o.as size ${as_.length} exceeds $colferListMax bytes');
      }
      buf.addByte(13);
      encodeVarint(buf, as_.length);
      for (final vi in as_) {
        encodeVarint(buf, vi.length);
        buf.add(vi);
      }
    }

    if (u8 != 0) {
      if (u8 > 255 || u8 < 0) {
        throw Exception('colfer: $u8 out of reach: u8');
      }
      buf.addByte(14);
      buf.addByte(u8);
    }

    if (u16 != 0) {
      if (u16 > 65535 || u16 < 0) {
        throw Exception('colfer: $u16 out of reach: u16');
      }
      if (u16 < 256) {
        buf.addByte(15 | 128);
        buf.addByte(u16);
      } else {
        buf.addByte(15);
        buf.addByte(u16 >> 8);
        buf.addByte(u16 & 255);
      }
    }

    if (f32s.isNotEmpty) {
      if (f32s.length > colferListMax) {
        throw Exception('colfer: gen.o.f32s size ${f32s.length} exceeds $colferListMax bytes');
      }
      buf.addByte(16);
      encodeVarint(buf, f32s.length);
      for (final vi in f32s) {
        if (vi.isNaN) {
          buf.add([0x7f, 0xc0, 0, 0]);
        } else {
          buf.add(Uint8List(4)..buffer.asByteData().setFloat32(0, vi));
        }
      }
    }

    if (f64s.isNotEmpty) {
      if (f64s.length > colferListMax) {
        throw Exception('colfer: gen.o.f64s size ${f64s.length} exceeds $colferListMax bytes');
      }
      buf.addByte(17);
      encodeVarint(buf, f64s.length);
      for (final vi in f64s) {
        if (vi.isNaN) {
          buf.add([0x7f, 0xf8, 0, 0, 0, 0, 0, 1]);
        } else {
          buf.add(Uint8List(8)..buffer.asByteData().setFloat64(0, vi));
        }
      }
    }

    buf.addByte(127);
    i = buf.length - i;
    if (i >= colferSizeMax) {
      throw Exception('colfer: gen.o size $i exceeds $colferSizeMax bytes');
    }
    return i;
  }

  /// Decodes data as Colfer and returns the number of bytes read.
  /// It May throw RangeError or colfer related Exception.
  int unmarshal(Uint8List data) {
    int header = 0, i = 0;
    var view = ByteData.view(data.buffer);
    int nextData() {
      int dataI = data[i];
      i++;
      return dataI;
    }

    void nextHeader() {
      header = nextData();
    }

    nextHeader();

    int readVarint() {
      int c = data[i];
      i++;
      if (c >= 0x80) {
        c &= 0x7f;
        for (int shift = 7;; shift += 7) {
          int b = data[i];
          i++;
          if (b < 0x80 || shift == 56) {
            c |= b << shift;
            break;
          }
          c |= (b & 0x7f) << shift;
        }
      }
      return c;
    }

    if (header == 0) {
      b = true;
      nextHeader();
    }

    if (header == 1) {
      u32 = readVarint();
      nextHeader();
    } else if (header == (1 | 128)) {
      u32 = view.getUint32(i);
      i += 4;
      nextHeader();
    }

    if (header == 2) {
      u64 = readVarint();
      nextHeader();
    } else if (header == (2 | 128)) {
      int v = view.getUint32(i) * 0x100000000;
      v += view.getUint32(i + 4);
      u64 = v;
      i += 8;
      nextHeader();
    }

    if (header == 3) {
      i32 = readVarint();
      nextHeader();
    } else if (header == (3 | 128)) {
      i32 = -1 * readVarint();
      nextHeader();
    }

    if (header == 4) {
      i64 = readVarint();
      nextHeader();
    } else if (header == (4 | 128)) {
      i64 = -1 * readVarint();
      nextHeader();
    }

    if (header == 5) {
      f32 = view.getFloat32(i);
      i += 4;
      nextHeader();
    }

    if (header == 6) {
      f64 = view.getFloat64(i);
      i += 8;
      nextHeader();
    }

    if (header == 7) {
      int s = view.getUint32(i);
      int us = view.getUint32(i + 4) ~/ 1000;
      t = DateTime.fromMicrosecondsSinceEpoch(s * 1000000 + us);
      i += 8;
      nextHeader();
    } else if (header == (7 | 128)) {
      int s = decodeInt64(data, i);
      int us = view.getUint32(i + 8) ~/ 1000;
      t = DateTime.fromMicrosecondsSinceEpoch(s * 1000000 + us);
      i += 12;
      nextHeader();
    }

    if (header == 8) {
      int size = readVarint();
      if (size < 0 || size > colferSizeMax) {
        throw Exception('colfer: gen.o.s size $size exceeds $colferSizeMax bytes');
      }

      int s = i;
      i += size;
      s_ = utf8.decode(data.sublist(s, i));
      nextHeader();
    }

    if (header == 9) {
      int size = readVarint();
      if (size < 0 || size > colferSizeMax) {
        throw Exception('colfer: gen.o.a size $size exceeds $colferSizeMax bytes');
      }

      int start = i;
      i += size;
      a_ = data.sublist(start, i);
      nextHeader();
    }

    if (header == 10) {
      var s = O();
      i += s.unmarshal(data.sublist(i));
      o = s;
      nextHeader();
    }

    if (header == 11) {
      int v = readVarint();
      if (v < 0 || v > colferListMax) {
        throw Exception('colfer: gen.o.os size $v exceeds $colferListMax bytes');
      }

      if (os.length != v) {
        os = List<O>.filled(v, O());
      }
      for (int vi = 0; vi < v; vi++) {
        i += os[vi].unmarshal(data.sublist(i));
      }
      nextHeader();
    }

    if (header == 12) {
      int v = readVarint();
      if (v < 0 || v > colferListMax) {
        throw Exception('colfer: gen.o.ss size $v exceeds $colferListMax bytes');
      }

      if (ss.length != v) {
        ss = List<String>.filled(v, '');
      }
      for (int vi = 0; vi < v; vi++) {
        int size = readVarint();
        if (size < 0 || size > colferSizeMax) {
          throw Exception('colfer: gen.o.ss size $size exceeds $colferSizeMax bytes');
        }

        int s = i;
        i += size;
        ss[vi] = utf8.decode(data.sublist(s, i));
      }
      nextHeader();
    }

    if (header == 13) {
      int v = readVarint();
      if (v < 0 || v > colferListMax) {
        throw Exception('colfer: gen.o.as size $v exceeds $colferListMax bytes');
      }

      if (as_.length != v) {
        as_ = List<Uint8List>.filled(v, Uint8List(0));
      }
      for (int vi = 0; vi < v; vi++) {
        int size = readVarint();
        if (size < 0 || size > colferSizeMax) {
          throw Exception('colfer: gen.o.as size $size exceeds $colferSizeMax bytes');
        }

        int s = i;
        i += size;
        as_[vi] = data.sublist(s, i);
      }
      nextHeader();
    }

    if (header == 14) {
      u8 = nextData();
      nextHeader();
    }

    if (header == 15) {
      u16 = (nextData() << 8) | nextData();
      nextHeader();
    } else if (header == (15 | 128)) {
      u16 = nextData();
      nextHeader();
    }

    if (header == 16) {
      int v = readVarint();
      if (v < 0 || v > colferListMax) {
        throw Exception('colfer: gen.o.f32s size $v exceeds $colferListMax bytes');
      }

      if (f32s.length != v) {
        f32s = Float32List(v);
      }
      for (int vi = 0; vi < v; vi++) {
        f32s[vi] = view.getFloat32(i);
        i += 4;
      }
      nextHeader();
    }

    if (header == 17) {
      int v = readVarint();
      if (v < 0 || v > colferListMax) {
        throw Exception('colfer: gen.o.f64s size $v exceeds $colferListMax bytes');
      }

      if (f64s.length != v) {
        f64s = Float64List(v);
      }
      for (int vi = 0; vi < v; vi++) {
        f64s[vi] = view.getFloat64(i);
        i += 8;
      }
      nextHeader();
    }

    if (header != 127) {
      throw Exception('colfer: unknown header $header at byte ${i - 1}');
    }
    if (i > colferSizeMax) {
      throw Exception('colfer: gen.o size $i exceeds $colferSizeMax bytes');
    }
    return i;
  }
}

// DromedaryCase oposes name casings.
class DromedaryCase {
  String pascalCase;

  @override
  bool operator ==(other) {
    if (!(other is DromedaryCase)) {
      print('compare type doesn\'t match $this != $other');
      return false;
    }
    if (other.pascalCase != pascalCase) {
      print('compare prop pascalCase doesn\'t match $pascalCase != ${other.pascalCase}');
      return false;
    }
    return true;
  }

  @override
  int get hashCode => pascalCase.hashCode;

  DromedaryCase({
    this.pascalCase = '',
  });

  /// Encodes as Colfer into buf and returns the number of bytes written.
  /// It May throw RangeError or colfer related Exception.
  int marshal(BytesBuilder buf) {
    int i = buf.length;

    if (pascalCase.isNotEmpty) {
      buf.addByte(0);
      var v = utf8.encode(pascalCase);
      encodeVarint(buf, v.length);
      buf.add(v);
    }

    buf.addByte(127);
    i = buf.length - i;
    if (i >= colferSizeMax) {
      throw Exception('colfer: gen.dromedaryCase size $i exceeds $colferSizeMax bytes');
    }
    return i;
  }

  /// Decodes data as Colfer and returns the number of bytes read.
  /// It May throw RangeError or colfer related Exception.
  int unmarshal(Uint8List data) {
    int header = 0, i = 0;
    int nextData() {
      int dataI = data[i];
      i++;
      return dataI;
    }

    void nextHeader() {
      header = nextData();
    }

    nextHeader();

    int readVarint() {
      int c = data[i];
      i++;
      if (c >= 0x80) {
        c &= 0x7f;
        for (int shift = 7;; shift += 7) {
          int b = data[i];
          i++;
          if (b < 0x80 || shift == 56) {
            c |= b << shift;
            break;
          }
          c |= (b & 0x7f) << shift;
        }
      }
      return c;
    }

    if (header == 0) {
      int size = readVarint();
      if (size < 0 || size > colferSizeMax) {
        throw Exception(
            'colfer: gen.dromedaryCase.PascalCase size $size exceeds $colferSizeMax bytes');
      }

      int s = i;
      i += size;
      pascalCase = utf8.decode(data.sublist(s, i));
      nextHeader();
    }

    if (header != 127) {
      throw Exception('colfer: unknown header $header at byte ${i - 1}');
    }
    if (i > colferSizeMax) {
      throw Exception('colfer: gen.dromedaryCase size $i exceeds $colferSizeMax bytes');
    }
    return i;
  }
}

// EmbedO has an inner object only.
// Covers regression of issue #66.
class EmbedO {
  O? inner;

  @override
  bool operator ==(other) {
    if (!(other is EmbedO)) {
      print('compare type doesn\'t match $this != $other');
      return false;
    }
    if (other.inner != inner) {
      print('compare prop inner doesn\'t match $inner != ${other.inner}');
      return false;
    }
    return true;
  }

  @override
  int get hashCode => inner.hashCode;

  EmbedO({
    this.inner,
  });

  /// Encodes as Colfer into buf and returns the number of bytes written.
  /// It May throw RangeError or colfer related Exception.
  int marshal(BytesBuilder buf) {
    int i = buf.length;

    if (inner != null) {
      buf.addByte(0);
      inner?.marshal(buf);
    }

    buf.addByte(127);
    i = buf.length - i;
    if (i >= colferSizeMax) {
      throw Exception('colfer: gen.EmbedO size $i exceeds $colferSizeMax bytes');
    }
    return i;
  }

  /// Decodes data as Colfer and returns the number of bytes read.
  /// It May throw RangeError or colfer related Exception.
  int unmarshal(Uint8List data) {
    int header = 0, i = 0;
    int nextData() {
      int dataI = data[i];
      i++;
      return dataI;
    }

    void nextHeader() {
      header = nextData();
    }

    nextHeader();
    if (header == 0) {
      var s = O();
      i += s.unmarshal(data.sublist(i));
      inner = s;
      nextHeader();
    }

    if (header != 127) {
      throw Exception('colfer: unknown header $header at byte ${i - 1}');
    }
    if (i > colferSizeMax) {
      throw Exception('colfer: gen.EmbedO size $i exceeds $colferSizeMax bytes');
    }
    return i;
  }
}
