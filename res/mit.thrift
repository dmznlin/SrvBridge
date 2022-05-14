namespace go mit
namespace delphi mit

struct ActionResult {
    1: required bool res; //执行结果
    2: optional i32 code = -1; //结果代码
    3: optional string data = ""; //结果数据
}

struct ActionParam {
    1: required string fname; //function name
    2: optional i32 ftype = 0; //function type
    3: optional string data = ""; //parameter's data
}

service Business {
    ActionResult Action(1:ActionParam param);
}