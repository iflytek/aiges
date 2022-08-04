# Copyright (c) OpenMMLab. All rights reserved.
import copy
import unittest.mock as mock

import numpy as np
import pytest
import torchvision.transforms as TF
from mmdet.core import BitmapMasks, PolygonMasks
from PIL import Image

import mmocr.datasets.pipelines.transforms as transforms


@mock.patch('%s.transforms.np.random.random_sample' % __name__)
@mock.patch('%s.transforms.np.random.randint' % __name__)
def test_random_crop_instances(mock_randint, mock_sample):

    img_gt = np.array([[0, 0, 0, 0, 0], [0, 0, 0, 0, 0], [0, 0, 1, 1, 1],
                       [0, 0, 1, 1, 1], [0, 0, 1, 1, 1]])
    # test target is bigger than img size in sample_offset
    mock_sample.side_effect = [1]
    rci = transforms.RandomCropInstances(6, instance_key='gt_kernels')
    (i, j) = rci.sample_offset(img_gt, (5, 5))
    assert i == 0
    assert j == 0

    # test the second branch in sample_offset

    rci = transforms.RandomCropInstances(3, instance_key='gt_kernels')
    mock_sample.side_effect = [1]
    mock_randint.side_effect = [1, 2]
    (i, j) = rci.sample_offset(img_gt, (5, 5))
    assert i == 1
    assert j == 2

    mock_sample.side_effect = [1]
    mock_randint.side_effect = [1, 2]
    rci = transforms.RandomCropInstances(5, instance_key='gt_kernels')
    (i, j) = rci.sample_offset(img_gt, (5, 5))
    assert i == 0
    assert j == 0

    # test the first bracnh is sample_offset

    rci = transforms.RandomCropInstances(3, instance_key='gt_kernels')
    mock_sample.side_effect = [0.1]
    mock_randint.side_effect = [1, 1]
    (i, j) = rci.sample_offset(img_gt, (5, 5))
    assert i == 1
    assert j == 1

    # test crop_img(img, offset, target_size)

    img = img_gt
    offset = [0, 0]
    target = [6, 6]
    crop = rci.crop_img(img, offset, target)
    assert np.allclose(img, crop[0])
    assert np.allclose(crop[1], [0, 0, 5, 5])

    target = [3, 2]
    crop = rci.crop_img(img, offset, target)
    assert np.allclose(np.array([[0, 0], [0, 0], [0, 0]]), crop[0])
    assert np.allclose(crop[1], [0, 0, 2, 3])

    # test crop_bboxes
    canvas_box = np.array([2, 3, 5, 5])
    bboxes = np.array([[2, 3, 4, 4], [0, 0, 1, 1], [1, 2, 4, 4],
                       [0, 0, 10, 10]])
    kept_bboxes, kept_idx = rci.crop_bboxes(bboxes, canvas_box)
    assert np.allclose(kept_bboxes,
                       np.array([[0, 0, 2, 1], [0, 0, 2, 1], [0, 0, 3, 2]]))
    assert kept_idx == [0, 2, 3]

    bboxes = np.array([[10, 10, 11, 11], [0, 0, 1, 1]])
    kept_bboxes, kept_idx = rci.crop_bboxes(bboxes, canvas_box)
    assert kept_bboxes.size == 0
    assert kept_bboxes.shape == (0, 4)
    assert len(kept_idx) == 0

    # test __call__
    rci = transforms.RandomCropInstances(3, instance_key='gt_kernels')
    results = {}
    gt_kernels = [img_gt, img_gt.copy()]
    results['gt_kernels'] = BitmapMasks(gt_kernels, 5, 5)
    results['img'] = img_gt.copy()
    results['mask_fields'] = ['gt_kernels']
    mock_sample.side_effect = [0.1]
    mock_randint.side_effect = [1, 1]
    output = rci(results)
    target = np.array([[0, 0, 0], [0, 1, 1], [0, 1, 1]])
    assert output['img_shape'] == (3, 3)

    assert np.allclose(output['img'], target)

    assert np.allclose(output['gt_kernels'].masks[0], target)
    assert np.allclose(output['gt_kernels'].masks[1], target)


@mock.patch('%s.transforms.np.random.random_sample' % __name__)
def test_scale_aspect_jitter(mock_random):
    img_scale = [(3000, 1000)]  # unused
    ratio_range = (0.5, 1.5)
    aspect_ratio_range = (1, 1)
    multiscale_mode = 'value'
    long_size_bound = 2000
    short_size_bound = 640
    resize_type = 'long_short_bound'
    keep_ratio = False
    jitter = transforms.ScaleAspectJitter(
        img_scale=img_scale,
        ratio_range=ratio_range,
        aspect_ratio_range=aspect_ratio_range,
        multiscale_mode=multiscale_mode,
        long_size_bound=long_size_bound,
        short_size_bound=short_size_bound,
        resize_type=resize_type,
        keep_ratio=keep_ratio)
    mock_random.side_effect = [0.5]

    # test sample_from_range

    result = jitter.sample_from_range([100, 200])
    assert result == 150

    # test _random_scale
    results = {}
    results['img'] = np.zeros((4000, 1000))
    mock_random.side_effect = [0.5, 1]
    jitter._random_scale(results)
    # scale1 0.5， scale2=1 scale =0.5  650/1000, w, h
    # print(results['scale'])
    assert results['scale'] == (650, 2600)


@mock.patch('%s.transforms.np.random.random_sample' % __name__)
def test_random_rotate(mock_random):

    mock_random.side_effect = [0.5, 0]
    results = {}
    img = np.random.rand(5, 5)
    results['img'] = img.copy()
    results['mask_fields'] = ['masks']
    gt_kernels = [results['img'].copy()]
    results['masks'] = BitmapMasks(gt_kernels, 5, 5)

    rotater = transforms.RandomRotateTextDet()

    results = rotater(results)
    assert np.allclose(results['img'], img)
    assert np.allclose(results['masks'].masks, img)


def test_color_jitter():
    img = np.ones((64, 256, 3), dtype=np.uint8)
    results = {'img': img}

    pt_official_color_jitter = TF.ColorJitter()
    output1 = pt_official_color_jitter(img)

    color_jitter = transforms.ColorJitter()
    output2 = color_jitter(results)

    assert np.allclose(output1, output2['img'])


def test_affine_jitter():
    img = np.ones((64, 256, 3), dtype=np.uint8)
    results = {'img': img}

    pt_official_affine_jitter = TF.RandomAffine(degrees=0)
    output1 = pt_official_affine_jitter(Image.fromarray(img))

    affine_jitter = transforms.AffineJitter(
        degrees=0,
        translate=None,
        scale=None,
        shear=None,
        resample=False,
        fillcolor=0)
    output2 = affine_jitter(results)

    assert np.allclose(np.array(output1), output2['img'])


def test_random_scale():
    h, w, c = 100, 100, 3
    img = np.ones((h, w, c), dtype=np.uint8)
    results = {'img': img, 'img_shape': (h, w, c)}

    polygon = np.array([0., 0., 0., 10., 10., 10., 10., 0.])

    results['gt_masks'] = PolygonMasks([[polygon]], *(img.shape[:2]))
    results['mask_fields'] = ['gt_masks']

    size = 100
    scale = (2., 2.)
    random_scaler = transforms.RandomScaling(size=size, scale=scale)

    results = random_scaler(results)

    out_img = results['img']
    out_poly = results['gt_masks'].masks[0][0]
    gt_poly = polygon * 2

    assert np.allclose(out_img.shape, (2 * h, 2 * w, c))
    assert np.allclose(out_poly, gt_poly)


@mock.patch('%s.transforms.np.random.randint' % __name__)
def test_random_crop_flip(mock_randint):
    img = np.ones((10, 10, 3), dtype=np.uint8)
    img[0, 0, :] = 0
    results = {'img': img, 'img_shape': img.shape}

    polygon = np.array([0., 0., 0., 10., 10., 10., 10., 0.])

    results['gt_masks'] = PolygonMasks([[polygon]], *(img.shape[:2]))
    results['gt_masks_ignore'] = PolygonMasks([], *(img.shape[:2]))
    results['mask_fields'] = ['gt_masks', 'gt_masks_ignore']

    crop_ratio = 1.1
    iter_num = 3
    random_crop_fliper = transforms.RandomCropFlip(
        crop_ratio=crop_ratio, iter_num=iter_num)

    # test crop_target
    pad_ratio = 0.1
    h, w = img.shape[:2]
    pad_h = int(h * pad_ratio)
    pad_w = int(w * pad_ratio)
    all_polys = results['gt_masks'].masks
    h_axis, w_axis = random_crop_fliper.generate_crop_target(
        img, all_polys, pad_h, pad_w)

    assert np.allclose(h_axis, (0, 11))
    assert np.allclose(w_axis, (0, 11))

    # test __call__
    polygon = np.array([1., 1., 1., 9., 9., 9., 9., 1.])
    results['gt_masks'] = PolygonMasks([[polygon]], *(img.shape[:2]))
    results['gt_masks_ignore'] = PolygonMasks([[polygon]], *(img.shape[:2]))

    mock_randint.side_effect = [0, 1, 2]
    results = random_crop_fliper(results)

    out_img = results['img']
    out_poly = results['gt_masks'].masks[0][0]
    gt_img = img
    gt_poly = polygon

    assert np.allclose(out_img, gt_img)
    assert np.allclose(out_poly, gt_poly)


@mock.patch('%s.transforms.np.random.random_sample' % __name__)
@mock.patch('%s.transforms.np.random.randint' % __name__)
def test_random_crop_poly_instances(mock_randint, mock_sample):
    results = {}
    img = np.zeros((30, 30, 3))
    poly_masks = PolygonMasks([[
        np.array([5., 5., 25., 5., 25., 10., 5., 10.])
    ], [np.array([5., 20., 25., 20., 25., 25., 5., 25.])]], 30, 30)
    results['img'] = img
    results['gt_masks'] = poly_masks
    results['gt_masks_ignore'] = PolygonMasks([], 30, 30)
    results['mask_fields'] = ['gt_masks', 'gt_masks_ignore']
    results['gt_labels'] = [1, 1]
    rcpi = transforms.RandomCropPolyInstances(
        instance_key='gt_masks', crop_ratio=1.0, min_side_ratio=0.3)

    # test sample_crop_box(img_size, results)
    mock_randint.side_effect = [0, 0, 0, 0, 30, 0, 0, 0, 15]
    crop_box = rcpi.sample_crop_box((30, 30), results)
    assert np.allclose(np.array(crop_box), np.array([0, 0, 30, 15]))

    # test __call__
    mock_randint.side_effect = [0, 0, 0, 0, 30, 0, 15, 0, 30]
    mock_sample.side_effect = [0.1]
    output = rcpi(results)
    target = np.array([5., 5., 25., 5., 25., 10., 5., 10.])
    assert len(output['gt_masks']) == 1
    assert len(output['gt_masks_ignore']) == 0
    assert np.allclose(output['gt_masks'].masks[0][0], target)
    assert output['img'].shape == (15, 30, 3)

    # test __call__ with blank instace_key masks
    mock_randint.side_effect = [0, 0, 0, 0, 30, 0, 15, 0, 30]
    mock_sample.side_effect = [0.1]
    rcpi = transforms.RandomCropPolyInstances(
        instance_key='gt_masks_ignore', crop_ratio=1.0, min_side_ratio=0.3)
    results['img'] = img
    results['gt_masks'] = poly_masks
    output = rcpi(results)
    assert len(output['gt_masks']) == 2
    assert np.allclose(output['gt_masks'].masks[0][0], poly_masks.masks[0][0])
    assert np.allclose(output['gt_masks'].masks[1][0], poly_masks.masks[1][0])
    assert output['img'].shape == (30, 30, 3)


@mock.patch('%s.transforms.np.random.random_sample' % __name__)
def test_random_rotate_poly_instances(mock_sample):
    results = {}
    img = np.zeros((30, 30, 3))
    poly_masks = PolygonMasks(
        [[np.array([10., 10., 20., 10., 20., 20., 10., 20.])]], 30, 30)
    results['img'] = img
    results['gt_masks'] = poly_masks
    results['mask_fields'] = ['gt_masks']
    rrpi = transforms.RandomRotatePolyInstances(rotate_ratio=1.0, max_angle=90)

    mock_sample.side_effect = [0., 1.]
    output = rrpi(results)
    assert np.allclose(output['gt_masks'].masks[0][0],
                       np.array([10., 20., 10., 10., 20., 10., 20., 20.]))
    assert output['img'].shape == (30, 30, 3)


@mock.patch('%s.transforms.np.random.random_sample' % __name__)
def test_square_resize_pad(mock_sample):
    results = {}
    img = np.zeros((15, 30, 3))
    polygon = np.array([10., 5., 20., 5., 20., 10., 10., 10.])
    poly_masks = PolygonMasks([[polygon]], 15, 30)
    results['img'] = img
    results['gt_masks'] = poly_masks
    results['mask_fields'] = ['gt_masks']
    srp = transforms.SquareResizePad(target_size=40, pad_ratio=0.5)

    # test resize with padding
    mock_sample.side_effect = [0.]
    output = srp(results)
    target = 4. / 3 * polygon
    target[1::2] += 10.
    assert np.allclose(output['gt_masks'].masks[0][0], target)
    assert output['img'].shape == (40, 40, 3)

    # test resize to square without padding
    results['img'] = img
    results['gt_masks'] = poly_masks
    mock_sample.side_effect = [1.]
    output = srp(results)
    target = polygon.copy()
    target[::2] *= 4. / 3
    target[1::2] *= 8. / 3
    assert np.allclose(output['gt_masks'].masks[0][0], target)
    assert output['img'].shape == (40, 40, 3)


def test_pyramid_rescale():
    img = np.random.randint(0, 256, size=(128, 100, 3), dtype=np.uint8)
    x = {'img': copy.deepcopy(img)}
    f = transforms.PyramidRescale()
    results = f(x)
    assert results['img'].shape == (128, 100, 3)

    # Test invalid inputs
    with pytest.raises(AssertionError):
        transforms.PyramidRescale(base_shape=(128))
    with pytest.raises(AssertionError):
        transforms.PyramidRescale(base_shape=128)
    with pytest.raises(AssertionError):
        transforms.PyramidRescale(factor=[])
    with pytest.raises(AssertionError):
        transforms.PyramidRescale(randomize_factor=[])
    with pytest.raises(AssertionError):
        f({})

    # Test factor = 0
    f_derandomized = transforms.PyramidRescale(
        factor=0, randomize_factor=False)
    results = f_derandomized({'img': copy.deepcopy(img)})
    assert np.all(results['img'] == img)
