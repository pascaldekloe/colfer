import gen.BaseTypes;
import gen.ListTypes;
import gen.ManyFlags;
import gen.OpaqueTypes;

import java.io.ByteArrayOutputStream;
import java.io.ByteArrayInputStream;
import java.io.ObjectInputStream;
import java.io.ObjectOutputStream;
import java.time.Instant;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Set;


public class test {

	private static boolean testSuccess = true;

	private static void failf(String format, Object... args) {
		System.err.printf(format + "\n", args);
		testSuccess = false;
	}

	public static void main(String[] args) {
		// core
		try {
			identity();
			marshaling();
			unmarshaling();
			bitFields();
		} catch (Exception e) {
			e.printStackTrace();
			System.exit(1);
		}

		// I/O streams
		try {
			streaming();
			serializable();
		} catch (Exception e) {
			e.printStackTrace();
			System.exit(1);
		}

		// serial boundaries
		try {
			marshalMax();
			// TODO: unmarshalMax();
			bufferOverflow();
		} catch (Exception e) {
			e.printStackTrace();
			System.exit(1);
		}

		if (! testSuccess) System.exit(2);
		System.err.println("all tests passed");
	}

	/** Instantiates test cases per serial hex as key. */
	static Map<String, BaseTypes> newGoldenBaseTypes() {
		return new LinkedHashMap<String, BaseTypes>() {{

			// all zero
			put("088110"
				+ "00" // uint8
				+ "00" // int8
				+ "01" // uint16
				+ "01" // int16
				+ "01" // uint32
				+ "01" // int32
				+ "01" // uint64
				+ "01" // int64
				+ "00000000" // float32
				+ "0000000000000000" // float64
				+ "0000000000000000" // timestamp
				+ "00" // text size
				+ "00", // bool
				new BaseTypes()
			);

			// small values
			put("508110"
				+ "01" // uint8
				+ "02" // int8
				+ "07" // uint16
				+ "11" // int16
				+ "0b" // uint32
				+ "19" // int32
				+ "0f" // uint64
				+ "21" // int64
				+ "00002041" // float32
				+ "0000000000002640" // float64
				+ "0d00000003000000" // timestamp
				+ "09" // text size
				+ "01" // bool
				+ "c280" // payload U+0080
				+ "e0a080" // payload U+0800
				+ "f0908080", // payload U+10000
				new BaseTypes()
					.withU8((byte)1)
					.withI8((byte)2)
					.withU16((short)3)
					.withI16((short)4)
					.withU32(5)
					.withI32(6)
					.withU64(7L)
					.withI64(8L)
					.withF32(10)
					.withF64(11)
					.withT(Instant.ofEpochSecond(12L, 13L))
					.withS("\u0080\u0800\ud800\udc00")
					.withB(true)
			);

			// large values
			put("388210"
				+ "ff" // uint8
				+ "7f" // int8
				+ "04" // uint16
				+ "04" // int16
				+ "10" // uint32
				+ "10" // int32
				+ "00" // uint64
				+ "00" // int64
				+ "ffff7f7f" // float32
				+ "ffffffffffffef7f" // float64
				+ "ffc99afbffffffff" // timestamp
				+ "0a" // text size
				+ "00" // bool (has no large value)
				+ "ffff" // overflow 65535
				+ "feff" // overflow 32767
				+ "ffffffff" // overflow 4294967295
				+ "feffffff" // overflow 2147483647
				+ "ffffffffffffffff" // overflow 18446744073709551615
				+ "feffffffffffffff" // overflow 9223372036854775807
				+ "7f" // payload U+007F
				+ "dfbf" // payload U+07FF
				+ "efbfbf" // payload U+FFFF
				+ "f48fbfbf", // payload U+10FFFF
				new BaseTypes()
					.withU8((byte)-1)
					.withI8((byte)127)
					.withU16((short)-1)
					.withI16((short)32767)
					.withU32(-1)
					.withI32(2147483647)
					.withU64(-1L)
					.withI64(9223372036854775807L)
					.withF32(Float.MAX_VALUE)
					.withF64(Double.MAX_VALUE)
					.withT(Instant.ofEpochSecond((1L << 34) - 1L, 999999999L))
					.withS("\u007f\u07ff\uffff\udbff\udfff")
			);

		}};
	}

	static void identity() {
		if (new BaseTypes().equals((Object) null))
			failf("BaseTypes equals null Object");
		if (new BaseTypes().equals((BaseTypes) null))
			failf("BaseTypes equals null BaseTypes");

		BaseTypes[] a = newGoldenBaseTypes().values().toArray(new BaseTypes[0]);
		BaseTypes[] b = newGoldenBaseTypes().values().toArray(new BaseTypes[0]);
		if (! Arrays.equals(a, b))
			failf("identity: golden BaseTypes not equal to self");
		if (Arrays.hashCode(a) != Arrays.hashCode(b))
			failf("identity: golden BaseTypes hash inconsistent");
	}

	static void marshaling() throws Exception {
		for (Entry<String, BaseTypes> e : newGoldenBaseTypes().entrySet()) {
			byte[] buf = new byte[BaseTypes.MARSHAL_MAX];
			int n = e.getValue().marshal(buf, 0);
			StringBuilder hex = new StringBuilder(n * 2);
			for (int i = 0; i < n; i++)
				hex.append(String.format("%02x", buf[i]));
			String got = hex.toString();
			String want = e.getKey();
			if (!got.equals(want))
				failf("marshaling: got serial 0x%s, want 0x%s", got, want);
		}
	}

	static void unmarshaling() {
		for (Entry<String, BaseTypes> golden : newGoldenBaseTypes().entrySet()) {
			BaseTypes want = golden.getValue();
			String hex = golden.getKey();
			byte[] buf = new byte[BaseTypes.COLFER_MAX];
			fromHex(buf, hex);

			BaseTypes got = new BaseTypes();
			int n = got.unmarshal(buf, 0, hex.length() / 2);
			if (n != hex.length() / 2)
				failf("unmarshaling: read %d bytes of serial 0x%s", n, hex);
			if (!got.equals(want)) {
				failf("unmarshaling: mismatch for serial 0x%s\ngot:", hex);
				dumpBaseTypes(got);
				failf("want:");
				dumpBaseTypes(want);
			}
		}
	}

	static void bitFields() {
		ManyFlags a = new ManyFlags();
		ManyFlags b = new ManyFlags();
		byte[] buf = new byte[ManyFlags.COLFER_MAX];

		// tests all possible boolean combinations
		for (int i = 0; i < ManyFlags.B17_FLAG; i++) {
			a._flags = i;

			int writen = a.marshal(buf, 0);
			if (writen < 4) {
				failf("bit fields: test abort on marshal error");
				return;
			}
			int readn = b.unmarshal(buf, 0, writen);
			if (readn != writen)
				failf("bit fields: marshal wrote %d bytes, unmarshal read %d bytes", writen, readn);
			if (b._flags != i)
				failf("bit fields: unmarshalled flags %x, marshalled %x", b._flags, i);
		}
	}

	static void streaming() throws Exception {
		ByteArrayOutputStream out = new ByteArrayOutputStream();

		byte[] buf = new byte[BaseTypes.COLFER_MAX];
		for (BaseTypes o : newGoldenBaseTypes().values()) {
			int n = o.marshal(buf, 0);
			if (n == 0) {
				failf("streaming: test abort on marshal error");
				return;
			}
			out.write(buf, 0, n);
		}

		BaseTypes.Unmarshaller unmarshaller = new BaseTypes.Unmarshaller(new ByteArrayInputStream(out.toByteArray()), 0);
		for (Entry<String, BaseTypes> golden : newGoldenBaseTypes().entrySet()) {
			BaseTypes got = unmarshaller.nextOrNull();
			if (got == null) {
				failf("streaming: unmarshal ended before serial 0x%s", golden.getKey());
				return;
			}
			if (! golden.getValue().equals(got)) {
				failf("streaming: unmarshal mismatch for serial 0x%s\ngot:", golden.getKey());
				dumpBaseTypes(got);
				failf("want:");
				dumpBaseTypes(golden.getValue());
				return;
			}
		}

		BaseTypes got = unmarshaller.nextOrNull();
		if (got != null) {
			failf("stream: unmarshal got an additional object:");
			dumpBaseTypes(got);
			return;
		}
	}

	static void serializable() throws Exception {
		Set<Entry<String, BaseTypes>> cases = newGoldenBaseTypes().entrySet();
		ByteArrayOutputStream buf = new ByteArrayOutputStream();

		ObjectOutputStream out = new ObjectOutputStream(buf);
		for (Entry<String, BaseTypes> e : cases)
			out.writeObject(e.getValue());
		out.close();

		ObjectInputStream in = new ObjectInputStream(new ByteArrayInputStream(buf.toByteArray()));
		for (Entry<String, BaseTypes> golden : cases) {
			BaseTypes got = (BaseTypes) in.readObject();
			BaseTypes want = golden.getValue();
			if (!got.equals(want)) {
				failf("serializable: mismatch for serial 0x%s\ngot:", golden.getKey());
				dumpBaseTypes(got);
				failf("want:");
				dumpBaseTypes(want);
			}
		}
	}

	static void marshalMax() {
		int n;

		n = new OpaqueTypes()
			.withA16l(new short[OpaqueTypes.COLFER_MAX / 2])
			.marshal(new byte[OpaqueTypes.COLFER_MAX], 0);
		if (n != 0)
			failf("marshal max: marshalled an oversized opaque16 binary into %d bytes", n);

		n = new BaseTypes()
			.withS(new String(new char[BaseTypes.COLFER_MAX]))
			.marshal(new byte[BaseTypes.COLFER_MAX + 1], 1);
		if (n != 0)
			failf("marshal max: marshalled an oversized text into %d bytes", n);

		n = new ListTypes()
			.withF32l(new float[ListTypes.COLFER_MAX / 4])
			.marshal(new byte[ListTypes.COLFER_MAX + 2], 2);
		if (n != 0)
			failf("marshal max: marshalled an oversized float32-list into %d bytes", n);
	}

	static void bufferOverflow() {
		int n;

		n = new BaseTypes()
			.withS(new String(new char[BaseTypes.COLFER_MAX / 2]))
			.marshal(new byte[BaseTypes.MARSHAL_BUF_MIN], 0);
		if (n != 0)
			failf("marshal max: marshaled an oversized text into %d bytes", n);

		// again with offset
		n = new BaseTypes()
			.withS(new String(new char[BaseTypes.COLFER_MAX / 2]))
			.marshal(new byte[BaseTypes.MARSHAL_BUF_MIN + 99], 99);
		if (n != 0)
			failf("marshal max: marshaled an oversized text into %d bytes", n);
	}

	private static void fromHex(byte[] buf, String s) {
		int len = s.length();
		if (len % 2 != 0)
			throw new IllegalArgumentException("odd number of hexadecimals");
		if (len / 2 > buf.length)
			throw new IllegalArgumentException("hex exceeds buffer capacity");
		for (int i = 0; i < len; i += 2) {
			int msn = Character.digit(s.charAt(i), 16);
			int lsn = Character.digit(s.charAt(i + 1), 16);
			buf[i / 2] = (byte) ((msn << 4) + lsn);
		}
	}

	private static void dumpBaseTypes(BaseTypes o) {
		System.err.printf("{ u8=%d i8=%d", o.u8, o.i8);
		System.err.printf(" u16=%d i16=%d", o.u16, o.i16);
		System.err.printf(" u32=%d i32=%d", o.u32, o.i32);
		System.err.printf(" u64=%d i64=%d", o.u64, o.i64);
		System.err.printf(" f32=%f f64=%f", o.f32, o.f64);
		System.err.printf(" t=%s s=0x", o.t);
		int utf16n = o.s.length();
		for (int i = 0; i < utf16n; i++)
			 System.err.printf("%04x", (short)o.s.charAt(i));
		System.err.printf(" b=%b }\n", o.getB());
	}

}
