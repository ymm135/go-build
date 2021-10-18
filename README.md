# go编译器和链接器
主要通过分析编译过程,学习go语法分析,中间代码与汇编的转换,最终破除Go的语法糖,从编译器与汇编角度真正理解Go语言.  

[Go编译器官方说明](https://github.com/golang/go/tree/master/src/cmd/compile)  
[Go编译器官方文档机翻](./go-compiler.md)  
[编译器知识:《自制编译器》](develop-compiler.md)  

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

## go 编译指令简介  

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

主要的编译指令有两个: **link, buildid**  , 对应源码分别为:
- src/cmd/link/main.go、src/cmd/link/doc.go
- src/cmd/buildid/buildid.go  

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

## go编译过程介绍  
1. [语法分析]逐行扫描源代码,将之转换为一系列的token，交给parser解析。
2. [语义分析]parser: 它将一系列token转换为AST(抽象语法树), 用于下一步生成代码。
3. [中间代码及生成]最后一步, 代码生成, 会利用上一步生成的AST并根据目标机器平台的不同，生成目标机器码。
> 注意：下面使用的代码包（go/scanner，go/parser，go/token，go/ast）主要是让我们可以方便地对 Go 代码进行解析和生成，做出更有趣的事情。但是 Go 本身的编译器并不是用这些代码包实现的。

[Go语法树分析](https://github.com/ymm135/go-ast-book)  
[Go程序编译成目标机器码](https://segmentfault.com/a/1190000016523685)  

### 语法分析  

```
package main

import (
	"fmt"
	"go/scanner"
	"go/token"
)

func main() {
	src := []byte(`
package main

import "fmt"

func main() {
	var s = "HelloWorld!"
	fmt.Println(s)
}
`)

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, src, nil, 0)

	for {
		pos, tok, lit := s.Scan()
		fmt.Printf("%-6s%-8s%q\n", fset.Position(pos), tok, lit)

		if tok == token.EOF {
			break
		}
	}
}

```

打印输出  
```
2:1   package "package"
2:9   IDENT   "main"
2:13  ;       "\n"
4:1   import  "import"
4:8   STRING  "\"fmt\""
4:13  ;       "\n"
6:1   func    "func"
6:6   IDENT   "main"
6:10  (       ""
6:11  )       ""
6:13  {       ""
7:2   var     "var"
7:6   IDENT   "s"
7:8   =       ""
7:10  STRING  "\"HelloWorld!\""
7:23  ;       "\n"
8:2   IDENT   "fmt"
8:5   .       ""
8:6   IDENT   "Println"
8:13  (       ""
8:14  IDENT   "s"
8:15  )       ""
8:16  ;       "\n"
9:1   }       ""
9:2   ;       "\n"
9:3   EOF     ""

```

### 语义分析 

```
package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
)

func main() {
	src := []byte(`
package main

import "fmt"

func main() {
	var s = "HelloWorld!"
	fmt.Println(s)
}
`)

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}

	ast.Print(fset, file)
}
```

打印输出  
```
     0  *ast.File {
     1  .  Package: 2:1
     2  .  Name: *ast.Ident {
     3  .  .  NamePos: 2:9
     4  .  .  Name: "main"
     5  .  }
     6  .  Decls: []ast.Decl (len = 2) {
     7  .  .  0: *ast.GenDecl {
     8  .  .  .  TokPos: 4:1
     9  .  .  .  Tok: import
    10  .  .  .  Lparen: -
    11  .  .  .  Specs: []ast.Spec (len = 1) {
    12  .  .  .  .  0: *ast.ImportSpec {
    13  .  .  .  .  .  Path: *ast.BasicLit {
    14  .  .  .  .  .  .  ValuePos: 4:8
    15  .  .  .  .  .  .  Kind: STRING
    16  .  .  .  .  .  .  Value: "\"fmt\""
    17  .  .  .  .  .  }
    18  .  .  .  .  .  EndPos: -
    19  .  .  .  .  }
    20  .  .  .  }
    21  .  .  .  Rparen: -
    22  .  .  }
    23  .  .  1: *ast.FuncDecl {
    24  .  .  .  Name: *ast.Ident {
    25  .  .  .  .  NamePos: 6:6
    26  .  .  .  .  Name: "main"
    27  .  .  .  .  Obj: *ast.Object {
    28  .  .  .  .  .  Kind: func
    29  .  .  .  .  .  Name: "main"
    30  .  .  .  .  .  Decl: *(obj @ 23)
    31  .  .  .  .  }
    32  .  .  .  }
    33  .  .  .  Type: *ast.FuncType {
    34  .  .  .  .  Func: 6:1
    35  .  .  .  .  Params: *ast.FieldList {
    36  .  .  .  .  .  Opening: 6:10
    37  .  .  .  .  .  Closing: 6:11
    38  .  .  .  .  }
    39  .  .  .  }
    40  .  .  .  Body: *ast.BlockStmt {
    41  .  .  .  .  Lbrace: 6:13
    42  .  .  .  .  List: []ast.Stmt (len = 2) {
    43  .  .  .  .  .  0: *ast.DeclStmt {
    44  .  .  .  .  .  .  Decl: *ast.GenDecl {
    45  .  .  .  .  .  .  .  TokPos: 7:2
    46  .  .  .  .  .  .  .  Tok: var
    47  .  .  .  .  .  .  .  Lparen: -
    48  .  .  .  .  .  .  .  Specs: []ast.Spec (len = 1) {
    49  .  .  .  .  .  .  .  .  0: *ast.ValueSpec {
    50  .  .  .  .  .  .  .  .  .  Names: []*ast.Ident (len = 1) {
    51  .  .  .  .  .  .  .  .  .  .  0: *ast.Ident {
    52  .  .  .  .  .  .  .  .  .  .  .  NamePos: 7:6
    53  .  .  .  .  .  .  .  .  .  .  .  Name: "s"
    54  .  .  .  .  .  .  .  .  .  .  .  Obj: *ast.Object {
    55  .  .  .  .  .  .  .  .  .  .  .  .  Kind: var
    56  .  .  .  .  .  .  .  .  .  .  .  .  Name: "s"
    57  .  .  .  .  .  .  .  .  .  .  .  .  Decl: *(obj @ 49)
    58  .  .  .  .  .  .  .  .  .  .  .  .  Data: 0
    59  .  .  .  .  .  .  .  .  .  .  .  }
    60  .  .  .  .  .  .  .  .  .  .  }
    61  .  .  .  .  .  .  .  .  .  }
    62  .  .  .  .  .  .  .  .  .  Values: []ast.Expr (len = 1) {
    63  .  .  .  .  .  .  .  .  .  .  0: *ast.BasicLit {
    64  .  .  .  .  .  .  .  .  .  .  .  ValuePos: 7:10
    65  .  .  .  .  .  .  .  .  .  .  .  Kind: STRING
    66  .  .  .  .  .  .  .  .  .  .  .  Value: "\"HelloWorld!\""
    67  .  .  .  .  .  .  .  .  .  .  }
    68  .  .  .  .  .  .  .  .  .  }
    69  .  .  .  .  .  .  .  .  }
    70  .  .  .  .  .  .  .  }
    71  .  .  .  .  .  .  .  Rparen: -
    72  .  .  .  .  .  .  }
    73  .  .  .  .  .  }
    74  .  .  .  .  .  1: *ast.ExprStmt {
    75  .  .  .  .  .  .  X: *ast.CallExpr {
    76  .  .  .  .  .  .  .  Fun: *ast.SelectorExpr {
    77  .  .  .  .  .  .  .  .  X: *ast.Ident {
    78  .  .  .  .  .  .  .  .  .  NamePos: 8:2
    79  .  .  .  .  .  .  .  .  .  Name: "fmt"
    80  .  .  .  .  .  .  .  .  }
    81  .  .  .  .  .  .  .  .  Sel: *ast.Ident {
    82  .  .  .  .  .  .  .  .  .  NamePos: 8:6
    83  .  .  .  .  .  .  .  .  .  Name: "Println"
    84  .  .  .  .  .  .  .  .  }
    85  .  .  .  .  .  .  .  }
    86  .  .  .  .  .  .  .  Lparen: 8:13
    87  .  .  .  .  .  .  .  Args: []ast.Expr (len = 1) {
    88  .  .  .  .  .  .  .  .  0: *ast.Ident {
    89  .  .  .  .  .  .  .  .  .  NamePos: 8:14
    90  .  .  .  .  .  .  .  .  .  Name: "s"
    91  .  .  .  .  .  .  .  .  .  Obj: *(obj @ 54)
    92  .  .  .  .  .  .  .  .  }
    93  .  .  .  .  .  .  .  }
    94  .  .  .  .  .  .  .  Ellipsis: -
    95  .  .  .  .  .  .  .  Rparen: 8:15
    96  .  .  .  .  .  .  }
    97  .  .  .  .  .  }
    98  .  .  .  .  }
    99  .  .  .  .  Rbrace: 9:1
   100  .  .  .  }
   101  .  .  }
   102  .  }
   103  .  Scope: *ast.Scope {
   104  .  .  Objects: map[string]*ast.Object (len = 1) {
   105  .  .  .  "main": *(obj @ 27)
   106  .  .  }
   107  .  }
   108  .  Imports: []*ast.ImportSpec (len = 1) {
   109  .  .  0: *(obj @ 12)
   110  .  }
   111  .  Unresolved: []*ast.Ident (len = 1) {
   112  .  .  0: *(obj @ 77)
   113  .  }
   114  }
```


### 中间代码及生成  

```
GOSSAFUNC=main GOOS=linux GOARCH=amd64 go build -gcflags -S hello.go 
```



  