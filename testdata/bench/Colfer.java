package testdata.bench;

/**
 * @author Commander Colfer
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class Colfer implements java.io.Serializable {

	private static final java.nio.charset.Charset utf8 = java.nio.charset.Charset.forName("UTF-8");

	public long Key;
	public String Host;
	public byte[] Addr;
	public int Port;
	public long Size;
	public long Hash;
	public double Ratio;
	public boolean Route;


	/**
	 * Writes in Colfer format.
	 * @param buf the data destination.
	 * @throws java.nio.BufferOverflowException when {@code buf} is too small.
	 */
	public final void marshal(java.nio.ByteBuffer buf) {
		buf.order(java.nio.ByteOrder.BIG_ENDIAN);
		buf.put((byte) 0x80);

		if (Key != 0) {
			byte header = 0;
			long x = Key;
			if (x < 0) {
				x = -x;
				header |= 0x80;
			}
			buf.put(header);
			putVarint(buf, x);
		}

		if (Host != null && ! Host.isEmpty()) {
			java.nio.ByteBuffer bytes = utf8.encode(Host);
			buf.put((byte) 0x08);
			putVarint(buf, bytes.limit());
			buf.put(bytes);
		}

		if (Addr != null && Addr.length != 0) {
			buf.put((byte) 0x09);
			putVarint(buf, Addr.length);
			buf.put(Addr);
		}

		if (Port != 0) {
			byte header = 3;
			int x = Port;
			if (x < 0) {
				x = -x;
				header |= 0x80;
			}
			buf.put(header);
			putVarint(buf, x);
		}

		if (Size != 0) {
			byte header = 4;
			long x = Size;
			if (x < 0) {
				x = -x;
				header |= 0x80;
			}
			buf.put(header);
			putVarint(buf, x);
		}

		if (Hash != 0) {
			buf.put((byte) 5);
			putVarint(buf, Hash);
		}

		if (Ratio != 0.0) {
			buf.put((byte) 6);
			buf.putDouble(Ratio);
		}

		if (Route) {
			buf.put((byte) 7);
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
			Key = getVarint64(buf);
			if ((header & 0x80) != 0)
			Key = (~Key) + 1;

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 1) {
			int length = getVarint32(buf);
			java.nio.ByteBuffer blob = java.nio.ByteBuffer.allocate(length);
			buf.get(blob.array());
			Host = utf8.decode(blob).toString();

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 2) {
			int length = getVarint32(buf);
			Addr = new byte[length];
			buf.get(Addr);

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 3) {
			Port = getVarint32(buf);
			if ((header & 0x80) != 0)
			Port = (~Port) + 1;

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 4) {
			Size = getVarint64(buf);
			if ((header & 0x80) != 0)
			Size = (~Size) + 1;

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 5) {
			Hash = getVarint64(buf);

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 6) {
			Ratio = buf.getDouble();

			if (! buf.hasRemaining()) return;
			header = buf.get() & 0xff;
			field = header & 0x7f;
		}

		if (field == 7) {
			Route = true;

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
