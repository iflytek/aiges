//
// Created by jlpan4 on 2021/4/27.
//

#ifndef H264BITSTREAM_H264_NALU_SPLIT_H
#define H264BITSTREAM_H264_NALU_SPLIT_H

#include "type.h"

pNaluList get_h264_nalu(uint8_t* buf, size_t sz);

#endif //H264BITSTREAM_H264_NALU_SPLIT_H
