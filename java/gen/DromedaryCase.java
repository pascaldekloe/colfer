package gen;

// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf.


/**
 * DromedaryCase mixes name conventions.
 * Its serial size has a natural limit.
 * @author generated by colf(1)
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
@SuppressWarnings(
value = "fallthrough"
)
public class DromedaryCase
implements java.io.Serializable {

	/** The lower boundary on output bytes. */
	public static int MARSHAL_MIN = 5;
	/** The upper boundary on output bytes. */
	public static int MARSHAL_MAX = 5 + 8;
	/** The lower boundary on input bytes. */
	public static int UNMARSHAL_MIN = 4;
	/** The upper boundary on input bytes. */
	public static int UNMARSHAL_MAX = 4096;
	/** The lower boundary for byte capacity on in and output buffers. */
	public static int BUF_MIN = (5 + 8 + 7) & ~7;

	/**
	 * title-case option
	 */
	@Deprecated()
	// @javax.validation.constraints.NotNull
	public int pascalCase;

	/**
	 * best-case scenario
	 */
	public byte withSnake;

	private static final long[] COLFER_MASKS = {
		0,
		0xffL,
		0xffffL,
		0xffffffL,
		0xffffffffL,
		0xffffffffffL,
		0xffffffffffffL,
		0xffffffffffffffL,
		0xffffffffffffffffL,
	};

	private static final sun.misc.Unsafe java_unsafe;

	static {
		try {
			java.lang.reflect.Field f = Class.class.forName("sun.misc.Unsafe").getDeclaredField("theUnsafe");
			f.setAccessible(true);
			java_unsafe = (sun.misc.Unsafe) f.get(null);
		} catch (Exception e) {
			throw new Error("Java unsafe API required", e);
		}
	}

	/** Default constructor. */
	public DromedaryCase() { }

	/** {@link java.io.InputStream} reader. */
	public static class Unmarshaller {

		/** The data source. */
		private final java.io.InputStream in;

		/** The read buffer. */
		private final byte[] buf;

		/** The start index in {@link #buf}. */
		private int off;

		/** The number of bytes in {@link #buf} (since {@link #off}). */
		private int len;


		/**
		 * Deserializes the following object.
		 * @param in the data source.
		 * @param bufn the buffer size in bytes.
		 */
		public Unmarshaller(java.io.InputStream in, int bufn) {
			this.in = in;
			this.buf = new byte[bufn < UNMARSHAL_MAX ? UNMARSHAL_MAX : bufn];
		}

		/**
		 * Unmarshals next in line.
		 * @return the result or {@code null} when EOF.
		 * @throws java.io.IOException from the {@code java.io.InputStream}.
		 * @throws java.io.EOFException on a partial record.
		 * @throws java.io.StreamCorruptedException when the data does not match this object's schema.
		 */
		public DromedaryCase nextOrNull() throws java.io.IOException {
			if (len == 0) {
				off = 0;
				if (!read()) return null; // EOF
			} else if (buf.length - off < BUF_MIN) {
				System.arraycopy(buf, off, buf, 0, len);
				off = 0;
			}

			DromedaryCase o = new DromedaryCase();
			while (true) {
				int size = o.unmarshal(buf, off, len);
				if (size > 3) {
					off += size;
					len -= size;
					return o;
				}
				if (size != 0)
					throw new java.io.StreamCorruptedException("illegal Colfer encoding");
				if (off != 0) {
					System.arraycopy(buf, off, buf, 0, len);
					off = 0;
				}
				if (!read())
					throw new java.io.EOFException("partial Colfer encoding");
			}
		}

		/** Buffer more data. The return is {@code false} on EOF. */
		private boolean read() throws java.io.IOException {
			int pos = this.off + this.len;
			int n = in.read(buf, pos, buf.length - pos);
			if (n < 0) return false;
			this.len += n;
			return true;
		}

	}

	/**
	 * Writes a Colfer encoding to the buffer. The serial size is guaranteed
	 * with {@link #MARSHAL_MIN} and {@link #MARSHAL_MAX}. Marshal may write
	 * anywhere beyond the offset—not limited to the serial size.
	 *
	 * @param buf the output buffer.
	 * @param off the start index [offset] in the buffer.
	 * @return the encoding size.
	 * @throws IllegalArgumentException when the buffer capacity since the
	 *         offset is less than {@link BUF_MIN}.
	 * @throws java.nio.BufferOverflowException when the data exceeds the
	 *         buffer capacity or {@link #MARSHAL_MAX}.
	 */
	public int marshal(byte[] buf, int off) {
		if (off < 0 || buf.length - off < BUF_MIN)
			throw new IllegalArgumentException("output buffer space less than BUF_MIN");

		int w = off + 5; // write index
		long word0 = 5 << 12;

		// pack .PascalCase int32
		long v0 = Integer.toUnsignedLong(this.pascalCase>>31 ^ this.pascalCase<<1);
		if (v0 < 128) {
			v0 = v0 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_LONG_BASE_OFFSET, v0);
			int bitCount = 64 - Long.numberOfLeadingZeros(v0);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v0 >>>= (tailSize << 3) - 1;
			v0 = (v0 | 1L) << tailSize;
		}
		word0 |= v0 << 24;

		// pack .with_snake opaque8
		word0 |= Byte.toUnsignedLong(this.withSnake) << 32;

		// write fixed positions
		int size = w - off;
		word0 |= size;
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_LONG_BASE_OFFSET + (0 * 8) + 0, (byte)(word0 >>> (0 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_LONG_BASE_OFFSET + (0 * 8) + 1, (byte)(word0 >>> (1 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_LONG_BASE_OFFSET + (0 * 8) + 2, (byte)(word0 >>> (2 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_LONG_BASE_OFFSET + (0 * 8) + 3, (byte)(word0 >>> (3 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_LONG_BASE_OFFSET + (0 * 8) + 4, (byte)(word0 >>> (4 * 8)));
		return size;
	}

	/**
	 * Reads a Colfer encoding from the buffer. Objects can be reused. All
	 * fields are initialized regardless of their value beforehand.
	 *
	 * The number of bytes read is guaranteed to lie within in the range of
	 * [{@link #UNMARSHAL_MIN}..{@link #UNMARSHAL_MAX}]. Return {@code 1}
	 * signals malformed data. Return {@code 0} signals incomplete data,
	 * a.k.a. end-of-file.
	 *
	 * Data selection within the buffer, including its exceptions, matches
	 * Java's standard {@link java.io.InputStream#read(byte[],int,int) read}
	 * and {@link java.io.OutputStream#write(byte[],int,int) write}.
	 *
	 * @param buf the input buffer.
	 * @param off the start index [offset] in the buffer.
	 * @param len the number of bytes available since the offset.
	 * @return either the encoding size, or 0 for EOF, or 1 for malformed.
	 * @throws IllegalArgumentException when the buffer capacity minus its
	 *         offset is less than {@link #BUF_MIN}.
	 * @throws IndexOutOfBoundsException when the buffer capacity does not
	 *         match the offset–length combination.
	 */
	public int unmarshal(byte[] buf, int off, int len) {
		if ((off | len) < 0 || buf.length - off < len)
			throw new IndexOutOfBoundsException("range beyond buffer dimensions");
		if (buf.length - off < BUF_MIN)
			throw new IllegalArgumentException("insufficient buffer capacity");
		if (len < 3) return 0;
		final long word0 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_LONG_BASE_OFFSET + (0L * 8L));

		final int size = (int)word0 & 0xfff;
		final int fixedSize = (int)(word0 >> 12) & 0xfff;
		if (size < fixedSize || fixedSize < 4) return 1;
		if (size > len) return 0;

		// read index at variable section
		int r = off + fixedSize;
		// unpack .PascalCase int32
		long v0 = word0 >> (24 + 1) & 0x7f;
		if ((1L << 24 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_LONG_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v0 | 0x80) + 1;
			r += tailSize;
			v0 <<= (tailSize << 3) - tailSize;
			v0 |= tail & COLFER_MASKS[tailSize];
		}
		this.pascalCase = (int)(v0 >>> 1) ^ -(int)(v0 & 1L);
		// unpack .with_snake opaque8
		this.withSnake = (byte)(word0 >> 32);

		// TODO: clear/undo absent fields

		return size;
	}

	/**
	 * {@link java.io.Serializable} version number reflects the fields present.
	 * Values in range [0, 127] belong to Colfer version 1.
	 */
	private static final long serialVersionUID = 5L << 7;

	/**
	 * {@link java.io.Serializable} as Colfer.
	 * @param out serial destination.
	 * @throws java.io.IOException a {@link java.io.WriteAbortedException}
	 * @throws java.io.IOException either an 
	 */
	private void writeObject(java.io.ObjectOutputStream out) throws java.io.IOException {
		byte[] buf = new byte[MARSHAL_MAX];
		int n = marshal(buf, 0);
		try {
			out.write(buf, 0, n);
		} catch (java.io.IOException e) {
			throw new java.io.WriteAbortedException("halt on Colfer payload", e);
		}
	}

	/**
	 * {@link java.io.Serializable} as Colfer.
	 * @param in serial source.
	 * @throws ClassNotFoundException never.
	 * @throws java.io.IOException either from {@code in} or a
	 *  {@link java.io.StreamCorruptedException}.
	 */
	private void readObject(java.io.ObjectInputStream in)
	throws ClassNotFoundException, java.io.IOException {
		byte[] buf = new byte[UNMARSHAL_MAX];
		in.readFully(buf, 0, UNMARSHAL_MIN);
		int size = (buf[0] & 0xff) | (buf[1] & 0xf) << 8;
		in.readFully(buf, UNMARSHAL_MIN, size - UNMARSHAL_MIN);
		if (unmarshal(buf, 0, size) != size)
			throw new java.io.StreamCorruptedException("not a dromedaryCase Colfer encoding");
	}

	/**
	 * Gets gen.dromedaryCase.PascalCase.
	 * @return the value.
	 */
	public int getPascalCase() {
		return this.pascalCase;
	}
	/**
	 * Sets gen.dromedaryCase.PascalCase.
	 * @param value the replacement.
	 */
	public void setPascalCase(int value) {
		this.pascalCase = value;
	}

	/**
	 * Sets gen.dromedaryCase.PascalCase.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public DromedaryCase withPascalCase(int value) {
		setPascalCase(value);
		return this;
	}

	/**
	 * Gets gen.dromedaryCase.with_snake.
	 * @return the value.
	 */
	public byte getWithSnake() {
		return this.withSnake;
	}
	/**
	 * Sets gen.dromedaryCase.with_snake.
	 * @param value the replacement.
	 */
	public void setWithSnake(byte value) {
		this.withSnake = value;
	}

	/**
	 * Sets gen.dromedaryCase.with_snake.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public DromedaryCase withWithSnake(byte value) {
		setWithSnake(value);
		return this;
	}

	/**
	 * Deep hash is consistent with {@link #equals(Object)}.
	 * @return the standard Java digest.
	 */
	@Override
	public final int hashCode() {
		int h = 1;
		h = h * 31 + this.pascalCase;
		h = h * 31 + (int)this.withSnake;
		return h;
	}

	/**
	 * Deep comparison is consistent with {@link #hashCode}.
	 * @param o anything, including {@code null}.
	 * @return the type and content match.
	 */
	@Override
	public final boolean equals(Object o) {
		return o instanceof DromedaryCase && equals((DromedaryCase)o);
	}

	/**
	 * Typed alternative to {@link #equals(Object)}.
	 * @param o same class or {@code null}.
	 * @return the content match.
	 */
	public final boolean equals(DromedaryCase o) {
		if (o == null) return false;
		if (o == this) return true;

		return this.pascalCase == o.pascalCase
			&& this.withSnake == o.withSnake;
	}
}
