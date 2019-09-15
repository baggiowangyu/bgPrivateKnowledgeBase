# POCO C++ 编译

## Windows平台

- 准备openssl项目，0.9.x均可。完成编译，生成 **```libeay32.lib```**、**```libeay32.dll```**、**```ssleay32.lib```**、**```ssleay32.dll```**；
- 准备openssl工程项结构：

```
    ┗openssl
      ┣build
      ┃┗include
      ┃  ┗openssl
      ┗win32
        ┗bin
          ┣debug
          ┗release
```

- 打开VS2008，首先编译Crypto，这个项目涉及到后面的安全性模块：

```
    编译过程中会报以下错误：\poco\Crypto\include\Poco/Crypto/EVPPKey.h(309) : error C3861: “EVP_PKEY_id”: 找不到标识符
    原因是Openssl已经将EVP_PKEY_id()函数移除了，我们可以在项目的的EVPPKey.h中手动添加该函数定义：

    #ifndef EVP_PKEY_id
    #define EVP_PKEY_id(evp_key) evp_key->type
    #endif
```

- 编译通过后其他项目基本不成问题了
