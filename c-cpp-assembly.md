# c/c++常用语句对应的汇编指令 
## gdb 设置及指令 
`help all`查看所有帮助，`help disassemble`查看汇编指令的详情  

`set disassembly-flavor intel` 将汇编指令格式 设置为intel格式，默认是att
```shell
(gdb) show disassembly-flavor
The disassembly flavor is “att”.
```

`disassemble`指令查看intel汇编样式
```shell
-exec disassemble

Dump of assembler code for function main:
   0x000000000040052d <+0>:	push   rbp
   0x000000000040052e <+1>:	mov    rbp,rsp
   0x0000000000400531 <+4>:	sub    rsp,0x10
=> 0x0000000000400535 <+8>:	mov    DWORD PTR [rbp-0xc],0x5
   0x000000000040053c <+15>:	lea    rax,[rbp-0xc]
   0x0000000000400540 <+19>:	mov    QWORD PTR [rbp-0x8],rax
   0x0000000000400544 <+23>:	mov    edi,0x4005f0
   0x0000000000400549 <+28>:	call   0x400410 <puts@plt>
   0x000000000040054e <+33>:	mov    eax,0x0
   0x0000000000400553 <+38>:	leave  
   0x0000000000400554 <+39>:	ret    
End of assembler dump.
```

可以使用`/m`或者`/r`与源码一起显示或者16进制
```shell
-exec help disass 
Disassemble a specified section of memory.
Default is the function surrounding the pc of the selected frame.
With a /m modifier, source lines are included (if available).
With a /r modifier, raw instructions in hex are included.
With a single argument, the function surrounding that address is dumped.
Two arguments (separated by a comma) are taken as a range of memory to dump,
  in the form of "start,end", or "start,+length".
```

`x`查看内存,`32c`显示32位字符，`32x`显示32位16进制
```shell
-exec x/32c 0x4005f0
0x4005f0:	72 'H'	101 'e'	108 'l'	108 'l'	111 'o'	32 ' '	87 'W'	111 'o'
0x4005f8:	114 'r'	108 'l'	100 'd'	0 '\000'	1 '\001'	27 '\033'	3 '\003'	59 ';'
0x400600:	48 '0'	0 '\000'	0 '\000'	0 '\000'	5 '\005'	0 '\000'	0 '\000'	0 '\000'
0x400608:	4 '\004'	-2 '\376'	-1 '\377'	-1 '\377'	124 '|'	0 '\000'	0 '\000'	0 '\000'


-exec x/32x 0x4005f0
0x4005f0:	0x48	0x65	0x6c	0x6c	0x6f	0x20	0x57	0x6f
0x4005f8:	0x72	0x6c	0x64	0x00	0x01	0x1b	0x03	0x3b
0x400600:	0x30	0x00	0x00	0x00	0x05	0x00	0x00	0x00
0x400608:	0x04	0xfe	0xff	0xff	0x7c	0x00	0x00	0x00
```

## 变量及结构体使用
```c
#include <stdio.h>

typedef struct Man
{
    char *name;
    int age;
} Man;

int main()
{
    int a = 5;
    int *b = &a;

    Man man;
    man.name = "xiaoming";
    man.age = 18;

    int c = add(a, *b);

    printf("Hello World %d\n", c);
    return 0;
}

int add(int a, int b)
{
    return a + b;
}
```

在vscode上使用`-exec disass /m`查看汇编指令
```shell
Dump of assembler code for function main:
10	int main() {
   0x000000000040052d <+0>:	push   rbp
   0x000000000040052e <+1>:	mov    rbp,rsp                           #rbp移动到rsp位置
   0x0000000000400531 <+4>:	sub    rsp,0x20                          #rsp下移32个字节

11	    int a = 5;
=> 0x0000000000400535 <+8>:	mov    DWORD PTR [rbp-0xc],0x5           #把0x5存放到[rbp-0xc]位置，[rbp-0xc]代表a变量地址

12	    int* b = &a;
   0x000000000040053c <+15>:	lea    rax,[rbp-0xc]                 #首先把[rbp-0xc]的地址存储到寄存器rax
   0x0000000000400540 <+19>:	mov    QWORD PTR [rbp-0x8],rax       #把rax存储的值赋给变量b [rbp-0x8]

13	
14	    Man man;
15	    man.name = "xiaoming";
   0x0000000000400544 <+23>:	mov    QWORD PTR [rbp-0x20],0x400600 #把"xiaoming"字符串地址0x400600 赋给[rbp-0x20]
                                                                     #[rbp-0x20]既是man.name的地址，也是变量man的地址
16	    man.age = 18; 
   0x000000000040054c <+31>:	mov    DWORD PTR [rbp-0x18],0x12     #age占用8个字节，把0x12赋给
17	
18	    int c = add(a, *b);
   0x0000000000400553 <+38>:	mov    rax,QWORD PTR [rbp-0x8]       
   0x0000000000400557 <+42>:	mov    edx,DWORD PTR [rax]
   0x0000000000400559 <+44>:	mov    eax,DWORD PTR [rbp-0x10]
   0x000000000040055c <+47>:	mov    esi,edx
   0x000000000040055e <+49>:	mov    edi,eax
   0x0000000000400560 <+51>:	mov    eax,0x0
   0x0000000000400565 <+56>:	call   0x400588 <add>               #调用方法<add>
   0x000000000040056a <+61>:	mov    DWORD PTR [rbp-0xc],eax

19	
20	    printf("Hello World %d\n", c);
=> 0x000000000040056d <+64>:	mov    eax,DWORD PTR [rbp-0xc]      #把变量c的地址赋给eax(16 bit)
   0x0000000000400570 <+67>:	mov    esi,eax                      
   0x0000000000400572 <+69>:	mov    edi,0x400639                 # 把"Hello World"的内存地址0x400639赋给edi
   0x0000000000400577 <+74>:	mov    eax,0x0
   0x000000000040057c <+79>:	call   0x400410 <printf@plt>        #传递两个参数，调用<printf@plt>，没有变量时调用puts@plt>

21	    return 0;
   0x0000000000400581 <+84>:	mov    eax,0x0

22	}
   0x0000000000400586 <+89>:	leave  
   0x0000000000400587 <+90>:	ret 

End of assembler dump.
```

查看两个字符串的值
```shell
-exec x/32c 0x400600
0x400600:	120 'x'	105 'i'	97 'a'	111 'o'	109 'm'	105 'i'	110 'n'	103 'g'
0x400608:	0 '\000'	72 'H'	101 'e'	108 'l'	108 'l'	111 'o'	32 ' '	87 'W'
0x400610:	111 'o'	114 'r'	108 'l'	100 'd'	0 '\000'	0 '\000'	0 '\000'	0 '\000'
0x400618:	1 '\001'	27 '\033'	3 '\003'	59 ';'	52 '4'	0 '\000'	0 '\000'	0 '\000'

-exec x/32c 0x400609
0x400609:	72 'H'	101 'e'	108 'l'	108 'l'	111 'o'	32 ' '	87 'W'	111 'o'
0x400611:	114 'r'	108 'l'	100 'd'	0 '\000'	0 '\000'	0 '\000'	0 '\000'	1 '\001'
0x400619:	27 '\033'	3 '\003'	59 ';'	52 '4'	0 '\000'	0 '\000'	0 '\000'	5 '\005'
0x400621:	0 '\000'	0 '\000'	0 '\000'	-24 '\350'	-3 '\375'	-1 '\377'	-1 '\377'	-128 '\200'
```

查看函数`<printf@plt> 0x400410`和`<add> 0x400588` 
```shell
-exec disass /m 0x400410,0x400430
Dump of assembler code from 0x400410 to 0x400430:
=> 0x0000000000400410 <printf@plt+0>:	jmp    QWORD PTR [rip+0x200c02]        # 0x601018
   0x0000000000400416 <printf@plt+6>:	push   0x0
   0x000000000040041b <printf@plt+11>:	jmp    0x400400
   0x0000000000400420 <__libc_start_main@plt+0>:	jmp    QWORD PTR [rip+0x200bfa]        # 0x601020
   0x0000000000400426 <__libc_start_main@plt+6>:	push   0x1
   0x000000000040042b <__libc_start_main@plt+11>:	jmp    0x400400
End of assembler dump.

-exec disass /m 0x400588,0x4005a0
Dump of assembler code from 0x400588 to 0x4005a0:
25	{
   0x0000000000400588 <add+0>:	push   rbp
   0x0000000000400589 <add+1>:	mov    rbp,rsp
   0x000000000040058c <add+4>:	mov    DWORD PTR [rbp-0x4],edi
   0x000000000040058f <add+7>:	mov    DWORD PTR [rbp-0x8],esi

26	    return a + b;
   0x0000000000400592 <add+10>:	mov    eax,DWORD PTR [rbp-0x8]
   0x0000000000400595 <add+13>:	mov    edx,DWORD PTR [rbp-0x4]
   0x0000000000400598 <add+16>:	add    eax,edx

27	}   0x000000000040059a <add+18>:	pop    rbp
   0x000000000040059b <add+19>:	ret    

End of assembler dump.
```

## [c++ 常用指令及汇编解析(正在更新)](https://github.com/ymm135/golang-cookbook/blob/master/md/c-cpp-golang/base-c++.md)   

