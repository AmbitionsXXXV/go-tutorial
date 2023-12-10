package main

import (
	"errors"       // 导入 errors 包，用于创建错误信息
	"fmt"          // 导入 fmt 包，用于格式化输出
	"unicode/utf8" // 导入 unicode/utf8 包，用于处理 UTF-8 字符串
)

// Reverse 函数用于反转一个 UTF-8 编码的字符串
// 如果输入的字符串不是有效的 UTF-8，则返回错误
func Reverse(s string) (string, error) {
	// 检查字符串是否为有效的 UTF-8
	if !utf8.ValidString(s) {
		return s, errors.New("input is not valid UTF-8")
	}

	fmt.Printf("input: %q\n", s)

	// 将字符串转换为 rune 切片
	// 使用 rune 而不是 byte，因为一个 Unicode 字符可能由多个字节组成
	// rune 能确保每个 Unicode 字符被完整表示，避免在字符中间切分造成的问题
	b := []rune(s)

	fmt.Printf("runes: %q\n", b)

	// 通过交换元素的方式反转 rune 切片
	for i, j := 0, len(b)-1; i < len(b)/2; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}

	// 将反转后的 rune 切片转换回字符串，并返回该字符串
	return string(b), nil
}

func main() {
	// 定义一个测试用的字符串
	input := "The quick brown fox jumped over the lazy dog"

	// 调用 Reverse 函数反转字符串
	rev, revErr := Reverse(input)
	// 再次调用 Reverse 函数对反转后的字符串进行反转
	doubleRev, doubleRevErr := Reverse(rev)

	// 输出原始字符串、第一次反转后的字符串以及第二次反转后的字符串
	fmt.Printf("original: %q\n", input)
	fmt.Printf("reversed: %q, err: %v\n", rev, revErr)
	fmt.Printf("reversed again: %q, err: %v\n", doubleRev, doubleRevErr)
}
