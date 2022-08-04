# Copyright (c) OpenMMLab. All rights reserved.
from mmcv.runner import BaseModule

from mmocr.models.builder import DECODERS


@DECODERS.register_module()
class BaseDecoder(BaseModule):
    """Base decoder class for text recognition."""

    def __init__(self, init_cfg=None, **kwargs):
        super().__init__(init_cfg=init_cfg)

    def forward_train(self, feat, out_enc, targets_dict, img_metas):
        raise NotImplementedError

    def forward_test(self, feat, out_enc, img_metas):
        raise NotImplementedError

    def forward(self,
                feat,
                out_enc,
                targets_dict=None,
                img_metas=None,
                train_mode=True):
        self.train_mode = train_mode
        if train_mode:
            return self.forward_train(feat, out_enc, targets_dict, img_metas)

        return self.forward_test(feat, out_enc, img_metas)
