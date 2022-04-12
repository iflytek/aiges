#include "pyCall.h"


char *callWrapperCreate(const char *usrTag, pParamList params,unsigned int psrIds[], int psrCnt, int *errNum, std::string sid)
{
    char *handle;

    PyGILState_STATE gstate = PyGILState_Ensure();
    PyObject *wrapperModule = PyImport_ImportModule(_wrapperName);
    PyObject *createFunc = PyObject_GetAttrString(wrapperModule, "wrapperCreate");
    Py_XDECREF(wrapperModule);
    if (!createFunc || !PyCallable_Check(createFunc))
    {
        PyGILState_Release(gstate);
        *errNum = WRAPPER::CError::NotImplementCreate;
        return NULL;
    }
    PyObject *pArgsT = PyTuple_New(5);
    std::vector<PyObject *> tmpPyObjectVec;
    try
    {
        //构建参数元组
        //构建请求句柄
        PyObject *pUsrTag = PyUnicode_FromString(usrTag);
        PyTuple_SetItem(pArgsT, 0, pUsrTag);

        //构建请求参数
        PyObject *pyParam = PyDict_New();
        for (pParamList p = params; p != NULL; p = p->next)
        {
            PyObject *tmpV = Py_BuildValue("s", p->value);
            tmpPyObjectVec.push_back(tmpV);
            PyDict_SetItemString(pyParam, p->key, tmpV);
            spdlog::debug("wrapper create param, key:{},value:{},sid:{}", p->key, p->value, sid);
        }
        PyTuple_SetItem(pArgsT, 1, pyParam);

        //构建个性化请求id
        int num = psrCnt;
        PyObject *pyPsrIds = PyTuple_New(num);
        if (num != 0)
        {
            for (int idx = 0; idx < num; idx++)
            {
                PyTuple_SetItem(pyPsrIds, idx, Py_BuildValue("i", psrIds[idx]));
                spdlog::debug("wrapper create psrId:{},sid:{}", psrIds[idx], sid);
            }
            PyTuple_SetItem(pArgsT, 2, pyPsrIds);
        }
        else
        {
            spdlog::debug("wrapper exec psrIds is empty.sid:{}", sid);
            PyTuple_SetItem(pArgsT, 2, pyPsrIds);
        }
        // //构建个性化请求个数
        PyTuple_SetItem(pArgsT, 3, Py_BuildValue("i", psrCnt));

        PyObject *pyErrNum = Py_BuildValue("i", 0);
        PyTuple_SetItem(pArgsT, 4, pyErrNum);

        spdlog::debug("wrapper create psrCnt .val :{},sid:{}", psrCnt, sid);
        PyObject *pRet = PyEval_CallObject(createFunc, pArgsT);
        if (pRet == NULL)
        {
            std::string errRlt = "";
            errRlt = log_python_exception();
            if (errRlt != "")
            {
                spdlog::error("wrapperCreate error:{},sid:{}", errRlt, sid);
            }
            Py_DECREF(pArgsT);
            Py_DECREF(createFunc);
            *errNum = WRAPPER::CError::innerError;
        }
        else
        {
            char *retHandle;
            int ret = 0;
            PyArg_Parse(pyErrNum, "i", &ret);
            if (ret == 0)
            {
                PyArg_Parse(pRet, "s", &retHandle);
                Py_DECREF(pRet);
                std::string tmHdlStr = retHandle;
                handle = (char *)malloc(strlen(tmHdlStr.c_str()));
                memcpy(handle, (char *)tmHdlStr.c_str(), strlen(tmHdlStr.c_str()));
                spdlog::debug("wrapper create handle .val :{},sid:{}", tmHdlStr, sid);
            }
            else
            {
                *errNum = ret;
                spdlog::error("wrapper create handle error .code :{},sid:{}", *errNum, sid);
            }
        }
    }
    catch (const std::exception &e)
    {
        std::string errRlt = "";
        errRlt = log_python_exception();
        if (errRlt != "")
        {
            spdlog::error("wrapperCreate error:{}, ret:{},sid:{}", errRlt, *errNum, sid);
        }
        *errNum = WRAPPER::CError::innerError;
    }
    spdlog::debug("wrapperCreate ret.{},tmpPyObjectVec size:{},sid:{}", *errNum, tmpPyObjectVec.size(), sid);
    for (auto &i : tmpPyObjectVec)
    {
        Py_XDECREF(i);
    }
    //Py_XDECREF(pyData);
    Py_DECREF(pArgsT);
    Py_DECREF(createFunc);
    PyGILState_Release(gstate);
    if (*errNum != 0)
    {
        return NULL;
    }
    else
    {
        return handle;
    }
}

int callWrapperWrite(char *handle, pDataList reqData,std::string sid)
{
    int ret = 0;
    PyGILState_STATE gstate = PyGILState_Ensure();
    PyObject *wrapperModule = PyImport_ImportModule(_wrapperName);
    PyObject *writeFunc = PyObject_GetAttrString(wrapperModule, "wrapperWrite");
    Py_XDECREF(wrapperModule);
    if (!writeFunc || !PyCallable_Check(writeFunc))
    {
        PyGILState_Release(gstate);
        ret = WRAPPER::CError::NotImplementWrite;
        return ret;
    }

    PyObject *pArgsT = PyTuple_New(2);
    std::vector<PyObject *> tmpPyObjectVec;
    try
    {
        //构建参数元组
        //构建请求句柄
        PyObject *pUsrTag = PyUnicode_FromString(handle);
        PyTuple_SetItem(pArgsT, 0, pUsrTag);

        //构建请求数据
        int dataNum = 0;
        for (pDataList tmpDataPtr = reqData; tmpDataPtr != NULL; tmpDataPtr = tmpDataPtr->next)
        {
            dataNum++;
        }
        spdlog::debug("call wrapper write,datanum:{}，sid:{}", dataNum, sid);
        PyObject *pyDataList = PyTuple_New(dataNum);
        if (dataNum > 0)
        {
            pDataList p = reqData;
            for (int tmpIdx = 0; tmpIdx < dataNum; tmpIdx++)
            {
                PyObject *tmp = PyDict_New();

                PyObject *pyKey = Py_BuildValue("s", p->key);
                PyDict_SetItemString(tmp, "key", pyKey);
                tmpPyObjectVec.push_back(pyKey);
                //std::string datas(static_cast<char*>(p->data),p->len);
                //PyObject *pyData = Py_BuildValue("O", p->data);
                PyObject *pyData = PyBytes_FromStringAndSize((char *)(p->data), p->len);
                PyDict_SetItemString(tmp, "data", pyData);
                tmpPyObjectVec.push_back(pyData);

                PyObject *pyDataLen = Py_BuildValue("i", int(p->len));
                PyDict_SetItemString(tmp, "len", pyDataLen);
                tmpPyObjectVec.push_back(pyDataLen);

                PyObject *pyStatus = Py_BuildValue("i", int(p->status));
                PyDict_SetItemString(tmp, "status", pyStatus);
                tmpPyObjectVec.push_back(pyStatus);

                PyObject *pyType = Py_BuildValue("i", int(p->type));
                PyDict_SetItemString(tmp, "type", pyType);

                tmpPyObjectVec.push_back(pyType);

                PyObject *tmpDesc = PyDict_New();
                tmpPyObjectVec.push_back(tmpDesc);
                for (pParamList descP = p->desc; descP != NULL; descP = descP->next)
                {
                    PyObject *tmpV = Py_BuildValue("s", descP->value);
                    tmpPyObjectVec.push_back(tmpV);
                    PyDict_SetItemString(tmpDesc, descP->key, tmpV);
                }
                PyDict_SetItemString(tmp, "desc", tmpDesc);

                PyTuple_SetItem(pyDataList, tmpIdx, tmp);
                p = p->next;
            }
        }
        PyTuple_SetItem(pArgsT, 1, pyDataList);

        PyObject *pRet = PyEval_CallObject(writeFunc, pArgsT);
        if (pRet == NULL)
        {
            std::string errRlt = "";
            errRlt = log_python_exception();
            if (errRlt != "")
            {
                spdlog::error("wrapperExec error:{},sid:{}", errRlt, sid);
            }
            Py_DECREF(pArgsT);
            Py_DECREF(writeFunc);
            ret = WRAPPER::CError::innerError;
        }
        else
        {
            PyArg_Parse(pRet, "i", &ret);
            Py_DECREF(pRet);
        }
    }
    catch (const std::exception &e)
    {
        std::string errRlt = "";
        errRlt = log_python_exception();
        if (errRlt != "")
        {
            spdlog::error("wrapperExec error:{}, ret:{},sid:{}", errRlt, ret, sid);
        }
        ret = WRAPPER::CError::innerError;
    }
    spdlog::debug("wrapperwite ret.{},tmpPyObjectVec size:{},sid:{}", ret, tmpPyObjectVec.size(), sid);
    for (auto &i : tmpPyObjectVec)
    {
        Py_XDECREF(i);
    }
    //Py_XDECREF(pyData);
    Py_DECREF(pArgsT);
    Py_DECREF(writeFunc);
    PyGILState_Release(gstate);
    return ret;
}

int callWrapperRead(char *handle, pDataList *respData,std::string sid)
{
    int ret = 0;
    PyGILState_STATE gstate = PyGILState_Ensure();
    PyObject *wrapperModule = PyImport_ImportModule(_wrapperName);
    PyObject *readFunc = PyObject_GetAttrString(wrapperModule, "wrapperRead");
    Py_XDECREF(wrapperModule);
    if (!readFunc || !PyCallable_Check(readFunc))
    {
        PyGILState_Release(gstate);
        ret = WRAPPER::CError::NotImplementRead;
        return ret;
    }

    //构建参数元组
    PyObject *pArgsT = PyTuple_New(2);
    std::vector<PyObject *> tmpPyObjectVec;
    try
    {
        //构建请求句柄
        PyObject *pUsrTag = PyUnicode_FromString(handle);
        PyTuple_SetItem(pArgsT, 0, pUsrTag);

        //构建响应数据体
        PyObject *pyRespData = PyList_New(0);
        PyTuple_SetItem(pArgsT, 1, pyRespData);

        PyObject *pRet = PyEval_CallObject(readFunc, pArgsT);
        if (pRet == NULL)
        {
            std::string errRlt = "";
            errRlt = log_python_exception();
            if (errRlt != "")
            {
                spdlog::error("wrappeRead error:{},sid:{}", errRlt, sid);
            }
            Py_DECREF(pArgsT);
            Py_DECREF(readFunc);
            ret = WRAPPER::CError::innerError;
        }
        else
        {
            PyArg_Parse(pRet, "i", &ret);
            Py_DECREF(pRet);
            if (ret == 0)
            {
                //读取响应
                int rltSize = PyList_Size(pyRespData);
                if (rltSize != 0)
                {
                    pDataList headPtr;
                    pDataList prePtr;
                    pDataList curPtr;
                    for (int idx = 0; idx < rltSize; idx++)
                    {
                        pDataList tmpData = new (DataList);

                        PyObject *tmpDict = PyList_GetItem(pyRespData, idx);
                        char *tmpRltKey = pyDictStrToChar(tmpDict, DATA_KEY, sid, ret);
                        if (ret != 0)
                        {
                            break;
                        }
                        else
                        {
                            tmpData->key = tmpRltKey;
                        }

                        int integerVal = 0;
                        ret = pyDictIntToInt(tmpDict, DATA_LEN, integerVal, sid);
                        if (ret != 0)
                        {
                            break;
                        }
                        else
                        {
                            tmpData->len = integerVal;
                        }

                        char *tmpRltData = pyDictStrToChar(tmpDict, DATA_DATA, sid, ret);
                        if (ret != 0)
                        {
                            break;
                        }
                        else
                        {
                            tmpData->data = (void *)tmpRltData;
                        }
                        int interValSta = 0;
                        ret = pyDictIntToInt(tmpDict, DATA_STATUS, interValSta, sid);
                        if (ret != 0)
                        {
                            break;
                        }
                        else
                        {
                            tmpData->status = DataStatus(interValSta);
                        }
                        int interValType = 0;
                        ret = pyDictIntToInt(tmpDict, DATA_TYPE, interValType, sid);
                        if (ret != 0)
                        {
                            break;
                        }
                        else
                        {
                            tmpData->type = DataType(interValType);
                        }
                        tmpData->next = NULL;
                        //检查下是否需要desc吧
                        tmpData->desc = pyDictToDesc(tmpDict, DATA_DESC, sid, ret);
                        if (ret != 0)
                        {
                            break;
                        }
                        if (idx == 0)
                        {
                            headPtr = tmpData;
                            prePtr = tmpData;
                            curPtr = tmpData;
                        }
                        else
                        {
                            curPtr = tmpData;
                            prePtr->next = curPtr;
                            prePtr = curPtr;
                        }
                        spdlog::debug("get result,key:{},data:{},len:{},type:{},status:{},sid:{}", 
                        tmpData->key, (char *)tmpData->data,tmpData->len, tmpData->type, tmpData->status, sid);
                    }
                    *respData = headPtr;
                }
            }
        }
    }
    catch (const std::exception &e)
    {
        std::string errRlt = "";
        errRlt = log_python_exception();
        if (errRlt != "")
        {
            spdlog::error("wrapperRead error:{}, ret:{},sid:{}", errRlt, ret, sid);
        }
        ret = WRAPPER::CError::innerError;
    }

    for (auto &i : tmpPyObjectVec)
    {
        Py_XDECREF(i);
    }
    //Py_XDECREF(pyData);
    Py_DECREF(pArgsT);
    Py_DECREF(readFunc);
    PyGILState_Release(gstate);
    return ret;
}

int callWrapperDestroy(char *handle,std::string sid)
{
    int ret = 0;
    PyGILState_STATE gstate = PyGILState_Ensure();
    PyObject *wrapperModule = PyImport_ImportModule(_wrapperName);
    PyObject *destoryFunc = PyObject_GetAttrString(wrapperModule, "wrapperDestory");
    Py_XDECREF(wrapperModule);
    if (!destoryFunc || !PyCallable_Check(destoryFunc))
    {
        PyGILState_Release(gstate);
        ret = WRAPPER::CError::NotImplementDestory;
        return ret;
    }

    //构建参数元组
    PyObject *pArgsT = PyTuple_New(1);
    try
    {
        //构建请求句柄
        PyObject *pUsrTag = PyUnicode_FromString(handle);
        PyTuple_SetItem(pArgsT, 0, pUsrTag);

        PyObject *pRet = PyEval_CallObject(destoryFunc, pArgsT);
        if (pRet == NULL)
        {
            std::string errRlt = "";
            errRlt = log_python_exception();
            if (errRlt != "")
            {
                spdlog::error("wrappeRead error:{},sid:{}", errRlt, sid);
            }
            Py_DECREF(pArgsT);
            Py_DECREF(destoryFunc);
            ret = WRAPPER::CError::innerError;
        }
        else
        {
            PyArg_Parse(pRet, "i", &ret);
            Py_DECREF(pRet);
        }
    }
    catch (const std::exception &e)
    {
        std::string errRlt = "";
        errRlt = log_python_exception();
        if (errRlt != "")
        {
            spdlog::error("wrapperRead error:{}, ret:{},sid:{}", errRlt, ret, sid);
        }
        ret = WRAPPER::CError::innerError;
    }
    //Py_XDECREF(pyData);
    Py_DECREF(pArgsT);
    Py_DECREF(destoryFunc);
    PyGILState_Release(gstate);
    return ret;
}
