thrift -gen -r go mit.thrift
:: 生成go接口

xcopy /y /e /i .\gen-go\MIT ..\src\MIT
:: 移动至源码目录

thrift -gen -r delphi mit.thrift
:: 生成delphi接口
