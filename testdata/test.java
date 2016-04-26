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
	goldenCases.put("7f", want);
	want = new O();
	want.b = true;
	goldenCases.put("007f", want);
	want = new O();
	want.u32 = 1;
	goldenCases.put("01017f", want);
	want = new O();
	want.u32 = -1;
	goldenCases.put("01ffffffff0f7f", want);
	want = new O();
	want.u64 = 1;
	goldenCases.put("02017f", want);
	want = new O();
	want.u64 = -1;
	goldenCases.put("02ffffffffffffffffff017f", want);
	want = new O();
	want.i32 = 1;
	goldenCases.put("03017f", want);
	want = new O();
	want.i32 = -1;
	goldenCases.put("83017f", want);
	want = new O();
	want.i32 = Integer.MAX_VALUE;
	goldenCases.put("03ffffffff077f", want);
	want = new O();
	want.i32 = Integer.MIN_VALUE;
	goldenCases.put("8380808080087f", want);
	want = new O();
	want.i64 = 1;
	goldenCases.put("04017f", want);
	want = new O();
	want.i64 = -1;
	goldenCases.put("84017f", want);
	want = new O();
	want.i64 = Long.MAX_VALUE;
	goldenCases.put("04ffffffffffffffff7f7f", want);
	want = new O();
	want.i64 = Long.MIN_VALUE;
	goldenCases.put("84808080808080808080017f", want);
	want = new O();
	want.f32 = Float.MIN_VALUE;
	goldenCases.put("05000000017f", want);
	want = new O();
	want.f32 = Float.MAX_VALUE;
	goldenCases.put("057f7fffff7f", want);
	want = new O();
	want.f32 = Float.NaN;
	goldenCases.put("057fc000007f", want);
	want = new O();
	want.f64 = Double.MIN_VALUE;
	goldenCases.put("0600000000000000017f", want);
	want = new O();
	want.f64 = Double.MAX_VALUE;
	goldenCases.put("067fefffffffffffff7f", want);
	want = new O();
	want.f64 = Double.NaN;
	goldenCases.put("067ff80000000000007f", want);
	want = new O();
	want.t = Instant.ofEpochSecond(1441739050, 0);
	goldenCases.put("070000000055ef312a7f", want);
	want = new O();
	want.t = Instant.ofEpochSecond(1441739050, 777888999);
	goldenCases.put("870000000055ef312a2e5da4e77f", want);
	want = new O();
	want.s = "A";
	goldenCases.put("0801417f", want);
	want = new O();
	want.s = "a\0";
	goldenCases.put("080261007f", want);
	want = new O();
	want.s = "\u0080\u0800\ud800\udc00";
	goldenCases.put("0809c280e0a080f09080807f", want);
	want = new O();
	want.a = new byte[]{(byte) 0xff};
	goldenCases.put("0901ff7f", want);
	want = new O();
	want.a = new byte[]{2, 0};
	goldenCases.put("090202007f", want);
	want = new O();
	want.o = new O();
	goldenCases.put("0a7f7f", want);

	want = new O();
	want.o = new O();
	want.o.b = true;
	goldenCases.put("0a007f7f", want);
	want = new O();
	want.os = new O[1];
	want.os[0] = new O();
	want.os[0].b = true;
	goldenCases.put("0b01007f7f", want);
	want = new O();
	want.os = new O[2];
	want.os[0] = new O();
	want.os[1] = new O();
	goldenCases.put("0b027f7f7f", want);
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
	String hex = new BigInteger(1, buf.array()).toString(16);
	while (buf.remaining() * 2 > hex.length())
		hex = "0" + hex;
	return hex;
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
