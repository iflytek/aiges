# Copyright (c) OpenMMLab. All rights reserved.
import warnings

import torch.nn as nn
from mmcv.runner import BaseModule

from mmocr.models.builder import HEADS
from .head_mixin import HeadMixin


@HEADS.register_module()
class TextSnakeHead(HeadMixin, BaseModule):
    """The class for TextSnake head: TextSnake: A Flexible Representation for
    Detecting Text of Arbitrary Shapes.

    TextSnake: `A Flexible Representation for Detecting Text of Arbitrary
    Shapes <https://arxiv.org/abs/1807.01544>`_.

    Args:
        in_channels (int): Number of input channels.
        out_channels (int): Number of output channels.
        downsample_ratio (float): Downsample ratio.
        loss (dict): Configuration dictionary for loss type.
        postprocessor (dict): Config of postprocessor for TextSnake.
        train_cfg, test_cfg: Depreciated.
        init_cfg (dict or list[dict], optional): Initialization configs.
    """

    def __init__(self,
                 in_channels,
                 out_channels=5,
                 downsample_ratio=1.0,
                 loss=dict(type='TextSnakeLoss'),
                 postprocessor=dict(
                     type='TextSnakePostprocessor', text_repr_type='poly'),
                 train_cfg=None,
                 test_cfg=None,
                 init_cfg=dict(
                     type='Normal',
                     override=dict(name='out_conv'),
                     mean=0,
                     std=0.01),
                 **kwargs):
        old_keys = ['text_repr_type', 'decoding_type']
        for key in old_keys:
            if kwargs.get(key, None):
                postprocessor[key] = kwargs.get(key)
                warnings.warn(
                    f'{key} is deprecated, please specify '
                    'it in postprocessor config dict. See '
                    'https://github.com/open-mmlab/mmocr/pull/640 '
                    'for details.', UserWarning)
        BaseModule.__init__(self, init_cfg=init_cfg)
        HeadMixin.__init__(self, loss, postprocessor)

        assert isinstance(in_channels, int)
        self.in_channels = in_channels
        self.out_channels = out_channels
        self.downsample_ratio = downsample_ratio
        self.train_cfg = train_cfg
        self.test_cfg = test_cfg

        self.out_conv = nn.Conv2d(
            in_channels=self.in_channels,
            out_channels=self.out_channels,
            kernel_size=1,
            stride=1,
            padding=0)

    def forward(self, inputs):
        """
        Args:
            inputs (Tensor): Shape :math:`(N, C_{in}, H, W)`, where
                :math:`C_{in}` is ``in_channels``. :math:`H` and :math:`W`
                should be the same as the input of backbone.

        Returns:
            Tensor: A tensor of shape :math:`(N, 5, H, W)`.
        """
        outputs = self.out_conv(inputs)
        return outputs
