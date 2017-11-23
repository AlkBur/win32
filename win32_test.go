package win32

import (
	"testing"
	"syscall"
	"unsafe"
	"log"
)

func TestWin32(t *testing.T)  {
	WinMain()
	return
}

func FatalErrNo(funcname string, err error) {
	errno, _ := err.(syscall.Errno)
	log.Fatalf("%s failed: %d %s\n", funcname, uint32(errno), err)
}

func WinMain() int {
	hInstance, err := GetModuleHandle(nil)
	if err != nil {
		FatalErrNo("GetModuleHandle", err)
	}
	// Get icon we're going to use.
	icon, err:= LoadIcon(0, IDI_APPLICATION)
	if err != nil {
		FatalErrNo("LoadIcon", err)
	}
	cursor, err := LoadCursor(0, IDC_ARROW)
	if err != nil {
		FatalErrNo("LoadCursor", err)
	}

	lpszClassName := syscall.StringToUTF16Ptr("myWindowClass")


	var wcex WNDCLASSEX
	wcex.Size            = uint32(unsafe.Sizeof(wcex))
	wcex.Style         = CS_HREDRAW | CS_VREDRAW
	wcex.WndProc       = syscall.NewCallback(WndProc)
	wcex.ClsExtra        = 0
	wcex.WndExtra        = 0
	wcex.Instance         = hInstance
	wcex.Icon         = icon
	wcex.Cursor       = cursor
	wcex.Background = COLOR_WINDOW + 11

	wcex.MenuName  = nil

	wcex.ClassName = lpszClassName
	wcex.IconSm       = icon

	if _, err = RegisterClassEx(&wcex); err != nil{
		FatalErrNo("RegisterClassEx", err)
	}

	hWnd, err := CreateWindowEx(
		0, lpszClassName, syscall.StringToUTF16Ptr("Simple Go Window!"),
		WS_OVERLAPPEDWINDOW | WS_VISIBLE,
		CW_USEDEFAULT, CW_USEDEFAULT, 400, 400, 0, 0, hInstance, 0)

	if err != nil {
		FatalErrNo("CreateWindowEx", err)
	}

	ShowWindow(hWnd, SW_SHOWDEFAULT)
	if err := UpdateWindow(hWnd); err != nil {
		FatalErrNo("CreateWindowEx", err)
	}

	var msg MSG
	for {
		r, err := GetMessage(&msg, 0, 0, 0)
		if  err != nil {
			break
		}
		if r == 0 {
			// WM_QUIT received -> get out
			break
		}
		TranslateMessage(&msg)
		DispatchMessage(&msg)
	}
	return int(msg.WParam)
}

func WndProc(hWnd HWND, msg uint32, wParam, lParam uintptr) (uintptr) {
	switch msg {
	case WM_DESTROY:
		PostQuitMessage(0)
	default:
		return DefWindowProc(hWnd, msg, wParam, lParam)
	}
	return 0
}
