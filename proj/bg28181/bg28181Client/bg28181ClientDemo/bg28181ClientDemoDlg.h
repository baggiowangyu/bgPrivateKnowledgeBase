
// bg28181ClientDemoDlg.h : 头文件
//

#pragma once
#include "afxwin.h"

#include "bg28181Client.h"


// Cbg28181ClientDemoDlg 对话框
class Cbg28181ClientDemoDlg : public CDialogEx
{
// 构造
public:
	Cbg28181ClientDemoDlg(CWnd* pParent = NULL);	// 标准构造函数

// 对话框数据
	enum { IDD = IDD_BG28181CLIENTDEMO_DIALOG };

	protected:
	virtual void DoDataExchange(CDataExchange* pDX);	// DDX/DDV 支持


// 实现
protected:
	HICON m_hIcon;

	// 生成的消息映射函数
	virtual BOOL OnInitDialog();
	afx_msg void OnSysCommand(UINT nID, LPARAM lParam);
	afx_msg void OnPaint();
	afx_msg HCURSOR OnQueryDragIcon();
	DECLARE_MESSAGE_MAP()

public:
	CEdit m_cLocalIp;
	CEdit m_cLocalPort;
	CEdit m_cLocalGBCode;
	CEdit m_cServerIp;
	CEdit m_cServerPort;
	CEdit m_cServerGBCode;
	CEdit m_cUsername;
	CEdit m_cPassword;

	CEdit m_cKeepAliveSendRate;

public:
	bg28181Client _28181_client_;

public:
	afx_msg void OnBnClickedBtnRegist();
	afx_msg void OnBnClickedBtnUnregist();
	
};
