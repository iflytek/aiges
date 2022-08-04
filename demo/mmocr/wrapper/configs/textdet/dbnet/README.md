# DBNet

> [Real-time Scene Text Detection with Differentiable Binarization](https://arxiv.org/abs/1911.08947)

<!-- [ALGORITHM] -->
## Abstract

Recently, segmentation-based methods are quite popular in scene text detection, as the segmentation results can more accurately describe scene text of various shapes such as curve text. However, the post-processing of binarization is essential for segmentation-based detection, which converts probability maps produced by a segmentation method into bounding boxes/regions of text. In this paper, we propose a module named Differentiable Binarization (DB), which can perform the binarization process in a segmentation network. Optimized along with a DB module, a segmentation network can adaptively set the thresholds for binarization, which not only simplifies the post-processing but also enhances the performance of text detection. Based on a simple segmentation network, we validate the performance improvements of DB on five benchmark datasets, which consistently achieves state-of-the-art results, in terms of both detection accuracy and speed. In particular, with a light-weight backbone, the performance improvements by DB are significant so that we can look for an ideal tradeoff between detection accuracy and efficiency. Specifically, with a backbone of ResNet-18, our detector achieves an F-measure of 82.8, running at 62 FPS, on the MSRA-TD500 dataset.

<div align=center>
<img src="https://user-images.githubusercontent.com/22607038/142791306-0da6db2a-20a6-4a68-b228-64ff275f67b3.png"/>
</div>

## Results and models

### ICDAR2015

|                                    Method                                     |                                                      Pretrained Model                                                      |  Training set   |    Test set    | #epochs | Test size | Recall | Precision | Hmean |                                                                                                                         Download                                                                                                                          |
| :---------------------------------------------------------------------------: | :------------------------------------------------------------------------------------------------------------------------: | :-------------: | :------------: | :-----: | :-------: | :----: | :-------: | :---: | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|     [DBNet_r18](/configs/textdet/dbnet/dbnet_r18_fpnc_1200e_icdar2015.py)     |                                                          ImageNet                                                          | ICDAR2015 Train | ICDAR2015 Test |  1200   |    736    | 0.731  |   0.871   | 0.795 |      [model](https://download.openmmlab.com/mmocr/textdet/dbnet/dbnet_r18_fpnc_sbn_1200e_icdar2015_20210329-ba3ab597.pth) \| [log](https://download.openmmlab.com/mmocr/textdet/dbnet/dbnet_r18_fpnc_sbn_1200e_icdar2015_20210329-ba3ab597.log.json)      |
| [DBNet_r50dcn](/configs/textdet/dbnet/dbnet_r50dcnv2_fpnc_1200e_icdar2015.py) | [Synthtext](https://download.openmmlab.com/mmocr/textdet/dbnet/dbnet_r50dcnv2_fpnc_sbn_2e_synthtext_20210325-aa96e477.pth) | ICDAR2015 Train | ICDAR2015 Test |  1200   |   1024    | 0.814  |   0.868   | 0.840 | [model](https://download.openmmlab.com/mmocr/textdet/dbnet/dbnet_r50dcnv2_fpnc_sbn_1200e_icdar2015_20211025-9fe3b590.pth) \| [log](https://download.openmmlab.com/mmocr/textdet/dbnet/dbnet_r50dcnv2_fpnc_sbn_1200e_icdar2015_20211025-9fe3b590.log.json) |


## Citation

```bibtex
@article{Liao_Wan_Yao_Chen_Bai_2020,
    title={Real-Time Scene Text Detection with Differentiable Binarization},
    journal={Proceedings of the AAAI Conference on Artificial Intelligence},
    author={Liao, Minghui and Wan, Zhaoyi and Yao, Cong and Chen, Kai and Bai, Xiang},
    year={2020},
    pages={11474-11481}}
```
