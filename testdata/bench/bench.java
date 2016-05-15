package testdata.bench;

import org.junit.Test;

import java.util.Arrays;


public class bench {

	private static Colfer[] newTestData() {
		Colfer c1 = new Colfer();
		c1.key = 1234567890L;
		c1.host = "db003lz12";
		c1.port = 389;
		c1.size = 452;
		c1.hash = 0x488b5c2428488918L;
		c1.ratio = 0.99;
		c1.route = true;

		Colfer c2 = new Colfer();
		c2.key = 1234567891L;
		c2.host = "localhost";
		c2.port = 22;
		c2.size = 4096;
		c2.hash = 0x243048899c24c824L;
		c2.ratio = 0.20;
		c2.route = false;

		Colfer c3 = new Colfer();
		c3.key = 1234567892L;
		c3.host = "kdc.local";
		c3.port = 88;
		c3.size = 1984;
		c3.hash = 0x000048891c24485cL;
		c3.ratio = 0.06;
		c3.route = false;

		Colfer c4 = new Colfer();
		c4.key = 1234567893L;
		c4.host = "vhost8.dmz.example.com";
		c4.port = 27017;
		c4.size = 59741;
		c4.hash = 0x5c2408488b9c2489L;
		c4.ratio = 0.0;
		c4.route = true;

		return new Colfer[] {c1, c2, c3, c4};
	}

	// prevent compiler optimization
	public static byte[] holdSerial;
	public static Colfer holdData;

	@Test
	public void benchMarshal() {
		Colfer[] testData = newTestData();
		final int n = 20000000;

		long start = System.nanoTime();
		for (int i = 0; i < n; i++) {
			holdSerial = new byte[200];
			testData[i % testData.length].marshal(holdSerial, 0);
		}
		long end = System.nanoTime();

		System.err.printf("%dM marshals avg %dns\n", n / 1000000, (end - start) / n);
	}

	@Test
	public void benchMarshalReuse() {
		Colfer[] testData = newTestData();
		holdSerial = new byte[Colfer.colferSizeMax];
		final int n = 20000000;

		long start = System.nanoTime();
		for (int i = 0; i < n; i++) {
			testData[i % testData.length].marshal(holdSerial, 0);
		}
		long end = System.nanoTime();

		System.err.printf("%dM marshals with buffer reuse avg %dns\n", n / 1000000, (end - start) / n);
	}

	@Test
	public void benchUnmarshal() {
		Colfer[] testData = newTestData();
		byte[][] serials = new byte[testData.length][];
		for (int i = 0; i < serials.length; i++) {
			byte[] buf = new byte[200];
			int n = testData[i].marshal(buf, 0);
			serials[i] = Arrays.copyOf(buf, n);
		}
		final int n = 20000000;

		long start = System.nanoTime();
		for (int i = 0; i < n; i++) {
			holdData = new Colfer();
			holdData.unmarshal(serials[i % serials.length], 0);
		}
		long end = System.nanoTime();

		System.err.printf("%dM unmarshals avg %dns\n", n / 1000000, (end - start) / n);
	}

}
