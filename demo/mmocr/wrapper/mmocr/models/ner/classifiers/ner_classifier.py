# Copyright (c) OpenMMLab. All rights reserved.
from mmocr.models.builder import (DETECTORS, build_convertor, build_decoder,
                                  build_encoder, build_loss)
from mmocr.models.textrecog.recognizer.base import BaseRecognizer


@DETECTORS.register_module()
class NerClassifier(BaseRecognizer):
    """Base class for NER classifier."""

    def __init__(self,
                 encoder,
                 decoder,
                 loss,
                 label_convertor,
                 train_cfg=None,
                 test_cfg=None,
                 init_cfg=None):
        super().__init__(init_cfg=init_cfg)
        self.label_convertor = build_convertor(label_convertor)

        self.encoder = build_encoder(encoder)

        decoder.update(num_labels=self.label_convertor.num_labels)
        self.decoder = build_decoder(decoder)

        loss.update(num_labels=self.label_convertor.num_labels)
        self.loss = build_loss(loss)

    def extract_feat(self, imgs):
        """Extract features from images."""
        raise NotImplementedError(
            'Extract feature module is not implemented yet.')

    def forward_train(self, imgs, img_metas, **kwargs):
        encode_out = self.encoder(img_metas)
        logits, _ = self.decoder(encode_out)
        loss = self.loss(logits, img_metas)
        return loss

    def forward_test(self, imgs, img_metas, **kwargs):
        encode_out = self.encoder(img_metas)
        _, preds = self.decoder(encode_out)
        pred_entities = self.label_convertor.convert_pred2entities(
            preds, img_metas['attention_masks'])
        return pred_entities

    def aug_test(self, imgs, img_metas, **kwargs):
        raise NotImplementedError('Augmentation test is not implemented yet.')

    def simple_test(self, img, img_metas, **kwargs):
        raise NotImplementedError('Simple test is not implemented yet.')
