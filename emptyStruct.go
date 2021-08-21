// https://devlights.hatenablog.com/entry/2019/12/03/114645
// の写経
package main

import (
	"fmt"
	"log"
	"time"
	"unsafe"
)

type (
	es3 struct{}
	es4 struct{}
)

func (es3) f1() string {
	return "hello es3"
}

func (es4) f2() string {
	return "hello es4"
}

func EmptyStruct() error {
	// --------------------------------------------------
	// 空の構造体 (Empty struct)
	//   - 空の構造体は struct{} で表す
	//   - 空の構造体はメモリサイズが 0
	//     - 型のメモリサイズは Unsafe.sizeof() で調べられる
	//   - 空の構造体は 属性 を持たない
	//   - 空の構造体は 同じアドレス を示す
	//     - なので、型の別名を付与してメソッドを作ることで簡易的な utilクラス みたいに出来る
	//   - メモリサイズが 0 なので、終了通知だけを送る done チャネルなどに便利
	//
	// ref: https://dave.cheney.net/2014/03/25/the-empty-struct
	// --------------------------------------------------
	var (
		emptyStruct    struct{}
		emptyInterface interface{}
	)

	emptyStruct = struct{}{}
	emptyInterface = emptyStruct

	memSize1 := unsafe.Sizeof(emptyStruct)
	memSize2 := unsafe.Sizeof(emptyInterface)
	fmt.Printf("メモリサイズ：\tEmptyStruct[%d] EmptyInteface[%d]\n", memSize1, memSize2)

	var (
		v1 = struct{}{}
		v2 = struct{}{}
	)

	fmt.Printf("アドレス確認：\tv1[%p]\tv2[%p]\taddr_equal? [%v]\n", &v1, &v2, v1 == v2)

	type (
		es1 struct{}
		es2 struct{}
	)

	var (
		v3 = es1{}
		v4 = es2{}
	)

	// 型の別名を付与した状態でも、アドレスは同じか？
	// （補足）v3 == v4 は、型が異なるためコンパイルエラーとなる
	fmt.Printf("別の型名、かつ、両方とも空のstruct{}の場合のアドレス：\tv3[%p(%T)]\tv4[%p(%T)]\n", &v3, v3, &v4, v4)

	// 型 es3 と es4 の定義は、本ソースコードの上部を参照。
	var (
		v5 = es3{}
		v6 = es4{}
	)

	// 型の別名を付与して、さらにメソッドを付与しても、アドレスは同じか？
	// （補足）v5 == v6 は、型が異なるためコンパイルエラーになる
	fmt.Printf("別の型名、かつ、メソッド付与した場合のアドレス：\tv5[%p(%T)]\tv6[%p(%T)]\n", &v5, v5, &v6, v6)

	// es3とes4は、元々は同じ struct{} だが、異なる型別名がついているので、Go内部では全く別の型として扱われる
	// なので、それぞれに定義したメソッドも、ちゃんと対象となる型の方に紐づく
	fmt.Println(v5.f1(), v6.f2())

	// done チャネルと空構造体の組み合わせ
	// done チャネルは、処理終了を通知したいだけなので、 close しかない
	type (
		nop struct{}
	)

	var (
		done = make(chan nop)
	)

	go func() {
		defer close(done)

		log.Println("\t==> goroutine begin")
		time.Sleep(2 * time.Second)
		log.Println("\t==> goroutine end")
	}()

	log.Println("main goroutine wait start")
	<-done
	log.Println("main goroutine wait end")

	return nil

}
