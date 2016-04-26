function getGoldenCases() {
	var unixCase = new Date();
	unixCase.setTime(1441739050000);
	var nanoCase = new Date();
	nanoCase.setTime(1441739050777);

	return {
		'7f': {},
		'007f': {b: true},
		'01017f': {u32: 1},
		'01ffffffff0f7f': {u32: 4294967295},
		'02017f': {u64: 1},
		'02ffffffffffffff0f7f': {u64: Number.MAX_SAFE_INTEGER},
		'03017f': {i32: 1},
		'83017f': {i32: -1},
		'03ffffffff077f': {i32: 2147483647},
		'8380808080087f': {i32: -2147483648},
		'04017f': {i64: 1},
		'84017f': {i64: -1},
		'04ffffffffffffff0f7f': {i64: Number.MAX_SAFE_INTEGER},
		'84ffffffffffffff0f7f': {i64: -Number.MAX_SAFE_INTEGER},
		'05000000017f': {f32: 1.401298464324817e-45},
		'057f7fffff7f': {f32: 3.4028234663852886e+38},
		'057fc000007f': {f32: NaN},
		'0600000000000000017f': {f64: Number.MIN_VALUE},
		'067fefffffffffffff7f': {f64: Number.MAX_VALUE},
		'067ff80000000000007f': {f64: NaN},
		'070000000055ef312a7f': {t: unixCase},
		'870000000055ef312a2e5da4e77f': {t: nanoCase, t_ns: 888999},
		'0801417f': {s: 'A'},
		'080261007f': {s: 'a\x00'},
		'0809c280e0a080f09080807f': {s: '\u0080\u0800\u{10000}'},
		'0901ff7f': {a: new Uint8Array([0xFF])},
		'090202007f': {a: new Uint8Array([2, 0])},
		'0a7f7f': {o: {}},
		'0a007f7f': {o: {b: true}},
		'0b01007f7f': {os: [{b: true}]},
		'0b027f7f7f': {os: [{}, {}]}
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
