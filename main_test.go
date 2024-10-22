package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"testing"
)

func TestBentoFiles(t *testing.T) {
	dir := "tests/"
	fileInfos, err := ioutil.ReadDir(dir)
	require.NoError(t, err)

	for _, fileInfo := range fileInfos {
		// This is useful for debugging a single file, so i'll leave it
		// commented out.
		//if fileInfo.Name() != "number.bento" {
		//	continue
		//}

		if !strings.HasSuffix(fileInfo.Name(), ".bento") {
			continue
		}

		t.Run(fileInfo.Name(), func(t *testing.T) {
			file, err := os.Open(dir + fileInfo.Name())
			require.NoError(t, err)

			parser := NewParser(file)
			program, err := parser.Parse()
			require.NoError(t, err)

			compiler := NewCompiler(program)
			compiledProgram := compiler.Compile()

			vm := NewVirtualMachine(compiledProgram)
			vm.out = bytes.NewBuffer(nil)
			err = vm.Run()
			require.NoError(t, err)

			// There can be a specific expectation file depending on the OS.
			expectedFilePath := dir + strings.Replace(fileInfo.Name(), ".bento", "."+runtime.GOOS+".txt", -1)
			expectedData, err := ioutil.ReadFile(expectedFilePath)

			if err != nil {
				// Fallback to generic expectation file.
				expectedFilePath = dir + strings.Replace(fileInfo.Name(), ".bento", ".txt", -1)
				expectedData, err = ioutil.ReadFile(expectedFilePath)
				require.NoError(t, err)
			}

			assert.Equal(t, string(expectedData), vm.out.(*bytes.Buffer).String())
		})
	}
}
