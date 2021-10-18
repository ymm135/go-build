# 自制编译器
[源码及pdf文档](res)  


cbc的的安装建议使用docker方式安装，安装指令为`docker pull leungwensen/cbc-ubuntu-64bit`, 运行`docker run -tid -v code:/root/code leungwensen/cbc-ubuntu-64bit`.可以在vscode中远程连接镜像.     
> 镜像里的cbc命令为cbc -Wa,--32 -Wl,-melf_i386的别名，可以直接执行。  

cbc安装及运行时如果出现问题, 可以通过`sh -x cbc xxxx`查看运行参数定位问题    
```
$ cbc hello.cb 
readlink: illegal option -- f
usage: readlink [-n] [file ...]
错误: 找不到或无法加载主类 net.loveruby.cflat.compiler.Compiler

$ sh -x cbc hello.cb 
+ JAVA=java
++ readlink -f cbc
readlink: illegal option -- f
usage: readlink [-n] [file ...]
+ cmd_path=
+++ dirname ''
++ dirname .
+ srcdir_root=.
+ java -classpath ./lib/cbc.jar net.loveruby.cflat.compiler.Compiler -I./import -L./lib hello.cb
错误: 找不到或无法加载主类 net.loveruby.cflat.compiler.Compiler

//错误显示. 不能再test文件夹下运行, 需要在根目录运行cbc命令, 包含lib、import文件夹  
```

- [第1章　开始制作编译器]()　　
- [第2章　CЬ和cbc]()　　
- [第3章　语法分析的概要]()　　
- [第4章　词法分析]()　　
- [第5章　基于JavaCC的解析器的描述]()　　
- [第6章　语法分析]()　　
- [第7章　JavaCC的action和抽象语法树]()　　
- [第8章　抽象语法树的生成]()　　
- [第9章　语义分析(1)引用的消解]()　　
- [第10章　语义分析(2)静态类型检查]()　　
- [第11章　中间代码的转换]()　　
- [第12章　x86 架构的概要]()　　
- [第13章　x86 汇编器编程]()　　
- [第14章　函数和变量]()　　
- [第15章　编译表达式和语句]()　　
- [第16章　分配栈帧]()　　
- [第17章　优化的方法]()　　
- [第18章　生成目标文件]()　　
- [第19章　链接和库]()　　
- [第20章　加载程序]()　　
- [第21章　生成地址无关代码]()　　
- [第22章　扩展阅读]()

 