# mmocr

# demo

### v1流程
1. 在https://www.iflyaicloud.com/新建能力

2. 编写`wrapper.py`文件，利用`AIservice`以及`xtest`工具完成验证
    - 插件的实现参考`mmocr/wrapper_v1.py`

3. 编写requirements.txt

4. 发布能力

## v2流程，基于`envd`构建镜像
1. 安装`envd`
    ```bash
    pip install --pre --upgrade envd
    envd bootstrap
   ```
2. 通过`envd`进行镜像构建
    ```python
    envd build --t IMAGE:TAG --f build.envd         
    ```

3. 快速开始一个项目
    ```python
    python -m aiges -n "project name"
    ```
4.  完成`wrapper.py`的编写以及本地调试，并发布能力
    - 插件的实现参考`mmocr/wrapper_v1.py`


   
