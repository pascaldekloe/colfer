function getGoldenCases() {
	var unixCase = new Date();
	unixCase.setTime(1441739050000);
	var nanoCase = new Date();
	nanoCase.setTime(1441739050777);

	return {
		'80': {},
		'8000': {b: true},
		'800101': {u32: 1},
		'8001ffffffff0f': {u32: 4294967295},
		'800201': {u64: 1},
		'800280808080808080808002': {u64: 18446744073709552000},
		'800301': {i32: 1},
		'808301': {i32: -1},
		'8003ffffffff07': {i32: 2147483647},
		'80838080808008': {i32: -2147483648},
		'800401': {i64: 1},
		'808401': {i64: -1},
		'800480808080808080808001': {i64: 9223372036854776000},
		'808480808080808080808001': {i64: -9223372036854775808},
		'800500000001': {f32: 1.401298464324817e-45},
		'80057f7fffff': {f32: 3.4028234663852886e+38},
		'80057fc00000': {f32: NaN},
		'80060000000000000001': {f64: Number.MIN_VALUE},
		'80067fefffffffffffff': {f64: Number.MAX_VALUE},
		'80067ff8000000000000': {f64: NaN},
		'80070000000055ef312a': {t: unixCase},
		'80870000000055ef312a2e5da4e7': {t: nanoCase, t_ns: 888999},
		'80080141': {s: 'A'},
		'8008026100': {s: 'a\x00'},
		'800809c280e0a080f0908080': {s: '\u0080\u0800\u{10000}'},
		'800901ff': {a: new Uint8Array([0xFF])},
		'8009020200': {a: new Uint8Array([2, 0])}
	}
}

QUnit.test('marshal', function(assert) {
	var golden = getGoldenCases()
	for (hex in golden) {
		var feed = golden[hex];
		var got = encodeHex(testdata.marshalO(feed));
		assert.equal(got, hex, hex + ': ' + JSON.stringify(feed));
	}
});

QUnit.test('unmarshal', function(assert) {
	var golden = getGoldenCases()
	for (hex in golden) {
		var want = golden[hex];
		var got = testdata.unmarshalO(decodeHex(hex));
		assert.deepEqual(got, want, hex + ': ' + JSON.stringify(want));
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
