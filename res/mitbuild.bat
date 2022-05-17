thrift -r --gen go mit.thrift
:: 生成go接口

xcopy /y /e /i .\gen-go\mit ..\src\mit
:: 移动至源码目录

thrift -r --gen delphi:async mit.thrift
:: 生成delphi接口,包括异步调用(双向通信)