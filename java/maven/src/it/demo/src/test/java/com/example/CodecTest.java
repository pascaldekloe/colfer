package com.example;

import com.example.demo.*;

import org.junit.Test;
import static org.junit.Assert.assertEquals;


/**
 * Tests the generated sources.
 * @author Pascal S. de Kloe
 */
public class CodecTest {

	/** Verifies the settings from the POM. */
	@Test
	public void pluginConfiguration() {
		assertEquals("size maximum", 2048, Course.colferSizeMax);
		assertEquals("list maximum", 99, Course.colferListMax);
	}

	/** Runs a full serialiazation cycle. */
	@Test
	public void codec() {
		Course a = new Course();
		a.name = "Koninklijke Haagsche Golf & Country Club";

		byte[] buf = new byte[100];
		int wrote = a.marshal(buf, 0);

		Course b = new Course();
		int read = b.unmarshal(buf, 0);

		assertEquals("write and read byte count", wrote, read);
		assertEquals("original and copy", a, b);
	}

}
