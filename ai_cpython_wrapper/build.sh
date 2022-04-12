#!/bin/bash
#export LD_LIBRARY_PATH=./lib:$LD_LIBRARY_PATH
echo $LD_LIBRARY_PATH
g++ -fPIC -shared -std=c++11 -Wno-attributes -g -O0 -I. -I./include/spdlog/include -I./include/ -o libpyCallCommon.so pyCall_common.cpp -L. -lpython
g++ -fPIC -shared -std=c++11 -Wno-attributes -g -O0 -I. -I./include/spdlog/include -I./include/ -o libpyCallOnce.so pyCall_once.cpp -L. -lpython
g++ -fPIC -shared -std=c++11 -Wno-attributes -g -O0 -I. -I./include/spdlog/include -I./include/ -o libpyCallStream.so pyCall_stream.cpp -L. -lpython
g++ -fPIC -shared -std=c++11 -Wno-attributes -g -O0 -I. -I./include/spdlog/include -I./include/ -L. -lpyCallCommon -lpyCallOnce -lpyCallStream -lboost_filesystem -lboost_system -o libwrapper.so wrapper.cpp
cp libwrapper.so libpyCallCommon.so libpyCallOnce.so libpyCallStream.so libpython.so ./wrapperlib/