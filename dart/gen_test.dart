import 'package:test/test.dart';
import 'dart:typed_data';
import './Colfer.dart';
import 'package:convert/convert.dart';

const MaxUint8 = (1 << 8) - 1;
const MaxUint16 = (1 << 16) - 1;
const MaxUint32 = (1 << 32) - 1;
const MaxUint64 = -1;
const MaxInt8 = MaxUint8 >> 1;
const MinInt8 = -MaxInt8 - 1;
const MaxInt16 = MaxUint16 >> 1;
const MinInt16 = -MaxInt16 - 1;
const MaxInt32 = MaxUint32 >> 1;
const MinInt32 = -MaxInt32 - 1;
const MaxInt64 = (1 << 63) - 1;
const MinInt64 = -MaxInt64 - 1;

class Golden {
  final String serial;
  final O object;

  Golden(this.serial, this.object);
}

List<Golden> newGoldenCases() {
  var tmp = Uint8List(192);
  tmp.fillRange(0, 192, 9);
  return [
    Golden('7f', O()),
    Golden('007f', O(b: true)),
    Golden('01017f', O(u32: 1)),
    Golden('01ff017f', O(u32: MaxUint8)),
    Golden('01ffff037f', O(u32: MaxUint16)),
    Golden('81ffffffff7f', O(u32: MaxUint32)),
    Golden('02017f', O(u64: 1)),
    Golden('02ff017f', O(u64: MaxUint8)),
    Golden('02ffff037f', O(u64: MaxUint16)),
    Golden('02ffffffff0f7f', O(u64: MaxUint32)),
    Golden('82ffffffffffffffff7f', O(u64: MaxUint64)),
    Golden('03017f', O(i32: 1)),
    Golden('83017f', O(i32: -1)),
    Golden('037f7f', O(i32: MaxInt8)),
    Golden('8380017f', O(i32: MinInt8)),
    Golden('03ffff017f', O(i32: MaxInt16)),
    Golden('838080027f', O(i32: MinInt16)),
    Golden('03ffffffff077f', O(i32: MaxInt32)),
    Golden('8380808080087f', O(i32: MinInt32)),
    Golden('04017f', O(i64: 1)),
    Golden('84017f', O(i64: -1)),
    Golden('047f7f', O(i64: MaxInt8)),
    Golden('8480017f', O(i64: MinInt8)),
    Golden('04ffff017f', O(i64: MaxInt16)),
    Golden('848080027f', O(i64: MinInt16)),
    Golden('04ffffffff077f', O(i64: MaxInt32)),
    Golden('8480808080087f', O(i64: MinInt32)),
    Golden('04ffffffffffffffff7f7f', O(i64: MaxInt64)),
    Golden('848080808080808080807f', O(i64: MinInt64)),
    Golden('84ffffffffffffffff7f7f', O(i64: MinInt64 + 1)),
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
    Golden('0901ff7f', O(a: Uint8List.fromList([MaxUint8]))),
    Golden('090202007f', O(a: Uint8List.fromList([2, 0]))),
    Golden(
        '09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f',
        O(a: tmp)),
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
    Golden('0eff7f', O(u8: MaxUint8)),
    Golden('8f017f', O(u16: 1)),
    Golden('0fffff7f', O(u16: MaxUint16)),
    Golden('1002000000003f8000007f', O(f32s: Float32List.fromList([0, 1]))),
    Golden('11014058c000000000007f', O(f64s: Float64List.fromList([99]))),
  ];
}

void main() {
  var cases = newGoldenCases();
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
        var obj = O();
        obj.unmarshal(bytes);
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
            fail('something else then expected RangeError');
          }
          fail('should break with RangeError, but it doesn\'t');
        }
      });
    }
  });
}
