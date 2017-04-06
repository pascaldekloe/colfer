package com.example;


/**
 * Go nuts with object-oriented programming.
 * @author Pascal S. de Kloe
 */
public abstract class BeanParent implements Colferable {

	@Override
	public String toString() {
		return "domain bean: " + getClass().getSimpleName();
	}

}
