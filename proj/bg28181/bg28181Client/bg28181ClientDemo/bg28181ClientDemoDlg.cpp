
// bg28181ClientDemoDlg.cpp : 实现文件
//

#include "stdafx.h"
#include "bg28181ClientDemo.h"
#include "bg28181ClientDemoDlg.h"
#include "afxdialogex.h"

#ifdef _DEBUG
#define new DEBUG_NEW
#endif


// 用于应用程序“关于”菜单项的 CAboutDlg 对话框

class CAboutDlg : public CDialogEx
{
public:
	CAboutDlg();

// 对话框数据
	enum { IDD = IDD_ABOUTBOX };

	protected:
	virtual void DoDataExchange(CDataExchange* pDX);    // DDX/DDV 支持

// 实现
protected:
	DECLARE_MESSAGE_MAP()
};

CAboutDlg::CAboutDlg() : CDialogEx(CAboutDlg::IDD)
{
}

void CAboutDlg::DoDataExchange(CDataExchange* pDX)
{
	CDialogEx::DoDataExchange(pDX);
}

BEGIN_MESSAGE_MAP(CAboutDlg, CDialogEx)
END_MESSAGE_MAP()


// Cbg28181ClientDemoDlg 对话框



Cbg28181ClientDemoDlg::Cbg28181ClientDemoDlg(CWnd* pParent /*=NULL*/)
	: CDialogEx(Cbg28181ClientDemoDlg::IDD, pParent)
{
	m_hIcon = AfxGetApp()->LoadIcon(IDR_MAINFRAME);
}

void Cbg28181ClientDemoDlg::DoDataExchange(CDataExchange* pDX)
{
	CDialogEx::DoDataExchange(pDX);
	DDX_Control(pDX, IDC_EDIT_LOCAL_IP, m_cLocalIp);
	DDX_Control(pDX, IDC_EDIT_LOCAL_PORT, m_cLocalPort);
	DDX_Control(pDX, IDC_EDIT_LOCAL_GBCODE, m_cLocalGBCode);
	DDX_Control(pDX, IDC_EDIT_SERVER_IP, m_cServerIp);
	DDX_Control(pDX, IDC_EDIT_SERVER_PORT, m_cServerPort);
	DDX_Control(pDX, IDC_EDIT_SERVER_GBCODE, m_cServerGBCode);
	DDX_Control(pDX, IDC_EDIT_AUTH_USERNAME, m_cUsername);
	DDX_Control(pDX, IDC_EDIT_AUTH_PASSWORD, m_cPassword);
}

BEGIN_MESSAGE_MAP(Cbg28181ClientDemoDlg, CDialogEx)
	ON_WM_SYSCOMMAND()
	ON_WM_PAINT()
	ON_WM_QUERYDRAGICON()
	ON_BN_CLICKED(IDC_BTN_REGIST, &Cbg28181ClientDemoDlg::OnBnClickedBtnRegist)
	ON_BN_CLICKED(IDC_BTN_UNREGIST, &Cbg28181ClientDemoDlg::OnBnClickedBtnUnregist)
END_MESSAGE_MAP()


// Cbg28181ClientDemoDlg 消息处理程序

BOOL Cbg28181ClientDemoDlg::OnInitDialog()
{
	CDialogEx::OnInitDialog();

	// 将“关于...”菜单项添加到系统菜单中。

	// IDM_ABOUTBOX 必须在系统命令范围内。
	ASSERT((IDM_ABOUTBOX & 0xFFF0) == IDM_ABOUTBOX);
	ASSERT(IDM_ABOUTBOX < 0xF000);

	CMenu* pSysMenu = GetSystemMenu(FALSE);
	if (pSysMenu != NULL)
	{
		BOOL bNameValid;
		CString strAboutMenu;
		bNameValid = strAboutMenu.LoadString(IDS_ABOUTBOX);
		ASSERT(bNameValid);
		if (!strAboutMenu.IsEmpty())
		{
			pSysMenu->AppendMenu(MF_SEPARATOR);
			pSysMenu->AppendMenu(MF_STRING, IDM_ABOUTBOX, strAboutMenu);
		}
	}

	// 设置此对话框的图标。  当应用程序主窗口不是对话框时，框架将自动
	//  执行此操作
	SetIcon(m_hIcon, TRUE);			// 设置大图标
	SetIcon(m_hIcon, FALSE);		// 设置小图标

	// TODO:  在此添加额外的初始化代码
	m_cLocalIp.SetWindowText(_T("192.168.231.1"));
	m_cLocalPort.SetWindowText(_T("5090"));
	m_cLocalGBCode.SetWindowText(_T("44000000002320000001"));

	m_cServerIp.SetWindowText(_T("192.168.231.131"));
	m_cServerPort.SetWindowText(_T("5060"));
	m_cServerGBCode.SetWindowText(_T("34020000002000000001"));
	m_cUsername.SetWindowText(_T("admin"));
	m_cPassword.SetWindowText(_T("12345678"));

	return TRUE;  // 除非将焦点设置到控件，否则返回 TRUE
}

void Cbg28181ClientDemoDlg::OnSysCommand(UINT nID, LPARAM lParam)
{
	if ((nID & 0xFFF0) == IDM_ABOUTBOX)
	{
		CAboutDlg dlgAbout;
		dlgAbout.DoModal();
	}
	else
	{
		CDialogEx::OnSysCommand(nID, lParam);
	}
}

// 如果向对话框添加最小化按钮，则需要下面的代码
//  来绘制该图标。  对于使用文档/视图模型的 MFC 应用程序，
//  这将由框架自动完成。

void Cbg28181ClientDemoDlg::OnPaint()
{
	if (IsIconic())
	{
		CPaintDC dc(this); // 用于绘制的设备上下文

		SendMessage(WM_ICONERASEBKGND, reinterpret_cast<WPARAM>(dc.GetSafeHdc()), 0);

		// 使图标在工作区矩形中居中
		int cxIcon = GetSystemMetrics(SM_CXICON);
		int cyIcon = GetSystemMetrics(SM_CYICON);
		CRect rect;
		GetClientRect(&rect);
		int x = (rect.Width() - cxIcon + 1) / 2;
		int y = (rect.Height() - cyIcon + 1) / 2;

		// 绘制图标
		dc.DrawIcon(x, y, m_hIcon);
	}
	else
	{
		CDialogEx::OnPaint();
	}
}

//当用户拖动最小化窗口时系统调用此函数取得光标
//显示。
HCURSOR Cbg28181ClientDemoDlg::OnQueryDragIcon()
{
	return static_cast<HCURSOR>(m_hIcon);
}



void Cbg28181ClientDemoDlg::OnBnClickedBtnRegist()
{
	// 先处理输入的数据
	CString str_local_ip;
	m_cLocalIp.GetWindowText(str_local_ip);

	CString str_local_port;
	m_cLocalPort.GetWindowText(str_local_port);
	int local_port = _ttoi(str_local_port.GetBuffer(0));

	CString str_local_gbcode;
	m_cLocalGBCode.GetWindowText(str_local_gbcode);

	CString str_server_ip;
	m_cServerIp.GetWindowText(str_server_ip);

	CString str_server_port;
	m_cServerPort.GetWindowText(str_server_port);
	int server_port = _ttoi(str_server_port.GetBuffer(0));

	CString str_server_gbcode;
	m_cServerGBCode.GetWindowText(str_server_gbcode);

	CString str_auth_username;
	m_cUsername.GetWindowText(str_auth_username);

	CString str_auth_password;
	m_cPassword.GetWindowText(str_auth_password);
	
	USES_CONVERSION;
	int errCode = _28181_client_.Initialize(T2A(str_local_ip.GetBuffer(0)), local_port, T2A(str_local_gbcode.GetBuffer(0)), NET_UDP);
	if (errCode != 0)
	{
		MessageBox(_T("初始化28181环境失败！"), _T("错误"), MB_OK | MB_ICONERROR);
		return ;
	}

	errCode = _28181_client_.Register(T2A(str_server_ip.GetBuffer(0)), server_port, T2A(str_server_gbcode.GetBuffer(0)), T2A(str_auth_username.GetBuffer(0)), T2A(str_auth_password.GetBuffer(0)), 3600);
	if (errCode != 0)
	{
		MessageBox(_T("注册失败！"), _T("错误"), MB_OK | MB_ICONERROR);
		return;
	}

	return;
}


void Cbg28181ClientDemoDlg::OnBnClickedBtnUnregist()
{
	// TODO:  在此添加控件通知处理程序代码
}
