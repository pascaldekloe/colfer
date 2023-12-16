import java.util.Arrays;

import bench.Colfer;


public class bench {

	public static void main(String[] args) {
		benchMarshal();
		benchUnmarshal();
		System.out.printf("checksum %x\n", sum);
	}

	private static Colfer[] newTestData() {
		return new Colfer[] {
			new Colfer()
				.withKey(1234567890L)
				.withHost("db003lz12")
				.withPort((short)389)
				.withSize(452)
				.withHash(0x488b5c2428488918L)
				.withRatio(0.99)
				.withRoute(true),

			new Colfer()
				.withKey(1234567891L)
				.withHost("localhost")
				.withPort((short)22)
				.withSize(4096)
				.withHash(0x243048899c24c824L)
				.withRatio(0.20)
				.withRoute(false),

			new Colfer()
				.withKey(1234567892L)
				.withHost("kdc.local")
				.withPort((short)88)
				.withSize(1984)
				.withHash(0x000048891c24485cL)
				.withRatio(0.06)
				.withRoute(false),

			new Colfer()
				.withKey(1234567893L)
				.withHost("vhost8.dmz.example.com")
				.withPort((short)27017)
				.withSize(59741)
				.withHash(0x5c2408488b9c2489L)
				.withRatio(0.0)
				.withRoute(true)

		};
	}

	public static int sum; // prevents compiler optimization

	static void benchMarshal() {
		Colfer[] testData = newTestData();
		byte[] buf = new byte[Colfer.MARSHAL_MAX];

		long start = System.nanoTime();
		final int n = 20000000 + (int)(start & 7);
		for (int i = 0; i < n; i++) {
			int size = testData[i % testData.length].marshal(buf, 0);
			if (size == 0) {
				System.out.println("exit on marshal error");
				return;
			}
			sum += size;
		}
		long end = System.nanoTime();

		System.out.printf("%d M marshals avg %d ns\n", n / 1000000, (end - start) / n);
		sum += Arrays.hashCode(buf);
	}

	static void benchUnmarshal() {
		Colfer[] testData = newTestData();
		byte[][] serials = new byte[testData.length][];
		for (int i = 0; i < serials.length; i++) {
			serials[i] = new byte[Colfer.COLFER_MAX];
			int size = testData[i].marshal(serials[i], 0);
			if (size == 0) {
				System.out.println("exit on marshal error");
				return;
			}
		}
		Colfer o = new Colfer();

		long start = System.nanoTime();
		final int n = 20000000 + (int)(start & 7);
		for (int i = 0; i < n; i++) {
			int size = o.unmarshal(serials[i % serials.length], 0, 200);
			if (size < Colfer.COLFER_MIN) {
				System.out.printf("exit on unmarshal error %d", size);
				return;
			}
			sum += size;
		}
		long end = System.nanoTime();

		System.out.printf("%d M unmarshals avg %d ns\n", n / 1000000, (end - start) / n);
		sum += o.hashCode();
	}

}
