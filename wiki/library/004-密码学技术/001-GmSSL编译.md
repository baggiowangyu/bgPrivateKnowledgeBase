# GmSSL编译

## Windows环境编译

### 编译32位Release动态链接库

```
cd <gmssl-project-dir>
perl Configure VC-WIN32 shared no-asm --prefix="E:/GmSSL/win32-release-shared-export" --openssldir="E:/GmSSL/win32-release-shared/ssl"
nmake
nmake test
nmake install
nmake clean
```

### 编译32位Release静态库

```
cd <gmssl-project-dir>
perl Configure VC-WIN32 no-asm no-shared --prefix="E:/GmSSL/win32-release-static-export" --openssldir="E:/GmSSL/win32-release-static/ssl"
nmake
nmake test
nmake install
nmake clean
```

### 编译32位Debug动态链接库

```
cd <gmssl-project-dir>
perl Configure VC-WIN32 shared no-asm --debug --prefix="E:/GmSSL/win32-debug-shared-export" --openssldir="E:/GmSSL/win32-debug-shared/ssl"
nmake
nmake test
nmake install
nmake clean
```

### 编译32位Debug静态库

```
cd <gmssl-project-dir>
perl Configure VC-WIN32 no-asm no-shared --debug --prefix="E:/GmSSL/win32-debug-static-export" --openssldir="E:/GmSSL/win32-debug-static/ssl"
nmake
nmake test
nmake install
nmake clean
```

### 编译64位Release动态链接库

win64配置参数需要根据自身系统确定：perl Configure { VC-WIN64A | VC-WIN64I }。
- VC-WIN64A：amd64
- VC-WIN64I：IA64

```
cd <gmssl-project-dir>
perl Configure VC-WIN64A shared no-asm --prefix="E:/GmSSL/win64-release-shared-export" --openssldir="E:/GmSSL/win64-release-shared/ssl"  
nmake  
nmake test  
nmake install  
nmake clean  
```

### 编译64位Release静态库

win64配置参数需要根据自身系统确定：perl Configure { VC-WIN64A | VC-WIN64I }。
- VC-WIN64A：amd64
- VC-WIN64I：IA64

```
cd <gmssl-project-dir>
perl Configure VC-WIN64A no-asm no-shared --prefix="E:/GmSSL/win64-release-static-export" --openssldir="E:/GmSSL/win64-release-static/ssl"  
nmake  
nmake test  
nmake install  
nmake clean  
```

### 编译64位Debug动态链接库

win64配置参数需要根据自身系统确定：perl Configure { VC-WIN64A | VC-WIN64I }。
- VC-WIN64A：amd64
- VC-WIN64I：IA64

```
cd <gmssl-project-dir>
perl Configure VC-WIN64A shared no-asm --debug --prefix="E:/GmSSL/win64-debug-shared-export" --openssldir="E:/GmSSL/win64-debug-shared/ssl"  
nmake  
nmake test  
nmake install  
nmake clean  
```

### 编译64位Debug静态库

win64配置参数需要根据自身系统确定：perl Configure { VC-WIN64A | VC-WIN64I }。
- VC-WIN64A：amd64
- VC-WIN64I：IA64

```
cd <gmssl-project-dir>
perl Configure VC-WIN64A no-asm no-shared --debug --prefix="E:/GmSSL/win64-debug-static-export" --openssldir="E:/GmSSL/win64-debug-static/ssl"  
nmake  
nmake test  
nmake install  
nmake clean  
```
