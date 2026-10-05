---
description: <script setup lang="go"> のブロックに書けるもの、変数と関数のテンプレートからの見え方、props の受け取り方。
---

# Go ブロックの書き方

<!--@include: ../_gosfc/readme-writing-the-go-block.md-->

## 状態とイベント

ブロックが宣言する変数はコンポーネントの状態です。テンプレートのイベントが Go の関数を
呼ぶとそれが実行され、テンプレートは新しい値を表示します。

```vue
<template>
  <button type="button" @click="Increment">{{ count }} 回クリックしました</button>
  <input :value="name" @input="Rename($event.target.value)" />
  <p>こんにちは、{{ name }}</p>
</template>

<script setup lang="go">
count := 0
name := "Gopher"

func Increment() {
	count++
}

func Rename(s string) {
	name = s
}
</script>
```

テンプレートが受け取るのは、各値を JavaScript に変換したコピーです (スライスは配列、
構造体はオブジェクト)。JavaScript 側でコピーを変えても Go の変数は変わりません。Go の状態は
Go の関数を呼んで変えてください。
