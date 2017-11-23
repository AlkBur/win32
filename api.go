package win32

import (
	"syscall"
	"unsafe"
	"fmt"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	moduser32   = syscall.NewLazyDLL("user32.dll")
	modgdi32 = syscall.NewLazyDLL("gdi32.dll")

	procGetModuleHandleW = modkernel32.NewProc("GetModuleHandleW")
	procRegisterClassExW = moduser32.NewProc("RegisterClassExW")
	procCreateWindowExW  = moduser32.NewProc("CreateWindowExW")
	procDefWindowProcW   = moduser32.NewProc("DefWindowProcW")
	procDestroyWindow    = moduser32.NewProc("DestroyWindow")
	procPostQuitMessage  = moduser32.NewProc("PostQuitMessage")
	procShowWindow       = moduser32.NewProc("ShowWindow")
	procUpdateWindow     = moduser32.NewProc("UpdateWindow")
	procGetMessageW      = moduser32.NewProc("GetMessageW")
	procTranslateMessage = moduser32.NewProc("TranslateMessage")
	procDispatchMessageW = moduser32.NewProc("DispatchMessageW")
	procLoadIconW        = moduser32.NewProc("LoadIconW")
	procLoadCursorW      = moduser32.NewProc("LoadCursorW")
	procSetCursor        = moduser32.NewProc("SetCursor")
	procSendMessageW     = moduser32.NewProc("SendMessageW")
	procPostMessageW     = moduser32.NewProc("PostMessageW")

	procGetDC            = moduser32.NewProc("GetDC")
	procGetStockObject   = modgdi32.NewProc("GetStockObject")
	procGetDeviceCaps    = modgdi32.NewProc("GetDeviceCaps")
	procUnregisterClass  = moduser32.NewProc("UnregisterClassW")
	procPeekMessage      = moduser32.NewProc("PeekMessageW")
	procGetDesktopWindow = moduser32.NewProc("GetDesktopWindow")
	procGetClientRect                 = moduser32.NewProc("GetClientRect")
)

var (
	// Some globally known cursors
	IDC_ARROW uint16 = 32512
	IDC_IBEAM uint16 = 32513
	IDC_WAIT  uint16 = 32514
	IDC_CROSS uint16 = 32515

	// Some globally known icons
	IDI_APPLICATION uint16 = 32512
	IDI_HAND        uint16 = 32513
	IDI_QUESTION    uint16 = 32514
	IDI_EXCLAMATION uint16 = 32515
	IDI_ASTERISK    uint16 = 32516
	IDI_WINLOGO     uint16 = 32517
	IDI_WARNING     = IDI_EXCLAMATION
	IDI_ERROR       = IDI_HAND
	IDI_INFORMATION = IDI_ASTERISK
)

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func GetModuleHandle(modname *uint16) (handle HINSTANCE, err error) {
	r0, _, e1 := syscall.Syscall(procGetModuleHandleW.Addr(), 1, uintptr(unsafe.Pointer(modname)), 0, 0)
	handle = HINSTANCE(r0)
	if handle == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func LoadIcon(instance HINSTANCE, id uint16) (icon HICON, err error) {
	iconName := MakeIntResource(id)
	ret, _, e1 := syscall.Syscall(procLoadIconW.Addr(), 2,
		uintptr(instance),
		uintptr(unsafe.Pointer(iconName)),
		0)
	icon = HICON(ret)
	if icon == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func LoadCursor(instance HINSTANCE, id uint16) (cursor HCURSOR, err error) {
	cursorname := MakeIntResource(id)
	r0, _, e1 := syscall.Syscall(procLoadCursorW.Addr(), 2, uintptr(instance), uintptr(unsafe.Pointer(cursorname)), 0)
	cursor = HCURSOR(r0)
	if cursor == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func PostQuitMessage(exitcode int32) {
	syscall.Syscall(procPostQuitMessage.Addr(), 1, uintptr(exitcode), 0, 0)
	return
}

func DefWindowProc(hwnd HWND, msg uint32, wparam uintptr, lparam uintptr) (lresult uintptr) {
	r0, _, _ := syscall.Syscall6(procDefWindowProcW.Addr(), 4, uintptr(hwnd), uintptr(msg), uintptr(wparam), uintptr(lparam), 0, 0)
	lresult = uintptr(r0)
	return
}

func RegisterClassEx(wndclass *WNDCLASSEX) (err error) {
	r0, _, e1 := syscall.Syscall(procRegisterClassExW.Addr(), 1, uintptr(unsafe.Pointer(wndclass)), 0, 0)
	atom := uint16(r0)
	if atom == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func CreateWindowEx(exstyle uint32, classname *uint16, windowname *uint16, style uint32, x int32, y int32, width int32, height int32, wndparent HWND, menu HMENU, instance HINSTANCE, param uintptr) (hwnd HWND, err error) {
	r0, _, e1 := syscall.Syscall12(procCreateWindowExW.Addr(), 12, uintptr(exstyle), uintptr(unsafe.Pointer(classname)), uintptr(unsafe.Pointer(windowname)), uintptr(style), uintptr(x), uintptr(y), uintptr(width), uintptr(height), uintptr(wndparent), uintptr(menu), uintptr(instance), uintptr(param))
	hwnd = HWND(r0)
	if hwnd == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ShowWindow(hwnd HWND, cmdshow int32) (wasvisible bool) {
	r0, _, _ := syscall.Syscall(procShowWindow.Addr(), 2, uintptr(hwnd), uintptr(cmdshow), 0)
	wasvisible = bool(r0 != 0)
	return
}

func UpdateWindow(hwnd HWND) (err error) {
	r1, _, e1 := syscall.Syscall(procUpdateWindow.Addr(), 1, uintptr(hwnd), 0, 0)
	if int(r1) == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func GetMessage(msg *MSG, hwnd syscall.Handle, MsgFilterMin uint32, MsgFilterMax uint32) (ret int32, err error) {
	r0, _, e1 := syscall.Syscall6(procGetMessageW.Addr(), 4, uintptr(unsafe.Pointer(msg)), uintptr(hwnd), uintptr(MsgFilterMin), uintptr(MsgFilterMax), 0, 0)
	ret = int32(r0)
	if ret == -1 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func TranslateMessage(msg *MSG) (done bool) {
	r0, _, _ := syscall.Syscall(procTranslateMessage.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	done = bool(r0 != 0)
	return
}

func DispatchMessage(msg *MSG) (ret int32) {
	r0, _, _ := syscall.Syscall(procDispatchMessageW.Addr(), 1, uintptr(unsafe.Pointer(msg)), 0, 0)
	ret = int32(r0)
	return
}

func GetDC(hwnd HWND) (ret HDC) {
	r0, _, _ := syscall.Syscall(procGetDC.Addr(), 1, uintptr(hwnd), 0, 0)
	ret = HDC(r0)
	return
}

func StringToUTF16(str string) *uint16 {
	return syscall.StringToUTF16Ptr(str)
}

func (wcex WNDCLASSEX) Sizeof() uint32 {
	return uint32(unsafe.Sizeof(wcex))
}

func GetStockObject(fnObject int) HGDIOBJ {
	ret, _, _ := procGetDeviceCaps.Call(
		uintptr(fnObject))

	return HGDIOBJ(ret)
}

func UnregisterClass(classname *uint16, instance HINSTANCE) bool {
	r0, _, _ := syscall.Syscall(procUnregisterClass.Addr(), 2,
		uintptr(unsafe.Pointer(classname)),
		uintptr(instance),
		0)
	return bool(r0 != 0)
}

func PeekMessage(lpMsg *MSG, hwnd HWND, wMsgFilterMin, wMsgFilterMax, wRemoveMsg uint32) bool {
	ret, _, _ := procPeekMessage.Call(
		uintptr(unsafe.Pointer(lpMsg)),
		uintptr(hwnd),
		uintptr(wMsgFilterMin),
		uintptr(wMsgFilterMax),
		uintptr(wRemoveMsg))

	return ret != 0
}

func GetDesktopWindow() HWND {
	ret, _, _ := procGetDesktopWindow.Call()
	return HWND(ret)
}

func GetClientRect(hwnd HWND) *RECT {
	var rect RECT
	ret, _, _ := procGetClientRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)))

	if ret == 0 {
		panic(fmt.Sprintf("GetClientRect(%d) failed", hwnd))
	}

	return &rect
}