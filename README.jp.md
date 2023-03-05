# IVIS NAS External Web Server
[![ko](https://img.shields.io/badge/lang-ko-red.svg)](https://github.com/picel/IVIS_NAS_Helper/blob/main/README.md)
[![ja](https://img.shields.io/badge/lang-ja-blue.svg)](https://github.com/picel/IVIS_NAS_Helper/blob/main/README.jp.md)

## 概要
- IVIS Labの外部ネットワークからNASにアクセスするためのWebサーバー。
- 外部ネットワークからNASにアクセスするためのUIを提供。

## 開発環境
- Go
- NGINX (Reverse Proxy)
- Redis (セッション管理)
- ZeroTier (NASと通信)

## ルート情報
- GET /
    - ログインページ
- POST /login
    - ログインリクエスト
    - id, pwをbodyに含めてリクエスト
    - ZeroTier Network (192.168.195.1) の /loginCheck ページにリクエスト
    - 成功時にJWTトークンを発行、失敗時にログインページにリダイレクト
- GET /logout
    - ログアウト
    - JWTトークンを削除
    - ログインページにリダイレクト
    ---
- GET /files/:path
    - ファイル一覧ページ
    - pathはファイルのパス
- GET /serve_file/:path
    - ファイルをブラウザで表示
    - pathはファイルのパス
- GET /download_file/:path　（バグあり）
    - ファイルをダウンロード
    - pathはファイルのパス
- GET /download_dir/:path
    - ディレクトリをダウンロード
    - pathはディレクトリのパス
    - ディレクトリをzipに圧縮してダウンロード

    ---
- POST /remote_login
    - [IVIS Admin Panel](https://github.com/picel/ivis_admin) からリクエストされるリモートログインAPI
    - 成功時にJWTトークンを発行、失敗時に401 Unauthorizedを返す
- POST /remove_token_verify
    - [IVIS Admin Panel](https://github.com/picel/ivis_admin) からリクエストされるJWTトークンの検証API
    - JWTトークンが有効な場合に200 OKを返す、無効な場合に401 Unauthorizedを返す

## 実行画面
- ログインページ
    ![login](https://user-images.githubusercontent.com/30901178/222891856-9b6833ec-d093-452b-8ebe-2d31ac5d89d3.png)
- ファイル一覧ページ
    ![files](https://user-images.githubusercontent.com/30901178/222891876-d88f0054-d227-4fab-8341-242232ded8ea.png)

## 注意事項
- keyファイルはgitに上げない。
    - keyファイルは/keyディレクトリに配置。
    - JWTVerifyKey []byte 定義
