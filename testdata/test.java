package testdata;

import org.junit.Test;

import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;


public class test {

    static Map<String, O> getGoldenCases() {
	Map<String, O> goldenCases = new HashMap<>();
	O want = new O();
	goldenCases.put("80", want);
	want = new O();
	want.b = true;
	goldenCases.put("8000", want);
	want = new O();
	want.u32 = 1;
	goldenCases.put("800101", want);
	want = new O();
	want.u32 = -1;
	goldenCases.put("8001ffffffff0f", want);
	want = new O();
	want.u64 = 1;
	goldenCases.put("800201", want);
	want = new O();
	want.u64 = -1;
	goldenCases.put("8002ffffffffffffffffff01", want);
	want = new O();
	want.i32 = 1;
	goldenCases.put("800301", want);
	want = new O();
	want.i32 = -1;
	goldenCases.put("808301", want);
	want = new O();
	want.i32 = Integer.MAX_VALUE;
	goldenCases.put("8003ffffffff07", want);
	want = new O();
	want.i32 = Integer.MIN_VALUE;
	goldenCases.put("80838080808008", want);
	want = new O();
	want.i64 = 1;
	goldenCases.put("800401", want);
	want = new O();
	want.i64 = -1;
	goldenCases.put("808401", want);
	want = new O();
	want.i64 = Long.MAX_VALUE;
	goldenCases.put("8004ffffffffffffffff7f", want);
	want = new O();
	want.i64 = Long.MIN_VALUE;
	goldenCases.put("808480808080808080808001", want);
	want = new O();
	want.f32 = Float.MIN_VALUE;
	goldenCases.put("800500000001", want);
	want = new O();
	want.f32 = Float.MAX_VALUE;
	goldenCases.put("80057f7fffff", want);
	want = new O();
	want.f32 = Float.NaN;
	goldenCases.put("80057fc00000", want);
	want = new O();
	want.f64 = Double.MIN_VALUE;
	goldenCases.put("80060000000000000001", want);
	want = new O();
	want.f64 = Double.MAX_VALUE;
	goldenCases.put("80067fefffffffffffff", want);
	want = new O();
	want.f64 = Double.NaN;
	goldenCases.put("80067ff8000000000000", want);
	want = new O();
	want.t = Instant.ofEpochSecond(1441739050, 0);
	goldenCases.put("80070000000055ef312a", want);
	want = new O();
	want.t = Instant.ofEpochSecond(1441739050, 777888999);
	goldenCases.put("80870000000055ef312a2e5da4e7", want);
	want = new O();
	want.s = "A";
	goldenCases.put("80080141", want);
	want = new O();
	want.s = "a\0";
	goldenCases.put("8008026100", want);
	want = new O();
	want.s = "\u0080\u0800\ud800\udc00";
	goldenCases.put("800809c280e0a080f0908080", want);
	want = new O();
	want.a = new byte[]{(byte) 0xff};
	goldenCases.put("800901ff", want);
	want = new O();
	want.a = new byte[]{2, 0};
	goldenCases.put("8009020200", want);
	return goldenCases;
    }

    @Test
    public void TestEncode() {
	for (Entry<String, O> e : getGoldenCases().entrySet()) {
	    try {
		ByteBuffer buf = ByteBuffer.allocate(e.getKey().length() / 2);
		e.getValue().marshal(buf);
		buf.flip();
		assertEquals("serial", e.getKey(), toHex(buf));
	    } catch (Exception ex) {
		assertNull("exception for serial " + e.getKey(), ex);
	    }
	}
    }

    @Test
    public void TestDecode() {
	for (String hex : getGoldenCases().keySet()) {
	    try {
		O o = new O();
		o.unmarshal(parseHex(hex));

		ByteBuffer buf = ByteBuffer.allocate(hex.length() / 2);
		o.marshal(buf);
		buf.flip();
		assertEquals(hex, toHex(buf));
	    } catch (Exception ex) {
		assertNull("exception for serial " + hex, ex);
	    }
	}
    }

    static String toHex(ByteBuffer buf) {
	return new BigInteger(1, buf.array()).toString(16);
    }

    static ByteBuffer parseHex(String s) {
	int len = s.length();
	byte[] data = new byte[len / 2];
	for (int i = 0; i < len; i += 2) {
	    int nibble0 = Character.digit(s.charAt(i), 16);
	    int nibble1 = Character.digit(s.charAt(i + 1), 16);
	    data[i / 2] = (byte) ((nibble0 << 4) + nibble1);
	}
	return ByteBuffer.wrap(data);
    }

}
