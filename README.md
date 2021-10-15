# go编译器和链接器
主要通过分析编译过程,学习go语法分析,中间代码与汇编的转换,最终破除Go的语法糖,从编译器与汇编角度真正理解Go语言.  

# helloword编译过程分析  

代码在根目录: hello.go  
```
package main

import "fmt"

func main() {
	var s = "HelloWorld!"
	fmt.Println(s)
}
```

go build 指令常用参数
```
$ go help build  

The build flags are shared by the build, clean, get, install, list, run,
and test commands:

	-a
		force rebuilding of packages that are already up-to-date.
	-n
		print the commands but do not run them.
	-p n
		the number of programs, such as build commands or
		test binaries, that can be run in parallel.
		The default is GOMAXPROCS, normally the number of CPUs available.
	-race
		enable data race detection.
		Supported only on linux/amd64, freebsd/amd64, darwin/amd64, windows/amd64,
		linux/ppc64le and linux/arm64 (only for 48-bit VMA).
	-msan
		enable interoperation with memory sanitizer.
		Supported only on linux/amd64, linux/arm64
		and only with Clang/LLVM as the host C compiler.
		On linux/arm64, pie build mode will be used.
	-v
		print the names of packages as they are compiled.
	-work
		print the name of the temporary work directory and
		do not delete it when exiting.
	-x
		print the commands.

```

输出编译过程及文件  
```
$ go build -x -work  hello.go  

WORK=/var/folders/1g/cxqzw10d5vz2c8npwkkm1cfr0000gn/T/go-build1523506390
mkdir -p $WORK/b001/
cat >$WORK/b001/importcfg.link << 'EOF' # internal
packagefile command-line-arguments=/Users/zero/Library/Caches/go-build/13/13db8eaff8e9308791b824dd51ba589fa4efe242ae932849962521d5b97677c3-d
packagefile fmt=/Users/zero/go/sdk/go1.16.9/pkg/darwin_amd64/fmt.a
...
EOF
mkdir -p $WORK/b001/exe/
cd .
/Users/zero/go/sdk/go1.16.9/pkg/tool/darwin_amd64/link -o $WORK/b001/exe/a.out -importcfg $WORK/b001/importcfg.link -buildmode=exe -buildid=EpmhAIqOPGpJ3Qm4BqAq/S1sH1pm351_b1Cjcp1jh/iibyKbRQoYoSEzgSenAu/EpmhAIqOPGpJ3Qm4BqAq -extld=clang /Users/zero/Library/Caches/go-build/13/13db8eaff8e9308791b824dd51ba589fa4efe242ae932849962521d5b97677c3-d
/Users/zero/go/sdk/go1.16.9/pkg/tool/darwin_amd64/buildid -w $WORK/b001/exe/a.out # internal
mv $WORK/b001/exe/a.out hello
```

work目录的文件结构如下:  
```
└── b001
    ├── exe              //存储exe文件 exe/a.out 
    └── importcfg.link   //存储packagefile文件路径
```

主要的编译指令有两个: **link, buildid**  , 对应源码分别为:src/cmd/link/main.go(doc.go使用文档) ,src/cmd/buildid/buildid.go  
```
/Users/zero/go/sdk/go1.16.9/pkg/tool/darwin_amd64/link -o $WORK/b001/exe/a.out -importcfg $WORK/b001/importcfg.link -buildmode=exe -buildid=EpmhAIqOPGpJ3Qm4BqAq/S1sH1pm351_b1Cjcp1jh/iibyKbRQoYoSEzgSenAu/EpmhAIqOPGpJ3Qm4BqAq -extld=clang /Users/zero/Library/Caches/go-build/13/13db8eaff8e9308791b824dd51ba589fa4efe242ae932849962521d5b97677c3-d
/Users/zero/go/sdk/go1.16.9/pkg/tool/darwin_amd64/buildid -w $WORK/b001/exe/a.out # internal

link 常用参数
├── -o file
|        write output to file
├── -importcfg file
|         read import configuration from file
├── -buildmode mode
|         set build mode
├── -buildid id
|        record id as Go toolchain build id 
查看buildid文件中,可看到: build id "S1sH1pm351_b1Cjcp1jh/iibyKbRQoYoSEzgSenAu"字样  

buildid 用法
usage: go tool buildid [-w] file
  -w    write build ID
```


