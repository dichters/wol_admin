# mydev-go开发说明
---
## 1. 项目需求说明
这个项目是一个代码模板，是一个fiber webservice后端。后续的代码都将基于这个样板工程完成。服务接口使用restful http接口提供。
http接口默认返回application/json。
## 2. 工程目录结构
项目父目录是mydev-go。包含两个子项目：
- mydev-service：web后端工程。
- mydev-api: 公共类放在这里（比如输入返回的类、异常码等等）
- mydev-sdk：通过提供sdk包，供外部直接调用（免于自己拼接http接口）。
- 我们人为约定，所有的子模块版本保持一致。
## 3. 指定版本
- go: 1.25.6
- fiber: v2.50.0

其余依赖如果需要，请结合当前截止时间（2026年1月），选用最新的稳定版或者长期支持版。
## 4. 功能细节
### 4.1. 提供基本的http接口样板
本项目的http接口，一般以application/json提供。格式为：
```json
{
  "code": 0,
  "msg": "Success.",
  "data": []
}
```
其中：
code：返回码，Integer类型。返回0为成功，其它为失败。
msg：返回信息。当code=0时返回"Success."，其它情况则返回具体错误原因所提供的信息。
data: 具体的数据，list格式。如果没有任何需要返回的数据，则返回空list。即使只有1个元素，也要返回list。
本服务的接口全部使用POST。
- 提供一个{host}/srv/v1/hc，返回：
```json
{
  "code": 0,
  "msg": "Success.",
  "data": []
}
```
- 提供一个{host}/srv/v1/divide，参数：
```json
{
  "a": 1,
  "b": 2
}
```
其中，a和b为数值。

成功返回：
```json
{
  "code": 0,
  "msg": "Success.",
  "data": [0.5]
}
```
data只有一个元素，为a/b的值。
如果b为0，不要做额外处理。异常处理交给拦截器。
### 4.2 异常处理相关
#### 异常类：
设计一个BaseException extends Exception。内置属性：code, msg。

另设计BizException extends BaseException 和 ServiceException extends BaseException。
#### 异常码code的枚举值：
放一个枚举类。随便写两个值，后续再完善。

### 4.3 关键拦截器

#### 异常处理拦截器Middleware
捕获所有未处理的错误（Error）和宕机（Panic），并返回标准化的 JSON 错误信息。
分类处理规则：
业务错误 (BizException)：
Code：使用异常对象中定义的 Code。
Msg：使用异常对象中定义的 Msg。
服务错误 (ServiceException)：
Code：使用异常对象中定义的 Code。
Msg：强制写死为 "Internal Server Error."（不暴露具体细节）。
未知错误 (其他 Error 或 Panic)：
Code：强制写死为 99999。
Msg：强制写死为 "Internal Server Error."。



#### 请求mdc id获取拦截器：
依照MDC规范，我们接收请求header中的traceId和spanId。header中对应的key分别应该为X-B3-TraceId何X-B3-spanId（依照B3 Propagation 的协议）。

如果请求的header中带有这两个id，则返回体中放入。

两个id不要显式出现在我们自定义的实体类中，而是要在HttpRequest和HttpResponse中用header的kv来存取。

#### 请求mdc id生成拦截器
如果请求header中没有两个mdc id，则自己生成。这一步需要第一时间完成，这样后续的日志才能写入这两个id。


两个id不要显式出现在我们自定义的实体类中，而是要在HttpRequest和HttpResponse中用header的kv来存取。


### 4.4 日志相关
使用slog。日志中请打印traceid和spanid。日志的默认级别为INFO。

除了stdout以外，还要输出到文件。输出在log目录。全部日志输出在service.log中；日志级别不小于ERROR的输出在error.log。所有日志在每天0点0分切片，保存7天，更早的删除。

## 5. 环境相关
$GOROOT: /Users/dichters/.env/go
$GOPATH: /Users/dichters/.go

