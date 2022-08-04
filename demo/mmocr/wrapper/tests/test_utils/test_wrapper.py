# Copyright (c) OpenMMLab. All rights reserved.
import numpy as np
import pytest
import torch

from mmocr.models.textdet.postprocess import (DBPostprocessor,
                                              FCEPostprocessor,
                                              TextSnakePostprocessor)
from mmocr.models.textdet.postprocess.utils import comps2boundaries, poly_nms


def test_db_boxes_from_bitmaps():
    """Test the boxes_from_bitmaps function in db_decoder."""
    pred = np.array([[[0.8, 0.8, 0.8, 0.8, 0], [0.8, 0.8, 0.8, 0.8, 0],
                      [0.8, 0.8, 0.8, 0.8, 0], [0.8, 0.8, 0.8, 0.8, 0],
                      [0.8, 0.8, 0.8, 0.8, 0]]])
    preds = torch.FloatTensor(pred).requires_grad_(True)
    db_decode = DBPostprocessor(text_repr_type='quad', min_text_width=0)
    boxes = db_decode(preds)
    assert len(boxes) == 1


def test_fcenet_decode():

    k = 1
    preds = []
    preds.append(torch.ones(1, 4, 10, 10))
    preds.append(torch.ones(1, 4 * k + 2, 10, 10))
    fcenet_decode = FCEPostprocessor(
        fourier_degree=k, num_reconstr_points=50, nms_thr=0.01)
    boundaries = fcenet_decode(preds=preds, scale=1)

    assert isinstance(boundaries, list)


def test_poly_nms():
    threshold = 0
    polygons = []
    polygons.append([10, 10, 10, 30, 30, 30, 30, 10, 0.95])
    polygons.append([15, 15, 15, 25, 25, 25, 25, 15, 0.9])
    polygons.append([40, 40, 40, 50, 50, 50, 50, 40, 0.85])
    polygons.append([5, 5, 5, 15, 15, 15, 15, 5, 0.7])

    keep_poly = poly_nms(polygons, threshold)
    assert isinstance(keep_poly, list)


def test_comps2boundaries():

    # test comps2boundaries
    x1 = np.arange(2, 18, 2)
    x2 = x1 + 2
    y1 = np.ones(8) * 2
    y2 = y1 + 2
    comp_scores = np.ones(8, dtype=np.float32) * 0.9
    text_comps = np.stack([x1, y1, x2, y1, x2, y2, x1, y2,
                           comp_scores]).transpose()
    comp_labels = np.array([1, 1, 1, 1, 1, 3, 5, 5])
    shuffle = [3, 2, 5, 7, 6, 0, 4, 1]
    boundaries = comps2boundaries(text_comps[shuffle], comp_labels[shuffle])
    assert len(boundaries) == 3

    # test comps2boundaries with blank inputs
    boundaries = comps2boundaries(text_comps[[]], comp_labels[[]])
    assert len(boundaries) == 0


def test_textsnake_decode():

    maps = torch.zeros((1, 6, 224, 224), dtype=torch.float)
    maps[:, 0:2, :, :] = -10.
    maps[:, 0, 60:100, 50:170] = 10.
    maps[:, 1, 75:85, 60:160] = 10.
    maps[:, 2, 75:85, 60:160] = 0.
    maps[:, 3, 75:85, 60:160] = 1.
    maps[:, 4, 75:85, 60:160] = 10.
    # test decoding with text center region of small area
    maps[:, 0:2, 150:152, 5:7] = 10.
    textsnake_decode = TextSnakePostprocessor()
    results = textsnake_decode(torch.squeeze(maps))
    assert len(results) == 1

    # test decoding with small radius
    maps.fill_(0.)
    maps[:, 0:2, :, :] = -10.
    maps[:, 0, 120:140, 20:40] = 10.
    maps[:, 1, 120:140, 20:40] = 10.
    maps[:, 2, 120:140, 20:40] = 0.
    maps[:, 3, 120:140, 20:40] = 1.
    maps[:, 4, 120:140, 20:40] = 0.5

    results = textsnake_decode(torch.squeeze(maps))
    assert len(results) == 0


def test_db_decode():
    pred = torch.zeros((1, 8, 8))
    pred[0, 2:7, 2:7] = 0.8
    expect_result_quad = [[
        1.0, 8.0, 1.0, 1.0, 8.0, 1.0, 8.0, 8.0, 0.800000011920929
    ]]
    expect_result_poly = [[
        8, 2, 8, 6, 6, 8, 2, 8, 1, 6, 1, 2, 2, 1, 6, 1, 0.800000011920929
    ]]
    with pytest.raises(AssertionError):
        DBPostprocessor(text_repr_type='dummpy')
    db_decode = DBPostprocessor(text_repr_type='quad', min_text_width=1)
    result_quad = db_decode(preds=pred)
    db_decode = DBPostprocessor(text_repr_type='poly', min_text_width=1)
    result_poly = db_decode(preds=pred)
    assert result_quad == expect_result_quad
    assert result_poly == expect_result_poly
