package catch

import "C"
import (
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

//export goSigHandler
func goSigHandler(sig C.int, cStack **C.char, cLen C.int) {
	convStack := visualStack(cStack, cLen)
	dump(convStack, debug.Stack(), sigErrTable[syscall.Signal(sig)])
}

// https://github.com/golang/go/wiki/cgo#Turning_C_arrays_into_Go_slices
func visualStack(cStack **C.char, cLen C.int) (conv []byte) {
	tmpSlice := (*[1 << 30]*C.char)(unsafe.Pointer(cStack))[:cLen:cLen]
	for _, s := range tmpSlice {
		var elf, sym, relAddr, absAddr string
		goStr := strings.ReplaceAll(C.GoString(s), " ", "") // 去除空格符
		seps := []string{"(", "+", ")", "[", "]"}
		sArrs := strSplits(goStr, seps)
		if len(sArrs) == 4 {
			elf, sym, relAddr, absAddr = sArrs[0], sArrs[1], sArrs[2], sArrs[3]
		} else if len(sArrs) == 3 {
			elf, relAddr, absAddr = sArrs[0], sArrs[1], sArrs[2]
		}
		if strings.HasPrefix(relAddr, "0x") {
			relAddr = relAddr[2:]
		}
		if strings.HasPrefix(absAddr, "0x") {
			absAddr = absAddr[2:]
		}

		if len(sym) == 0 {
			// 取绝对地址取址, 防止cgo相关符号获取异常;
			absAddr = relAddr
		} else if elf != os.Args[0] {
			absAddr = dynSymAddr(elf, sym, relAddr)
		}
		cmd := exec.Command("addr2line", "-Cfe", elf, absAddr)
		if op, err := cmd.Output(); err == nil {
			goStr = string(op)
		}
		goStr = strings.Replace(goStr, "\n", "\n        ", 1)
		conv = append(conv, []byte(goStr)...)
	}
	return
}

func dynSymAddr(elf string, sym string, offset string) (address string) {
	// 调用库取符号及相对地址 (若sym为空,相对地址即最终地址), 取动态符号表数据;
	cmd := exec.Command("readelf", "--dyn-syms", elf)
	if op, err := cmd.Output(); err == nil {
		symtab := strings.Split(string(op), "\n")
		for _, v := range symtab {
			if strings.Contains(v, sym) {
				symLine := strings.Fields(v)
				if len(symLine) == 8 {
					if strings.HasPrefix(symLine[7], sym+"@") || strings.Compare(symLine[7], sym) == 0 {
						hex, err := strconv.ParseInt(symLine[1], 16, 64)
						hex2, err2 := strconv.ParseInt(offset, 16, 64)
						if err == nil && err2 == nil {
							address = fmt.Sprintf("%x", hex+hex2)
							break
						}
					}
				}
			}
		}
	}
	return
}

func strSplits(src string, seps []string) (arrs []string) {
	arrs = append(arrs, src)
	for _, sep := range seps {
		var tmp []string
		for _, str := range arrs {
			arr := strings.Split(str, sep)
			for _, v := range arr {
				if len(v) > 0 {
					tmp = append(tmp, v)
				}
			}
		}
		arrs = tmp
	}
	return
}
