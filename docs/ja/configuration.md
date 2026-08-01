# DMXBOX設定リファレンス

[English](../en/configuration.md) | [ドキュメント一覧](README.md)

DMXBOXはバックエンドのカレントディレクトリにある`config.json`を読み込みます。ソースツリー内の開発用ファイルは`backend/config.json`です。ブラウザから保存すると設定全体が置き換えられ、バックエンドモジュールが再読み込みされます。

手作業で編集する前に、正常に動作している設定をバックアップしてください。

## 全体例

```json
{
    "output": {
        "target": ["console"],
        "ftdi": { "port": "COM3" },
        "artnet": {
            "addr": "127.0.0.1",
            "universe": 0,
            "subuni": 0,
            "net": 0
        },
        "osc": {
            "ip": "127.0.0.1",
            "port": 8765,
            "format": "/yosc:req/set/MIXER:Current/InCh/Fader/On/{}/1",
            "type": "int",
            "inverse": true,
            "channels": [1, 2, 3, 4]
        }
    },
    "input": {
        "modules": ["http"],
        "http": {
            "ip": "127.0.0.1",
            "port": 8080,
            "accepts": [
                "http://localhost:8080",
                "http://127.0.0.1:8080"
            ]
        },
        "tcp": {
            "ip": "127.0.0.1",
            "port": 50000
        }
    },
    "dmx": {
        "groups": {
            "stage": {
                "name": "Stage",
                "devices": [
                    { "model": "dimmer", "channel": 1, "max": [255] },
                    { "model": "wclight", "channel": 10, "max": [255, 255, 0] }
                ]
            }
        },
        "fadeInterval": 0.7,
        "delay": 0,
        "fps": 30
    }
}
```

## `input`

### `input.modules`

有効にする入力モジュールの一覧です。

- `http`：HTTP APIと本番用フロントエンド
- `tcp`：プレーンTCPコマンドサーバー

通常のローカル構成では`http`を含めてください。

### `input.http`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `ip` | string | 待受アドレス。ローカル限定では`127.0.0.1`を使用します。 |
| `port` | integer | HTTPポート。通常は`8080`です。 |
| `accepts` | string array | CORSで許可するOriginです。 |

ループバック以外のアドレスへバインドすると、APIと設定画面がネットワークへ公開されます。DMXBOXは現在アプリケーションレベルの認証を提供していないため、信頼できるネットワークと適切なホストファイアウォールがある場合に限って使用してください。

### `input.tcp`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `ip` | string | TCP待受アドレスです。 |
| `port` | integer | TCP待受ポート。通常は`50000`です。 |

コマンドの区切り、応答、コマンド一覧、互換エイリアスについては[TCPプロトコル仕様](tcp-protocol.md)を参照してください。

## `output`

### `output.target`

有効にする出力の一覧です。

- `console`：ハードウェアを使わずDMX出力をログ表示
- `ftdi`：シリアルDMX出力
- `artnet`：Art-Net出力
- `osc`：OSCのMute/Unmute出力

複数の出力を同時に有効化できます。DMXモジュールには少なくとも1つのDMXレンダラー（`console`、`ftdi`、`artnet`）が必要です。

### `output.ftdi`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `port` | string | `COM3`や`/dev/ttyUSB0`などのシリアルデバイスです。 |

### `output.artnet`

| 項目 | 型 | 範囲 | 説明 |
| --- | --- | --- | --- |
| `addr` | string | — | 宛先IPv4アドレスです。ローカルインターフェースから到達可能である必要があります。 |
| `port` | integer | 省略可 | UDPポート。省略時は6454です。 |
| `net` | integer | 0～127 | Art-NetのNetフィールドです。 |
| `subuni` | integer | 0～15 | Art-NetのSub-Netフィールドです。 |
| `universe` | integer | 0～15 | Art-NetのUniverseフィールドです。 |

### `output.osc`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `ip` | string | OSC受信先アドレスです。 |
| `port` | integer | OSC受信先UDPポートです。 |
| `format` | string | OSCアドレスのテンプレートです。すべての`{}`がチャンネル値へ置換されます。 |
| `type` | string | `int`は`int32`、`float`は`float32`として送信します。 |
| `inverse` | boolean | 受信先へ送るMute値を反転します。 |
| `channels` | integer array | アドレステンプレートへ代入する値です。 |

## `dmx`

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `groups` | object | 固定のグループIDからグループ定義へのマップです。 |
| `fadeInterval` | number | デフォルトのフェード時間（秒）です。 |
| `delay` | number | 設定上のエフェクト遅延（秒）です。 |
| `fps` | number | DMX更新レートです。フロントエンドは0～99を許容しますが、正の値を使用してください。 |

### グループ

`dmx.groups`の各キーはHTTP・TCP制御APIから使用するグループIDです。各グループは次を持ちます。

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `name` | string | ブラウザに表示する名前です。 |
| `devices` | array | グループとして同時に制御するデバイスです。 |

### デバイス

| 項目 | 型 | 説明 |
| --- | --- | --- |
| `model` | string | `dimmer`または`wclight`です。 |
| `channel` | integer | 1始まりのDMX開始チャンネル。1～512です。 |
| `max` | integer array | Fade In時のチャンネル別最大値。通常は0～255です。 |

デバイスのチャンネル構成：

- `dimmer`：1チャンネル。`max`は1要素必要です。
- `wclight`：寒色、暖色、フラッシュの3チャンネルです。現在フラッシュは常に0です。

デバイス全体が512チャンネルの範囲内に収まる必要があります。例えば`wclight`を開始できるのは510チャンネルまでです。

シリアルDMX、Art-Net、デバイス動作、フェード、終了時の挙動については[照明機材・出力仕様](lighting-specifications.md)を参照してください。

## 開発用フロントエンド設定

`frontend/public/config.json`はバックエンド設定とは別のファイルです。

```json
{
    "backendPort": 8080
}
```

Viteで起動する開発用フロントエンドが、ローカルバックエンドのポートを知るために使用します。本番リリースでは、フロントエンドとバックエンドの両方をGoのHTTPサーバーが配信します。
