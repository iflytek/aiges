#include "pyCall.h"


int callWrapperExec(const char *usrTag, pParamList params, pDataList reqData, pDataList *respData, unsigned int psrIds[], int psrCnt, std::string sid)
{
    int ret = 0;
    PyGILState_STATE gstate = PyGILState_Ensure();
    PyObject *wrapperModule = PyImport_ImportModule(_wrapperName);
    PyObject *execFunc = PyObject_GetAttrString(wrapperModule, "wrapperOnceExec");
    Py_XDECREF(wrapperModule);
    if (!execFunc || !PyCallable_Check(execFunc))
    {
        PyGILState_Release(gstate);
        return WRAPPER::CError::NotImplementExec;
    }
    PyObject *pArgsT = PyTuple_New(6);
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
            if(NULL==p->key){
                continue;
            }
            PyObject *tmpV = Py_BuildValue("s", p->value);
            tmpPyObjectVec.push_back(tmpV);
            PyDict_SetItemString(pyParam, p->key, tmpV);
            spdlog::debug("wrapper exec param, key:{},value:{},sid:{}", p->key, p->value, sid);
        }
        PyTuple_SetItem(pArgsT, 1, pyParam);
        //构建请求数据
        int dataNum = 0;
        for (pDataList tmpDataPtr = reqData; tmpDataPtr != NULL; tmpDataPtr = tmpDataPtr->next)
        {
            dataNum++;
        }
        spdlog::debug("call wrapper exec,datanum:{}，sid:{}", dataNum, sid);
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
                    if(NULL==descP->key){
                        continue;
                    }   
                    PyObject *tmpV = Py_BuildValue("s", descP->value);
                    tmpPyObjectVec.push_back(tmpV);
                    PyDict_SetItemString(tmpDesc, descP->key, tmpV);
                }
                PyDict_SetItemString(tmp, "desc", tmpDesc);

                PyTuple_SetItem(pyDataList, tmpIdx, tmp);
                p = p->next;
            }
        }
        PyTuple_SetItem(pArgsT, 2, pyDataList);
        //构建响应数据体
        PyObject *pyRespData = PyList_New(0);
        PyTuple_SetItem(pArgsT, 3, pyRespData);
        //构建个性化请求id
        int num = psrCnt;
        PyObject *pyPsrIds = PyTuple_New(num);
        if (num != 0)
        {
            for (int idx = 0; idx < num; idx++)
            {
                PyTuple_SetItem(pyPsrIds, idx, Py_BuildValue("i", psrIds[idx]));
                spdlog::debug("wrapper exec psrId:{},sid:{}", psrIds[idx], sid);
            }
            PyTuple_SetItem(pArgsT, 4, pyPsrIds);
        }
        else
        {
            spdlog::debug("wrapper exec psrIds is empty.sid:{}", sid);
            PyTuple_SetItem(pArgsT, 4, pyPsrIds);
        }
        // //构建个性化请求个数
        PyTuple_SetItem(pArgsT, 5, Py_BuildValue("i", psrCnt));
        spdlog::debug("wrapper exec psrCnt .val :{},sid:{}", psrCnt, sid);
        PyObject *pRet = PyEval_CallObject(execFunc, pArgsT);
        if (pRet == NULL)
        {
            std::string errRlt = "";
            errRlt = log_python_exception();
            if (errRlt != "")
            {
                spdlog::error("wrapperExec error:{},sid:{}", errRlt, sid);
            }
            Py_DECREF(pArgsT);
            Py_DECREF(execFunc);
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
            spdlog::error("wrapperExec error:{}, ret:{},sid:{}", errRlt, ret, sid);
        }
        ret = WRAPPER::CError::innerError;
    }
    spdlog::debug("wrapperExec ret.{},tmpPyObjectVec size:{},sid:{}", ret, tmpPyObjectVec.size(), sid);
    for (auto &i : tmpPyObjectVec)
    {
        Py_XDECREF(i);
    }
    //Py_XDECREF(pyData);
    Py_DECREF(pArgsT);
    Py_DECREF(execFunc);
    PyGILState_Release(gstate);
    return ret;
}

