#ifndef __AIGES_TYPE_H__
#define __AIGES_TYPE_H__


///**
// * C/C++ wchar_t support
// */
//#ifdef __cplusplus
//# include <cwchar>
//#else  /* c */
//# include <wchar.h>
//#endif /* wchar */


///**
// * Wrapper API type
// */
//#if defined(_MSC_VER)            /* Microsoft Visual C++ */
//  #if !defined(WrapperAPI)
//    #define WrapperAPI __stdcall
//  #endif
//  #pragma pack(push, 8)
//#elif defined(__BORLANDC__)      /* Borland C++ */
//  #if !defined(WrapperAPI)
//    #define WrapperAPI __stdcall
//  #endif
//  #pragma option -a8
//#elif defined(__WATCOMC__)       /* Watcom C++ */
//  #if !defined(WrapperAPI)
//    #define WrapperAPI __stdcall
//  #endif
//  #pragma pack(push, 8)
//#else                            /* Any other including Unix */
//  #if !defined(WrapperAPI)
    #define WrapperAPI __attribute__ ((visibility("default")))
//  #endif
//#endif


///**
// * True and false
// */
//#ifndef FALSE
//#define FALSE		0
//#endif	/* FALSE */
//
//#ifndef TRUE
//#define TRUE		1
//#endif	/* TRUE */

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
    char*   key;            // 数据标识
    void*   data;           // 数据实体
    unsigned int len;       // 数据长度
    DataType    type;       // 数据类型
    DataStatus status;      // 数据状态
    pDescList desc;         // 数据描述
    struct DataList* next;  // 链表指针
}*  pDataList;


///* Reset the structure packing alignments for different compilers. */
//#if defined(_MSC_VER)            /* Microsoft Visual C++ */
//  #pragma pack(pop)
//#elif defined(__BORLANDC__)      /* Borland C++ */
//  #pragma option -a.
//#elif defined(__WATCOMC__)       /* Watcom C++ */
//  #pragma pack(pop)
//#endif

#endif