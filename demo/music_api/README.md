# 利用第三方API实现Python加载器插件

### Python加载器插件V1
- 在[引擎托管平台](https://docs.iflyaicloud.com/aipaas-doc/)中，加载器插件连接了引擎和云原生组件，以实现云服务化的目的。用Python实现加载器插件，对使用者来说更易上手，能快速高效地实现服务请求

- 三方音乐识别API的文档[请参考](https://docs.acrcloud.cn/api-reference/identification-api/)

    - 音乐识别API可以接受用户的音频输入，文件格式为`mp3`，`wav`，`wma`，`amr`，`ogg`，`ape`，`acc`，`spx`，`m4a`，`mp4`，`FLAC`等
    
    - 音频可以为清晰的音乐或者是音乐哼唱，本例中分别使用了`music`和`humming`表示两种模式
    
    - 用户的请求数据是音频的二进制流格式，三方返回的识别结果是JSON对象，包含了`cost_time`、`status`、`metadata`三个字段，JSON对象的具体格式[请参考](https://docs.acrcloud.cn/metadata/music-broadcast)

    - 音乐识别可以应用于听歌识曲等

## :star:Python加载器插件V2

### 关于Python实现加载器的介绍可以参考[V1版本](#python加载器插件v1)，下面介绍了相比于上一版本，这一版本有哪些不同和改进。   

1. 类似于V1版本的Python加载器插件，实现的函数同样为`wrapperInit`、`wrapperFini`、`wrapperOnceExec`和`wrapperError`，不同的是，由于继承自`WrapperBase`，基类里说明了必须实现的接口，否则会出现`NotImplementedError`错误

2. 运行中用到的参数，V1版本是将变量声明为全局变量，在`wrapperInit`初始化后，其余函数体内将其声明为`global`；V2版本目前是将变量声明为类变量，实例变量同样可选

3. 需要注意的是， `wrapperOnceExec`函数执行返回的类型是`Response`对象，而不是前一版本表示错误码的`int`类型，意味着**无论结果正常与否**，均需实例化`Response`对象并返回
    
    - 未出现异常时，`Response`对象是是由一个或多个`ResponseData`对象构成的列表，其中`ResponseData`类有`key`、`data`、`len`、`type`和`status`五个成员变量

    - 出现异常时，直接调用`Response`对象的`response_err`方法返回错误码

4. 实现`Wrapper`类时，必须**继承**`WrapperBase`类，前三个成员函数的实现可以参考[V1版本实现](https://xfyun.github.io/athena_website/blog/music/api)
    
    ```python
         class Wrapper(WrapperBase):
            def wrapperInit(cls, config: {}) -> int:
               ...
            
            def wrapperFini(cls) -> int:
               ...
            
            def wrapperError(cls, ret: int) -> str:
               ...
            
            # 这里需要注意返回的类型是 Response 对象
            def wrapperOnceExec(self, params: {}, reqData: DataListCls) -> Response:
               res = Response()
               # 调用三方API的过程
               ...
               # 拿到返回的结果
               
               # 如果发生错误
               if error_occur:
                  return res.response_err(error_code)

               l = ResponseData()
               l.key = "output_text"
               l.type = 0
               l.status = 3
               l.data = r.text
               l.len = len(r.text.encode())

               res.list = [l]
               return res
     ```

5. 对于本地调试运行，需要注意下列几点

      - 额外声明用户请求和用户响应两个类

         ```python
         class UserRequest(object):
            '''
            定义请求类:
            params:  params 开头的属性代表最终HTTP协议中的功能参数parameters部分， 对应的是xtest.toml中的parameter字段
                     params Field支持 StringParamField，
                     NumberParamField，BooleanParamField,IntegerParamField，每个字段均支持枚举
                     params 属性多用于协议中的控制字段，请求body字段不属于params范畴

            input:    input字段多用与请求数据段，即body部分，当前支持 ImageBodyField、 StringBodyField和AudioBodyField
            '''
            params1 = StringParamField(key="mode", enums=["music", "humming"], value='humming')

            input1 = AudioBodyField(key="data", path="/home/wrapper/test.wav")
            
         class UserResponse(object):
            '''
            定义响应类:
            accepts:  accepts代表响应中包含哪些字段, 以及数据类型

            input:    input字段多用与请求数据段，即body部分，当前支持 ImageBodyField, StringBodyField, 和AudioBodyField
            '''
            accept1 = StringBodyField(key="ouput_text")
         ```
      - 在`Wrapper`类中声明一个类成员变量的字典类型config，模拟`wrapperInit`函数中传递参数，后期选择注释即可，在本例中如下
         ```python
         class Wrapper(WrapperBase):
            # 实例化用户请求类和用户响应类
            requestCls = UserRequest()
            responseCls = UserResponse()
            
            # 用于模拟aiges读入参数的字典
            config = {}
            config = {
            "requrl" : ...,
            "http_method" : ...,
            "http_uri" : ...,
            "access_key_music" : ...,
            "access_secret_music" : ...,
            "access_key_humming" : ...,
            "access_secret_humming" : ...
            }
         ```

      - 声明`main`函数，实例化`Wrapper`对象，运行程序
         ```python
            if __name__ == '__main__':
               m = Wrapper()
               m.schema()
               m.run()
         ```



<details>

<summary> 加载器插件的实现可以参考</summary>

- [v1版本加载器插件](https://iflytek.github.io/athena_website/blog/music/api/) 

- [v2版本加载器插件](https://iflytek.github.io/athena_website/docs/%E5%8A%A0%E8%BD%BD%E5%99%A8/Python%E6%8F%92%E4%BB%B6)
</details>
