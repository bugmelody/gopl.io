// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 180.

// Package tempconv performs Celsius and Fahrenheit temperature computations.
package tempconv

import (
	"flag"
	"fmt"
)

type Celsius float64
type Fahrenheit float64

func CToF(c Celsius) Fahrenheit { return Fahrenheit(c*9.0/5.0 + 32.0) }
func FToC(f Fahrenheit) Celsius { return Celsius((f - 32.0) * 5.0 / 9.0) }

func (c Celsius) String() string { return fmt.Sprintf("%g°C", c) }

/*
//!+flagvalue
package flag

// Value is the interface to the value stored in a flag.
type Value interface {
	String() string
	Set(string) error
}
//!-flagvalue
*/

/**
Notice that celsiusFlag embeds a Celsius, there by getting a String method for free.
To satisfy flag.Value, we need only declare the Set method
*/
//!+celsiusFlag
// *celsiusFlag satisfies the flag.Value interface.
type celsiusFlag struct{ Celsius }

// 为了修改 receiver,需要使用指针作为 receiver
func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	/**
	The call to fmt.Sscanf parses a floating-point number (value) and a string (unit) from the
	input s. Although one must usually check Sscanf’s error result, in this case we don’t need
	to because if there was a problem, no switch case will match.
	*/
	fmt.Sscanf(s, "%f%s", &value, &unit) // no error check needed
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Fahrenheit(value))
		return nil
	}
	return fmt.Errorf("invalid temperature %q", s)
}

//!-celsiusFlag

//!+CelsiusFlag

/**
The CelsiusFlag function below wraps it all up. To the caller, it returns a pointer to the Celsius
field embedded within the celsiusFlag variable f. The Celsius field is the variable that will be
updated by the Set method during flags processing. The call to Var adds the flag to the
application’s set of command-line flags, the global variable flag.CommandLine.
Programs with unusually complex command-line interfaces may have several variables of this type. The
call to Var assigns a *celsiusFlag argument to a flag.Value parameter, causing the compiler to check
that *celsiusFlag has the necessary methods.
*/

// CelsiusFlag defines a Celsius flag with the specified name,
// default value, and usage, and returns the address of the flag variable.
// The flag argument must have a quantity and a unit, e.g., "100C".
func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	
	/**
	$ go doc flag.CommandLine
	var CommandLine = NewFlagSet(os.Args[0], ExitOnError)
			CommandLine is the default set of command-line flags, parsed from os.Args.
			The top-level functions such as BoolVar, Arg, and so on are wrappers for the
			methods of CommandLine.
	*/
	
	/**
	$ go doc flag.Flagset.Var
	func (f *FlagSet) Var(value Value, name string, usage string)
			Var defines a flag with the specified name and usage string. The type and
			value of the flag are represented by the first argument, of type Value,
			which typically holds a user-defined implementation of Value. For instance,
			the caller could create a flag that turns a comma-separated string into a
			slice of strings by giving the slice the methods of Value; in particular,
			Set would decompose the comma-separated string into the slice.
	 */
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}

//!-CelsiusFlag
