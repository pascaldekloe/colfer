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
	goldenCases.put("807f", want);
	want = new O();
	want.b = true;
	goldenCases.put("80007f", want);
	want = new O();
	want.u32 = 1;
	goldenCases.put("8001017f", want);
	want = new O();
	want.u32 = -1;
	goldenCases.put("8001ffffffff0f7f", want);
	want = new O();
	want.u64 = 1;
	goldenCases.put("8002017f", want);
	want = new O();
	want.u64 = -1;
	goldenCases.put("8002ffffffffffffffffff017f", want);
	want = new O();
	want.i32 = 1;
	goldenCases.put("8003017f", want);
	want = new O();
	want.i32 = -1;
	goldenCases.put("8083017f", want);
	want = new O();
	want.i32 = Integer.MAX_VALUE;
	goldenCases.put("8003ffffffff077f", want);
	want = new O();
	want.i32 = Integer.MIN_VALUE;
	goldenCases.put("808380808080087f", want);
	want = new O();
	want.i64 = 1;
	goldenCases.put("8004017f", want);
	want = new O();
	want.i64 = -1;
	goldenCases.put("8084017f", want);
	want = new O();
	want.i64 = Long.MAX_VALUE;
	goldenCases.put("8004ffffffffffffffff7f7f", want);
	want = new O();
	want.i64 = Long.MIN_VALUE;
	goldenCases.put("8084808080808080808080017f", want);
	want = new O();
	want.f32 = Float.MIN_VALUE;
	goldenCases.put("8005000000017f", want);
	want = new O();
	want.f32 = Float.MAX_VALUE;
	goldenCases.put("80057f7fffff7f", want);
	want = new O();
	want.f32 = Float.NaN;
	goldenCases.put("80057fc000007f", want);
	want = new O();
	want.f64 = Double.MIN_VALUE;
	goldenCases.put("800600000000000000017f", want);
	want = new O();
	want.f64 = Double.MAX_VALUE;
	goldenCases.put("80067fefffffffffffff7f", want);
	want = new O();
	want.f64 = Double.NaN;
	goldenCases.put("80067ff80000000000007f", want);
	want = new O();
	want.t = Instant.ofEpochSecond(1441739050, 0);
	goldenCases.put("80070000000055ef312a7f", want);
	want = new O();
	want.t = Instant.ofEpochSecond(1441739050, 777888999);
	goldenCases.put("80870000000055ef312a2e5da4e77f", want);
	want = new O();
	want.s = "A";
	goldenCases.put("800801417f", want);
	want = new O();
	want.s = "a\0";
	goldenCases.put("80080261007f", want);
	want = new O();
	want.s = "\u0080\u0800\ud800\udc00";
	goldenCases.put("800809c280e0a080f09080807f", want);
	want = new O();
	want.a = new byte[]{(byte) 0xff};
	goldenCases.put("800901ff7f", want);
	want = new O();
	want.a = new byte[]{2, 0};
	goldenCases.put("80090202007f", want);
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
