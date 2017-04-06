package com.example;


/**
 * Core marshalling API of Colfer.
 * @author Pascal S. de Kloe
 */
public interface Colferable {

	/**
	 * Serializes the object.
	 * @param buf the data destination.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferOverflowException when {@code buf} is too small.
	 * @throws IllegalStateException on an upper limit breach defined by either {@link #colferSizeMax} or {@link #colferListMax}.
	 */
	int marshal(byte[] buf, int offset);

	/**
	 * Deserializes the object.
	 * @param buf the data source.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferUnderflowException when {@code buf} is incomplete. (EOF)
	 * @throws SecurityException on an upper limit breach defined by either {@code colferSizeMax} or {@code colferListMax}.
	 * @throws InputMismatchException when the data does not match this object's schema.
	 */
	int unmarshal(byte[] buf, int offset);

	/**
	 * Deserializes the object.
	 * @param buf the data source.
	 * @param offset the initial index for {@code buf}, inclusive.
	 * @param end the index limit for {@code buf}, exclusive.
	 * @return the final index for {@code buf}, exclusive.
	 * @throws BufferUnderflowException when {@code buf} is incomplete. (EOF)
	 * @throws SecurityException on an upper limit breach defined by either {@code colferSizeMax} or {@code colferListMax}.
	 * @throws InputMismatchException when the data does not match this object's schema.
	 */
	public int unmarshal(byte[] buf, int offset, int end);

}
