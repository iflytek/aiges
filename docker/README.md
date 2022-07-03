# docker 镜像build工具

## 依赖

```
pip install jinja2
pip install plumbum  
```

## 生成 aiges Dockerfile

***当前python, golang,distro版本固定为: 3.9.13, 1.17, ubuntu1804***

Github构建:
```
python docker/build.py generate  --all --use_github  
```

内部构建:

```
python docker/build.py generate  --all  

```





 