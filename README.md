## 何？

ReactiveXのGo版である、[RxGo](https://github.com/ReactiveX/RxGo)を使ってみたサンプル集です

## サンプル

| メイン| 概要 |
|---|---|
| [csv_parser](https://github.com/ngiyshhk/rxgo-examples#csv_parser) | 入力csvの指定したカラムのみ標準出力に出力する |
| [web_watcher](https://github.com/ngiyshhk/rxgo-examples#web_watcher) | 指定URLの指定要素を定期的に監視して、更新分をslackに通知する |

## csv_parser
### 概要
入力csvの指定したカラムのみ標準出力に出力する

### 使い方
| オプション | 説明 |
|---|---|
| C | 出力カラム(カンマ区切り) |
| t | tsvの場合、指定する |
```shell
% bin/mac/csv_parser -h
Usage of bin/mac/csv_parser:
  -C string
        display columns
  -t    tsv
```

### 実行例(osx)
```shell
curl http://hojin.ctot.jp/markets/CSV/01_USDJPY_D.csv | \
nkf -w | \
bin/mac/csv_parser -C 1,3
```

## web_watcher
### 概要
指定URLの指定要素を定期的に監視して、更新分をslackに通知する

### 使い方
| オプション | 説明 | デフォルト |
|---|---|---|
| t | 監視するページ | - |
| c | ページの文字コード | utf-8 |
| e | 監視する要素 | - |
| a | 取得する属性 | text |
| s | 監視する間隔(秒) | 10 |
| p | postするslackのincoming-webhook url | - |

```shell
% bin/mac/web_watcher -h
Usage of bin/mac/web_watcher:
  -a string
        target attribute (default "text")
  -c string
        target charset (default "utf-8")
  -e string
        target element
  -p string
        slack incoming hook url
  -s int
        access duration[second] (default 10)
  -t string
        target url
```

### 実行例(osx)
```shell
# ライブドアニュースのタイトルを300秒に一回取得し、更新分をslackに通知する
bin/mac/web_watcher \
-t "http://news.livedoor.com/" \
-c "euc-jp" \
-e "ul.topicsList > li" \
-s 300 \
-p "https://hooks.slack.com/services/hogehoge"
```
