#!/usr/bin/env python
class EngineBase:   # clase EngineBase(object):  c call exception

    '''
    aiges engine wrapper
    '''
    version = '1.0.1'

    def wrapperInit(self, kwargs):  # **kwargs: c map convert to python dict? how??
        print(kwargs)
        return 0

    def wrapperFini(self):
        return 0

    def wrapperError(self, code):
        return "unknown err code"

    def wrapperDebugInfo(self, inst):
        return "unknown debug info"

    def wrapperExec(self, kwargs, msg):
        print(kwargs)
	print(msg)
        return 10010


if __name__ == "__main__":
    engine = EngineBase()
    print(engine.version)
