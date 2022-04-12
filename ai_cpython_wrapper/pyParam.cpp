//PythonParam.cpp
#include "PythonParam.h"


PythonParamBuilder::PythonParamBuilder()
:mArgs(nullptr){

}

PythonParamBuilder::~PythonParamBuilder(){
    for ( std::vector<PyObject*>::iterator it = mTupleObjects.begin();
        it != mTupleObjects.end();
	 ++it ){
            if( *it )
		Py_DECREF( *it );
	}

    if ( mArgs )
        Py_DECREF(mArgs);
}

PyObject* PythonParamBuilder::Build(){
    if ( mTupleObjects.empty() )
        return nullptr;

    mArgs = PyTuple_New( mTupleObjects.size() );

    for ( int i=0; i<mTupleObjects.size(); i++){
         PyTuple_SetItem(mArgs, i, mTupleObjects[i] );
    }

    return mArgs;
}

bool PythonParamBuilder::AddPyObject( PyObject* obj ){
     if ( nullptr == obj )
        return false;
     mTupleObjects.push_back(obj);
     return true;
}

bool PythonParamBuilder::AddString( const std::string& str ){
     return AddPyObject( ParseStringToPyStr( str ) );
}

bool PythonParamBuilder::AddList( const std::vector<std::string>& params ){
    return AddPyObject( ParseListToPyList(params) );
}

bool PythonParamBuilder::AddMap( const std::map<std::string, std::string>& mapParam){
    return AddPyObject( ParseMapToPyDict( mapParam ) );
}

PyObject* PythonParamBuilder::ParseMapToPyDict(const std::map<std::string, std::string>& mapParam){
    PyObject* pyDict = PyDict_New();
    for (std::map<std::string, std::string>::const_iterator it = mapParam.begin();
          it != mapParam.end();
          ++it) {
              PyObject* pyValue = PyString_FromString(it->second.c_str());
              if (pyValue == NULL) {
                   printf( "Parse param:[%s] to PyStringObject failed.", it->second.c_str());
		   return NULL;
	      }
	      if (PyDict_SetItemString(pyDict, it->first.c_str(), pyValue) < 0) {
		   printf( "Parse key:[%s] value:[%s] failed.", it->first.c_str(), it->second.c_str());
		   return NULL;
	      }
    } 
    return pyDict;
}

PyObject* PythonParamBuilder::ParseListToPyList(const std::vector<std::string>& params){
    size_t size = params.size();
    PyObject* paramList =  PyList_New(size);
    for (size_t i = 0; i < size; i++) {
	PyObject* pyStr = PyString_FromStringAndSize(params[i].data(), params[i].size());
	if (pyStr == NULL) {
  	    printf( "Parse param:[%s] to PyStringObject failed.", params[i].c_str());
	    break;
	}
	if (PyList_SetItem(paramList, i, pyStr) != 0) {
	    printf( "param:[%s] append to PyParamList failed.", params[i].c_str());
	    break;
	}
    }
    return paramList;
}

PyObject* PythonParamBuilder::ParseStringToPyStr( const std::string& str ){
    return PyString_FromStringAndSize( str.data(), str.size() );
}