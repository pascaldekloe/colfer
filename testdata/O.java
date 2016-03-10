package testdata;

/**
 * @author Commander Colfer
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class O implements java.io.Serializable {

	private static final java.nio.charset.Charset utf8 = java.nio.charset.Charset.forName("UTF-8");

	public boolean B;
	public int U32;
	public long U64;
	public int I32;
	public long I64;
	public float F32;
	public double F64;
	public java.time.Instant T;
	public String S;
	public byte[] A;


	/**
	 * Writes in Colfer format.
	 * @param buf the data destination.
	 * @throws java.nio.BufferOverflowException when {@code buf} is too small.
	 */
	public final void marshal(java.nio.ByteBuffer buf) {
		buf.order(java.nio.ByteOrder.BIG_ENDIAN);
		buf.put((byte) 0x80);

		if (B) {
			buf.put((byte) 0);
		}

		if (U32 != 0) {
			buf.put((byte) 1);
			putVarint(buf, U32);
		}

		if (U64 != 0) {
			buf.put((byte) 2);
			putVarint(buf, U64);
		}

		if (I32 != 0) {
			byte header = 3;
			int x = I32;
			if (x < 0) {
				x = -x;
				header |= 0x80;
			}
			buf.put(header);
			putVarint(buf, x);
		}

		if (I64 != 0) {
			byte header = 4;
			long x = I64;
			if (x < 0) {
				x = -x;
				header |= 0x80;
			}
			buf.put(header);
			putVarint(buf, x);
		}

		if (F32 != 0.0f) {
			buf.put((byte) 5);
			buf.putFloat(F32);
		}

		if (F64 != 0.0) {
			buf.put((byte) 6);
			buf.putDouble(F64);
		}

		if (T != null) {
			long s = T.getEpochSecond();
			int ns = T.getNano();
			if (! (s == 0 && ns == 0)) {
				byte header = 7;
				if (ns != 0) header |= 0x80;
				buf.put(header);
				buf.putLong(s);
				if (ns != 0) buf.putInt(ns);
			}
		}

		if (S != null && ! S.isEmpty()) {
			java.nio.ByteBuffer bytes = utf8.encode(S);
			buf.put((byte) 0x08);
			putVarint(buf, bytes.limit());
			buf.put(bytes);
		}

		if (A != null && A.length != 0) {
			buf.put((byte) 0x09);
			putVarint(buf, A.length);
			buf.put(A);
		}

	}

	/**
	 * Reads in Colfer format.
	 * @param buf the data source.
	 * @throws java.nio.BufferUnderflowException when {@code buf} is incomplete.
	 */
	public final void unmarshal(java.nio.ByteBuffer buf) {
		int header = buf.get() & 0xff;
		if (header != 0x80)
			throw new IllegalArgumentException("magic header mismatch");

		if (! buf.hasRemaining()) return;
		header = buf.get() & 0xff;
		int field = header & 0x7f;

		if (field == 0) {
			B = true;

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 1) {
			U32 = getVarint32(buf);

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 2) {
			U64 = getVarint64(buf);

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 3) {
			I32 = getVarint32(buf);
			if ((header & 0x80) != 0)
			I32 = (~I32) + 1;

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 4) {
			I64 = getVarint64(buf);
			if ((header & 0x80) != 0)
			I64 = (~I64) + 1;

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 5) {
			F32 = buf.getFloat();

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 6) {
			F64 = buf.getDouble();

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 7) {
			long s = buf.getLong();
			if ((header & 0x80) == 0) {
				T = java.time.Instant.ofEpochSecond(s);
			} else {
				int ns = buf.getInt();
				T = java.time.Instant.ofEpochSecond(s, ns);
			}

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 8) {
			int length = getVarint32(buf);
			java.nio.ByteBuffer blob = java.nio.ByteBuffer.allocate(length);
			buf.get(blob.array());
			S = utf8.decode(blob).toString();

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 9) {
			int length = getVarint32(buf);
			A = new byte[length];
			buf.get(A);

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		throw new IllegalArgumentException("pending data");
	}

	/**
	 * Serializes an integer.
	 * @param buf the data destination.
	 * @param x the value.
	 */
	private static void putVarint(java.nio.ByteBuffer buf, int x) {
		while ((x & 0xffffff80) != 0) {
			buf.put((byte) (x | 0x80));
			x >>>= 7;
		}
		buf.put((byte) x);
	}

	/**
	 * Serializes an integer.
	 * @param buf the data destination.
	 * @param x the value.
	 */
	private static void putVarint(java.nio.ByteBuffer buf, long x) {
		while ((x & 0xffffffffffffff80L) != 0) {
			buf.put((byte) (x | 0x80));
			x >>>= 7;
		}
		buf.put((byte) x);
	}

	/**
	 * Deserializes a 32-bit integer.
	 * @param buf the data source.
	 * @return the value.
	 */
	private static int getVarint32(java.nio.ByteBuffer buf) {
		int x = 0;
		int shift = 0;
		while (true) {
			int b = buf.get() & 0xff;
			x |= (b & 0x7f) << shift;
			if (b < 0x80) return x;
			shift += 7;
		}
	}

	/**
	 * Deserializes a 64-bit integer.
	 * @param buf the data source.
	 * @return the value.
	 */
	private static long getVarint64(java.nio.ByteBuffer buf) {
		long x = 0;
		int shift = 0;
		while (true) {
			long b = buf.get() & 0xffL;
			x |= (b & 0x7f) << shift;
			if (b < 0x80) return x;
			shift += 7;
		}
	}

}
