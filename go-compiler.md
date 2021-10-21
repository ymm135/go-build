
# [Go 编译器简介](https://github.com/golang/go/tree/master/src/cmd/compile)
cmd/compile包含构成Go编译器的主要包。编译器在逻辑上可以分为四个阶段，我们将在包含其代码的包列表旁边简要描述这些阶段。  

在提到编译器时，您有时可能会听到“前端”和“后端”这两个术语。粗略地说，这些转化为我们将在此处列出的前两个和最后两个阶段。第三个术语“中端”通常是指第二阶段发生的大部分工作。    

请注意，`go/*`包系列，例如`go/parserand`和`go/types`，与编译器无关。由于编译器最初是用C编写的，因此开发了`go/*`这些包以支持编写工具使用Go代码，例如`gofmt`和`vet`。    

需要说明的是，“gc”代表“Go compiler”，与大写的“GC”关系不大，GC代表垃圾收集。  

## 1.解析(Parsing)  
- `cmd/compile/internal/syntax` （词法分析器，解析器，语法树）  
<br>
在编译的第一阶段，对源代码进行标记（词法分析）、解析（语法分析），并为每个源文件构建语法树。   
  
每个语法树都是相应源文件的精确表示，其节点对应于源文件的各种元素，例如表达式、声明和语句。语法树还包括用于错误报告和创建调试信息的位置信息。  

## 2. 类型检查和AST转换(Type-checking and AST transformations)
- `cmd/compile/internal/gc` （创建编译器AST、类型检查、AST转换）  
<br>
从它用C编写时开始，gc包包含一个AST定义。它的所有代码都是根据它编写的，所以gc包必须做的第一件事是将语法包的语法树转换为编译器的AST表示. 这个额外的步骤将来可能会被重构掉。  

然后对AST进行类型检查。第一步是名称解析和类型推断，它们确定哪个对象属于哪个标识符，以及每个表达式具有什么类型。类型检查包括某些额外的检查，例如“声明和未使用”以及确定函数是否终止。  
  
某些转换也在AST上完成。一些节点根据类型信息进行细化，例如字符串加法从算术加法节点类型中分离出来。其他示例如死代码消除、函数调用内联和转义分析。  


## 3. 通用 SSA(Generic SSA)
- `cmd/compile/internal/gc` （转换为 SSA）
- `cmd/compile/internal/ssa` （SSA 通行证和规则）  
<br>
在此阶段，AST被转换为静态单分配 (SSA) 形式，这是一种具有特定属性的较低级别的中间表示，可以更轻松地实现优化并最终从中生成机器代码。 
 
在此转换期间，将应用内在函数。这些是编译器被教导要根据具体情况用高度优化的代码替换的特殊函数。  

在AST到SSA转换期间，某些节点也被降低为更简单的组件，以便编译器的其余部分可以使用它们。例如，内建的copy被内存移动取代，范围循环被重写为for循环。由于历史原因，其中一些目前在转换为SSA之前发生，但长期计划是将它们全部移到这里。  
  
然后，应用一系列与机器无关的通行证和规则。这些不涉及任何单一的计算机体系结构，因此可以在所有GOARCH变体上运行。  
  
这些通用传递的一些示例包括死代码消除、删除不需要的nil检查和删除未使用的分支。通用重写规则主要涉及表达式，例如将某些表达式替换为常量值，以及优化乘法和浮点运算。  
  

## 4. 生成机器码(Generating machine code)  
- `cmd/compile/internal/ssa` (SSA 降低和特定于拱门的通行证)(SSA lowering and arch-specific passes)  
- `cmd/internal/obj` （机器码生成）(machine code generation)  
<br>
编译器的机器相关阶段从“较低”阶段开始，它将通用值重写为特定于机器的变体。例如，在 md6 内存操作数上是可能的，因此可以组合许多加载-存储操作。  
  
请注意，lower pass 运行所有特定于机器的重写规则，因此它目前也应用了许多优化。  
  
一旦 SSA 被“降低”并且更加特定于目标架构，就会运行最终的代码优化通道。这包括另一个死代码消除过程，将值移动到更接近它们的用途，删除从未读取的局部变量以及寄存器分配。  
  
作为此步骤的一部分完成的其他重要工作包括堆栈帧布局，它将堆栈偏移分配给局部变量，以及指针活跃度分析，它计算每个GC安全点上哪些堆栈上指针是活跃的。  
  
在 SSA 生成阶段结束时，Go 函数已转换为一系列 obj.Prog 指令。这些被传递给汇编器 ( `cmd/internal/obj`)，汇编器将它们转换为机器代码并写出最终的目标文件。目标文件还将包含反射数据、导出数据和调试信息。  
  

## 进一步阅读  
要深入了解 SSA 包的工作原理，包括其传递和规则，请前往[cmd/compile/internal/ssa/README.md](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/README.md) 

# Go 编译器的SSA后端介绍  

这个包包含编译器的静态单赋值表单组件。如果您不熟悉SSA，其 [Wikipedia文章](https://zh.wikipedia.org/wiki/%E9%9D%99%E6%80%81%E5%8D%95%E8%B5%8B%E5%80%BC%E5%BD%A2%E5%BC%8F) 是一个很好的起点。  
如果您还不熟悉 Go 编译器，建议您先阅读 [cmd/compile/README.md](https://github.com/golang/go/blob/master/src/cmd/compile/README.md) 。该文档概述了编译器，并解释了SSA在其中的作用和目的。  

## Key concepts  
下面描述的名称可能与 Go 对应的名称松散相关，但请注意，它们并不等效。例如，Go 块语句有一个变量作用域，而 SSA 没有变量和变量作用域的概念。
  
值和块以其唯一的顺序ID命名也可能令人惊讶。它们很少对应原始代码中的命名实体，例如变量或函数参数。顺序 ID 还允许编译器避免映射，并且始终可以使用调试和位置信息将值追溯到 Go 代码。  

## Values  
值是SSA的基本组成部分。根据SSA的定义，一个值只定义一次，但可以使用任意次。值主要由唯一标识符、运算符、类型和一些参数组成。  

运算符或Op描述计算值的操作。每个运算符的语义可以在 中找到`gen/*Ops.go`。例如，OpAdd8 采用两个包含 8 位整数的值参数并将其结果相加。这是两个uint8值相加的可能 SSA 表示：  

```
// var c uint8 = a + b
v4 = Add8 <uint8> v2 v3
```  

值的类型通常是Go类型。例如，上面例子中的值有一个uint8类型，一个常量布尔值将有一个bool类型。然而，某些类型不是来自Go并且是特殊的；下面我们将介绍memory其中最常见的。    
有关更多信息，请参阅 [value.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/value.go)   

## 内存类型(Memory types)
`memory`表示全局内存状态。`Op`需要一个存储器参数取决于存储器状态，和一个Op具有存储型影响存储器的状态。这确保内存操作保持正确的顺序。例如：  

```
// *a = 3
// *b = *a
v10 = Store <mem> {int} v6 v8 v1
v14 = Store <mem> {int} v7 v8 v10
```

在这里，Store将它的第二个参数（类型int）存储到第一个参数（类型*int）中。最后一个参数是内存状态；由于第二个存储取决于第一个存储定义的内存值，因此两个存储不能重新排序。  

有关更多信息，请参阅 [cmd/compile/internal/types/type.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/types/type.go) 。  

  
## 块(Blocks)  
块代表函数控制流图中的基本块。它本质上是一个值列表，用于定义此块的操作。除了值列表之外，块主要由唯一标识符、种类和后继块列表组成。  

最简单的一种是`plain`块；它只是将控制流交给另一个块，因此它的后继列表包含一个块。  

另一种常见的块类型是`exit`块。它们有一个最终值，称为控制值，它必须返回一个内存状态。这对于函数返回一些值是必要的，例如 - 调用者需要一些内存状态来依赖，以确保它正确接收这些返回值。  

我们将提到的最后一个重要的块类型是`if`块。它有一个必须是布尔值的单个控制值，并且它正好有两个后继块。如果 `bool` 为真，则将控制流交给第一个后继，否则交给第二个。  

下面是一个用基本块表示的 if-else 控制流示例：  

```
// func(b bool) int {
// 	if b {
// 		return 2
// 	}
// 	return 3
// }
b1:
  v1 = InitMem <mem>
  v2 = SP <uintptr>
  v5 = Addr <*int> {~r1} v2
  v6 = Arg <bool> {b}
  v8 = Const64 <int> [2]
  v12 = Const64 <int> [3]
  If v6 -> b2 b3
b2: <- b1
  v10 = VarDef <mem> {~r1} v1
  v11 = Store <mem> {int} v5 v8 v10
  Ret v11
b3: <- b1
  v14 = VarDef <mem> {~r1} v1
  v15 = Store <mem> {int} v5 v12 v14
  Ret v15
```

有关更多信息，请参阅 [block.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/block.go) 。  

## 函数(Functions)  
函数表示函数声明及其主体。它主要由名称、类型（其签名）、构成其主体的块列表以及所述列表中的入口块组成。  

当一个函数被调用时，控制流被交给它的入口块。如果函数终止，控制流最终将到达退出块，从而结束函数调用。  

请注意，一个函数可能有零个或多个退出块，就像 Go 函数可以有任意数量的返回点一样，但它必须只有一个入口点块。  

另请注意，某些 SSA 函数是自动生成的，例如用作映射键的每种类型的哈希函数。  

例如，这是一个空函数在 SSA 中的样子，只有一个退出块，返回一个无趣的内存状态： 
```
foo func()
  b1:
    v1 = InitMem <mem>
    Ret v1
``` 

有关更多信息，请参阅[func.go]()。  

## 编译通过(Compiler passes)  
拥有 SSA 形式的程序本身并不是很有用。它的优势在于编写优化程序来修改程序以使其更好是多么容易。Go 编译器完成此操作的方式是通过传递列表。  

每次通过都以某种方式转换 SSA 函数。例如，死代码消除过程将删除它可以证明永远不会执行的块和值，而 `nil` 检查消除过程将删除它可以证明是多余的 `nil` 检查。  

编译器一次传递一个函数的工作，默认情况下按顺序运行一次。    

该`lower`通行证是特殊的; 它将 SSA 表示从机器无关转换为机器相关。也就是说，一些抽象运算符被替换为它们的非通用对应物，可能会减少或增加最终值的数量。  

有关更多信息，请参阅 [compile.gopasses](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/compile.go) 中定义的列表。  

## 玩转SSA(Playing with SSA)  
查看并习惯编译器的 SSA 运行的一个好方法是通过 GOSSAFUNC. 例如，要查看 funcFoo的初始 SSA 形式和最终生成的程序集，可以运行：  
```
GOSSAFUNC=Foo go build
```  

生成的ssa.html文件还将包含每个编译阶段的 SSA 函数，以便于查看每个阶段对特定程序的作用。您还可以单击值和块以突出显示它们，以帮助遵循控制流和值。  

GOSSAFUNC 中指定的值也可以是包限定的函数名，例如  
```
GOSSAFUNC=blah.Foo go build
```

这将匹配最终后缀为“blah”的包中任何名为“Foo”的函数（例如，something/blah.Foo、anotherthing/extra/blah.Foo）。  

如果需要非 HTML 转储，请在 GOSSAFUNC 值后附加一个“+”，转储将写入标准输出：  

```
GOSSAFUNC=Bar+ go build
```

## Hacking on SSA
虽然大多数编译器通过直接在 Go 代码中实现，但其他一些是代码生成的。这是目前通过重写规则完成的，这些规则有自己的语法并在`gen/*.rules`. 可以通过这种方式轻松快速地编写更简单的优化，但重写规则不适用于更复杂的优化。  

要阅读有关重写规则的更多信息，请查看 [gen/generic.rules](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/gen/generic.rules) 和 [gen/rulegen.go](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/gen/rulegen.go) 中的顶级评论 。  

同样，管理运算符的代码也是从 生成的代码 `gen/*Ops.go`，因为维护几个表比维护大量代码更容易。更改规则或运算符后，请参阅 [gen/README](https://github.com/golang/go/blob/master/src/cmd/compile/internal/ssa/gen/README) 以获取有关如何再次生成 Go 代码的说明。    






  

  