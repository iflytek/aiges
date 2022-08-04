label_convertor = dict(
    type='SegConvertor', dict_type='DICT36', with_unknown=True, lower=True)

model = dict(
    type='SegRecognizer',
    backbone=dict(
        type='ResNet31OCR',
        layers=[1, 2, 5, 3],
        channels=[32, 64, 128, 256, 512, 512],
        out_indices=[0, 1, 2, 3],
        stage4_pool_cfg=dict(kernel_size=2, stride=2),
        last_stage_pool=True),
    neck=dict(
        type='FPNOCR', in_channels=[128, 256, 512, 512], out_channels=256),
    head=dict(
        type='SegHead',
        in_channels=256,
        upsample_param=dict(scale_factor=2.0, mode='nearest')),
    loss=dict(
        type='SegLoss', seg_downsample_ratio=1.0, seg_with_loss_weight=True),
    label_convertor=label_convertor)
