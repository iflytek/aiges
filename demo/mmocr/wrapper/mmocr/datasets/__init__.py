# Copyright (c) OpenMMLab. All rights reserved.
from mmdet.datasets.builder import DATASETS, build_dataloader, build_dataset

from . import utils
from .base_dataset import BaseDataset
from .icdar_dataset import IcdarDataset
from .kie_dataset import KIEDataset
from .ner_dataset import NerDataset
from .ocr_dataset import OCRDataset
from .ocr_seg_dataset import OCRSegDataset
from .openset_kie_dataset import OpensetKIEDataset
from .pipelines import CustomFormatBundle, DBNetTargets, FCENetTargets
from .text_det_dataset import TextDetDataset
from .uniform_concat_dataset import UniformConcatDataset
from .utils import *  # NOQA

__all__ = [
    'DATASETS', 'IcdarDataset', 'build_dataloader', 'build_dataset',
    'BaseDataset', 'OCRDataset', 'TextDetDataset', 'CustomFormatBundle',
    'DBNetTargets', 'OCRSegDataset', 'KIEDataset', 'FCENetTargets',
    'NerDataset', 'UniformConcatDataset', 'OpensetKIEDataset'
]

__all__ += utils.__all__
