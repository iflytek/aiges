g++ -fPIC -shared -Wno-attributes -std=c++0x -g -O0 -rdynamic -o libwrapper.so wrapper.cpp
g++ -fPIC -shared -Wno-attributes -std=c++0x -g -O0 -rdynamic -o libwrapper-catch.so wrapper_catch.cpp
