
// bg28181ClientDemo.h : PROJECT_NAME Ӧ�ó������ͷ�ļ�
//

#pragma once

#ifndef __AFXWIN_H__
	#error "�ڰ������ļ�֮ǰ������stdafx.h�������� PCH �ļ�"
#endif

#include "resource.h"		// ������


// Cbg28181ClientDemoApp: 
// �йش����ʵ�֣������ bg28181ClientDemo.cpp
//

class Cbg28181ClientDemoApp : public CWinApp
{
public:
	Cbg28181ClientDemoApp();

// ��д
public:
	virtual BOOL InitInstance();

// ʵ��

	DECLARE_MESSAGE_MAP()
};

extern Cbg28181ClientDemoApp theApp;