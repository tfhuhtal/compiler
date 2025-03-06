package assembler

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func Assemble(assemblyCode, outputFile string) ([]byte, error) {
	tempDir, err := ioutil.TempDir("", "compiler_")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	stdlibS := filepath.Join(tempDir, "stdlib.s")
	stdlibO := filepath.Join(tempDir, "stdlib.o")
	programS := filepath.Join(tempDir, "program.s")
	programO := filepath.Join(tempDir, "program.o")
	outputExe := filepath.Join(tempDir, "a.out")

	if err := ioutil.WriteFile(stdlibS, []byte(STDLIB_ASM_CODE), 0644); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(programS, []byte(assemblyCode), 0644); err != nil {
		return nil, err
	}

	if err := exec.Command("as", "-g", "-o", stdlibO, stdlibS).Run(); err != nil {
		return nil, err
	}
	if err := exec.Command("as", "-g", "-o", programO, programS).Run(); err != nil {
		return nil, err
	}
	if err := exec.Command("ld", "-o", outputExe, "-static", stdlibO, programO).Run(); err != nil {
		return nil, err
	}

	if outputFile != "" {
		return nil, os.Rename(outputExe, outputFile)
	}
	return ioutil.ReadFile(outputExe)
}

const STDLIB_ASM_CODE = `
	.global _start
	.global print_int
	.global print_bool
	.global read_int
	.extern main
	.section .text

# BEGIN START (we skip this part when linking with C)
# ***** Function '_start' *****
# Calls function 'main' and halts the program

_start:
	call main
	movq $60, %rax
	xorq %rdi, %rdi
	syscall
# END START

# ***** Function 'print_int' *****
print_int:
	pushq %rbp
	movq %rsp, %rbp
	movq %rdi, %r10
	decq %rsp
	movb $10, (%rsp)
	decq %rsp
	xorq %r9, %r9
	xorq %rax, %rax
	cmpq $0, %rdi
	je .Ljust_zero
	jge .Ldigit_loop
	incq %r9
.Ldigit_loop:
	cmpq $0, %rdi
	je .Ldigits_done
	movq %rdi, %rax
	movq $10, %rcx
	cqto
	idivq %rcx
	movq %rax, %rdi
	cmpq $0, %rdx
	jge .Lnot_negative
	negq %rdx
.Lnot_negative:
	addq $48, %rdx
	movb %dl, (%rsp)
	decq %rsp
	jmp .Ldigit_loop
.Ljust_zero:
	movb $48, (%rsp)
	decq %rsp
.Ldigits_done:
	cmpq $0, %r9
	je .Lminus_done
	movb $45, (%rsp)
	decq %rsp
.Lminus_done:
	movq $1, %rax
	movq $1, %rdi
	movq %rsp, %rsi
	incq %rsi
	movq %rbp, %rdx
	subq %rsp, %rdx
	decq %rdx
	syscall
	movq %rbp, %rsp
	popq %rbp
	movq %r10, %rax
	ret

# ***** Function 'print_bool' *****
print_bool:
	pushq %rbp
	movq %rsp, %rbp
	movq %rdi, %r10
	cmpq $0, %rdi
	jne .Ltrue
	movq $false_str, %rsi
	movq $false_str_len, %rdx
	jmp .Lwrite
.Ltrue:
	movq $true_str, %rsi
	movq $true_str_len, %rdx
.Lwrite:
	movq $1, %rax
	movq $1, %rdi
	syscall
	movq %rbp, %rsp
	popq %rbp
	movq %r10, %rax
	ret

true_str:
	.ascii "true"
true_str_len = . - true_str
false_str:
	.ascii "false"
false_str_len = . - false_str

# ***** Function 'read_int' *****
read_int:
	pushq %rbp
	movq %rsp, %rbp
	pushq %r12
	pushq $0
	xorq %r9, %r9
	xorq %r10, %r10
	xorq %r12, %r12
.Lloop:
	xorq %rax, %rax
	xorq %rdi, %rdi
	movq %rsp, %rsi
	movq $1, %rdx
	syscall
	cmpq $0, %rax
	jg .Lno_error
	je .Lend_of_input
	jmp .Lerror
.Lend_of_input:
	cmpq $0, %r12
	je .Lerror
	jmp .Lend
.Lno_error:
	incq %r12
	movq (%rsp), %r8
	cmpq $10, %r8
	je .Lend
	cmpq $45, %r8
	jne .Lnegation_done
	xorq $1, %r9
.Lnegation_done:
	cmpq $48, %r8
	jl .Lloop
	cmpq $57, %r8
	jg .Lloop
	subq $48, %r8
	imulq $10, %r10
	addq %r8, %r10
	jmp .Lloop
.Lend:
	cmpq $0, %r9
	je .Lfinal_negation_done
	neg %r10
.Lfinal_negation_done:
	popq %r12
	movq %rbp, %rsp
	popq %rbp
	movq %r10, %rax
	ret
.Lerror:
	movq $1, %rax
	movq $2, %rdi
	movq $read_int_error_str, %rsi
	movq $read_int_error_str_len, %rdx
	syscall
	movq $60, %rax
	movq $1, %rdi
	syscall

read_int_error_str:
	.ascii "Error: read_int() failed to read input\\n"
read_int_error_str_len = . - read_int_error_str
`
