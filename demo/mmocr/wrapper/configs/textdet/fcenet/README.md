# FCENet

> [Fourier Contour Embedding for Arbitrary-Shaped Text Detection](https://arxiv.org/abs/2104.10442)

<!-- [ALGORITHM] -->

## Abstract

One of the main challenges for arbitrary-shaped text detection is to design a good text instance representation that allows networks to learn diverse text geometry variances. Most of existing methods model text instances in image spatial domain via masks or contour point sequences in the Cartesian or the polar coordinate system. However, the mask representation might lead to expensive post-processing, while the point sequence one may have limited capability to model texts with highly-curved shapes. To tackle these problems, we model text instances in the Fourier domain and propose one novel Fourier Contour Embedding (FCE) method to represent arbitrary shaped text contours as compact signatures. We further construct FCENet with a backbone, feature pyramid networks (FPN) and a simple post-processing with the Inverse Fourier Transformation (IFT) and Non-Maximum Suppression (NMS). Different from previous methods, FCENet first predicts compact Fourier signatures of text instances, and then reconstructs text contours via IFT and NMS during test. Extensive experiments demonstrate that FCE is accurate and robust to fit contours of scene texts even with highly-curved shapes, and also validate the effectiveness and the good generalization of FCENet for arbitrary-shaped text detection. Furthermore, experimental results show that our FCENet is superior to the state-of-the-art (SOTA) methods on CTW1500 and Total-Text, especially on challenging highly-curved text subset.

<div align=center>
<img src="https://user-images.githubusercontent.com/22607038/142791859-1b0ebde4-b151-4c25-ba1b-f354bd8ddc8c.png"/>
</div>


## Results and models

### CTW1500

|                                 Method                                 |     Backbone     | Pretrained Model | Training set  |   Test set   | #epochs |  Test size  | Recall | Precision | Hmean |                                                                                                  Download                                                                                                   |
| :--------------------------------------------------------------------: | :--------------: | :--------------: | :-----------: | :----------: | :-----: | :---------: | :----: | :-------: | :---: | :---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
| [FCENet](/configs/textdet/fcenet/fcenet_r50dcnv2_fpn_1500e_ctw1500.py) | ResNet50 + DCNv2 |     ImageNet     | CTW1500 Train | CTW1500 Test |  1500   | (736, 1080) | 0.828  |   0.875   | 0.851 | [model](https://download.openmmlab.com/mmocr/textdet/fcenet/fcenet_r50dcnv2_fpn_1500e_ctw1500_20211022-e326d7ec.pth) \| [log](https://download.openmmlab.com/mmocr/textdet/fcenet/20210511_181328.log.json) |

### ICDAR2015

|                               Method                                | Backbone | Pretrained Model | Training set | Test set  | #epochs |  Test size   | Recall | Precision | Hmean |                                                                                                 Download                                                                                                 |
| :-----------------------------------------------------------------: | :------: | :--------------: | :----------: | :-------: | :-----: | :----------: | :----: | :-------: | :---: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
| [FCENet](/configs/textdet/fcenet/fcenet_r50_fpn_1500e_icdar2015.py) | ResNet50 |     ImageNet     |  IC15 Train  | IC15 Test |  1500   | (2260, 2260) | 0.819  |   0.880   | 0.849 | [model](https://download.openmmlab.com/mmocr/textdet/fcenet/fcenet_r50_fpn_1500e_icdar2015_20211022-daefb6ed.pth) \| [log](https://download.openmmlab.com/mmocr/textdet/fcenet/20210601_222655.log.json) |

## Citation

```bibtex
@InProceedings{zhu2021fourier,
      title={Fourier Contour Embedding for Arbitrary-Shaped Text Detection},
      author={Yiqin Zhu and Jianyong Chen and Lingyu Liang and Zhanghui Kuang and Lianwen Jin and Wayne Zhang},
      year={2021},
      booktitle = {CVPR}
      }
```
