# microservice-monitor

microservice-monitorは、K8s上で動いているマイクロサービスのCPU/Memory使用状況を収集するサービスです。


# 概要

microservice-monitorは、一定期間ごとにK8sのmertics-serverから各マイクロサービスのCPU/Memory使用状況を取得します。  
取得するCPU/Memoryは、Pod単位とContainer単位での取得が可能です。  
取得したCPU/Memory使用状況は、JSON形式でlogとして出力されます。

また、microservice-monitorにはアラート機能があり、一定の閾値を超えた場合にアラートを出すことができます。

# 動作環境

microservice-monitorは、K8sが動作する環境であればどこでも動作可能です。
ただし、前提としてK8s上にmetrics-serverがデプロイされている必要があります。

# セットアップ

以下のコマンドでDockerイメージの作成ができます。

```
$ make docker-build
```

# 起動方法

以下のコマンドでK8s上にデプロイできます。

```
$ kubectl apply -f deployments
```

# 環境変数

|環境変数名|値|
|----------|--|
|INTERVAL|5(s)|
|WINDOW_SIZE|5|
|CONFIG_DIR|path to config file dir|
|SLACK_URL|slack incoming url|

また、アラート通知を行いたいマイクロサービスをYAMLファイルで定義することができます。詳細はexampl/alert-setting.ymlを参照してください

# CPU/Memory使用状況のログ

以下の形式でマイクロサービス単位のCPU/Memory使用状況のログが出力されます。

```
{"timestamp":"2021-04-15T08:25:58Z","name":"coredns-f9fd979d6-fd2fg","containers":[{"name":"coredns","metrics":{"cpu":16,"memory":18096128}}],"metrics":{"cpu":16,"memory":18096128}}
```
