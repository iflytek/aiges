### 利用第三方API实现Python加载器插件
- 在[引擎托管平台](https://docs.iflyaicloud.com/aipaas-doc/)中，加载器插件连接了引擎和云原生组件，以实现云服务化的目的。用Python实现加载器插件，对使用者来说更易上手，能快速高效地实现服务请求。

- 三方音乐识别API的文档请[参考](https://docs.acrcloud.cn/api-reference/identification-api/)
    
    - 音乐识别API可以接受用户的音频输入，文件格式为`mp3`，`wav`，`wma`，`amr`，`ogg`，`ape`，`acc`，`spx`，`m4a`，`mp4`，`FLAC`等
    
    - 音频可以为清晰的音乐或者是音乐哼唱
    
    - 用户的请求数据是音频的二进制流格式，三方返回的识别结果是JSON对象，包含了`cost_time`、`status`、`metadata`三个字段，JSON对象的具体格式[请参考](https://docs.acrcloud.cn/metadata/music-broadcast)

- 加载器插件的实现介绍可以[参考](https://xfyun.github.io/athena_website/blog/music/api/)