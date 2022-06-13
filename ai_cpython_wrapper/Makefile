PYLIB=$(shell python3-config  --embed --ldflags)
PYINC=$(shell python3-config --cflags --embed)
wrapper:
	echo $(PYLIB)
	g++ -fPIC -shared -std=c++11 -Wno-attributes  $(PYINC)  -g -O0 -I. -I./include/spdlog/include -I./include/ -o libpyCallCommon.so pyCall_common.cpp -L. $(PYLIB)
	g++ -fPIC -shared -std=c++11 -Wno-attributes  $(PYINC) -g -O0 -I. -I./include/spdlog/include -I./include/ -o libpyCallOnce.so pyCall_once.cpp -L. -L /opt/conda/envs/loader/lib  $(PYLIB)
	g++ -fPIC -shared -std=c++11 -Wno-attributes  $(PYINC) -g -O0 -I. -I./include/spdlog/include -I./include/ -o libpyCallStream.so pyCall_stream.cpp -L. -L /opt/conda/envs/loader/lib   $(PYLIB)
	g++  -Wall -pedantic -Wextra -fPIC -shared -std=c++11 -Wno-attributes  $(PYINC) -g -O0 -I. -I./include/spdlog/include -I./include/ -o libwrapper.so wrapper.cpp -L. -L /opt/conda/envs/loader/lib -lboost_filesystem -lboost_system  -lpyCallCommon -lpyCallOnce -lpyCallStream $(PYLIB)
	mkdir -p wrapper_lib
	cp libwrapper.so libpyCallCommon.so libpyCallOnce.so libpyCallStream.so ./wrapper_lib
