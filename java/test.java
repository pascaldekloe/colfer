import gen.O;

import java.io.ByteArrayOutputStream;
import java.io.ByteArrayInputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.time.Instant;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;

import static java.nio.charset.StandardCharsets.UTF_8;


public class test {

	static boolean testSuccess = true;


	public static void main(String[] args) {
		try {
			identity();

			marshal();
			unmarshal();
			stream();

			marshalMax();
			marshalTextMax();
			marshalBinaryMax();
			marshalListMax();

			unmarshalMax();
			unmarshalTextMax();
			unmarshalBinaryMax();
			unmarshalListMax();

			serializable();
		} catch (Exception e) {
			e.printStackTrace();
			System.exit(1);
		}

		if (! testSuccess) System.exit(2);
	}

	static void fail(String format, Object... args) {
		format += "\n";
		System.err.printf(format, args);

		testSuccess = false;
	}

	static Map<String, O> newGoldenCases() {
		Map<String, O> goldenCases = new LinkedHashMap<>();
		newCase(goldenCases, "7f");
		newCase(goldenCases, "007f").b = true;
		newCase(goldenCases, "01017f").u32 = 1;
		newCase(goldenCases, "01ff017f").u32 = 255;
		newCase(goldenCases, "01ffff037f").u32 = 65535;
		newCase(goldenCases, "81ffffffff7f").u32 = -1;
		newCase(goldenCases, "02017f").u64 = 1L;
		newCase(goldenCases, "02ff017f").u64 = 255L;
		newCase(goldenCases, "02ffff037f").u64 = 65535L;
		newCase(goldenCases, "02ffffffff0f7f").u64 = 4294967295L;
		newCase(goldenCases, "82ffffffffffffffff7f").u64 = -1L;
		newCase(goldenCases, "03017f").i32 = 1;
		newCase(goldenCases, "83017f").i32 = -1;
		newCase(goldenCases, "037f7f").i32 = Byte.MAX_VALUE;
		newCase(goldenCases, "8380017f").i32 = Byte.MIN_VALUE;
		newCase(goldenCases, "03ffff017f").i32 = Short.MAX_VALUE;
		newCase(goldenCases, "838080027f").i32 = Short.MIN_VALUE;
		newCase(goldenCases, "03ffffffff077f").i32 = Integer.MAX_VALUE;
		newCase(goldenCases, "8380808080087f").i32 = Integer.MIN_VALUE;
		newCase(goldenCases, "04017f").i64 = 1;
		newCase(goldenCases, "84017f").i64 = -1;
		newCase(goldenCases, "047f7f").i64 = Byte.MAX_VALUE;
		newCase(goldenCases, "8480017f").i64 = Byte.MIN_VALUE;
		newCase(goldenCases, "04ffff017f").i64 = Short.MAX_VALUE;
		newCase(goldenCases, "848080027f").i64 = Short.MIN_VALUE;
		newCase(goldenCases, "04ffffffff077f").i64 = Integer.MAX_VALUE;
		newCase(goldenCases, "8480808080087f").i64 = Integer.MIN_VALUE;
		newCase(goldenCases, "04ffffffffffffffff7f7f").i64 = Long.MAX_VALUE;
		newCase(goldenCases, "848080808080808080807f").i64 = Long.MIN_VALUE;
		newCase(goldenCases, "05000000017f").f32 = Float.MIN_VALUE;
		newCase(goldenCases, "057f7fffff7f").f32 = Float.MAX_VALUE;
		newCase(goldenCases, "057fc000007f").f32 = Float.NaN;
		newCase(goldenCases, "0600000000000000017f").f64 = Double.MIN_VALUE;
		newCase(goldenCases, "067fefffffffffffff7f").f64 = Double.MAX_VALUE;
		newCase(goldenCases, "067ff80000000000007f").f64 = Double.NaN;
		newCase(goldenCases, "0755ef312a2e5da4e77f").t = Instant.ofEpochSecond(1441739050L, 777888999);
		newCase(goldenCases, "870000000100000000000000007f").t = Instant.ofEpochSecond(1L << 32, 0);
		newCase(goldenCases, "87ffffffffffffffff2e5da4e77f").t = Instant.ofEpochSecond(-1L, 777888999);
		newCase(goldenCases, "87fffffff14f443f00000000007f").t = Instant.ofEpochSecond(-63094636800L, 0);
		newCase(goldenCases, "0801417f").s = "A";
		newCase(goldenCases, "080261007f").s = "a\0";
		newCase(goldenCases, "0809c280e0a080f09080807f").s = "\u0080\u0800\ud800\udc00";
		newCase(goldenCases, "08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f").s = "                                                                                                                                ";
		newCase(goldenCases, "0901ff7f").a = new byte[]{-1};
		newCase(goldenCases, "090202007f").a = new byte[]{2, 0};
		newCase(goldenCases, "09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f").a = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t".getBytes(UTF_8);
		newCase(goldenCases, "0a7f7f").o = new O();
		O inner = new O();
		inner.b = true;
		newCase(goldenCases, "0a007f7f").o = inner;
		O element = new O();
		element.b = true;
		newCase(goldenCases, "0b01007f7f").os = new O[] {element};
		newCase(goldenCases, "0b027f7f7f").os = new O[] {new O(), new O()};
		newCase(goldenCases, "0c0300016101627f").ss = new String[] {"", "a", "b"};
		newCase(goldenCases, "0d0201000201027f").as = new byte[][]{new byte[]{0}, new byte[]{1, 2}};
		newCase(goldenCases, "0e017f").u8 = 1;
		newCase(goldenCases, "0eff7f").u8 = -1;
		newCase(goldenCases, "8f017f").u16 = 1;
		newCase(goldenCases, "0fffff7f").u16 = -1;
		newCase(goldenCases, "1002000000003f8000007f").f32s = new float[] {0, 1};
		newCase(goldenCases, "11014058c000000000007f").f64s = new double[] {99};
		return goldenCases;
	}

	static O newCase(Map<String, O> cases, String hex) {
		O o = new O();
		cases.put(hex, o);
		return o;
	}


	static void identity() {
		if (new O().equals((Object) null))
			fail("equals null Object");
		if (new O().equals((O) null))
			fail("equals null O");

		Object[] a = newGoldenCases().values().toArray();
		Object[] b = newGoldenCases().values().toArray();
		if (! Arrays.equals(a, b))
			fail("golden cases not equal");
		if (Arrays.hashCode(a) != Arrays.hashCode(b))
			fail("golden cases hash not equal");
	}

	static void marshal() throws Exception {
		for (Entry<String, O> e : newGoldenCases().entrySet()) {
			byte[] buf = new byte[O.colferSizeMax];
			int n = e.getValue().marshal(buf, 0);
			if (n != e.getKey().length() / 2)
				fail("marshal: got write index %d for serial 0x%s", n, e.getKey());
			String got = toHex(Arrays.copyOf(buf, n));
			if (! got.equals(e.getKey()))
				fail("marshal: got serial 0x%s, want %s", got, e.getKey());
		}
	}

	static void unmarshal() {
		for (Entry<String, O> e : newGoldenCases().entrySet()) {
			O o = new O();
			byte[] serial = parseHex(e.getKey());
			int i = o.unmarshal(serial, 0);

			if (i != serial.length)
				fail("unmarshal: 0x%s: got read index %d for serial 0x%s", i, e.getKey());
			if (! e.getValue().equals(o))
				fail("unmarshal: mismatch for serial 0x%s", e.getKey());
		}
	}

	static void stream() throws Exception {
		ByteArrayOutputStream out = new ByteArrayOutputStream();

		byte[] buf = new byte[1];
		for (O o : newGoldenCases().values()) {
			buf = o.marshal(out, buf);
		}

		O.Unmarshaller unmarshaller = new O.Unmarshaller(new ByteArrayInputStream(out.toByteArray()), new byte[1]);
		for (Entry<String, O> e : newGoldenCases().entrySet()) {
			O got = unmarshaller.next();
			if (got == null) {
				fail("stream: missing as of serial 0x%s", e.getKey());
				return;
			}
			if (! e.getValue().equals(got))
				fail("stream: mismatch for serial 0x%s", e.getKey());
		}
		if (unmarshaller.next() != null)
			fail("stream: data tail");
	}

	static void marshalMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			O o = new O();
			o.u64 = 1;
			o.marshal(new byte[O.colferSizeMax], 0);
			fail("no marshal max exception");
		} catch (IllegalStateException e) {
			String want = "colfer: gen.o exceeds 2 bytes";
			if (! want.equals(e.getMessage()))
				fail("marshal max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	static void marshalTextMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			O o = new O();
			o.s = "AAA";
			o.marshal(new byte[6], 0);
			fail("no marshal text max exception");
		} catch (IllegalStateException e) {
			String want = "colfer: gen.o.s size 3 exceeds 2 UTF-8 bytes";
			if (! want.equals(e.getMessage()))
				fail("marshal text max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	static void marshalBinaryMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			O o = new O();
			o.a = new byte[]{0, 1, 2};
			o.marshal(new byte[O.colferSizeMax], 0);
			fail("no marshal binary max exception");
		} catch (IllegalStateException e) {
			String want = "colfer: gen.o.a size 3 exceeds 2 bytes";
			if (! want.equals(e.getMessage()))
				fail("marshal binary max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	static void marshalListMax() {
		int origMax = O.colferListMax;
		O.colferListMax = 9;
		try {
			O o = new O();
			o.os = new O[10];
			o.marshal(new byte[O.colferSizeMax], 0);
			fail("no marshal list max exception");
		} catch (IllegalStateException e) {
			String want = "colfer: gen.o.os length 10 exceeds 9 elements";
			if (! want.equals(e.getMessage()))
				fail("marshal list max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferListMax = origMax;
		}
	}

	static void unmarshalMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			byte[] serial = parseHex("02017f");
			new O().unmarshal(serial, 0);
			fail("no unmarshal max exception");
		} catch (SecurityException e) {
			String want = "colfer: gen.o exceeds 2 bytes";
			if (! want.equals(e.getMessage()))
				fail("unmarshal max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	static void unmarshalTextMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 9;
		try {
			byte[] serial = parseHex("080a414141");
			new O().unmarshal(serial, 0);
			fail("no unmarshal text max exception");
		} catch (SecurityException e) {
			String want = "colfer: gen.o.s size 10 exceeds 9 UTF-8 bytes";
			if (! want.equals(e.getMessage()))
				fail("unmarshal text max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	static void unmarshalBinaryMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 9;
		try {
			byte[] serial = parseHex("090a414141");
			new O().unmarshal(serial, 0);
			fail("no unmarshal binary max exception");
		} catch (SecurityException e) {
			String want = "colfer: gen.o.a size 10 exceeds 9 bytes";
			if (! want.equals(e.getMessage()))
				fail("unmarshal binary max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	static void unmarshalListMax() {
		int origMax = O.colferListMax;
		O.colferListMax = 9;
		try {
			byte[] serial = parseHex("0b0a7f7f7f");
			new O().unmarshal(serial, 0);
			fail("no unmarshal list max exception");
		} catch (SecurityException e) {
			String want = "colfer: gen.o.os length 10 exceeds 9 elements";
			if (! want.equals(e.getMessage()))
				fail("unmarshal list max error: %s\nwant: %s", e.getMessage(), want);
		} finally {
			O.colferListMax = origMax;
		}
	}

	static void serializable() throws Exception {
		Set<Entry<String, O>> cases = newGoldenCases().entrySet();
		ByteArrayOutputStream buf = new ByteArrayOutputStream();

		ObjectOutputStream out = new ObjectOutputStream(buf);
		for (Entry<String, O> e : cases)
			out.writeObject(e.getValue());
		out.close();

		ObjectInputStream in = new ObjectInputStream(new ByteArrayInputStream(buf.toByteArray()));
		for (Entry<String, O> e : cases) {
			O got = (O) in.readObject();
			O want = e.getValue();
			if (want.equals(got)) continue;
			byte[] serial = new byte[O.colferSizeMax];
			int n = got.marshal(serial, 0);
			fail("got 0x%s, want 0x%s", toHex(Arrays.copyOf(serial, n)), e.getKey());
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
