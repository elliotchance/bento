package main

import (
	"fmt"
	"math/big"
)

// System defines all of the inbuilt functions.
var System = map[string]func(vm *VirtualMachine, args []int){
	"display ?":                display,
	"set ? to ?":               setVariable,
	"add ? and ? into ?":       add,
	"subtract ? from ? into ?": subtract,
	"multiply ? and ? into ?":  multiply,
	"divide ? by ? into ?":     divide,
}

func display(vm *VirtualMachine, args []int) {
	switch value := vm.GetArg(args[0]).(type) {
	case *string: // text
		_, _ = fmt.Fprintf(vm.out, "%v\n", *value)

	case *big.Rat: // number
		_, _ = fmt.Fprintf(vm.out, "%v\n", value.FloatString(6))

	default:
		panic(value)
	}
}

func setVariable(vm *VirtualMachine, args []int) {
	var newValue interface{}

	switch value := vm.GetArg(args[1]).(type) {
	case *string: // text
		newValue = NewText(*value)

	case *big.Rat: // number
		newValue = big.NewRat(0, 1).Set(value)

	default:
		panic(value)
	}

	vm.SetArg(args[0], newValue)
}

func add(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	vm.SetArg(args[2], big.NewRat(0, 1).Add(a, b))
}

func subtract(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	vm.SetArg(args[2], big.NewRat(0, 1).Sub(b, a))
}

func multiply(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	vm.SetArg(args[2], big.NewRat(0, 1).Mul(a, b))
}

func divide(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	vm.SetArg(args[2], big.NewRat(0, 1).Quo(a, b))
}
