package gen;

// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf.


/**
 * ListTypes contains each type supported in list form.
 * @author generated by colf(1)
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class ListTypes
implements java.io.Serializable {

	/** The lower boundary on output bytes. */
	public static int MARSHAL_MIN = 11;
	/** The upper boundary on output bytes. */
	public static int MARSHAL_MAX = 4096;
	/** The lower boundary on input bytes. */
	public static int UNMARSHAL_MIN = 4;
	/** The upper boundary on input bytes. */
	public static int UNMARSHAL_MAX = 4096;
	/** The lower boundary for byte capacity on in and output buffers. */
	public static int BUF_MIN = (11 + 64 + 7) & ~7;

	/**
	 * Test 8-bit values.
	 */
	public byte[] a8l = zero_a8l;

	/**
	 * Test 16-bit values.
	 */
	public short[] a16l = zero_a16l;

	/**
	 * Test 32-bit values.
	 */
	public int[] a32l = zero_a32l;

	/**
	 * Test 64-bit values.
	 */
	public long[] a64l = zero_a64l;

	/**
	 * Test single precision–floating points.
	 */
	public float[] f32l = zero_f32l;

	/**
	 * Test double precision–floating points.
	 */
	public double[] f64l = zero_f64l;

	/**
	 * Test timestamps (with nanosecond precision).
	 */
	public java.time.Instant[] tl = zero_tl;

	/**
	 * Test Unicode strings of variable size.
	 */
	public String[] sl = zero_sl;
	private static final byte[] zero_a8l = new byte[0];
	private static final short[] zero_a16l = new short[0];
	private static final int[] zero_a32l = new int[0];
	private static final long[] zero_a64l = new long[0];
	private static final float[] zero_f32l = new float[0];
	private static final double[] zero_f64l = new double[0];
	private static final java.time.Instant[] zero_tl = new java.time.Instant[0];
	private static final String[] zero_sl = new String[0];

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
	public ListTypes() { }

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
		public ListTypes nextOrNull() throws java.io.IOException {
			if (len == 0) {
				off = 0;
				if (!read()) return null; // EOF
			} else if (buf.length - off < BUF_MIN) {
				System.arraycopy(buf, off, buf, 0, len);
				off = 0;
			}

			ListTypes o = new ListTypes();
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
	public int marshalWithBounds(byte[] buf, int off) {
		if (off < 0 || buf.length - off < BUF_MIN)
			throw new IllegalArgumentException("output buffer space less than BUF_MIN");

		int w = off + 11; // write index
		long word0 = 11 << 12;

		// pack .a8l []opaque8
		long v0 = Integer.toUnsignedLong(this.a8l.length);
		if (v0 < 128) {
			v0 = v0 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v0);
			int bitCount = 64 - Long.numberOfLeadingZeros(v0);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v0 >>>= (tailSize << 3) - 1;
			v0 = (v0 | 1L) << tailSize;
		}
		word0 |= v0 << 24;

		// pack .a16l []opaque16
		long v1 = Integer.toUnsignedLong(this.a16l.length);
		if (v1 < 128) {
			v1 = v1 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v1);
			int bitCount = 64 - Long.numberOfLeadingZeros(v1);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v1 >>>= (tailSize << 3) - 1;
			v1 = (v1 | 1L) << tailSize;
		}
		word0 |= v1 << 32;

		// pack .a32l []opaque32
		long v2 = Integer.toUnsignedLong(this.a32l.length);
		if (v2 < 128) {
			v2 = v2 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v2);
			int bitCount = 64 - Long.numberOfLeadingZeros(v2);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v2 >>>= (tailSize << 3) - 1;
			v2 = (v2 | 1L) << tailSize;
		}
		word0 |= v2 << 40;

		// pack .a64l []opaque64
		long v3 = Integer.toUnsignedLong(this.a64l.length);
		if (v3 < 128) {
			v3 = v3 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v3);
			int bitCount = 64 - Long.numberOfLeadingZeros(v3);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v3 >>>= (tailSize << 3) - 1;
			v3 = (v3 | 1L) << tailSize;
		}
		word0 |= v3 << 48;

		// pack .f32l []float32
		long v4 = Integer.toUnsignedLong(this.f32l.length);
		if (v4 < 128) {
			v4 = v4 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v4);
			int bitCount = 64 - Long.numberOfLeadingZeros(v4);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v4 >>>= (tailSize << 3) - 1;
			v4 = (v4 | 1L) << tailSize;
		}
		word0 |= v4 << 56;

		// pack .f64l []float64
		long v5 = Integer.toUnsignedLong(this.f64l.length);
		if (v5 < 128) {
			v5 = v5 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v5);
			int bitCount = 64 - Long.numberOfLeadingZeros(v5);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v5 >>>= (tailSize << 3) - 1;
			v5 = (v5 | 1L) << tailSize;
		}
		long word1 = v5;

		// pack .tl []timestamp
		long v6 = Integer.toUnsignedLong(this.tl.length);
		if (v6 < 128) {
			v6 = v6 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v6);
			int bitCount = 64 - Long.numberOfLeadingZeros(v6);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v6 >>>= (tailSize << 3) - 1;
			v6 = (v6 | 1L) << tailSize;
		}
		word1 |= v6 << 8;

		// pack .sl []text
		long v7 = Integer.toUnsignedLong(this.sl.length);
		if (v7 < 128) {
			v7 = v7 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v7);
			int bitCount = 64 - Long.numberOfLeadingZeros(v7);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v7 >>>= (tailSize << 3) - 1;
			v7 = (v7 | 1L) << tailSize;
		}
		word1 |= v7 << 16;

		// write payloads
		// TODO: implement text list
		if (buf.length - w < this.tl.length << 3)
			throw new java.nio.BufferOverflowException();
		for (java.time.Instant t : this.tl) {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, t.getEpochSecond() << 30 | Integer.toUnsignedLong(t.getNano()));
			w += 8;
		}
		if (buf.length - w < this.f64l.length << 3)
			throw new java.nio.BufferOverflowException();
		for (double d : this.f64l) {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, Double.doubleToRawLongBits(d));
			w += 8;
		}
		if (buf.length - w < this.f32l.length << 2)
			throw new java.nio.BufferOverflowException();
		for (float f : this.f32l) {
			java_unsafe.putInt(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, Float.floatToRawIntBits(f));
			w += 4;
		}
		if (buf.length - w < this.a64l.length << 3)
			throw new java.nio.BufferOverflowException();
		for (long b : this.a64l) {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, b);
			w += 8;
		}
		if (buf.length - w < this.a32l.length << 2)
			throw new java.nio.BufferOverflowException();
		for (int b : this.a32l) {
			java_unsafe.putInt(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, b);
			w += 4;
		}
		if (buf.length - w < this.a16l.length << 1)
			throw new java.nio.BufferOverflowException();
		for (short b : this.a16l) {
			java_unsafe.putShort(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, b);
			w += 2;
		}
		if (buf.length - w < this.a8l.length)
			throw new java.nio.BufferOverflowException();
		System.arraycopy(this.a8l, 0, buf, w, this.a8l.length);
		w += this.a8l.length;

		// write fixed positions
		int size = w - off;
		if (size > MARSHAL_MAX)
			throw new java.nio.BufferOverflowException();
		word0 |= size;
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (0 * 8), word0);
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (1 * 8) + 0, (byte)(word1 >>> (0 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (1 * 8) + 1, (byte)(word1 >>> (1 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (1 * 8) + 2, (byte)(word1 >>> (2 * 8)));
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
		final long word0 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (0L * 8L));
		final long word1 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (1L * 8L));

		final int size = (int)word0 & 0xfff;
		final int fixedSize = (int)(word0 >> 12) & 0xfff;
		if (size < fixedSize || fixedSize < 4) return 1;
		if (size > len) return 0;
		// payloads read backwards
		int end = off + size;

		// read index at variable section
		int r = off + fixedSize;
		// unpack .a8l []opaque8
		long v0 = word0 >> (24 + 1) & 0x7f;
		if ((1L << 24 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v0 | 0x80) + 1;
			r += tailSize;
			v0 <<= (tailSize << 3) - tailSize;
			v0 |= tail & COLFER_MASKS[tailSize];
		}
		final long v0Size = v0 * 1;
		if ((long)(end - r) < v0Size)
			return 1;
		this.a8l = new byte[(int) v0];
		end -= this.a8l.length;
		System.arraycopy(buf, end, this.a8l, 0, this.a8l.length);
		// unpack .a16l []opaque16
		long v1 = word0 >> (32 + 1) & 0x7f;
		if ((1L << 32 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v1 | 0x80) + 1;
			r += tailSize;
			v1 <<= (tailSize << 3) - tailSize;
			v1 |= tail & COLFER_MASKS[tailSize];
		}
		final long v1Size = v1 * 1;
		if ((long)(end - r) < v1Size)
			return 1;
		this.a16l = new short[(int) v1];
		// unpack .a32l []opaque32
		long v2 = word0 >> (40 + 1) & 0x7f;
		if ((1L << 40 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v2 | 0x80) + 1;
			r += tailSize;
			v2 <<= (tailSize << 3) - tailSize;
			v2 |= tail & COLFER_MASKS[tailSize];
		}
		final long v2Size = v2 * 1;
		if ((long)(end - r) < v2Size)
			return 1;
		this.a32l = new int[(int) v2];
		// unpack .a64l []opaque64
		long v3 = word0 >> (48 + 1) & 0x7f;
		if ((1L << 48 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v3 | 0x80) + 1;
			r += tailSize;
			v3 <<= (tailSize << 3) - tailSize;
			v3 |= tail & COLFER_MASKS[tailSize];
		}
		final long v3Size = v3 * 1;
		if ((long)(end - r) < v3Size)
			return 1;
		this.a64l = new long[(int) v3];
		// unpack .f32l []float32
		long v4 = word0 >> (56 + 1) & 0x7f;
		if ((1L << 56 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v4 | 0x80) + 1;
			r += tailSize;
			v4 <<= (tailSize << 3) - tailSize;
			v4 |= tail & COLFER_MASKS[tailSize];
		}
		final long v4Size = v4 * 1;
		if ((long)(end - r) < v4Size)
			return 1;
		this.f32l = new float[(int) v4];
		// unpack .f64l []float64
		long v5 = word1 >> (0 + 1) & 0x7f;
		if ((1L << 0 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v5 | 0x80) + 1;
			r += tailSize;
			v5 <<= (tailSize << 3) - tailSize;
			v5 |= tail & COLFER_MASKS[tailSize];
		}
		final long v5Size = v5 * 1;
		if ((long)(end - r) < v5Size)
			return 1;
		this.f64l = new double[(int) v5];
		// unpack .tl []timestamp
		long v6 = word1 >> (8 + 1) & 0x7f;
		if ((1L << 8 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v6 | 0x80) + 1;
			r += tailSize;
			v6 <<= (tailSize << 3) - tailSize;
			v6 |= tail & COLFER_MASKS[tailSize];
		}
		final long v6Size = v6 * 1;
		if ((long)(end - r) < v6Size)
			return 1;
		this.tl = new java.time.Instant[(int) v6];
		// unpack .sl []text
		long v7 = word1 >> (16 + 1) & 0x7f;
		if ((1L << 16 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v7 | 0x80) + 1;
			r += tailSize;
			v7 <<= (tailSize << 3) - tailSize;
			v7 |= tail & COLFER_MASKS[tailSize];
		}
		final long v7Size = v7 * 1;
		if ((long)(end - r) < v7Size)
			return 1;
		this.sl = new String[(int) v7];

		// TODO: clear/undo absent fields

		return size;
	}

	/**
	 * {@link java.io.Serializable} version number reflects the fields present.
	 * Values in range [0, 127] belong to Colfer version 1.
	 */
	private static final long serialVersionUID = 11L << 7;

	/**
	 * {@link java.io.Serializable} as Colfer.
	 * @param out serial destination.
	 * @throws java.io.IOException a {@link java.io.WriteAbortedException}
	 *         or an {@link java.io.InvalidObjectException} when encoding
	 *         would exceed {@link #MARSHAL_MAX}.
	 * @throws java.io.IOException either an 
	 */
	private void writeObject(java.io.ObjectOutputStream out) throws java.io.IOException {
		byte[] buf = new byte[MARSHAL_MAX];
		int n = marshalWithBounds(buf, 0);
		if (n == 0) throw new java.io.InvalidObjectException("MARSHAL_MAX reached");
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
			throw new java.io.StreamCorruptedException("not a ListTypes Colfer encoding");
	}

	/**
	 * Gets gen.ListTypes.a8l.
	 * @return the value.
	 */
	public byte[] getA8l() {
		return this.a8l;
	}
	/**
	 * Sets gen.ListTypes.a8l.
	 * @param value the replacement.
	 */
	public void setA8l(byte[] value) {
		this.a8l = value;
	}

	/**
	 * Sets gen.ListTypes.a8l.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withA8l(byte[] value) {
		setA8l(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.a16l.
	 * @return the value.
	 */
	public short[] getA16l() {
		return this.a16l;
	}
	/**
	 * Sets gen.ListTypes.a16l.
	 * @param value the replacement.
	 */
	public void setA16l(short[] value) {
		this.a16l = value;
	}

	/**
	 * Sets gen.ListTypes.a16l.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withA16l(short[] value) {
		setA16l(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.a32l.
	 * @return the value.
	 */
	public int[] getA32l() {
		return this.a32l;
	}
	/**
	 * Sets gen.ListTypes.a32l.
	 * @param value the replacement.
	 */
	public void setA32l(int[] value) {
		this.a32l = value;
	}

	/**
	 * Sets gen.ListTypes.a32l.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withA32l(int[] value) {
		setA32l(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.a64l.
	 * @return the value.
	 */
	public long[] getA64l() {
		return this.a64l;
	}
	/**
	 * Sets gen.ListTypes.a64l.
	 * @param value the replacement.
	 */
	public void setA64l(long[] value) {
		this.a64l = value;
	}

	/**
	 * Sets gen.ListTypes.a64l.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withA64l(long[] value) {
		setA64l(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.f32l.
	 * @return the value.
	 */
	public float[] getF32l() {
		return this.f32l;
	}
	/**
	 * Sets gen.ListTypes.f32l.
	 * @param value the replacement.
	 */
	public void setF32l(float[] value) {
		this.f32l = value;
	}

	/**
	 * Sets gen.ListTypes.f32l.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withF32l(float[] value) {
		setF32l(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.f64l.
	 * @return the value.
	 */
	public double[] getF64l() {
		return this.f64l;
	}
	/**
	 * Sets gen.ListTypes.f64l.
	 * @param value the replacement.
	 */
	public void setF64l(double[] value) {
		this.f64l = value;
	}

	/**
	 * Sets gen.ListTypes.f64l.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withF64l(double[] value) {
		setF64l(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.tl.
	 * @return the value.
	 */
	public java.time.Instant[] getTl() {
		return this.tl;
	}
	/**
	 * Sets gen.ListTypes.tl.
	 * @param value the replacement.
	 */
	public void setTl(java.time.Instant[] value) {
		this.tl = value;
	}

	/**
	 * Sets gen.ListTypes.tl.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withTl(java.time.Instant[] value) {
		setTl(value);
		return this;
	}

	/**
	 * Gets gen.ListTypes.sl.
	 * @return the value.
	 */
	public String[] getSl() {
		return this.sl;
	}
	/**
	 * Sets gen.ListTypes.sl.
	 * @param value the replacement.
	 */
	public void setSl(String[] value) {
		this.sl = value;
	}

	/**
	 * Sets gen.ListTypes.sl.
	 * @param value the replacement.
	 * @return {@code this}.
	 */
	public ListTypes withSl(String[] value) {
		setSl(value);
		return this;
	}

	/**
	 * Deep hash is consistent with {@link #equals(Object)}.
	 * @return the standard Java digest.
	 */
	@Override
	public final int hashCode() {
		int h = 1;
		h = h * 31 + java.util.Arrays.hashCode(this.a8l);
		h = h * 31 + java.util.Arrays.hashCode(this.a16l);
		h = h * 31 + java.util.Arrays.hashCode(this.a32l);
		h = h * 31 + java.util.Arrays.hashCode(this.a64l);
		h = h * 31 + java.util.Arrays.hashCode(this.f32l);
		h = h * 31 + java.util.Arrays.hashCode(this.f64l);
		h = h * 31 + java.util.Arrays.hashCode(this.tl);
		h = h * 31 + java.util.Arrays.hashCode(this.sl);
		return h;
	}

	/**
	 * Deep comparison is consistent with {@link #hashCode}.
	 * Two not-a-number values compare equal.
	 * @param o anything, including {@code null}.
	 * @return the type and content match.
	 */
	@Override
	public final boolean equals(Object o) {
		return o instanceof ListTypes && equals((ListTypes)o);
	}

	/**
	 * Typed alternative to {@link #equals(Object)}.
	 * @param o same class or {@code null}.
	 * @return the content match.
	 */
	public final boolean equals(ListTypes o) {
		if (o == null) return false;
		if (o == this) return true;

		return java.util.Arrays.equals(this.a8l, o.a8l)
			&& java.util.Arrays.equals(this.a16l, o.a16l)
			&& java.util.Arrays.equals(this.a32l, o.a32l)
			&& java.util.Arrays.equals(this.a64l, o.a64l)
			&& java.util.Arrays.equals(this.f32l, o.f32l)
			&& java.util.Arrays.equals(this.f64l, o.f64l)
			&& java.util.Arrays.equals(this.tl, o.tl)
			&& java.util.Arrays.equals(this.sl, o.sl);
	}
}
