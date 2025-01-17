package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
)

// System defines all of the inbuilt functions.
var System = map[string]func(vm *VirtualMachine, args []int){
	// This is a really dodgy hack until we can properly support varargs. Each
	// of the arguments will be printed with no space between them and a single
	// newline will be written after any (including zero) arguments.
	"display":                       display,
	"display ?":                     display,
	"display ? ?":                   display,
	"display ? ? ?":                 display,
	"display ? ? ? ?":               display,
	"display ? ? ? ? ?":             display,
	"display ? ? ? ? ? ?":           display,
	"display ? ? ? ? ? ? ?":         display,
	"display ? ? ? ? ? ? ? ? ?":     display,
	"display ? ? ? ? ? ? ? ? ? ?":   display,
	"display ? ? ? ? ? ? ? ? ? ? ?": display,

	// The other built-in functions.
	"set ? to ?":                                            setVariable,
	"add ? and ? into ?":                                    add,
	"subtract ? from ? into ?":                              subtract,
	"multiply ? and ? into ?":                               multiply,
	"divide ? by ? into ?":                                  divide,
	"run system command ?":                                  system,
	"run system command ? output into ?":                    systemOutput,
	"run system command ? status code into ?":               systemStatus,
	"run system command ? output into ? status code into ?": systemOutputStatus,
}

func display(vm *VirtualMachine, args []int) {
	for _, arg := range args {
		// TODO: Convert this switch into an interface.
		switch value := vm.GetArg(arg).(type) {
		case *string: // text
			_, _ = fmt.Fprintf(vm.out, "%v", *value)

		case *Number:
			_, _ = fmt.Fprintf(vm.out, "%v", value.String())

		case nil: // blackhole

		case *Backend:
			response, err := value.send(&BackendRequest{
				Sentence: "display ?",
				Args:     []string{fmt.Sprintf("%v", value)},
			})
			if err != nil {
				panic(err)
			}

			_, _ = fmt.Fprintf(vm.out, "%v", response.Text)

		default:
			panic(value)
		}
	}

	_, _ = fmt.Fprint(vm.out, "\n")
}

func setVariable(vm *VirtualMachine, args []int) {
	switch value := vm.GetArg(args[1]).(type) {
	case *string: // text
		vm.SetArg(args[0], NewText(*value))

	case *Number:
		vm.GetNumber(args[0]).Set(value)
	}
}

func add(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	c := vm.GetNumber(args[2])
	c.Add(a, b)
}

func subtract(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	c := vm.GetNumber(args[2])

	// Notice there are in reverse order because the language is
	// "subtract a from b".
	c.Sub(b, a)
}

func multiply(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	c := vm.GetNumber(args[2])
	c.Mul(a, b)
}

func divide(vm *VirtualMachine, args []int) {
	a := vm.GetNumber(args[0])
	b := vm.GetNumber(args[1])
	c := vm.GetNumber(args[2])
	c.Quo(a, b)
}

func runSystemCommand(rawCommand string) (output []byte, status int) {
	cmd := exec.Command("sh", "-c", rawCommand)
	var err error
	output, err = cmd.CombinedOutput()

	if msg, ok := err.(*exec.ExitError); ok {
		status = msg.Sys().(syscall.WaitStatus).ExitStatus()
	}

	return
}

func system(vm *VirtualMachine, args []int) {
	rawCommand := vm.GetText(args[0])
	output, _ := runSystemCommand(*rawCommand)
	_, _ = vm.out.Write(output)
}

func systemOutput(vm *VirtualMachine, args []int) {
	rawCommand := vm.GetText(args[0])
	output, _ := runSystemCommand(*rawCommand)
	vm.SetArg(args[1], NewText(string(output)))
}

func systemStatus(vm *VirtualMachine, args []int) {
	rawCommand := vm.GetText(args[0])
	_, status := runSystemCommand(*rawCommand)
	vm.SetArg(args[1], NewNumber(strconv.Itoa(status), 0))
}

func systemOutputStatus(vm *VirtualMachine, args []int) {
	rawCommand := vm.GetText(args[0])
	output, status := runSystemCommand(*rawCommand)
	vm.SetArg(args[1], NewText(string(output)))
	vm.SetArg(args[2], NewNumber(strconv.Itoa(status), 0))
}
