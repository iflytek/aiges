# Text Recognition Training set, including:
# Synthetic Datasets: SynthText (with character level boxes)

train_img_root = 'data/mixture'

train_img_prefix = f'{train_img_root}/SynthText'

train_ann_file = f'{train_img_root}/SynthText/instances_train.txt'

train = dict(
    type='OCRSegDataset',
    img_prefix=train_img_prefix,
    ann_file=train_ann_file,
    loader=dict(
        type='AnnFileLoader',
        repeat=1,
        file_format='txt',
        parser=dict(
            type='LineJsonParser', keys=['file_name', 'annotations', 'text'])),
    pipeline=None,
    test_mode=False)

train_list = [train]
