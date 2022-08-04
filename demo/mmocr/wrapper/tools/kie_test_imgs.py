#!/usr/bin/env python
# Copyright (c) OpenMMLab. All rights reserved.
import argparse
import ast
import os
import os.path as osp

import mmcv
import numpy as np
import torch
from mmcv import Config
from mmcv.image import tensor2imgs
from mmcv.parallel import MMDataParallel
from mmcv.runner import load_checkpoint

from mmocr.datasets import build_dataloader, build_dataset
from mmocr.models import build_detector


def save_results(model, img_meta, gt_bboxes, result, out_dir):
    assert 'filename' in img_meta, ('Please add "filename" '
                                    'to "meta_keys" in config.')
    assert 'ori_texts' in img_meta, ('Please add "ori_texts" '
                                     'to "meta_keys" in config.')

    out_json_file = osp.join(out_dir,
                             osp.basename(img_meta['filename']) + '.json')

    idx_to_cls = {}
    if model.module.class_list is not None:
        for line in mmcv.list_from_file(model.module.class_list):
            class_idx, class_label = line.strip().split()
            idx_to_cls[int(class_idx)] = class_label

    json_result = [{
        'text':
        text,
        'box':
        box,
        'pred':
        idx_to_cls.get(
            pred.argmax(-1).cpu().item(),
            pred.argmax(-1).cpu().item()),
        'conf':
        pred.max(-1)[0].cpu().item()
    } for text, box, pred in zip(img_meta['ori_texts'], gt_bboxes,
                                 result['nodes'])]

    mmcv.dump(json_result, out_json_file)


def test(model, data_loader, show=False, out_dir=None):
    model.eval()
    results = []
    dataset = data_loader.dataset
    prog_bar = mmcv.ProgressBar(len(dataset))
    for i, data in enumerate(data_loader):
        with torch.no_grad():
            result = model(return_loss=False, rescale=True, **data)

        batch_size = len(result)
        if show or out_dir:
            img_tensor = data['img'].data[0]
            img_metas = data['img_metas'].data[0]
            if np.prod(img_tensor.shape) == 0:
                imgs = [mmcv.imread(m['filename']) for m in img_metas]
            else:
                imgs = tensor2imgs(img_tensor, **img_metas[0]['img_norm_cfg'])
            assert len(imgs) == len(img_metas)
            gt_bboxes = [data['gt_bboxes'].data[0][0].numpy().tolist()]

            for i, (img, img_meta) in enumerate(zip(imgs, img_metas)):
                if 'img_shape' in img_meta:
                    h, w, _ = img_meta['img_shape']
                    img_show = img[:h, :w, :]
                else:
                    img_show = img

                if out_dir:
                    out_file = osp.join(out_dir,
                                        osp.basename(img_meta['filename']))
                else:
                    out_file = None

                model.module.show_result(
                    img_show,
                    result[i],
                    gt_bboxes[i],
                    show=show,
                    out_file=out_file)

                if out_dir:
                    save_results(model, img_meta, gt_bboxes[i], result[i],
                                 out_dir)

        for _ in range(batch_size):
            prog_bar.update()
    return results


def parse_args():
    parser = argparse.ArgumentParser(
        description='MMOCR visualize for kie model.')
    parser.add_argument('config', help='Test config file path.')
    parser.add_argument('checkpoint', help='Checkpoint file.')
    parser.add_argument('--show', action='store_true', help='Show results.')
    parser.add_argument(
        '--out-dir',
        help='Directory where the output images and results will be saved.')
    parser.add_argument('--local_rank', type=int, default=0)
    parser.add_argument(
        '--device',
        help='Use int or int list for gpu. Default is cpu',
        default=None)
    args = parser.parse_args()
    if 'LOCAL_RANK' not in os.environ:
        os.environ['LOCAL_RANK'] = str(args.local_rank)

    return args


def main():
    args = parse_args()
    assert args.show or args.out_dir, ('Please specify at least one '
                                       'operation (show the results / save )'
                                       'the results with the argument '
                                       '"--show" or "--out-dir".')
    device = args.device
    if device is not None:
        device = ast.literal_eval(f'[{device}]')
    cfg = Config.fromfile(args.config)
    # import modules from string list.
    if cfg.get('custom_imports', None):
        from mmcv.utils import import_modules_from_strings
        import_modules_from_strings(**cfg['custom_imports'])
    # set cudnn_benchmark
    if cfg.get('cudnn_benchmark', False):
        torch.backends.cudnn.benchmark = True

    distributed = False

    # build the dataloader
    dataset = build_dataset(cfg.data.test)
    data_loader = build_dataloader(
        dataset,
        samples_per_gpu=1,
        workers_per_gpu=cfg.data.workers_per_gpu,
        dist=distributed,
        shuffle=False)

    # build the model and load checkpoint
    cfg.model.train_cfg = None
    model = build_detector(cfg.model, test_cfg=cfg.get('test_cfg'))
    load_checkpoint(model, args.checkpoint, map_location='cpu')

    model = MMDataParallel(model, device_ids=device)
    test(model, data_loader, args.show, args.out_dir)


if __name__ == '__main__':
    main()
