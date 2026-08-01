# DMXBOX TCPプロトコル仕様

[English](../en/tcp-protocol.md) | [ドキュメント一覧](README.md) | [設定リファレンス](configuration.md)

## 目的

オプションのTCPモジュールは、ローカル制御ソフトウェアや専用制御機器から短いコマンドを受信し、DMX FadeまたはOSC Mute操作へルーティングします。

有効化する設定：

```json
{
    "input": {
        "modules": ["http", "tcp"],
        "tcp": {
            "ip": "127.0.0.1",
            "port": 50000
        }
    }
}
```

## 通信方式

| 項目 | 現在の動作 |
| --- | --- |
| トランスポート | TCP |
| デフォルトアドレス | `127.0.0.1:50000` |
| 通信方向 | クライアントがコマンドを送信し、サーバーが確認応答を返す |
| コマンド区切り | CRLF（`\r\n`）、LF（`\n`）、CR（`\r`） |
| 読み取りバッファ | 512バイト |
| 確認応答 | ASCIIの`ack\r\n` |
| 接続の寿命 | 1接続で複数回読み取れる。1コマンド1接続にも対応し、推奨する |

コマンドと識別子にはASCII互換の文字列を想定しています。実装は明示的な文字コード変換を行いません。

## コマンド文法

```text
fade-command = ("fadeIn" | "fadeOut") SP group [SP options]
options      = option *("," option)
option       = name ":" value
mute-command = "mute" [SP boolean]
boolean      = "true" | "false"
```

Fadeの意味が定義されているオプションは`duration`と`interval`です。不明なオプション名は内部メッセージへ渡されますが、DMXハンドラーでは無視されます。

## Fadeコマンド

```text
fadeIn <group>
fadeOut <group>
```

- `<group>`は`dmx.groups`のキーです。
- `fadeIn`は各デバイスを設定済みの`max`値へ遷移させます。
- `fadeOut`は各デバイスを0へ遷移させます。
- オプションを省略すると、`dmx.fadeInterval`を所要時間として使用します。

例：

```text
fadeIn stage
fadeOut stage
```

### コマンド単位のオプション

```text
fadeIn <group> duration:<秒>,interval:<秒>
fadeOut <group> duration:<秒>,interval:<秒>
```

| オプション | 意味 |
| --- | --- |
| `duration` | 遷移時間（秒）。`0`は即時Cutです。 |
| `interval` | 遷移を開始するまでの待ち時間（秒）。0より大きい値で開始を遅延します。 |

例：

```text
fadeIn stage duration:1.5,interval:0
fadeOut stage duration:0,interval:0
```

値は32ビット浮動小数点数として解析されます。不正な値は無視され、設定済みのデフォルト値が使用されます。

## Muteコマンド

```text
mute
mute true
mute false
```

- `mute`と`mute true`はMuteを要求します。
- `mute false`はUnmuteを要求します。
- 受信機器へ送るには、`output.target`でOSCモジュールを有効にする必要があります。

## 互換エイリアス

エイリアスは、単一トークンの完全なコマンドとして受信した場合だけ展開されます。

| エイリアス | 展開後 | 利用条件 |
| --- | --- | --- |
| `fi` | `fadeIn stg` | `stg`グループが存在する場合 |
| `fo` | `fadeOut stg` | `stg`グループが存在する場合 |
| `ci` | `fadeIn stg interval:0` | `stg`グループが存在する場合 |
| `co` | `fadeOut stg interval:0` | `stg`グループが存在する場合 |
| `fai` | `fadeIn aud` | `aud`グループが存在する場合 |
| `fao` | `fadeOut aud` | `aud`グループが存在する場合 |
| `mute` | `mute true` | 常時 |
| `unmute` | `mute false` | 常時 |

現在の`ci`と`co`が設定するのは`interval:0`だけで、`duration:0`は設定しません。そのため必ずしも即時Cutにはなりません。即時遷移が必要な場合は、`duration:0`を明示してください。

## 応答

各行断片を処理した後、サーバーは次を返します。

```text
ack\r\n
```

現在の確認応答が示すのは、入力を読み取り、コマンドパーサーへ渡したことだけです。次の成功は保証しません。

- コマンド名を認識したこと
- グループが存在したこと
- DMXまたはOSCモジュールが内部メッセージを受理したこと
- 物理出力機器が操作を完了したこと

現在、否定応答や構造化されたエラー応答はありません。

## クライアント例

Python：

```python
import socket

with socket.create_connection(("127.0.0.1", 50000), timeout=2) as client:
    client.sendall(b"fadeIn stage\r\n")
    response = client.recv(512)
    assert response == b"ack\r\n"
```

## 既知の制約

- TCPはアプリケーション上のメッセージ境界を保持しません。現在の実装は完全な行までバッファリングせず、ソケットの読み取り単位で解析します。コマンドを短くし、一度に1コマンドを送信し、次を送る前に`ack\r\n`を待ってください。
- 1コマンドは512バイトの読み取りバッファに収める必要があります。
- 空コマンドや不明なコマンドにも`ack\r\n`を返す場合があります。
- 空白には意味があります。コマンド、グループ、オプションの間には半角スペースを1つ使用してください。
- 設定を再読み込みするとTCPリスナーが再生成され、使用中の接続が閉じたり中断したりする場合があります。
- 現在のプロトコルに認証と暗号化はありません。信頼できるネットワークとホストファイアウォールで保護しない限り、`127.0.0.1`へバインドしてください。
