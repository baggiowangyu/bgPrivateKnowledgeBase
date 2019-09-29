# 从汇编角度来看数值比较与字符串比较的性能差异

## 研究环境

- 开发环境：Visual Studio 2008
- 运行环境：Windows 10 x64

## 代码片段

```
#include "stdafx.h"
#include <iostream>


int _tmain(int argc, _TCHAR* argv[])
{
	int a = 123;
	bool bret = a == 123;

	char *b = "123";
	strcmp(b, "123");

	return 0;
}
```

## 汇编分析

在进行分析之前，我们假定，一条汇编语句消耗CPU时间单位为1。那么我们通过一个算法占用的CPU时间单位来评估算法的消耗成本，由此来简单评估几个算法之间的效能差异。

```
bool bret = a == 123;
```

可翻译为汇编语言：

```
xor         eax,eax 
cmp         dword ptr [a],7Bh 
sete        al   
mov         byte ptr [bret],al 
```

其中

```
mov         byte ptr [bret],al 
```

是将比较结果赋值给变量bret，可以忽略。

那么我们可以初步认定，采用数值比较的方式，消耗的CPU时间单位为3.

接下来，我们来对字符串比较进行效能评估：

```
strcmp(b, "123");
```

可翻译为汇编语言：

```
000913D8  push        offset string "123" (95800h) 
000913DD  mov         eax,dword ptr [b] 
000913E0  push        eax  
000913E1  call        @ILT+165(_strcmp) (910AAh) 
```

其中strcmp的汇编实现为：

```
strcmp  proc \
        str1:ptr byte, \
        str2:ptr byte

        OPTION PROLOGUE:NONE, EPILOGUE:NONE

        .FPO    ( 0, 2, 0, 0, 0, 0 )

        mov     edx,[esp + 4]   ; edx = src
        mov     ecx,[esp + 8]   ; ecx = dst

        test    edx,3
        jnz     short dopartial

        align   4
dodwords:
        mov     eax,[edx]

        cmp     al,[ecx]
        jne     short donene
        or      al,al
        jz      short doneeq
        cmp     ah,[ecx + 1]
        jne     short donene
        or      ah,ah
        jz      short doneeq

        shr     eax,16

        cmp     al,[ecx + 2]
        jne     short donene
        or      al,al
        jz      short doneeq
        cmp     ah,[ecx + 3]
        jne     short donene
        add     ecx,4
        add     edx,4
        or      ah,ah
        jnz     short dodwords

        align   4
doneeq:
        xor     eax,eax
        ret

        align   4
donene:
        ; The instructions below should place -1 in eax if src < dst,
        ; and 1 in eax if src > dst.

        sbb     eax,eax
        sal     eax,1
        add     eax,1
        ret

        align   4
dopartial:
        test    edx,1
        jz      short doword

        mov     al,[edx]
        add     edx,1
        cmp     al,[ecx]
        jne     short donene
        add     ecx,1
        or      al,al
        jz      short doneeq

        test    edx,2
        jz      short dodwords


        align   4
doword:
        mov     ax,[edx]
        add     edx,2
        cmp     al,[ecx]
        jne     short donene
        or      al,al
        jz      short doneeq
        cmp     ah,[ecx + 1]
        jne     short donene
        or      ah,ah
        jz      short doneeq
        add     ecx,2
        jmp     short dodwords

strcmp  endp

        end
```

在此环境实际测试执行了23条汇编语句，那么粗略评估CPU时间单位为26

简单评估，当字符串长度为3的时候，字符串比较消耗CPU资源是数值比较的9倍左右。
当字符串长度为10时，字符串比较消耗的CPU时间单位为63，是数值比较的21倍。

可以粗略判断，当字符串越长，字符串比较的CPU消耗时间会线性增加，数值比较的时间消耗固定不变。
