function getGoldenCases() {
	var unixCase = new Date();
	unixCase.setTime(1441739050000);
	var nanoCase = new Date();
	nanoCase.setTime(1441739050777);

	return {
		'807f': {},
		'80007f': {b: true},
		'8001017f': {u32: 1},
		'8001ffffffff0f7f': {u32: 4294967295},
		'8002017f': {u64: 1},
		'8002ffffffffffffff0f7f': {u64: Number.MAX_SAFE_INTEGER},
		'8003017f': {i32: 1},
		'8083017f': {i32: -1},
		'8003ffffffff077f': {i32: 2147483647},
		'808380808080087f': {i32: -2147483648},
		'8004017f': {i64: 1},
		'8084017f': {i64: -1},
		'8004ffffffffffffff0f7f': {i64: Number.MAX_SAFE_INTEGER},
		'8084ffffffffffffff0f7f': {i64: -Number.MAX_SAFE_INTEGER},
		'8005000000017f': {f32: 1.401298464324817e-45},
		'80057f7fffff7f': {f32: 3.4028234663852886e+38},
		'80057fc000007f': {f32: NaN},
		'800600000000000000017f': {f64: Number.MIN_VALUE},
		'80067fefffffffffffff7f': {f64: Number.MAX_VALUE},
		'80067ff80000000000007f': {f64: NaN},
		'80070000000055ef312a7f': {t: unixCase},
		'80870000000055ef312a2e5da4e77f': {t: nanoCase, t_ns: 888999},
		'800801417f': {s: 'A'},
		'80080261007f': {s: 'a\x00'},
		'800809c280e0a080f09080807f': {s: '\u0080\u0800\u{10000}'},
		'800901ff7f': {a: new Uint8Array([0xFF])},
		'80090202007f': {a: new Uint8Array([2, 0])}
	}
}

QUnit.test('marshal', function(assert) {
	var golden = getGoldenCases()
	for (hex in golden) {
		var feed = golden[hex];
		var desc = hex + ': ' + JSON.stringify(feed)
		try {
			var got = encodeHex(testdata.marshalO(feed));
			assert.equal(got, hex, desc);
		} catch (err) {
			assert.equal(err, 'no error', desc);
		}
	}
});

QUnit.test('unmarshal', function(assert) {
	var golden = getGoldenCases()
	for (hex in golden) {
		var want = golden[hex];
		var desc = hex + ': ' + JSON.stringify(want)
		try {
			var got = testdata.unmarshalO(decodeHex(hex));
			assert.deepEqual(got, want, desc);
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
