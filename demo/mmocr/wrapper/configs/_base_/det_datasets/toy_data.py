root = 'tests/data/toy_dataset'

# dataset with type='TextDetDataset'
train1 = dict(
    type='TextDetDataset',
    img_prefix=f'{root}/imgs',
    ann_file=f'{root}/instances_test.txt',
    loader=dict(
        type='AnnFileLoader',
        repeat=4,
        file_format='txt',
        parser=dict(
            type='LineJsonParser',
            keys=['file_name', 'height', 'width', 'annotations'])),
    pipeline=None,
    test_mode=False)

# dataset with type='IcdarDataset'
train2 = dict(
    type='IcdarDataset',
    ann_file=f'{root}/instances_test.json',
    img_prefix=f'{root}/imgs',
    pipeline=None)

test = dict(
    type='TextDetDataset',
    img_prefix=f'{root}/imgs',
    ann_file=f'{root}/instances_test.txt',
    loader=dict(
        type='AnnFileLoader',
        repeat=1,
        file_format='txt',
        parser=dict(
            type='LineJsonParser',
            keys=['file_name', 'height', 'width', 'annotations'])),
    pipeline=None,
    test_mode=True)

train_list = [train1, train2]

test_list = [test]
