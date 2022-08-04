# Copyright (c) OpenMMLab. All rights reserved.
import math

import torch
import torch.nn as nn

from mmocr.models.builder import LOSSES


@LOSSES.register_module()
class CTCLoss(nn.Module):
    """Implementation of loss module for CTC-loss based text recognition.

    Args:
        flatten (bool): If True, use flattened targets, else padded targets.
        blank (int): Blank label. Default 0.
        reduction (str): Specifies the reduction to apply to the output,
            should be one of the following: ('none', 'mean', 'sum').
        zero_infinity (bool): Whether to zero infinite losses and
            the associated gradients. Default: False.
            Infinite losses mainly occur when the inputs
            are too short to be aligned to the targets.
    """

    def __init__(self,
                 flatten=True,
                 blank=0,
                 reduction='mean',
                 zero_infinity=False,
                 **kwargs):
        super().__init__()
        assert isinstance(flatten, bool)
        assert isinstance(blank, int)
        assert isinstance(reduction, str)
        assert isinstance(zero_infinity, bool)

        self.flatten = flatten
        self.blank = blank
        self.ctc_loss = nn.CTCLoss(
            blank=blank, reduction=reduction, zero_infinity=zero_infinity)

    def forward(self, outputs, targets_dict, img_metas=None):
        """
        Args:
            outputs (Tensor): A raw logit tensor of shape :math:`(N, T, C)`.
            targets_dict (dict): A dict with 3 keys ``target_lengths``,
                ``flatten_targets`` and ``targets``.

                - | ``target_lengths`` (Tensor): A tensor of shape :math:`(N)`.
                    Each item is the length of a word.

                - | ``flatten_targets`` (Tensor): Used if ``self.flatten=True``
                    (default). A tensor of shape
                    (sum(targets_dict['target_lengths'])). Each item is the
                    index of a character.

                - | ``targets`` (Tensor): Used if ``self.flatten=False``. A
                    tensor of :math:`(N, T)`. Empty slots are padded with
                    ``self.blank``.

            img_metas (dict): A dict that contains meta information of input
                images. Preferably with the key ``valid_ratio``.

        Returns:
            dict: The loss dict with key ``loss_ctc``.
        """
        valid_ratios = None
        if img_metas is not None:
            valid_ratios = [
                img_meta.get('valid_ratio', 1.0) for img_meta in img_metas
            ]

        outputs = torch.log_softmax(outputs, dim=2)
        bsz, seq_len = outputs.size(0), outputs.size(1)
        outputs_for_loss = outputs.permute(1, 0, 2).contiguous()  # T * N * C

        if self.flatten:
            targets = targets_dict['flatten_targets']
        else:
            targets = torch.full(
                size=(bsz, seq_len), fill_value=self.blank, dtype=torch.long)
            for idx, tensor in enumerate(targets_dict['targets']):
                valid_len = min(tensor.size(0), seq_len)
                targets[idx, :valid_len] = tensor[:valid_len]

        target_lengths = targets_dict['target_lengths']
        target_lengths = torch.clamp(target_lengths, min=1, max=seq_len).long()

        input_lengths = torch.full(
            size=(bsz, ), fill_value=seq_len, dtype=torch.long)
        if not self.flatten and valid_ratios is not None:
            input_lengths = [
                math.ceil(valid_ratio * seq_len)
                for valid_ratio in valid_ratios
            ]
            input_lengths = torch.Tensor(input_lengths).long()

        loss_ctc = self.ctc_loss(outputs_for_loss, targets, input_lengths,
                                 target_lengths)

        losses = dict(loss_ctc=loss_ctc)

        return losses
