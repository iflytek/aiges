// PythonParam.h
#ifndef __PYTHON_CPP_PYTHON_PARAM_H__
#define __PYTHON_CPP_PYTHON_PARAM_H__

#include <map>
#include <string>
#include <vector>

#include "include/python/Python.h"
std::string log_python_exception();
class PythonParamBuilder{
public:
    PythonParamBuilder();
    ~PythonParamBuilder();
    PyObject* Build();

    bool AddString( const std::string& );
    bool AddList( const std::vector<std::string>& params );
    bool AddMap( const std::map<std::string, std::string>& mapParam);
    template<typename T>
    bool AddCTypeParam( const std::string& typeStr, T v ){
		PyObject * param = Py_BuildValue( typeStr, v );
		return AddPyObject( param );
    }
private:
    bool AddPyObject( PyObject* obj );
    PyObject* ParseMapToPyDict(const std::map<std::string, std::string>& mapParam);
    PyObject* ParseListToPyList(const std::vector<std::string>& params);
    PyObject* ParseStringToPyStr( const std::string& str );
    
    PyObject* PasreBytesToPyBytes(char *value,unsigned int len);
private:
    PyObject* mArgs;
    std::vector<PyObject*> mTupleObjects;
};

#endif