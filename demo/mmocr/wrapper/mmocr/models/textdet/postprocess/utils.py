# Copyright (c) OpenMMLab. All rights reserved.
import functools
import operator

import cv2
import numpy as np
import pyclipper
from numpy.fft import ifft
from numpy.linalg import norm
from shapely.geometry import Polygon

from mmocr.core.evaluation.utils import boundary_iou


def filter_instance(area, confidence, min_area, min_confidence):
    return bool(area < min_area or confidence < min_confidence)


def box_score_fast(bitmap, _box):
    h, w = bitmap.shape[:2]
    box = _box.copy()
    xmin = np.clip(np.floor(box[:, 0].min()).astype(np.int32), 0, w - 1)
    xmax = np.clip(np.ceil(box[:, 0].max()).astype(np.int32), 0, w - 1)
    ymin = np.clip(np.floor(box[:, 1].min()).astype(np.int32), 0, h - 1)
    ymax = np.clip(np.ceil(box[:, 1].max()).astype(np.int32), 0, h - 1)

    mask = np.zeros((ymax - ymin + 1, xmax - xmin + 1), dtype=np.uint8)
    box[:, 0] = box[:, 0] - xmin
    box[:, 1] = box[:, 1] - ymin
    cv2.fillPoly(mask, box.reshape(1, -1, 2).astype(np.int32), 1)
    return cv2.mean(bitmap[ymin:ymax + 1, xmin:xmax + 1], mask)[0]


def unclip(box, unclip_ratio=1.5):
    poly = Polygon(box)
    distance = poly.area * unclip_ratio / poly.length
    offset = pyclipper.PyclipperOffset()
    offset.AddPath(box, pyclipper.JT_ROUND, pyclipper.ET_CLOSEDPOLYGON)
    expanded = np.array(offset.Execute(distance))
    return expanded


def fill_hole(input_mask):
    h, w = input_mask.shape
    canvas = np.zeros((h + 2, w + 2), np.uint8)
    canvas[1:h + 1, 1:w + 1] = input_mask.copy()

    mask = np.zeros((h + 4, w + 4), np.uint8)

    cv2.floodFill(canvas, mask, (0, 0), 1)
    canvas = canvas[1:h + 1, 1:w + 1].astype(np.bool)

    return ~canvas | input_mask


def centralize(points_yx,
               normal_sin,
               normal_cos,
               radius,
               contour_mask,
               step_ratio=0.03):

    h, w = contour_mask.shape
    top_yx = bot_yx = points_yx
    step_flags = np.ones((len(points_yx), 1), dtype=np.bool)
    step = step_ratio * radius * np.hstack([normal_sin, normal_cos])
    while np.any(step_flags):
        next_yx = np.array(top_yx + step, dtype=np.int32)
        next_y, next_x = next_yx[:, 0], next_yx[:, 1]
        step_flags = (next_y >= 0) & (next_y < h) & (next_x > 0) & (
            next_x < w) & contour_mask[np.clip(next_y, 0, h - 1),
                                       np.clip(next_x, 0, w - 1)]
        top_yx = top_yx + step_flags.reshape((-1, 1)) * step
    step_flags = np.ones((len(points_yx), 1), dtype=np.bool)
    while np.any(step_flags):
        next_yx = np.array(bot_yx - step, dtype=np.int32)
        next_y, next_x = next_yx[:, 0], next_yx[:, 1]
        step_flags = (next_y >= 0) & (next_y < h) & (next_x > 0) & (
            next_x < w) & contour_mask[np.clip(next_y, 0, h - 1),
                                       np.clip(next_x, 0, w - 1)]
        bot_yx = bot_yx - step_flags.reshape((-1, 1)) * step
    centers = np.array((top_yx + bot_yx) * 0.5, dtype=np.int32)
    return centers


def merge_disks(disks, disk_overlap_thr):
    xy = disks[:, 0:2]
    radius = disks[:, 2]
    scores = disks[:, 3]
    order = scores.argsort()[::-1]

    merged_disks = []
    while order.size > 0:
        if order.size == 1:
            merged_disks.append(disks[order])
            break
        i = order[0]
        d = norm(xy[i] - xy[order[1:]], axis=1)
        ri = radius[i]
        r = radius[order[1:]]
        d_thr = (ri + r) * disk_overlap_thr

        merge_inds = np.where(d <= d_thr)[0] + 1
        if merge_inds.size > 0:
            merge_order = np.hstack([i, order[merge_inds]])
            merged_disks.append(np.mean(disks[merge_order], axis=0))
        else:
            merged_disks.append(disks[i])

        inds = np.where(d > d_thr)[0] + 1
        order = order[inds]
    merged_disks = np.vstack(merged_disks)

    return merged_disks


def poly_nms(polygons, threshold):
    assert isinstance(polygons, list)

    polygons = np.array(sorted(polygons, key=lambda x: x[-1]))

    keep_poly = []
    index = [i for i in range(polygons.shape[0])]

    while len(index) > 0:
        keep_poly.append(polygons[index[-1]].tolist())
        A = polygons[index[-1]][:-1]
        index = np.delete(index, -1)

        iou_list = np.zeros((len(index), ))
        for i in range(len(index)):
            B = polygons[index[i]][:-1]

            iou_list[i] = boundary_iou(A, B, 1)
        remove_index = np.where(iou_list > threshold)
        index = np.delete(index, remove_index)

    return keep_poly


def fourier2poly(fourier_coeff, num_reconstr_points=50):
    """ Inverse Fourier transform
        Args:
            fourier_coeff (ndarray): Fourier coefficients shaped (n, 2k+1),
                with n and k being candidates number and Fourier degree
                respectively.
            num_reconstr_points (int): Number of reconstructed polygon points.
        Returns:
            Polygons (ndarray): The reconstructed polygons shaped (n, n')
        """

    a = np.zeros((len(fourier_coeff), num_reconstr_points), dtype='complex')
    k = (len(fourier_coeff[0]) - 1) // 2

    a[:, 0:k + 1] = fourier_coeff[:, k:]
    a[:, -k:] = fourier_coeff[:, :k]

    poly_complex = ifft(a) * num_reconstr_points
    polygon = np.zeros((len(fourier_coeff), num_reconstr_points, 2))
    polygon[:, :, 0] = poly_complex.real
    polygon[:, :, 1] = poly_complex.imag
    return polygon.astype('int32').reshape((len(fourier_coeff), -1))


class Node:

    def __init__(self, ind):
        self.__ind = ind
        self.__links = set()

    @property
    def ind(self):
        return self.__ind

    @property
    def links(self):
        return set(self.__links)

    def add_link(self, link_node):
        self.__links.add(link_node)
        link_node.__links.add(self)


def graph_propagation(edges, scores, text_comps, edge_len_thr=50.):
    """Propagate edge score information and construct graph. This code was
    partially adapted from https://github.com/GXYM/DRRG licensed under the MIT
    license.

    Args:
        edges (ndarray): The edge array of shape N * 2, each row is a node
            index pair that makes up an edge in graph.
        scores (ndarray): The edge score array.
        text_comps (ndarray): The text components.
        edge_len_thr (float): The edge length threshold.

    Returns:
        vertices (list[Node]): The Nodes in graph.
        score_dict (dict): The edge score dict.
    """
    assert edges.ndim == 2
    assert edges.shape[1] == 2
    assert edges.shape[0] == scores.shape[0]
    assert text_comps.ndim == 2
    assert isinstance(edge_len_thr, float)

    edges = np.sort(edges, axis=1)
    score_dict = {}
    for i, edge in enumerate(edges):
        if text_comps is not None:
            box1 = text_comps[edge[0], :8].reshape(4, 2)
            box2 = text_comps[edge[1], :8].reshape(4, 2)
            center1 = np.mean(box1, axis=0)
            center2 = np.mean(box2, axis=0)
            distance = norm(center1 - center2)
            if distance > edge_len_thr:
                scores[i] = 0
        if (edge[0], edge[1]) in score_dict:
            score_dict[edge[0], edge[1]] = 0.5 * (
                score_dict[edge[0], edge[1]] + scores[i])
        else:
            score_dict[edge[0], edge[1]] = scores[i]

    nodes = np.sort(np.unique(edges.flatten()))
    mapping = -1 * np.ones((np.max(nodes) + 1), dtype=np.int)
    mapping[nodes] = np.arange(nodes.shape[0])
    order_inds = mapping[edges]
    vertices = [Node(node) for node in nodes]
    for ind in order_inds:
        vertices[ind[0]].add_link(vertices[ind[1]])

    return vertices, score_dict


def connected_components(nodes, score_dict, link_thr):
    """Conventional connected components searching. This code was partially
    adapted from https://github.com/GXYM/DRRG licensed under the MIT license.

    Args:
        nodes (list[Node]): The list of Node objects.
        score_dict (dict): The edge score dict.
        link_thr (float): The link threshold.

    Returns:
        clusters (List[list[Node]]): The clustered Node objects.
    """
    assert isinstance(nodes, list)
    assert all([isinstance(node, Node) for node in nodes])
    assert isinstance(score_dict, dict)
    assert isinstance(link_thr, float)

    clusters = []
    nodes = set(nodes)
    while nodes:
        node = nodes.pop()
        cluster = {node}
        node_queue = [node]
        while node_queue:
            node = node_queue.pop(0)
            neighbors = set([
                neighbor for neighbor in node.links if
                score_dict[tuple(sorted([node.ind, neighbor.ind]))] >= link_thr
            ])
            neighbors.difference_update(cluster)
            nodes.difference_update(neighbors)
            cluster.update(neighbors)
            node_queue.extend(neighbors)
        clusters.append(list(cluster))
    return clusters


def clusters2labels(clusters, num_nodes):
    """Convert clusters of Node to text component labels. This code was
    partially adapted from https://github.com/GXYM/DRRG licensed under the MIT
    license.

    Args:
        clusters (List[list[Node]]): The clusters of Node objects.
        num_nodes (int): The total node number of graphs in an image.

    Returns:
        node_labels (ndarray): The node label array.
    """
    assert isinstance(clusters, list)
    assert all([isinstance(cluster, list) for cluster in clusters])
    assert all(
        [isinstance(node, Node) for cluster in clusters for node in cluster])
    assert isinstance(num_nodes, int)

    node_labels = np.zeros(num_nodes)
    for cluster_ind, cluster in enumerate(clusters):
        for node in cluster:
            node_labels[node.ind] = cluster_ind
    return node_labels


def remove_single(text_comps, comp_pred_labels):
    """Remove isolated text components. This code was partially adapted from
    https://github.com/GXYM/DRRG licensed under the MIT license.

    Args:
        text_comps (ndarray): The text components.
        comp_pred_labels (ndarray): The clustering labels of text components.

    Returns:
        filtered_text_comps (ndarray): The text components with isolated ones
            removed.
        comp_pred_labels (ndarray): The clustering labels with labels of
            isolated text components removed.
    """
    assert text_comps.ndim == 2
    assert text_comps.shape[0] == comp_pred_labels.shape[0]

    single_flags = np.zeros_like(comp_pred_labels)
    pred_labels = np.unique(comp_pred_labels)
    for label in pred_labels:
        current_label_flag = (comp_pred_labels == label)
        if np.sum(current_label_flag) == 1:
            single_flags[np.where(current_label_flag)[0][0]] = 1
    keep_ind = [i for i in range(len(comp_pred_labels)) if not single_flags[i]]
    filtered_text_comps = text_comps[keep_ind, :]
    filtered_labels = comp_pred_labels[keep_ind]

    return filtered_text_comps, filtered_labels


def norm2(point1, point2):
    return ((point1[0] - point2[0])**2 + (point1[1] - point2[1])**2)**0.5


def min_connect_path(points):
    """Find the shortest path to traverse all points. This code was partially
    adapted from https://github.com/GXYM/DRRG licensed under the MIT license.

    Args:
        points(List[list[int]]): The point sequence [[x0, y0], [x1, y1], ...].

    Returns:
        shortest_path(List[list[int]]): The shortest index path.
    """
    assert isinstance(points, list)
    assert all([isinstance(point, list) for point in points])
    assert all([isinstance(coord, int) for point in points for coord in point])

    points_queue = points.copy()
    shortest_path = []
    current_edge = [[], []]

    edge_dict0 = {}
    edge_dict1 = {}
    current_edge[0] = points_queue[0]
    current_edge[1] = points_queue[0]
    points_queue.remove(points_queue[0])
    while points_queue:
        for point in points_queue:
            length0 = norm2(point, current_edge[0])
            edge_dict0[length0] = [point, current_edge[0]]
            length1 = norm2(current_edge[1], point)
            edge_dict1[length1] = [current_edge[1], point]
        key0 = min(edge_dict0.keys())
        key1 = min(edge_dict1.keys())

        if key0 <= key1:
            start = edge_dict0[key0][0]
            end = edge_dict0[key0][1]
            shortest_path.insert(0, [points.index(start), points.index(end)])
            points_queue.remove(start)
            current_edge[0] = start
        else:
            start = edge_dict1[key1][0]
            end = edge_dict1[key1][1]
            shortest_path.append([points.index(start), points.index(end)])
            points_queue.remove(end)
            current_edge[1] = end

        edge_dict0 = {}
        edge_dict1 = {}

    shortest_path = functools.reduce(operator.concat, shortest_path)
    shortest_path = sorted(set(shortest_path), key=shortest_path.index)

    return shortest_path


def in_contour(cont, point):
    x, y = point
    is_inner = cv2.pointPolygonTest(cont, (int(x), int(y)), False) > 0.5
    return is_inner


def fix_corner(top_line, bot_line, start_box, end_box):
    """Add corner points to predicted side lines. This code was partially
    adapted from https://github.com/GXYM/DRRG licensed under the MIT license.

    Args:
        top_line (List[list[int]]): The predicted top sidelines of text
            instance.
        bot_line (List[list[int]]): The predicted bottom sidelines of text
            instance.
        start_box (ndarray): The first text component box.
        end_box (ndarray): The last text component box.

    Returns:
        top_line (List[list[int]]): The top sidelines with corner point added.
        bot_line (List[list[int]]): The bottom sidelines with corner point
            added.
    """
    assert isinstance(top_line, list)
    assert all(isinstance(point, list) for point in top_line)
    assert isinstance(bot_line, list)
    assert all(isinstance(point, list) for point in bot_line)
    assert start_box.shape == end_box.shape == (4, 2)

    contour = np.array(top_line + bot_line[::-1])
    start_left_mid = (start_box[0] + start_box[3]) / 2
    start_right_mid = (start_box[1] + start_box[2]) / 2
    end_left_mid = (end_box[0] + end_box[3]) / 2
    end_right_mid = (end_box[1] + end_box[2]) / 2
    if not in_contour(contour, start_left_mid):
        top_line.insert(0, start_box[0].tolist())
        bot_line.insert(0, start_box[3].tolist())
    elif not in_contour(contour, start_right_mid):
        top_line.insert(0, start_box[1].tolist())
        bot_line.insert(0, start_box[2].tolist())
    if not in_contour(contour, end_left_mid):
        top_line.append(end_box[0].tolist())
        bot_line.append(end_box[3].tolist())
    elif not in_contour(contour, end_right_mid):
        top_line.append(end_box[1].tolist())
        bot_line.append(end_box[2].tolist())
    return top_line, bot_line


def comps2boundaries(text_comps, comp_pred_labels):
    """Construct text instance boundaries from clustered text components. This
    code was partially adapted from https://github.com/GXYM/DRRG licensed under
    the MIT license.

    Args:
        text_comps (ndarray): The text components.
        comp_pred_labels (ndarray): The clustering labels of text components.

    Returns:
        boundaries (List[list[float]]): The predicted boundaries of text
            instances.
    """
    assert text_comps.ndim == 2
    assert len(text_comps) == len(comp_pred_labels)
    boundaries = []
    if len(text_comps) < 1:
        return boundaries
    for cluster_ind in range(0, int(np.max(comp_pred_labels)) + 1):
        cluster_comp_inds = np.where(comp_pred_labels == cluster_ind)
        text_comp_boxes = text_comps[cluster_comp_inds, :8].reshape(
            (-1, 4, 2)).astype(np.int32)
        score = np.mean(text_comps[cluster_comp_inds, -1])

        if text_comp_boxes.shape[0] < 1:
            continue

        elif text_comp_boxes.shape[0] > 1:
            centers = np.mean(
                text_comp_boxes, axis=1).astype(np.int32).tolist()
            shortest_path = min_connect_path(centers)
            text_comp_boxes = text_comp_boxes[shortest_path]
            top_line = np.mean(
                text_comp_boxes[:, 0:2, :], axis=1).astype(np.int32).tolist()
            bot_line = np.mean(
                text_comp_boxes[:, 2:4, :], axis=1).astype(np.int32).tolist()
            top_line, bot_line = fix_corner(top_line, bot_line,
                                            text_comp_boxes[0],
                                            text_comp_boxes[-1])
            boundary_points = top_line + bot_line[::-1]

        else:
            top_line = text_comp_boxes[0, 0:2, :].astype(np.int32).tolist()
            bot_line = text_comp_boxes[0, 2:4:-1, :].astype(np.int32).tolist()
            boundary_points = top_line + bot_line

        boundary = [p for coord in boundary_points for p in coord] + [score]
        boundaries.append(boundary)

    return boundaries
