# Mastodon yukibot - GO

> [!WARNING]
> It's MyGO🥒🥒🥒🥒🥒

> 我从来没觉得写代码开心过

## 开始之前

在 `MS_HOST/settings/applications` (MS_HOST 为你的`Mastodon`地址) 创建一个应用

之后你可以得到

应用 ID "xxx"

应用密钥 "xxx"

你的访问令牌 "xxx"

前往 https://aistudio.google.com/app/apikey 获取你的`GEMINI_API_KEY`

## 配置

在项目目录新建`.env`文件

```.env
MS_HOST="https://xxx.xxx"
MS_CID="xxx"
MS_SECRET="xxx"
MS_TOKEN="xxx"
GEMINI_API_KEY="xxx"
```

-   MS_HOST 为你的`Mastodon`地址
-   MS_CID 为你获取到的`Mastodon`应用 ID
-   MS_SECRET 为你获取到的`Mastodon`应用密钥
-   MS_TOKEN 为你获取到的`Mastodon`你的访问令牌
-   GEMINI_API_KEY 为你在`aistudio`获取到的 apikey

## 运行

终端下输入

```shell
go mod tidy && go build
./MyGO
```
