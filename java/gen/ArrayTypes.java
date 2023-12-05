package gen;

// Code generated by colf(1); DO NOT EDIT.
// The compiler used schema file test.colf.


/**
 * ArrayTypes contains each BaseType supported in array form,
 * which is all but bool(ean).
 * @author generated by colf(1)
 * @see <a href="https://github.com/pascaldekloe/colfer">Colfer's home</a>
 */
public class ArrayTypes
implements java.io.Serializable {

	/** The lower boundary on output bytes. */
	public static int MARSHAL_MIN = 61;
	/** The upper boundary on output bytes. */
	public static int MARSHAL_MAX = 4096;
	/** The lower boundary on input bytes. */
	public static int UNMARSHAL_MIN = 5;
	/** The upper boundary on input bytes. */
	public static int UNMARSHAL_MAX = 4096;
	/** The lower boundary for byte capacity on in and output buffers. */
	public static int BUF_MIN = (61 + 112 + 7) & ~7;

	/**
	 * Test 8 bit–unsigned integers. Two elements set the
	 * minimium size to 5.
	 */
	public final byte[] u8n2 = new byte[2];

	/**
	 * Test 8 bit–signed integers.
	 */
	public final byte[] i8n2 = new byte[2];

	/**
	 * Test 16 bit–unsigned integers.
	 */
	public final int[] u16n2 = new int[2];

	/**
	 * Test 16 bit–signed integers.
	 */
	public final int[] i16n2 = new int[2];

	/**
	 * Test 32 bit–unsigned integers.
	 */
	public final short[] u32n2 = new short[2];

	/**
	 * Test 32 bit–signed integers.
	 */
	public final short[] i32n2 = new short[2];

	/**
	 * Test 64 bit–unsigned integers.
	 */
	public final long[] u64n2 = new long[2];

	/**
	 * Test 64 bit–signed integers.
	 */
	public final long[] i64n2 = new long[2];

	/**
	 * Test single precision–floating points.
	 */
	public final float[] f32n2 = new float[2];

	/**
	 * Test double precision–floating points.
	 */
	public final double[] f64n2 = new double[2];

	/**
	 * Test timestamps (with nanosecond precision).
	 */
	public final java.time.Instant[] tn2 = new java.time.Instant[2];

	/**
	 * Test Unicode strings of variable size.
	 */
	public final String[] sn2 = new String[2];

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
	public ArrayTypes() { }

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
		public ArrayTypes nextOrNull() throws java.io.IOException {
			if (len == 0) {
				off = 0;
				if (!read()) return null; // EOF
			} else if (buf.length - off < BUF_MIN) {
				System.arraycopy(buf, off, buf, 0, len);
				off = 0;
			}

			ArrayTypes o = new ArrayTypes();
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
	 * Writes a Colfer encoding to the buffer. The output size is guaranteed
	 * with {@link #MARSHAL_MIN} and {@link #MARSHAL_MAX}. Return 0 signals
	 * that encoding would exceed {@link #MARSHAL_MAX}.
	 *
	 * @param buf the output buffer.
	 * @param off the start index [offset] in the buffer.
	 * @return the encoding size or 0.
	 * @throws IllegalArgumentException when the buffer capacity since the
	 *         offset is less than {@link BUF_MIN}.
	 * @throws java.nio.BufferOverflowException is prevented with a buffer
	 *         capacity since the offset of at least {@link #MARSHAL_MAX}.
	 */
	public int marshalWithBounds(byte[] buf, int off) {
		if (off < 0 || buf.length - off < BUF_MIN)
			throw new IllegalArgumentException("output buffer space less than BUF_MIN");

		int w = off + 61; // write index
		long word0 = 61 << 12;
		long variableSize = 0;

		// pack .u8n2 uint8
		word0 |= Byte.toUnsignedLong(this.u8n2[0]) << 24;
		word0 |= Byte.toUnsignedLong(this.u8n2[1]) << 32;

		// pack .i8n2 int8
		word0 |= Byte.toUnsignedLong(this.i8n2[0]) << 40;
		word0 |= Byte.toUnsignedLong(this.i8n2[1]) << 48;

		// pack .u16n2 uint32
		long v4 = Integer.toUnsignedLong(this.u16n2[0]);
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
		long v5 = Integer.toUnsignedLong(this.u16n2[1]);
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

		// pack .i16n2 int32
		long v6 = Integer.toUnsignedLong(this.i16n2[0]>>31 ^ this.i16n2[0]<<1);
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
		long v7 = Integer.toUnsignedLong(this.i16n2[1]>>31 ^ this.i16n2[1]<<1);
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

		// pack .u32n2 uint16
		long v8 = Short.toUnsignedLong(this.u32n2[0]);
		if (v8 < 128) {
			v8 = v8 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v8);
			int bitCount = 64 - Long.numberOfLeadingZeros(v8);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v8 >>>= (tailSize << 3) - 1;
			v8 = (v8 | 1L) << tailSize;
		}
		word1 |= v8 << 24;
		long v9 = Short.toUnsignedLong(this.u32n2[1]);
		if (v9 < 128) {
			v9 = v9 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v9);
			int bitCount = 64 - Long.numberOfLeadingZeros(v9);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v9 >>>= (tailSize << 3) - 1;
			v9 = (v9 | 1L) << tailSize;
		}
		word1 |= v9 << 32;

		// pack .i32n2 int16
		long v10 = Integer.toUnsignedLong(this.i32n2[0]>>15 ^ this.i32n2[0]<<1);
		if (v10 < 128) {
			v10 = v10 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v10);
			int bitCount = 64 - Long.numberOfLeadingZeros(v10);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v10 >>>= (tailSize << 3) - 1;
			v10 = (v10 | 1L) << tailSize;
		}
		word1 |= v10 << 40;
		long v11 = Integer.toUnsignedLong(this.i32n2[1]>>15 ^ this.i32n2[1]<<1);
		if (v11 < 128) {
			v11 = v11 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v11);
			int bitCount = 64 - Long.numberOfLeadingZeros(v11);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v11 >>>= (tailSize << 3) - 1;
			v11 = (v11 | 1L) << tailSize;
		}
		word1 |= v11 << 48;

		// pack .u64n2 uint64
		long v12 = this.u64n2[0];
		if (v12 < 128) {
			v12 = v12 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v12);
			int bitCount = 64 - Long.numberOfLeadingZeros(v12);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v12 >>>= (tailSize << 3) - 1;
			v12 = (v12 | 1L) << tailSize;
		}
		word1 |= v12 << 56;
		long v13 = this.u64n2[1];
		if (v13 < 128) {
			v13 = v13 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v13);
			int bitCount = 64 - Long.numberOfLeadingZeros(v13);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v13 >>>= (tailSize << 3) - 1;
			v13 = (v13 | 1L) << tailSize;
		}
		long word2 = v13;

		// pack .i64n2 int64
		long v14 = this.i64n2[0]>>63 ^ this.i64n2[0]<<1;
		if (v14 < 128) {
			v14 = v14 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v14);
			int bitCount = 64 - Long.numberOfLeadingZeros(v14);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v14 >>>= (tailSize << 3) - 1;
			v14 = (v14 | 1L) << tailSize;
		}
		word2 |= v14 << 8;
		long v15 = this.i64n2[1]>>63 ^ this.i64n2[1]<<1;
		if (v15 < 128) {
			v15 = v15 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v15);
			int bitCount = 64 - Long.numberOfLeadingZeros(v15);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v15 >>>= (tailSize << 3) - 1;
			v15 = (v15 | 1L) << tailSize;
		}
		word2 |= v15 << 16;

		// pack .f32n2 float32
		long v16 = Integer.toUnsignedLong(Float.floatToRawIntBits(this.f32n2[0]));
		word2 |= v16 << 24;
		long v17 = Integer.toUnsignedLong(Float.floatToRawIntBits(this.f32n2[1]));
		word2 |= v17 << 56;
		long word3 = v17 >> (64-56);

		// pack .f64n2 float64
		long v18 = Double.doubleToRawLongBits(this.f64n2[0]);
		word3 |= v18 << 24;
		long word4 = v18 >> (64-24);
		long v19 = Double.doubleToRawLongBits(this.f64n2[1]);
		word4 |= v19 << 24;
		long word5 = v19 >> (64-24);

		// pack .tn2 timestamp
		long v20 = this.tn2[0].getEpochSecond() << 30 | Integer.toUnsignedLong(this.tn2[0].getNano());
		word5 |= v20 << 24;
		long word6 = v20 >> (64-24);
		long v21 = this.tn2[1].getEpochSecond() << 30 | Integer.toUnsignedLong(this.tn2[1].getNano());
		word6 |= v21 << 24;
		long word7 = v21 >> (64-24);

		// pack .sn2 text
		long v22 = this.sn2[0].length();
		for (int i = 0, end = (int) v22; i < end; i++) {
			char c = this.sn2[0].charAt(i);
			if (c < '\u0080') continue; // 1 char to 1 UTF-8 byte
			if (c < '\u0800') v22++; // 1 char to 2 UTF-8 bytes
			else if (! Character.isHighSurrogate(c)) v22 += 2; // 1 char to 3 UTF-8 bytes
			else if (i + 1 < end && Character.isLowSurrogate(this.sn2[0].charAt(i + 1))) {
				v22 += 2; // 2 chars to 4 UTF-8 bytes
				i++;
			} // else broken surrogate replaced with '?'
		}
		variableSize += v22;
		if (v22 < 128) {
			v22 = v22 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v22);
			int bitCount = 64 - Long.numberOfLeadingZeros(v22);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v22 >>>= (tailSize << 3) - 1;
			v22 = (v22 | 1L) << tailSize;
		}
		word7 |= v22 << 24;
		long v23 = this.sn2[1].length();
		for (int i = 0, end = (int) v23; i < end; i++) {
			char c = this.sn2[1].charAt(i);
			if (c < '\u0080') continue; // 1 char to 1 UTF-8 byte
			if (c < '\u0800') v23++; // 1 char to 2 UTF-8 bytes
			else if (! Character.isHighSurrogate(c)) v23 += 2; // 1 char to 3 UTF-8 bytes
			else if (i + 1 < end && Character.isLowSurrogate(this.sn2[1].charAt(i + 1))) {
				v23 += 2; // 2 chars to 4 UTF-8 bytes
				i++;
			} // else broken surrogate replaced with '?'
		}
		variableSize += v23;
		if (v23 < 128) {
			v23 = v23 << 1 | 1L;
		} else {
			java_unsafe.putLong(buf, w + java_unsafe.ARRAY_BYTE_BASE_OFFSET, v23);
			int bitCount = 64 - Long.numberOfLeadingZeros(v23);
			int tailSize = (((bitCount - 1) >>> 3) + bitCount) >>> 3;
			w += tailSize;
			v23 >>>= (tailSize << 3) - 1;
			v23 = (v23 | 1L) << tailSize;
		}
		word7 |= v23 << 32;
		long size = (long)(w - off) + variableSize;
		// boundary checks
		if (size > MARSHAL_MAX)
			return 0;
		if (buf.length - w < (int)variableSize)
			throw new java.nio.BufferOverflowException();

		// write payloads
		for (int i = 0, end = this.sn2[0].length(); i < end; i++) {
			char c = this.sn2[0].charAt(i);
			if (c < '\u0080') {
				buf[w++] = (byte) c;
			} else if (c < '\u0800') {
				buf[w++] = (byte) (c >> 6 | 0xc0);
				buf[w++] = (byte) (c & 0x3f | 0x80);
			} else if (! Character.isHighSurrogate(c)) {
				buf[w++] = (byte) (c >> 12 | 0xe0);
				buf[w++] = (byte) (c >> 6 & 0x3f | 0xc0);
				buf[w++] = (byte) (c & 0x3f | 0xc0);
			} else if (i + 1 >= end) {
				buf[w++] = '?'; // incomplete pair
			} else {
				char low = this.sn2[0].charAt(++i);
				if (!Character.isLowSurrogate(low)) {
					buf[w++] = '?'; // broken pair
					i--;
				} else {
					int cp = Character.toCodePoint(c, low);
					buf[w++] = (byte) (cp >> 18 | 0xf0);
					buf[w++] = (byte) (c >> 12 & 0x3f | 0xc0);
					buf[w++] = (byte) (c >> 6 & 0x3f | 0xc0);
					buf[w++] = (byte) (c & 0x3f | 0xc0);
				}
			}
		}
		for (int i = 0, end = this.sn2[1].length(); i < end; i++) {
			char c = this.sn2[1].charAt(i);
			if (c < '\u0080') {
				buf[w++] = (byte) c;
			} else if (c < '\u0800') {
				buf[w++] = (byte) (c >> 6 | 0xc0);
				buf[w++] = (byte) (c & 0x3f | 0x80);
			} else if (! Character.isHighSurrogate(c)) {
				buf[w++] = (byte) (c >> 12 | 0xe0);
				buf[w++] = (byte) (c >> 6 & 0x3f | 0xc0);
				buf[w++] = (byte) (c & 0x3f | 0xc0);
			} else if (i + 1 >= end) {
				buf[w++] = '?'; // incomplete pair
			} else {
				char low = this.sn2[1].charAt(++i);
				if (!Character.isLowSurrogate(low)) {
					buf[w++] = '?'; // broken pair
					i--;
				} else {
					int cp = Character.toCodePoint(c, low);
					buf[w++] = (byte) (cp >> 18 | 0xf0);
					buf[w++] = (byte) (c >> 12 & 0x3f | 0xc0);
					buf[w++] = (byte) (c >> 6 & 0x3f | 0xc0);
					buf[w++] = (byte) (c & 0x3f | 0xc0);
				}
			}
		}

		// write fixed positions
		word0 |= size;
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (0 * 8), word0);
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (1 * 8), word1);
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (2 * 8), word2);
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (3 * 8), word3);
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (4 * 8), word4);
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (5 * 8), word5);
		java_unsafe.putLong(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (6 * 8), word6);
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (7 * 8) + 0, (byte)(word7 >>> (0 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (7 * 8) + 1, (byte)(word7 >>> (1 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (7 * 8) + 2, (byte)(word7 >>> (2 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (7 * 8) + 3, (byte)(word7 >>> (3 * 8)));
		java_unsafe.putByte(buf, off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (7 * 8) + 4, (byte)(word7 >>> (4 * 8)));

		return (int)size;
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
		final long word2 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (2L * 8L));
		final long word3 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (3L * 8L));
		final long word4 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (4L * 8L));
		final long word5 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (5L * 8L));
		final long word6 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (6L * 8L));
		final long word7 = java_unsafe.getLong(buf, (long)off + java_unsafe.ARRAY_BYTE_BASE_OFFSET + (7L * 8L));

		final int size = (int)word0 & 0xfff;
		final int fixedSize = (int)(word0 >> 12) & 0xfff;
		if (size < fixedSize || fixedSize < 4) return 1;
		if (size > len) return 0;
		// payloads read backwards
		int end = off + size;

		// read index at variable section
		int r = off + fixedSize;
		// unpack .u8n2 uint8
		this.u8n2[0] = (byte)(word0 >> 24);
		this.u8n2[1] = (byte)(word0 >> 32);
		// unpack .i8n2 int8
		this.i8n2[0] = (byte)(word0 >> 40);
		this.i8n2[1] = (byte)(word0 >> 48);
		// unpack .u16n2 uint32
		long v4 = word0 >> (56 + 1) & 0x7f;
		if ((1L << 56 & word0) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v4 | 0x80) + 1;
			r += tailSize;
			v4 <<= (tailSize << 3) - tailSize;
			v4 |= tail & COLFER_MASKS[tailSize];
		}
		this.u16n2[0] = (int)v4;
		long v5 = word1 >> (0 + 1) & 0x7f;
		if ((1L << 0 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v5 | 0x80) + 1;
			r += tailSize;
			v5 <<= (tailSize << 3) - tailSize;
			v5 |= tail & COLFER_MASKS[tailSize];
		}
		this.u16n2[1] = (int)v5;
		// unpack .i16n2 int32
		long v6 = word1 >> (8 + 1) & 0x7f;
		if ((1L << 8 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v6 | 0x80) + 1;
			r += tailSize;
			v6 <<= (tailSize << 3) - tailSize;
			v6 |= tail & COLFER_MASKS[tailSize];
		}
		this.i16n2[0] = (int)(v6 >>> 1) ^ -(int)(v6 & 1L);
		long v7 = word1 >> (16 + 1) & 0x7f;
		if ((1L << 16 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v7 | 0x80) + 1;
			r += tailSize;
			v7 <<= (tailSize << 3) - tailSize;
			v7 |= tail & COLFER_MASKS[tailSize];
		}
		this.i16n2[1] = (int)(v7 >>> 1) ^ -(int)(v7 & 1L);
		// unpack .u32n2 uint16
		long v8 = word1 >> (24 + 1) & 0x7f;
		if ((1L << 24 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v8 | 0x80) + 1;
			r += tailSize;
			v8 <<= (tailSize << 3) - tailSize;
			v8 |= tail & COLFER_MASKS[tailSize];
		}
		this.u32n2[0] = (short)v8;
		long v9 = word1 >> (32 + 1) & 0x7f;
		if ((1L << 32 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v9 | 0x80) + 1;
			r += tailSize;
			v9 <<= (tailSize << 3) - tailSize;
			v9 |= tail & COLFER_MASKS[tailSize];
		}
		this.u32n2[1] = (short)v9;
		// unpack .i32n2 int16
		long v10 = word1 >> (40 + 1) & 0x7f;
		if ((1L << 40 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v10 | 0x80) + 1;
			r += tailSize;
			v10 <<= (tailSize << 3) - tailSize;
			v10 |= tail & COLFER_MASKS[tailSize];
		}
		this.i32n2[0] = (short)((short)(v10 >>> 1) ^ -(short)(v10 & 1L));
		long v11 = word1 >> (48 + 1) & 0x7f;
		if ((1L << 48 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v11 | 0x80) + 1;
			r += tailSize;
			v11 <<= (tailSize << 3) - tailSize;
			v11 |= tail & COLFER_MASKS[tailSize];
		}
		this.i32n2[1] = (short)((short)(v11 >>> 1) ^ -(short)(v11 & 1L));
		// unpack .u64n2 uint64
		long v12 = word1 >> (56 + 1) & 0x7f;
		if ((1L << 56 & word1) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v12 | 0x80) + 1;
			r += tailSize;
			v12 <<= (tailSize << 3) - tailSize;
			v12 |= tail & COLFER_MASKS[tailSize];
		}
		this.u64n2[0] = v12;
		long v13 = word2 >> (0 + 1) & 0x7f;
		if ((1L << 0 & word2) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v13 | 0x80) + 1;
			r += tailSize;
			v13 <<= (tailSize << 3) - tailSize;
			v13 |= tail & COLFER_MASKS[tailSize];
		}
		this.u64n2[1] = v13;
		// unpack .i64n2 int64
		long v14 = word2 >> (8 + 1) & 0x7f;
		if ((1L << 8 & word2) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v14 | 0x80) + 1;
			r += tailSize;
			v14 <<= (tailSize << 3) - tailSize;
			v14 |= tail & COLFER_MASKS[tailSize];
		}
		this.i64n2[0] = v14 >>> 1 ^ -(v14 & 1L);
		long v15 = word2 >> (16 + 1) & 0x7f;
		if ((1L << 16 & word2) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v15 | 0x80) + 1;
			r += tailSize;
			v15 <<= (tailSize << 3) - tailSize;
			v15 |= tail & COLFER_MASKS[tailSize];
		}
		this.i64n2[1] = v15 >>> 1 ^ -(v15 & 1L);
		// unpack .f32n2 float32
		int v16 = (int)(word2 >>> 24);
		this.f32n2[0] = Float.intBitsToFloat(v16);
		int v17 = (int)(word2>>>56 | word3<<(64-56));
		this.f32n2[1] = Float.intBitsToFloat(v17);
		// unpack .f64n2 float64
		long v18 = word3>>>24 | word4<<(64-24);
		this.f64n2[0] = Double.longBitsToDouble(v18);
		long v19 = word4>>>24 | word5<<(64-24);
		this.f64n2[1] = Double.longBitsToDouble(v19);
		// unpack .tn2 timestamp
		long v20 = word5>>>24 | word6<<(64-24);
		this.tn2[0] = java.time.Instant.ofEpochSecond(v20 >>> 30, (int) v20 & (1 << 30) - 1);
		long v21 = word6>>>24 | word7<<(64-24);
		this.tn2[1] = java.time.Instant.ofEpochSecond(v21 >>> 30, (int) v21 & (1 << 30) - 1);
		// unpack .sn2 text
		long v22 = word7 >> (24 + 1) & 0x7f;
		if ((1L << 24 & word7) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v22 | 0x80) + 1;
			r += tailSize;
			v22 <<= (tailSize << 3) - tailSize;
			v22 |= tail & COLFER_MASKS[tailSize];
		}
		if ((long)(end - r) < v22)
			return 1;
		end -= (int)v22;
		this.sn2[0] = new String(buf, end, (int)v22, java.nio.charset.StandardCharsets.UTF_8);
		long v23 = word7 >> (32 + 1) & 0x7f;
		if ((1L << 32 & word7) == 0) {
			long tail = java_unsafe.getLong(buf, r + java_unsafe.ARRAY_BYTE_BASE_OFFSET);
			int tailSize = Long.numberOfTrailingZeros(v23 | 0x80) + 1;
			r += tailSize;
			v23 <<= (tailSize << 3) - tailSize;
			v23 |= tail & COLFER_MASKS[tailSize];
		}
		if ((long)(end - r) < v23)
			return 1;
		end -= (int)v23;
		this.sn2[1] = new String(buf, end, (int)v23, java.nio.charset.StandardCharsets.UTF_8);

		// TODO: clear/undo absent fields

		return size;
	}

	/**
	 * {@link java.io.Serializable} version number reflects the fields present.
	 * Values in range [0, 127] belong to Colfer version 1.
	 */
	private static final long serialVersionUID = 61L << 7;

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
			throw new java.io.StreamCorruptedException("not a ArrayTypes Colfer encoding");
	}

	/**
	 * Gets gen.ArrayTypes.u8n2.
	 * @return the value.
	 */
	public byte[] getU8n2() {
		return this.u8n2;
	}

	/**
	 * Gets gen.ArrayTypes.i8n2.
	 * @return the value.
	 */
	public byte[] getI8n2() {
		return this.i8n2;
	}

	/**
	 * Gets gen.ArrayTypes.u16n2.
	 * @return the value.
	 */
	public int[] getU16n2() {
		return this.u16n2;
	}

	/**
	 * Gets gen.ArrayTypes.i16n2.
	 * @return the value.
	 */
	public int[] getI16n2() {
		return this.i16n2;
	}

	/**
	 * Gets gen.ArrayTypes.u32n2.
	 * @return the value.
	 */
	public short[] getU32n2() {
		return this.u32n2;
	}

	/**
	 * Gets gen.ArrayTypes.i32n2.
	 * @return the value.
	 */
	public short[] getI32n2() {
		return this.i32n2;
	}

	/**
	 * Gets gen.ArrayTypes.u64n2.
	 * @return the value.
	 */
	public long[] getU64n2() {
		return this.u64n2;
	}

	/**
	 * Gets gen.ArrayTypes.i64n2.
	 * @return the value.
	 */
	public long[] getI64n2() {
		return this.i64n2;
	}

	/**
	 * Gets gen.ArrayTypes.f32n2.
	 * @return the value.
	 */
	public float[] getF32n2() {
		return this.f32n2;
	}

	/**
	 * Gets gen.ArrayTypes.f64n2.
	 * @return the value.
	 */
	public double[] getF64n2() {
		return this.f64n2;
	}

	/**
	 * Gets gen.ArrayTypes.tn2.
	 * @return the value.
	 */
	public java.time.Instant[] getTn2() {
		return this.tn2;
	}

	/**
	 * Gets gen.ArrayTypes.sn2.
	 * @return the value.
	 */
	public String[] getSn2() {
		return this.sn2;
	}

	/**
	 * Deep hash is consistent with {@link #equals(Object)}.
	 * @return the standard Java digest.
	 */
	@Override
	public final int hashCode() {
		int h = 1;
		h = h * 31 + (int)this.u8n2[0];
		h = h * 31 + (int)this.u8n2[1];
		h = h * 31 + (int)this.i8n2[0];
		h = h * 31 + (int)this.i8n2[1];
		h = h * 31 + this.u16n2[0];
		h = h * 31 + this.u16n2[1];
		h = h * 31 + this.i16n2[0];
		h = h * 31 + this.i16n2[1];
		h = h * 31 + (int)this.u32n2[0];
		h = h * 31 + (int)this.u32n2[1];
		h = h * 31 + (int)this.i32n2[0];
		h = h * 31 + (int)this.i32n2[1];
		h = h * 31 + Long.hashCode(this.u64n2[0]);
		h = h * 31 + Long.hashCode(this.u64n2[1]);
		h = h * 31 + Long.hashCode(this.i64n2[0]);
		h = h * 31 + Long.hashCode(this.i64n2[1]);
		h = h * 31 + Float.hashCode(this.f32n2[0]);
		h = h * 31 + Float.hashCode(this.f32n2[1]);
		h = h * 31 + Double.hashCode(this.f64n2[0]);
		h = h * 31 + Double.hashCode(this.f64n2[1]);
		h = h * 31 + this.tn2[0].hashCode();
		h = h * 31 + this.tn2[1].hashCode();
		h = h * 31 + this.sn2[0].hashCode();
		h = h * 31 + this.sn2[1].hashCode();
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
		return o instanceof ArrayTypes && equals((ArrayTypes)o);
	}

	/**
	 * Typed alternative to {@link #equals(Object)}.
	 * @param o same class or {@code null}.
	 * @return the content match.
	 */
	public final boolean equals(ArrayTypes o) {
		if (o == null) return false;
		if (o == this) return true;

		return this.u8n2[0] == o.u8n2[0]
			&& this.u8n2[1] == o.u8n2[1]
			&& this.i8n2[0] == o.i8n2[0]
			&& this.i8n2[1] == o.i8n2[1]
			&& this.u16n2[0] == o.u16n2[0]
			&& this.u16n2[1] == o.u16n2[1]
			&& this.i16n2[0] == o.i16n2[0]
			&& this.i16n2[1] == o.i16n2[1]
			&& this.u32n2[0] == o.u32n2[0]
			&& this.u32n2[1] == o.u32n2[1]
			&& this.i32n2[0] == o.i32n2[0]
			&& this.i32n2[1] == o.i32n2[1]
			&& this.u64n2[0] == o.u64n2[0]
			&& this.u64n2[1] == o.u64n2[1]
			&& this.i64n2[0] == o.i64n2[0]
			&& this.i64n2[1] == o.i64n2[1]
			&& (this.f32n2[0] == o.f32n2[0] || (this.f32n2[0] != this.f32n2[0] && o.f32n2[0] != o.f32n2[0]))
			&& (this.f32n2[1] == o.f32n2[1] || (this.f32n2[1] != this.f32n2[1] && o.f32n2[1] != o.f32n2[1]))
			&& (this.f64n2[0] == o.f64n2[0] || (this.f64n2[0] != this.f64n2[0] && o.f64n2[0] != o.f64n2[0]))
			&& (this.f64n2[1] == o.f64n2[1] || (this.f64n2[1] != this.f64n2[1] && o.f64n2[1] != o.f64n2[1]))
			&& this.tn2[0].equals(o.tn2[0])
			&& this.tn2[1].equals(o.tn2[1])
			&& this.sn2[0].equals(o.sn2[0])
			&& this.sn2[1].equals(o.sn2[1]);
	}
}
