namespace go MIT
namespace delphi MIT

struct Result {
    1: bool result
    2: i32 code;
    3: string data;
}

service Business {
    Result Action(1:string nFunName, 2:string nData);
}