<h1 style="text-align: center;">
  システム設計演習 レポート
</h1>
<p style="text-align: right;">学籍番号: 20B30790 &nbsp;&nbsp;藤井 &nbsp; 一喜 &nbsp;&nbsp;&nbsp;</p>

## 必須課題

### 1. HTTP 通信の基本

#### 1.1. HTTP Request Methods

参考資料: MDN Web Docs - [HTTP メソッド](https://developer.mozilla.org/ja/docs/Web/HTTP/Methods), Real World HTTP (オライリー・ジャパン)

- GET

  - 指定したリソースの情報を取得する。

- HEAD

  - GET と同様にリソースの情報を取得するが、レスポンスボディは要求しない

- POST

  - 指定したリソースに、情報を送信する際に使用するメソッド。リクエストボディにデータを含めることで、サーバーにデータを送信する。

- PUT

  - 指定したリソースに情報を送信することで、更新を行う目的で使用される

- DELETE
  - 指定したリソースを削除するために使用される

#### 1.2 HTTP Header

参考資料: MDN Web Docs - [HTTP ヘッダー](https://developer.mozilla.org/ja/docs/Web/HTTP/Headers), Real World HTTP (オライリー・ジャパン)

- Host

  - リクエスト先のホスト名とポート番号を指定する
  - 例: Tokyo Tech Portal マトリクスコードログインでは `Host: portal.nap.gsic.titech.ac.jp`となっている

- User-Agent

  - リクエストを送信したクライアントの情報を指定する
  - サーバーはこの情報を元に、クライアントの種類やバージョンを判断する

    (ブラウザ版からのアクセスなのか、モバイルアプリ版からのアクセスなのかを判別など)

    例: `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36` (Tokyo Tech Portal)

- Accept

  - クライアントが受け付けるレスポンスのメディアタイプを指定する
  - サーバーはこの情報を元に、クライアントが受け付けるレスポンスのメディアタイプを判断する

    例: `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9` (Tokyo Tech Portal)

- Referer

  - リクエストを送信する前に表示していたページの URL を指定する
  - サーバーはこの情報を元に、リファラログを記録することができる

    例: `https://portal.nap.gsic.titech.ac.jp/GetAccess/Login?Template=userpass_key&AUTHMETHOD=UserPassword`

- Accept-Encoding

  - クライアントが受け付けるレスポンスのエンコーディングを指定する
  - サーバーはこの情報を元に、クライアントが受け付けるレスポンスのエンコーディングを判断する (gzip など)

    例: Tokyo Tech Portal では`Accept-Encoding: gzip, deflate, br`

- Accept-Language

  - クライアントが受け付けるレスポンスの言語を指定する
  - サーバーはこの情報を元に、クライアントが受け付けるレスポンスの言語を判断する (Tokyo Tech Portal では `Accept-Language: en-US,en;q=0.9,ja;q=0.8`)

- Content-Type

  - リクエストボディのメディアタイプを指定する
  - サーバーはこの情報を元に、リクエストボディのメディアタイプを判断する (JSON など)

    例: Tokyo Tech Portal `Content-Type: text/html;charset=UTF-8` (Response Headers)

#### 1.3 HTTP Method Idempotence and Safe

- GET

  - Idempotence: Yes
  - Safe: Yes

- HEAD

  - Idempotence: Yes
  - Safe: Yes

- POST

  - Idempotence: No
  - Safe: No

- PUT

  - Idempotence: Yes
  - Safe: No

- DELETE
  - Idempotence: Yes
  - Safe: No

### 2. タスク管理アプリケーション

GitHub リポジトリ: [Link](https://github.com/okoge-kaz/todo-application)

#### 2.1 基本仕様

#### 2.2 シークエンス図

#### 2.3 データベース設計

## 加点課題

### 3. Web アプリケーション開発

#### 3.1. Stateless Type and Session Type

- Stateless Type

  - サーバーはクライアントの状態を保持しない
  - クライアントはリクエストごとにサーバーに対して認証情報を送信する必要がある

- Session Type
  - サーバーはクライアントの状態を保持する
  - クライアントはリクエストごとにサーバーに対して認証情報を送信する必要がない

#### 3.2 SQL Injection

- SQL Injection
  - SQL インジェクションとは、Web アプリケーションの SQL クエリに対して、意図しない SQL 文を注入する攻撃手法のこと
  - SQL インジェクションを行うと、データベースの情報を改ざんしたり、データベースの情報を盗み出したりすることができる

#### 3.3 XSS (Cross Site Scripting) and CSRF (Cross Site Request Forgery)

- XSS (Cross Site Scripting)

  - XSS とは、Web アプリケーションに対して、意図しない JavaScript を注入する攻撃手法のこと
  - XSS を行うと、クライアントの Cookie を盗み出したり、クライアントの Cookie を改ざんしたりすることができる

- CSRF (Cross Site Request Forgery)

  - CSRF とは、Web アプリケーションに対して、意図しないリクエストを送信する攻撃手法のこと
  - CSRF を行うと、クライアントの Cookie を盗み出したり、クライアントの Cookie を改ざんしたりすることができる

- CSRF 対策

  - CSRF 対策とは、CSRF 攻撃を防ぐための対策のこと
  - CSRF 対策としては、CSRF トークンを使用することがある
  - CSRF トークンとは、クライアントがリクエストを送信する際に、サーバーから発行されるランダムな文字列のこと
  - サーバーは、クライアントがリクエストを送信する際に、CSRF トークンをリクエストボディに含めるように要求する
  - サーバーは、クライアントがリクエストを送信する際に、CSRF トークンを Cookie に含めるように要求する
  - サーバーは、クライアントがリクエストを送信する際に、CSRF トークンをリクエストヘッダに含めるように要求する

- XSS 対策

  - XSS 対策とは、XSS 攻撃を防ぐための対策のこと
  - XSS 対策としては、HTML エスケープを使用することがある
  - HTML エスケープとは、HTML の特殊文字をエスケープすることのこと
  - HTML エスケープとしては、HTML エンティティを使用することがある
  - HTML エンティティとは、HTML の特殊文字をエスケープするための文字列のこと
  - HTML エンティティとしては、`&` を `&amp;` に、`<` を `&lt;` に、`>` を `&gt;` に、`"` を `&quot;` に、`'` を `&#x27;` に、`/` を `&#x2F;` に変換することがある

- CORS (Cross Origin Resource Sharing)
  - CORS とは、クロスオリジン間でリソースを共有するための仕組みのこと

#### 3.4 Server Side Rendering

- SSR (Server Side Rendering)
  - SSR とは、クライアント側で HTML を生成するのではなく、サーバー側で HTML を生成することのこと
  - SSR を行うと、クライアント側で HTML を生成するのではなく、サーバー側で HTML を生成することができる

#### 3.5 自由記述欄
