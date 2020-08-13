# hubcontrib tweeter
24時間の総コミット数と"芝生"をツイートします

## 使いかた
同一ディレクトリに`.env`ファイルを追加し, トークンを記述します.

```.env
GITHUB_TOKEN={GitHubのトークン}
CONSUMER_KEY={Twitterのトークン1}
CONSUMER_SECRET={Twitterのトークン2}
ACCESS_TOKEN={Twitterのトークン3}
ACCESS_SECRET={Twitterのトークン4}
```

※ GitHubのトークンは [ここ](https://github.com/settings/tokens) の `Personal access tokens` からいただいてください

起動すると24時間の総コミット数と"芝生"をツイートします
