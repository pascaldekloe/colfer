QUnit.test('constructor', function(assert) {
	assert.deepEqual(new testdata.O(), new testdata.O({}), 'absent and empty init');

	var o = new testdata.O({s: 'hello', i32: 42});
	assert.deepEqual(new testdata.O(o), o, 'clone');
});

function newGoldenCases() {
	var date1 = new Date(), date2 = new Date(), date3 = new Date(), date4 = new Date();
	date1.setTime(1441739050777);
	date2.setTime(0x100000000 * 1000);
	date3.setTime(-223);
	date4.setTime(-63094636800000);

	return {
		'7f': {},
		'007f': {b: true},
		'01017f': {u32: 1},
		'01ff017f': {u32: 255},
		'01ffff037f': {u32: 65535},
		'81ffffffff7f': {u32: 4294967295},
		'02017f': {u64: 1},
		'02ff017f': {u64: 255},
		'02ffff037f': {u64: 65535},
		'02ffffffff0f7f': {u64: 4294967295},
		'82001fffffffffffff7f': {u64: Number.MAX_SAFE_INTEGER},
		'03017f': {i32: 1},
		'83017f': {i32: -1},
		'037f7f': {i32: 127},
		'8380017f': {i32: -128},
		'03ffff017f': {i32: 32767},
		'838080027f': {i32: -32768},
		'03ffffffff077f': {i32: 2147483647},
		'8380808080087f': {i32: -2147483648},
		'04017f': {i64: 1},
		'84017f': {i64: -1},
		'047f7f': {i64: 127},
		'8480017f': {i64: -128},
		'04ffff017f': {i64: 32767},
		'848080027f': {i64: -32768},
		'04ffffffff077f': {i64: 2147483647},
		'8480808080087f': {i64: -2147483648},
		'04ffffffffffffff0f7f': {i64: Number.MAX_SAFE_INTEGER},
		'84ffffffffffffff0f7f': {i64: -Number.MAX_SAFE_INTEGER},
		'05000000017f': {f32: 1.401298464324817e-45},
		'057f7fffff7f': {f32: 3.4028234663852886e+38},
		'057fc000007f': {f32: NaN},
		'0600000000000000017f': {f64: Number.MIN_VALUE},
		'067fefffffffffffff7f': {f64: Number.MAX_VALUE},
		'067ff80000000000007f': {f64: NaN},
		'0755ef312a2e5da4e77f': {t: date1, t_ns: 888999},
		'870000000100000000000000007f': {t: date2, t_ns: 0},
		'87ffffffffffffffff2e5da4e77f': {t: date3, t_ns: 888999},
		'87fffffff14f443f00000000007f': {t: date4, t_ns: 0},
		'0801417f': {s: 'A'},
		'080261007f': {s: 'a\x00'},
		'0809c280e0a080f09080807f': {s: '\u0080\u0800\u{10000}'},
		'08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f': {s: '                                                                                                                                '},
		'0901ff7f': {a: new Uint8Array([0xFF])},
		'090202007f': {a: new Uint8Array([2, 0])},
		'09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f': {a: new Uint8Array(192).fill(9)},
		'0a7f7f': {o: new testdata.O()},
		'0a007f7f': {o: new testdata.O({b: true})},
		'0b01007f7f': {os: [new testdata.O({b: true})]},
		'0b027f7f7f': {os: [new testdata.O(), new testdata.O()]},
		'0c0300016101627f': {ss: ["", "a", "b"]},
		'0d0201000201027f': {as: [new Uint8Array([0]), new Uint8Array([1, 2])]},
		'0e017f': {u8: 1},
		'0eff7f': {u8: 255},
		'8f017f': {u16: 1},
		'0fffff7f': {u16: 65535}
	}
}

QUnit.test('marshal', function(assert) {
	var golden = newGoldenCases();
	for (hex in golden) {
		var feed = golden[hex];
		var desc = hex + ': ' + JSON.stringify(feed)
		try {
			var o = new testdata.O(feed);
			var got = encodeHex(o.marshal());
			assert.equal(got, hex, desc);
		} catch (err) {
			assert.equal(err, 'no error', desc);
		}
	}
});

QUnit.test('unmarshal', function(assert) {
	var golden = newGoldenCases();
	for (hex in golden) {
		var want = golden[hex];
		var desc = hex + ': ' + JSON.stringify(want)
		try {
			var got = new testdata.O();
			got.unmarshal(decodeHex(hex));
			assert.deepEqual(got, new testdata.O(want), desc);
		} catch (err) {
			assert.equal(err, 'no error', desc);
		}
	}
});

function encodeHex(bytes) {
	var s = '';
	if (!bytes) return s;

	for (var i = 0; i < bytes.length; i++) {
		var hex = (bytes[i] & 0xff).toString(16);
		if (hex.length == 1) hex = '0' + hex;
		s += hex;
	}
	return s;
}

function decodeHex(s) {
	if (!s) return new Uint8Array();;

	var a = [];
	for (var i = 0; i < s.length; i += 2)
		a.push(parseInt(s.substr(i, 2), 16));
	return new Uint8Array(a);
}
