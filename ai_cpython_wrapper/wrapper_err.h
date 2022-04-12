#ifndef _WRAPPER_ERROR_H
#define _WRAPPER_ERROR_H

namespace WRAPPER
{
    enum CError
    {
        innerError = -1,
        NotImplementError = -1000,
        NotImplementInit = -1001,
        NotImplementExec = -1002,
        NotImplementCreate = -1003,
        NotImplementWrite = -1004,
        NotImplementRead = -1005,
        NotImplementDestory = -1006,
        NotImplementFini = -1007,
        RltDataKeyInvalid = -1010,
        RltDataDataInvalid = -1011,
        RltDataLenInvalid = -1012,
        RltDataStatusInvalid = -1013,
        RltDataTypeInvalid = -1014,
    };

};
#endif