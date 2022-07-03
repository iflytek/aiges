#!/bin/bash
# Copyright © 2021 iflytek.com 
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "### Usage" >> Note.md
echo "
***当前支持的cuda-go-python基础镜像列表(包含cuda go python编译环境)***
#### Public 公网可拉取
***当前支持的python cuda golang基础编译镜像***

| repo                                                                                     | tag                         | python | cuda | os           |
|------------------------------------------------------------------------------------------|-----------------------------|--------|------|--------------|
| public.ecr.aws/iflytek-open/cuda-go-python-base:10.1-1.17-3.9.13-ubuntu1804 | 10.1-1.17-3.9.13-ubuntu1804 | 3.9.13 | 10.1 | ubuntu 18.04 |
| public.ecr.aws/iflytek-open/cuda-go-python-base:10.2-1.17-3.9.13-ubuntu1804 | 10.2-1.17-3.9.13-ubuntu1804 | 3.9.13 | 10.2 | ubuntu 18.04 |
| public.ecr.aws/iflytek-open/cuda-go-python-base:11.2-1.17-3.9.13-ubuntu1804 | 10.1-1.17-3.9.13-ubuntu1804 | 3.9.13 | 11.2 | ubuntu 18.04 |
| public.ecr.aws/iflytek-open/cuda-go-python-base:11.6-1.17-3.9.13-ubuntu1804 | 10.2-1.17-3.9.13-ubuntu1804 | 3.9.13 | 11.6 | ubuntu 18.04 |

***当前支持的aiges基础镜像列表***

| repo                                                                      | tag                    | python | cuda | os           |
|---------------------------------------------------------------------------|------------------------|--------|------|--------------|
| public.ecr.aws/iflytek-open/aiges-gpu:10.1-3.9.13-ubuntu1804 | 10.1-3.9.13-ubuntu1804 | 3.9.13 | 10.1 | ubuntu 18.04 |
| public.ecr.aws/iflytek-open/aiges-gpu:10.2-3.9.13-ubuntu1804 | 10.2-3.9.13-ubuntu1804 | 3.9.13 | 10.2 | ubuntu 18.04 |
| public.ecr.aws/iflytek-open/aiges-gpu:11.2-1.17-3.9.13-ubuntu1804 | 11.2-1.17-3.9.13-ubuntu1804 | 3.9.13 | 11.2 | ubuntu 18.04 |
| public.ecr.aws/iflytek-open/aiges-gpu:11.6-1.17-3.9.13-ubuntu1804 | 11.6-1.17-3.9.13-ubuntu1804 | 3.9.13 | 11.6 | ubuntu 18.04 |



#### 内部仓库

***当前支持的python cuda golang基础编译镜像***

| repo                                                                                     | tag                         | python | cuda | os           |
|------------------------------------------------------------------------------------------|-----------------------------|--------|------|--------------|
| artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.1-1.17-3.9.13-ubuntu1804 | 10.1-1.17-3.9.13-ubuntu1804 | 3.9.13 | 10.1 | ubuntu 18.04 |
| artifacts.iflytek.com/docker-private/atp/cuda-go-python-base:10.2-1.17-3.9.13-ubuntu1804 | 10.2-1.17-3.9.13-ubuntu1804 | 3.9.13 | 10.2 | ubuntu 18.04 |

***当前支持的aiges基础镜像列表***

| repo                                                                      | tag                    | python | cuda | os           |
|---------------------------------------------------------------------------|------------------------|--------|------|--------------|
| artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.1-3.9.13-ubuntu1804 | 10.1-3.9.13-ubuntu1804 | 3.9.13 | 10.1 | ubuntu 18.04 |
| artifacts.iflytek.com/docker-private/atp/aiges-gpu:10.2-3.9.13-ubuntu1804 | 10.2-3.9.13-ubuntu1804 | 3.9.13 | 10.2 | ubuntu 18.04 |

" >> Note.md
