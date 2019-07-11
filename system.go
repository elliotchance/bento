package main

import "fmt"

// System defines all of the inbuilt functions.
var System = map[string]func(vm *VirtualMachine, args []int){
	"display ?":  display,
	"set ? to ?": setVariable,
}

func display(vm *VirtualMachine, args []int) {
	_, _ = fmt.Fprintf(vm.out, "%v\n", *vm.GetArg(args[0]).(*string))
}

func setVariable(vm *VirtualMachine, args []int) {
	vm.SetArg(args[0], NewText(*vm.GetArg(args[1]).(*string)))
}
