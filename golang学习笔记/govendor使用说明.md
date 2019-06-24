1. 下载

```bash
go get -u github.com/kardianos/govendor
```

2. 初始化项目为vendor管理的项目

```bash
# 1.进入项目根目录
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>chdir
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod

# 2.初始化，无非就是在当前项目中创建vendor目录和vendor/vendor.json文件
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>govendor init

E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>

# 3.应用外部包，  +external (e) referenced packages in GOPATH but not in current project
# 这样项目就会去GOPATH中去找本项目依赖的包到vendor目录中
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>govendor add +external

E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>
```

3 .govendor常用操作

```bash
# 查看当前项目依赖包：包括列出项目本身的包和vendor中的外部包
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>govendor list
 v  github.com/gorilla/websocket
 v  github.com/tidwall/gjson
 v  github.com/tidwall/match
 v  gopkg.in/olivere/elastic.v3
 v  gopkg.in/olivere/elastic.v3/backoff
 v  gopkg.in/olivere/elastic.v3/uritemplates
 l  openshift-statistics-analysis-delete-pod
 l  openshift-statistics-analysis-delete-pod/pkg/elasticsearch
 l  openshift-statistics-analysis-delete-pod/pkg/utils

E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>


# Look at what is using a package
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>govendor list -v fmt
 s  fmt
    ├──  v  openshift-statistics-analysis-delete-pod/vendor/gopkg.in/olivere/elastic.v3
    ├──  v  openshift-statistics-analysis-delete-pod/vendor/gopkg.in/olivere/elastic.v3/uritemplates
    ├──  l  openshift-statistics-analysis-delete-pod
    ├──  l  openshift-statistics-analysis-delete-pod/pkg/elasticsearch
    └──  l  openshift-statistics-analysis-delete-pod/pkg/utils

E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>

# Test your repository only
E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>govendor test +local
?       openshift-statistics-analysis-delete-pod        [no test files]
?       openshift-statistics-analysis-delete-pod/pkg/elasticsearch      [no test files]
ok      openshift-statistics-analysis-delete-pod/pkg/utils      1.328s

E:\2017\000programming\go\code\src\openshift-statistics-analysis-delete-pod>


```

