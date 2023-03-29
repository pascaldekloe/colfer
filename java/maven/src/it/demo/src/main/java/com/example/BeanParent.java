package com.example;


/**
 * Go nuts with object-oriented programming.
 * @author Pascal S. de Kloe
 */
public abstract class BeanParent implements Colferable {

	@Override
	public boolean equals(Object o) {
		return o instanceof BeanParent;
	}

	@Override
	public int hashCode() {
		return 42;
	}

	@Override
	public String toString() {
		return "domain bean: " + getClass().getSimpleName();
	}

}
