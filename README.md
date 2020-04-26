## AWS Athena，Kinesis Hands-on 

###前提条件
AWS GLOBAL账号, 建议配置AWS Cli工具

安装配置 aws cli , version > 1.16.312

```bash
 #利用pip安装
 pip3 install awscli --upgrade --user

 #利用awscli-bundle安装 linux / macOS
  curl "https://s3.amazonaws.com/aws-cli/awscli-bundle.zip" -o "awscli-bundle.zip"
  unzip awscli-bundle.zip
  sudo ./awscli-bundle/install -i /usr/local/aws -b /usr/local/bin/aws
  
  #查看aws cli版本
  aws --version
  aws-cli/1.16.312 Python/3.7.6 Darwin/18.7.0 botocore/
```

配置aws cli的用户权限
```bash
 #使用aws configure配置aws cli的AccessKey/SecrectAccessKey
 aws configure
 AWS Access Key ID :
 AWS Secret Access Key :
 Default region name:
 Default output format [None]:

 #测试AK/SK是否生效
 aws sts get-caller-identity

 #如果可以正常返回以下内容(包含account id),则表示已经正确设置角色权限
 {
    "Account": "<your account id, etc.11111111>", 
    "UserId": "AIDAIG42GHSYU2TYCMCZW", 
    "Arn": "arn:aws:iam::<your account id, etc.11111111>:user/<iam user>"
 }
```


###Lab1: 使用athena 访问 s3 数据
进入ahtena服务

1.创建 nginx_access_logs 和app01_user_login 表
 步骤1分别创建2个表，支持两种不同的格式, **请注意**把s3://aws-hands-on-athena换成你自己的s3 bucket.

1.1 创建Athena nginx_access_logs表,使用正则表达式匹配nginx access log.

```bash
  #nginx access log 格式:
  127.0.0.1 - - [19/Jun/2012:09:16:22 +0100] "GET /GO.jpg HTTP/1.1" 499 0 "http://domain.com/htm_data/7/1206/758536.html" "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; .NET CLR 1.1.4322; .NET CLR 2.0.50727; .NET CLR 3.0.4506.2152; .NET CLR 3.5.30729; SE 2.X MetaSr 1.0)"

```
打开Athena 选择数据库(默认是default),在Query Editor中输入建表语句,点击运行
```sql
-- nginx_access_log 建表语句
CREATE EXTERNAL TABLE `nginx_access_logs`(
  `remote_addr` string COMMENT '', 
  `request_time` string COMMENT '', 
  `request_method` string COMMENT '', 
  `request_url` string COMMENT '', 
  `http_protocol` string COMMENT '',
  `http_code` int, 
  `response_size` string COMMENT '', 
  `http_referer` string COMMENT '', 
  `http_user_agent` string COMMENT '')
ROW FORMAT SERDE 
  'org.apache.hadoop.hive.serde2.RegexSerDe' 
WITH SERDEPROPERTIES ( 
  'input.regex'='([.0-9]*) - - ([^\"]*)\"([^\ ]*) ([^\ ]*)(.*?)\" (-|[0-9]*) (-|[0-9]*) (\".*?\") (\".*?\")') 
LOCATION
  's3://aws-hands-on-athena/nginx'
```

-------

1.2创建Athena app01_user_login表 配置user_login的json格式数据.

```bash
 #user_login json格式
 {"id":1,"first_name":"Stanford","last_name":"Wasmuth","email":"swasmuth0@stanford.edu","gender":"Male","ip_address":"27.14.197.121","lastlogin":"2019-12-23 07:46:51"},
 {"id":2,"first_name":"Eve","last_name":"Maeer","email":"emaeer1@shop-pro.jp","gender":"Female","ip_address":"223.213.166.71","lastlogin":"2019-06-03 07:59:29"}
```
运行以下建表语句
```sql
-- app01_user_login 建表语句
CREATE EXTERNAL TABLE `app01_user_login`(
  `id` string COMMENT '', 
  `first_name` string COMMENT '', 
  `last_name` string COMMENT '', 
  `email` string COMMENT '', 
  `gender` string COMMENT '', 
  `ip_address` string COMMENT '', 
  `lastlogin` TIMESTAMP COMMENT '')
ROW FORMAT SERDE 
  'org.openx.data.jsonserde.JsonSerDe' 
LOCATION
  's3://aws-hands-on-athena/app01'
```


-------

2.使用web console 或者aws cli 上传data目录中的nginx-access-example.log.gz 和user-login-example.json.gz文件.

```bash
 aws s3 cp nginx-access-example.log.gz  s3://aws-hands-on-athena/nginx/2020/04/25/nginx-access-example.log.gz
 aws s3 cp user-login-example.json.gz  s3://aws-hands-on-athena/app01/2020/04/25/user-login-example.json.gz
```
-------
3.测试athena query

打开athena Query Editor 
查询nginx 日志

```sql
 select * from nginx_access_logs
```

查询用户登陆日志
```sql
 select * from app01_user_login where lastlogin BETWEEN date '2019-01-01' AND date '2019-06-01'
```
-------
4.使用golang访问athena(可选)
 
 ```bash
  #编译golang项目
  cd golang-example
  go build
  
  #测试,请将-out替换成你自己的bucket
  ./athena-go-query -sql "select * from app01_user_login limit 1" -out s3://aws-hands-on-athena/query/
 ``` 

