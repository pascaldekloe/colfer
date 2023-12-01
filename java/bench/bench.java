import java.util.Arrays;

import bench.Colfer;


public class bench {

	public static void main(String[] args) {
		benchMarshal();
		benchUnmarshal();
		System.out.printf("checksum %x\n", sum);
	}

	private static Colfer[] newTestData() {
		Colfer c1 = new Colfer();
		c1.key = 1234567890L;
		c1.host = "db003lz12";
		c1.port = 389;
		c1.size = 452;
		c1.hash = 0x488b5c2428488918L;
		c1.ratio = 0.99;
		c1.bools = Colfer.ROUTE_FLAG;

		Colfer c2 = new Colfer();
		c2.key = 1234567891L;
		c2.host = "localhost";
		c2.port = 22;
		c2.size = 4096;
		c2.hash = 0x243048899c24c824L;
		c2.ratio = 0.20;
		c1.bools = 0;

		Colfer c3 = new Colfer();
		c3.key = 1234567892L;
		c3.host = "kdc.local";
		c3.port = 88;
		c3.size = 1984;
		c3.hash = 0x000048891c24485cL;
		c3.ratio = 0.06;
		c1.bools = 0;

		Colfer c4 = new Colfer();
		c4.key = 1234567893L;
		c4.host = "vhost8.dmz.example.com";
		c4.port = 27017;
		c4.size = 59741;
		c4.hash = 0x5c2408488b9c2489L;
		c4.ratio = 0.0;
		c1.bools = Colfer.ROUTE_FLAG;

		return new Colfer[] {c1, c2, c3, c4};
	}

	public static int sum; // prevents compiler optimization

	static void benchMarshal() {
		Colfer[] testData = newTestData();
		byte[] buf = new byte[Colfer.MARSHAL_MAX];
		final int n = 20000000;

		long start = System.nanoTime();
		for (int i = 0; i < n; i++) {
			int size = testData[i % testData.length].marshalWithBounds(buf, 0);
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
			serials[i] = new byte[200];
			int size = testData[i].marshalWithBounds(serials[i], 0);
			if (size == 0) {
				System.out.println("exit on marshal error");
				return;
			}
		}
		Colfer o = new Colfer();
		final int n = 20000000;

		long start = System.nanoTime();
		for (int i = 0; i < n; i++) {
			int size = o.unmarshal(serials[i % serials.length], 0, 200);
			if (size < Colfer.MARSHAL_MIN) {
				System.out.println("exit on unmarshal error");
				return;
			}
			sum += size;
		}
		long end = System.nanoTime();

		System.out.printf("%d M unmarshals avg %d ns\n", n / 1000000, (end - start) / n);
		sum += o.hashCode();
	}

}
