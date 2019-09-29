// bgCompareEfficient.cpp : 定义控制台应用程序的入口点。
//
// 本例用于研究几种条件判断的效率问题

#include "stdafx.h"
#include <iostream>


int _tmain(int argc, _TCHAR* argv[])
{
	int a = 123;
	bool bret = a == 123;

	char *b = "1234567890";
	strcmp(b, "1234567890");

	return 0;
}

