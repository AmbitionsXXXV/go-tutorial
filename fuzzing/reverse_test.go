package main

import (
	"testing"      // 导入 testing 包，用于编写测试用例
	"unicode/utf8" // 导入 unicode/utf8 包，用于检查字符串是否为有效的 UTF-8
)

// FuzzReverse 是一个模糊测试函数，用于测试 Reverse 函数
func FuzzReverse(f *testing.F) {
	// 定义一些初始的测试用例
	testcases := []string{"Hello, world", " ", "!12345"}

	// 通过 f.Add 向模糊测试引擎添加这些测试用例
	// 这些用例将作为种子语料，帮助引擎开始测试
	for _, tc := range testcases {
		f.Add(tc) // 使用 f.Add 提供种子语料
	}

	// f.Fuzz 接收一个匿名函数，用于定义具体的测试逻辑
	f.Fuzz(func(t *testing.T, orig string) {
		// 调用 Reverse 函数反转原始字符串
		rev, err1 := Reverse(orig)

		// 如果反转时产生错误，则终止这次测试迭代
		if err1 != nil {
			return
		}

		// 再次调用 Reverse 函数，对反转后的字符串进行反转
		doubleRev, err2 := Reverse(rev)

		// 如果第二次反转时产生错误，则终止这次测试迭代
		if err2 != nil {
			return
		}

		// 检查两次反转后的字符串是否与原始字符串相同
		// 如果不同，则报告错误
		if orig != doubleRev {
			t.Errorf("Before: %q, after: %q", orig, doubleRev)
		}

		// 检查原始字符串是否为有效的 UTF-8
		// 并确保反转后的字符串同样是有效的 UTF-8
		// 如果不是，则报告错误
		if utf8.ValidString(orig) && !utf8.ValidString(rev) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", rev)
		}
	})
}
