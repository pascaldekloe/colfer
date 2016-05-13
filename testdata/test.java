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
		newCase(goldenCases, "7f");
		newCase(goldenCases, "007f").b = true;
		newCase(goldenCases, "01017f").u32 = 1;
		newCase(goldenCases, "01ffffffff0f7f").u32 = -1;
		newCase(goldenCases, "02017f").u64 = 1;
		newCase(goldenCases, "02ffffffffffffffffff7f").u64 = -1;
		newCase(goldenCases, "03017f").i32 = 1;
		newCase(goldenCases, "83017f").i32 = -1;
		newCase(goldenCases, "03ffffffff077f").i32 = Integer.MAX_VALUE;
		newCase(goldenCases, "8380808080087f").i32 = Integer.MIN_VALUE;
		newCase(goldenCases, "04017f").i64 = 1;
		newCase(goldenCases, "84017f").i64 = -1;
		newCase(goldenCases, "04ffffffffffffffff7f7f").i64 = Long.MAX_VALUE;
		newCase(goldenCases, "848080808080808080807f").i64 = Long.MIN_VALUE;
		newCase(goldenCases, "05000000017f").f32 = Float.MIN_VALUE;
		newCase(goldenCases, "057f7fffff7f").f32 = Float.MAX_VALUE;
		newCase(goldenCases, "057fc000007f").f32 = Float.NaN;
		newCase(goldenCases, "0600000000000000017f").f64 = Double.MIN_VALUE;
		newCase(goldenCases, "067fefffffffffffff7f").f64 = Double.MAX_VALUE;
		newCase(goldenCases, "067ff80000000000007f").f64 = Double.NaN;
		newCase(goldenCases, "070000000055ef312a7f").t = Instant.ofEpochSecond(1441739050, 0);
		newCase(goldenCases, "870000000055ef312a2e5da4e77f").t = Instant.ofEpochSecond(1441739050, 777888999);
		newCase(goldenCases, "0801417f").s = "A";
		newCase(goldenCases, "080261007f").s = "a\0";
		newCase(goldenCases, "0809c280e0a080f09080807f").s = "\u0080\u0800\ud800\udc00";
		newCase(goldenCases, "0901ff7f").a = new byte[]{-1};
		newCase(goldenCases, "090202007f").a = new byte[]{2, 0};
		newCase(goldenCases, "0a7f7f").o = new O();
		O inner = new O();
		inner.b = true;
		newCase(goldenCases, "0a007f7f").o = inner;
		O element = new O();
		element.b = true;
		newCase(goldenCases, "0b01007f7f").os = new O[] {element};
		newCase(goldenCases, "0b027f7f7f").os = new O[] {new O(), new O()};
		return goldenCases;
	}

	private static O newCase(Map<String, O> cases, String hex) {
		O o = new O();
		cases.put(hex, o);
		return o;
	}

	@Test
	public void testEncode() {
		for (Entry<String, O> e : getGoldenCases().entrySet()) {
			try {
				ByteBuffer buf = ByteBuffer.allocate(e.getKey().length() / 2);
				e.getValue().marshal(buf);
				assertEquals("serial", e.getKey(), toHex(buf.array()));
			} catch (Exception ex) {
				assertNull("exception for serial " + e.getKey(), ex);
			}
		}
	}

	@Test
	public void testDecode() {
		for (Entry<String, O> e : getGoldenCases().entrySet()) {
			try {
				O o = new O();
				o.unmarshal(ByteBuffer.wrap(parseHex(e.getKey())));
				assertEquals(e.getKey(), e.getValue(), o);
			} catch (Exception ex) {
				assertNull("exception for serial " + e.getKey(), ex);
			}
		}
	}

	static String toHex(byte[] bytes) {
		String hex = new BigInteger(1, bytes).toString(16);
		while (bytes.length * 2 > hex.length())
			hex = "0" + hex;
		return hex;
	}

	static byte[] parseHex(String s) {
		int len = s.length();
		byte[] data = new byte[len / 2];
		for (int i = 0; i < len; i += 2) {
			int nibble0 = Character.digit(s.charAt(i), 16);
			int nibble1 = Character.digit(s.charAt(i + 1), 16);
			data[i / 2] = (byte) ((nibble0 << 4) + nibble1);
		}
		return data;
	}

}
