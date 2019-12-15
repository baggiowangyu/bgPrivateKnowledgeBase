
// bg28181ClientDemoDlg.h : ͷ�ļ�
//

#pragma once
#include "afxwin.h"

#include "bg28181Client.h"


// Cbg28181ClientDemoDlg �Ի���
class Cbg28181ClientDemoDlg : public CDialogEx
{
// ����
public:
	Cbg28181ClientDemoDlg(CWnd* pParent = NULL);	// ��׼���캯��

// �Ի�������
	enum { IDD = IDD_BG28181CLIENTDEMO_DIALOG };

	protected:
	virtual void DoDataExchange(CDataExchange* pDX);	// DDX/DDV ֧��


// ʵ��
protected:
	HICON m_hIcon;

	// ���ɵ���Ϣӳ�亯��
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
