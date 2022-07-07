# New Design of wrapper.py

## 背景

1. 当前的wrapper.py 由[C项目](https://github.com/xfyun/aiges_c_python_wrapper)
   实现了 [wrapper接口](https://github.com/xfyun/aiges_c_python_wrapper/blob/master/include/aiges/wrapper.h)实现。
   aiges_c_python_wrapper编译成 libwrapper.so，由aiges统一加载。
2. 当前如果python用户需要实现推理插件， 只需要参考 [wrapper.py](https://github.com/xfyun/aiges_c_python_wrapper/blob/master/wrapper.py)
   实现对应接口后，即可实现python推理。
3. 当前用户实现 wrapper.py后， 无法直接调试运行，且不太了解 aiges如何调用 wrapper.py 以及传递到 wrapper.py
   对应的参数是什么类型都非常疑惑，造成python版本的AI推理插件集成方式并不那么pythonic。

## 新版wrapper.py集成方式优化目标

1. 用户可以定义AI能力输入的数据字段，控制字段列表
2. 用户可以按需定义AI能力输出的字段列表
3. 平台工具可以通过wrapper.py 自动导出用户schema并配置到webgate，对用户屏蔽schema概念
4. 平台工具可以提供用户直接Run wrapper.py ，并按照平台真实加载 wrapper.py方式传递对应参数，方便用户在任何环境快速Debug，发现一些基础问题。
5. 尽可能简化用户输入，并且在有限的用户输入下，获取平台需要的信息

## wrapper.py 新设计

![img_1.png](img_1.png)
1. 提供 python sdk:  python sdk将发布到 pypi，方便用户随时更新安装
2. [为什么?](###为什么) 新wrapper要求用户 实现 `Wrapper` 类，并将原有 函数式 wrapper开头的函数放入到 `Wrapper` （类方法|对象方法？待讨论 todo)中去
3. [新wrapper设计](https://github.com/xfyun/aiges_python/blob/master/aiges_python/v2/wrapper.py)，要求用户在Wrapper类中除了要实现 原有的 wrapperInit WrapperExec 等实现之外，需要额外定义能力的输入，输出，最终生成的HTTP接口基于此信息生成

### 为什么

我们希望用户只需要定义关键的实现，而不必care背后
wrapper.py如何被调用的细节，但是这块背后逻辑其实是复杂的，我们不希望在wrapper.py中让用户过多的定义一些因为平台要求而必须要的一些设置，我们希望在SDK的基类中实现定义好这些默认行为，
比如wrapper.py真实调用顺序 为 `WrapperInit -> WrapperExec -> WrapperFin`

在基类中定义这个行为的好处是， 用户继承基类并实现必要方法后，可以直接 Run运行，并且调试拿到结果。

至于为什么希望用户在 Wrapper类中实现 对应方法，原因也是可以在基类行为中做一些 更Pythonic的魔法，简化用户的输入。

