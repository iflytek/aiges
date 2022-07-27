#ifndef __DSP_TYPE_H__
#define __DSP_TYPE_H__


/**
 * C/C++ wchar_t support
 */
#ifdef __cplusplus
# include <cwchar>
#else  /* c */
# include <wchar.h>
#endif /* wchar */


typedef enum{
    CTMeterCustom =   0,      // 自定义计量接口
    CTMetricsLog  =   1,      // 自定义metrics日志接口
    CTTraceLog    =   2,      // 自定义trace日志接口
} CtrlType;

typedef enum{
    DataText    =   0,      // 文本数据
    DataAudio   =   1,      // 音频数据
    DataImage   =   2,      // 图像数据
    DataVideo   =   3,      // 视频数据
    DataPer     =   4,      // 个性化数据
} DataType;

typedef enum{
    DataBegin   =   0,      // 首数据
    DataContinue =  1,      // 中间数据
    DataEnd     =   2,      // 尾数据
    DataOnce    =   3,      // 非会话单次输入输出
} DataStatus;

typedef struct ParamList{
    char* key;
    char* value;
    unsigned int vlen;
    struct ParamList* next;
}* pParamList, *pConfig, *pDescList;     // 配置对复用该结构定义

typedef struct DataList{
    void*   data;           // 数据实体
    unsigned int len;       // 数据长度
    int status;      // 数据状态
    int duration;      // 数据状态
    struct DataList* next;  // 链表指针
}*  pDataList;

typedef struct Framedata {
	void* data;
	int status; // 1:contine, 2:end
	int duration;
	struct Framedata* next;  // 链表指针
}* fDataList;

typedef struct AudioTag {
	char* encode;
	int rate;
	int spx;
}* pAudioTag;

typedef struct AVStreamTag {
    char* encoding;
    unsigned int quality;
    unsigned int sample_rate;
    unsigned int channels;
    unsigned int bit_depth;
}AVStreamTag, * pAVStreamTag;


typedef struct NaluNode {
    int offset;
    int start_code_len;
    int size;
    int nalu_type;
    struct NaluNode* next;
}NaluNode, * pNaluList;

typedef struct AudioFrameList{
    unsigned int start;       // 开始位置
    unsigned int end;       // 结束位置
    unsigned int len;       // 数据长度
    int status;      // 数据状态
    int duration;      // 数据时长 ms
    struct DataList* next;  // 链表指针
}* pAudioFrameList;

typedef struct Tag {
    int a;
    int b;
}* pTag;

typedef struct InfoMp4Decoder {
    int code;
    int video_encoding; // 0-pcm
    int video_frame_rate;
    int video_height;
    int video_width;
    int audio_encoding; // 0-h264
    int audio_sample_rate;
    int audio_channels;
    int audio_bitdepth;

//    struct InfoMp4Decoder* next;
}InfoMp4Decoder, * pInfoMp4DecoderList;


#endif
