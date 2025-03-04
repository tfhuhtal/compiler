.extern print_int
.extern print_bool
.extern read_int
.section .text


.global main
.type main, @function
main:
    pushq %rbp
    movq %rsp, %rbp
    subq $32, %rsp

# Label(L0)
.L0:

# Call(read_int, [], x0)
subq $8, %rsp
callq read_int
movq %rax, -8(%rbp)
add $8, %rsp

# Copy(x0, x1)
movq -8(%rbp), %rax
movq %rax, -16(%rbp)


# Call(print_int, [x1], x2)
subq $8, %rsp
movq -16(%rbp), %rdi
callq print_int
movq %rax, -24(%rbp)
add $8, %rsp

    movq $0, %rax
    movq %rbp, %rsp
    popq %rbp
    ret