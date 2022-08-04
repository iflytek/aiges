prefix = 'tests/data/ocr_char_ann_toy_dataset/'

train = dict(
    type='OCRSegDataset',
    img_prefix=f'{prefix}/imgs',
    ann_file=f'{prefix}/instances_train.txt',
    loader=dict(
        type='AnnFileLoader',
        repeat=100,
        file_format='txt',
        parser=dict(
            type='LineJsonParser', keys=['file_name', 'annotations', 'text'])),
    pipeline=None,
    test_mode=True)

test = dict(
    type='OCRDataset',
    img_prefix=f'{prefix}/imgs',
    ann_file=f'{prefix}/instances_test.txt',
    loader=dict(
        type='AnnFileLoader',
        repeat=1,
        file_format='txt',
        parser=dict(
            type='LineStrParser',
            keys=['filename', 'text'],
            keys_idx=[0, 1],
            separator=' ')),
    pipeline=None,
    test_mode=True)

train_list = [train]

test_list = [test]
