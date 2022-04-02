gcc -fPIC -shared -Wno-attributes -g -O0 -rdynamic -o libwrapper.so wrapper.cpp
gcc -fPIC -shared -Wno-attributes -std=c++0x -g -O0 -rdynamic -o libwrapper-catch.so wrapper_catch.cpp
