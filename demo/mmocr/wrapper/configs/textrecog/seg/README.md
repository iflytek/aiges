# SegOCR

<!-- [ALGORITHM] -->
## Abstract

Just a simple Seg-based baseline for text recognition tasks.


## Dataset

### Train Dataset

| trainset  | instance_num | repeat_num | source |
| :-------: | :----------: | :--------: | :----: |
| SynthText |   7266686    |     1      | synth  |

### Test Dataset

| testset | instance_num |   type    |
| :-----: | :----------: | :-------: |
| IIIT5K  |     3000     |  regular  |
|   SVT   |     647      |  regular  |
|  IC13   |     1015     |  regular  |
|  CT80   |     288      | irregular |

## Results and Models

| Backbone |  Neck  | Head  |       |        | Regular Text |       |       | Irregular Text |                                                                                           download                                                                                           |
| :------: | :----: | :---: | :---: | :----: | :----------: | :---: | :---: | :------------: | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: |
|          |        |       |       | IIIT5K |     SVT      | IC13  |       |      CT80      |
| R31-1/16 | FPNOCR |  1x   |       |  90.9  |     81.8     | 90.7  |       |      80.9      | [model](https://download.openmmlab.com/mmocr/textrecog/seg/seg_r31_1by16_fpnocr_academic-72235b11.pth) \| [log](https://download.openmmlab.com/mmocr/textrecog/seg/20210325_112835.log.json) |

:::{note}

-   `R31-1/16` means the size (both height and width ) of feature from backbone is 1/16 of input image.
-   `1x` means the size (both height and width) of feature from head is the same with input image.
:::

## Citation

```bibtex
@unpublished{key,
  title={SegOCR Simple Baseline.},
  author={},
  note={Unpublished Manuscript},
  year={2021}
}
```
