# Copyright (c) OpenMMLab. All rights reserved.
import torch.nn.functional as F
from mmcv.runner import BaseModule, ModuleList
from torch import nn

from mmocr.models.builder import NECKS


class FPEM(BaseModule):
    """FPN-like feature fusion module in PANet.

    Args:
        in_channels (int): Number of input channels.
        init_cfg (dict or list[dict], optional): Initialization configs.
    """

    def __init__(self, in_channels=128, init_cfg=None):
        super().__init__(init_cfg=init_cfg)
        self.up_add1 = SeparableConv2d(in_channels, in_channels, 1)
        self.up_add2 = SeparableConv2d(in_channels, in_channels, 1)
        self.up_add3 = SeparableConv2d(in_channels, in_channels, 1)
        self.down_add1 = SeparableConv2d(in_channels, in_channels, 2)
        self.down_add2 = SeparableConv2d(in_channels, in_channels, 2)
        self.down_add3 = SeparableConv2d(in_channels, in_channels, 2)

    def forward(self, c2, c3, c4, c5):
        """
        Args:
            c2, c3, c4, c5 (Tensor): Each has the shape of
                :math:`(N, C_i, H_i, W_i)`.

        Returns:
            list[Tensor]: A list of 4 tensors of the same shape as input.
        """
        # upsample
        c4 = self.up_add1(self._upsample_add(c5, c4))  # c4 shape
        c3 = self.up_add2(self._upsample_add(c4, c3))
        c2 = self.up_add3(self._upsample_add(c3, c2))

        # downsample
        c3 = self.down_add1(self._upsample_add(c3, c2))
        c4 = self.down_add2(self._upsample_add(c4, c3))
        c5 = self.down_add3(self._upsample_add(c5, c4))  # c4 / 2
        return c2, c3, c4, c5

    def _upsample_add(self, x, y):
        return F.interpolate(x, size=y.size()[2:]) + y


class SeparableConv2d(BaseModule):

    def __init__(self, in_channels, out_channels, stride=1, init_cfg=None):
        super().__init__(init_cfg=init_cfg)

        self.depthwise_conv = nn.Conv2d(
            in_channels=in_channels,
            out_channels=in_channels,
            kernel_size=3,
            padding=1,
            stride=stride,
            groups=in_channels)
        self.pointwise_conv = nn.Conv2d(
            in_channels=in_channels, out_channels=out_channels, kernel_size=1)
        self.bn = nn.BatchNorm2d(out_channels)
        self.relu = nn.ReLU()

    def forward(self, x):
        x = self.depthwise_conv(x)
        x = self.pointwise_conv(x)
        x = self.bn(x)
        x = self.relu(x)
        return x


@NECKS.register_module()
class FPEM_FFM(BaseModule):
    """This code is from https://github.com/WenmuZhou/PAN.pytorch.

    Args:
        in_channels (list[int]): A list of 4 numbers of input channels.
        conv_out (int): Number of output channels.
        fpem_repeat (int): Number of FPEM layers before FFM operations.
        align_corners (bool): The interpolation behaviour in FFM operation,
            used in :func:`torch.nn.functional.interpolate`.
        init_cfg (dict or list[dict], optional): Initialization configs.
    """

    def __init__(self,
                 in_channels,
                 conv_out=128,
                 fpem_repeat=2,
                 align_corners=False,
                 init_cfg=dict(
                     type='Xavier', layer='Conv2d', distribution='uniform')):
        super().__init__(init_cfg=init_cfg)
        # reduce layers
        self.reduce_conv_c2 = nn.Sequential(
            nn.Conv2d(
                in_channels=in_channels[0],
                out_channels=conv_out,
                kernel_size=1), nn.BatchNorm2d(conv_out), nn.ReLU())
        self.reduce_conv_c3 = nn.Sequential(
            nn.Conv2d(
                in_channels=in_channels[1],
                out_channels=conv_out,
                kernel_size=1), nn.BatchNorm2d(conv_out), nn.ReLU())
        self.reduce_conv_c4 = nn.Sequential(
            nn.Conv2d(
                in_channels=in_channels[2],
                out_channels=conv_out,
                kernel_size=1), nn.BatchNorm2d(conv_out), nn.ReLU())
        self.reduce_conv_c5 = nn.Sequential(
            nn.Conv2d(
                in_channels=in_channels[3],
                out_channels=conv_out,
                kernel_size=1), nn.BatchNorm2d(conv_out), nn.ReLU())
        self.align_corners = align_corners
        self.fpems = ModuleList()
        for _ in range(fpem_repeat):
            self.fpems.append(FPEM(conv_out))

    def forward(self, x):
        """
        Args:
            x (list[Tensor]): A list of four tensors of shape
                :math:`(N, C_i, H_i, W_i)`, representing C2, C3, C4, C5
                features respectively. :math:`C_i` should matches the number in
                ``in_channels``.

        Returns:
            list[Tensor]: Four tensors of shape
            :math:`(N, C_{out}, H_0, W_0)` where :math:`C_{out}` is
            ``conv_out``.
        """
        c2, c3, c4, c5 = x
        # reduce channel
        c2 = self.reduce_conv_c2(c2)
        c3 = self.reduce_conv_c3(c3)
        c4 = self.reduce_conv_c4(c4)
        c5 = self.reduce_conv_c5(c5)

        # FPEM
        for i, fpem in enumerate(self.fpems):
            c2, c3, c4, c5 = fpem(c2, c3, c4, c5)
            if i == 0:
                c2_ffm = c2
                c3_ffm = c3
                c4_ffm = c4
                c5_ffm = c5
            else:
                c2_ffm = c2_ffm + c2
                c3_ffm = c3_ffm + c3
                c4_ffm = c4_ffm + c4
                c5_ffm = c5_ffm + c5

        # FFM
        c5 = F.interpolate(
            c5_ffm,
            c2_ffm.size()[-2:],
            mode='bilinear',
            align_corners=self.align_corners)
        c4 = F.interpolate(
            c4_ffm,
            c2_ffm.size()[-2:],
            mode='bilinear',
            align_corners=self.align_corners)
        c3 = F.interpolate(
            c3_ffm,
            c2_ffm.size()[-2:],
            mode='bilinear',
            align_corners=self.align_corners)
        outs = [c2_ffm, c3, c4, c5]
        return tuple(outs)
