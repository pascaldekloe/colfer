import 'dart:typed_data';
import 'package:convert/convert.dart';
import 'package:test/test.dart';
import 'gen/Colfer.dart';

const maxUint8 = (1 << 8) - 1;
const maxUint16 = (1 << 16) - 1;
const maxUint32 = (1 << 32) - 1;
const maxUint64 = -1;
const maxInt8 = maxUint8 >> 1;
const minInt8 = -maxInt8 - 1;
const maxInt16 = maxUint16 >> 1;
const minInt16 = -maxInt16 - 1;
const maxInt32 = maxUint32 >> 1;
const minInt32 = -maxInt32 - 1;
const maxInt64 = (1 << 63) - 1;
const minInt64 = -maxInt64 - 1;

class Golden {
  Golden(this.serial, this.object);

  final String serial;
  final O object;
}

List<Golden> newGoldenCases() {
  return [
    Golden('7f', O()),
    Golden('007f', O(b: true)),
    Golden('01017f', O(u32: 1)),
    Golden('01ff017f', O(u32: maxUint8)),
    Golden('01ffff037f', O(u32: maxUint16)),
    Golden('81ffffffff7f', O(u32: maxUint32)),
    Golden('02017f', O(u64: 1)),
    Golden('02ff017f', O(u64: maxUint8)),
    Golden('02ffff037f', O(u64: maxUint16)),
    Golden('02ffffffff0f7f', O(u64: maxUint32)),
    Golden('82ffffffffffffffff7f', O(u64: maxUint64)),
    Golden('03017f', O(i32: 1)),
    Golden('83017f', O(i32: -1)),
    Golden('037f7f', O(i32: maxInt8)),
    Golden('8380017f', O(i32: minInt8)),
    Golden('03ffff017f', O(i32: maxInt16)),
    Golden('838080027f', O(i32: minInt16)),
    Golden('03ffffffff077f', O(i32: maxInt32)),
    Golden('8380808080087f', O(i32: minInt32)),
    Golden('04017f', O(i64: 1)),
    Golden('84017f', O(i64: -1)),
    Golden('047f7f', O(i64: maxInt8)),
    Golden('8480017f', O(i64: minInt8)),
    Golden('04ffff017f', O(i64: maxInt16)),
    Golden('848080027f', O(i64: minInt16)),
    Golden('04ffffffff077f', O(i64: maxInt32)),
    Golden('8480808080087f', O(i64: minInt32)),
    Golden('04ffffffffffffffff7f7f', O(i64: maxInt64)),
    Golden('848080808080808080807f', O(i64: minInt64)),
    Golden('84ffffffffffffffff7f7f', O(i64: minInt64 + 1)),
    Golden('05000000017f', O(f32: 1.401298464324817e-45)),
    Golden('057f7fffff7f', O(f32: 3.4028234663852886e+38)),
    Golden('057fc000007f', O(f32: double.nan)),
    Golden('0600000000000000017f', O(f64: double.minPositive)),
    Golden('067fefffffffffffff7f', O(f64: double.maxFinite)),
    Golden('067ff80000000000017f', O(f64: double.nan)),
    // the following 4 DateTime tests differ from the ones in other languages because of the nanosecond precision loss
    Golden('0755ef312a2e5da1007f', O(t: DateTime.fromMicrosecondsSinceEpoch(1441739050777888))),
    Golden('87000007d954159c00000003e87f',
        O(t: DateTime.fromMicrosecondsSinceEpoch(8630000000000000001))),
    Golden('87fffff82457de8000000003e87f',
        O(t: DateTime.fromMicrosecondsSinceEpoch(-8639999999999999999))),
    Golden('87ffffffffffffffff2e5da1007f', O(t: DateTime.fromMicrosecondsSinceEpoch(-222112))),
    Golden('0801417f', O(s: 'A')),
    Golden('080261007f', O(s: 'a\x00')),
    Golden('0809c280e0a080f09080807f', O(s: '\u0080\u0800\u{10000}')),
    Golden(
        '08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f',
        O(s: ' ' * 128)),
    Golden('0901ff7f', O(a: Uint8List.fromList([maxUint8]))),
    Golden('090202007f', O(a: Uint8List.fromList([2, 0]))),
    Golden(
        '09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f',
        O(a: Uint8List(192)..fillRange(0, 192, 9))),
    Golden('0a7f7f', O(o: O())),
    Golden('0a007f7f', O(o: O(b: true))),
    Golden('0b01007f7f', O(os: [O(b: true)])),
    Golden('0b027f7f7f', O(os: [O(), O()])),
    Golden('0c0300016101627f', O(ss: ['', 'a', 'b'])),
    Golden(
        '0d0201000201027f',
        O(as_0: [
          Uint8List.fromList([0]),
          Uint8List.fromList([1, 2])
        ])),
    Golden('0e017f', O(u8: 1)),
    Golden('0eff7f', O(u8: maxUint8)),
    Golden('8f017f', O(u16: 1)),
    Golden('0fffff7f', O(u16: maxUint16)),
    Golden('1002000000003f8000007f', O(f32s: Float32List.fromList([0, 1]))),
    Golden('11014058c000000000007f', O(f64s: Float64List.fromList([99]))),
  ];
}

void main() {
  var cases = newGoldenCases();

  group('Identity', () {
    for (int i = 0; i < cases.length; i++) {
      for (int j = 0; j < cases.length; j++) {
        var a = cases[i];
        var b = cases[j];
        if (i == j && a.hashCode != b.hashCode) {
          fail('inconsistent hash on object $a');
        } else if (i != j && a.hashCode == b.hashCode) {
          fail('hash collision on object $a & $b');
        }
      }
    }
  });

  group('Marshal', () {
    for (final gold in cases) {
      var buf = Uint8List(gold.object.marshalLen());
      int size = gold.object.marshalTo(buf);
      test('marshaled binary should match golden mapping', () {
        List<int> list = hex.decode(gold.serial);
        Uint8List bytes = Uint8List.fromList(list);
        buf = Uint8List.view(buf.buffer, 0, size);
        expect(buf, bytes);
      });
    }
  });

  group('Unmarshal', () {
    // getting separate cases because of the local changes
    for (final gold in newGoldenCases()) {
      test('struct should match golden mapping', () {
        List<int> list = hex.decode(gold.serial);
        Uint8List bytes = Uint8List.fromList(list);
        var obj = O()..unmarshal(bytes);
        // work around NaN != NaN
        if (obj.f32.isNaN && gold.object.f32.isNaN) {
          obj.f32 = 0;
          gold.object.f32 = 0;
        }
        if (obj.f64.isNaN && gold.object.f64.isNaN) {
          obj.f64 = 0;
          gold.object.f64 = 0;
        }
        expect(obj, equals(gold.object));
      });
    }
  });

  group('Unmarshal Incomplete', () {
    for (final gold in cases) {
      test('should fail with RangeError', () {
        List<int> list = hex.decode(gold.serial);
        var obj = O();
        for (int i = 0; i < list.length; i++) {
          Uint8List bytes = Uint8List.fromList(list.sublist(0, i));
          try {
            obj.unmarshal(bytes);
          } on RangeError {
            continue;
          } catch (e) {
            fail('want a RangeError, got: $e');
          }
          fail("should break with RangeError, but it doesn't");
        }
      });
    }
  });
}
