#ifndef __METEOR_AHSC_H__
#define __METEOR_AHSC_H__

#ifndef MTR_API
# ifdef WIN32
#  define MTR_API __declspec(dllexport)
# else
#  define MTR_API __attribute__ ((visibility ("default")))
# endif
#endif

/**
 * @brief   SyncState
 *          数据同步的参数值.
 */

 #define SyncState int
//enum SyncState {
//    /* 发送数据到数据同步组件 */
//    SYNC = 0,
//    /* 不发送同步数据到数据同步组件 */
//    SYNC_NO = 1
//};

/**
 * @brief   HbaseKeys
 *          用于定位HBase存储单元的key
 */
typedef struct HbaseKeys {
    /* HBase中的表名 */
    const char *table_;
    /* HBase表中的行键(row key) */
    const char *row_key_;
    /**
     * HBase表中的列族(family), 为节省存储空间, 尽量使用单个字符, 
     * 如: 个性化表中列族为"p", 由业务层在建表时约定 
     */
    const char *family_;
    /**
     * HBase表中列族之下的列(qualifer), 每个列族下可支持多个列, 
     * 不同行的多个列可存储为稀疏矩阵, 不推荐每一行均使用不同的列
     */
    const char *qualifier_;
} HbaseKeys;

#ifdef __cplusplus
extern  "C" {
#endif

/** 
 * @brief   ahsc_initialize
 *          初始化系统的第一个函数，仅需要调用一次
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const char *config   [in]
 *          接口传入的配置文件路径, 支持绝对路径和相对路径
 * @param   void* reserved       [in]
 *          保留以后使用
 * @see		ahsc_uninitialize
 */
MTR_API int ahsc_initialize(const char *config, void *reserved);
typedef int (*Proc_ahsc_initialize)(const char *config, void *reserved);

/** 
 * @brief   ahsc_uninitialize
 *          初始化系统的第一个函数，仅需要调用一次
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   void* reserved       [in]
 *          保留以后使用
 * @see		ahsc_initialize
 */
MTR_API int ahsc_uninitialize(void *reserved);
typedef int (*Proc_ahsc_uninitialize)(void *reserved);

/** 
 * @brief   ahsc_timeout
 *          设置接口调用超时时间, 以毫秒为单位, 缺省超时时间为3000ms, 可重复设置,
 *          以最后一次设置的超时时间为准
 * @return  void
 *          无调用失败情况
 * @param   int timeout          [in]
 *          接口调用超时时间, 以毫秒为单位
 * @see		ahsc_initialize
 */
MTR_API void ahsc_timeout(int timeout);
typedef void (*Proc_ahsc_timeout)(int timeout);

/** 
 * @brief   ahsc_put
 *          将数据写入指定的HBase数据单元, 同步接口, 完成写入后返回
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const void *buffer                 [in]
 *          用于存储数据的缓冲区
 * @param   int length                         [in]
 *          缓冲区长度
 * @param   SyncState sync                     [in]
 *          是否进行数据同步
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_put_async
 */
MTR_API int ahsc_put(const HbaseKeys *keys, const void *buffer, int length, 
    SyncState sync, const char *params, void *reserved);
typedef int (*Proc_ahsc_put)(const HbaseKeys *keys, const void *buffer, 
    int length, SyncState sync, const char *params, void *reserved);

/** 
 * @brief   Proc_ahsc_put_callback
 *          ahsc_put_async接口的回掉函数指针类型
 * @return  void
 *          无
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const void *buffer                 [in]
 *          用于存储数据的缓冲区
 * @param   int length                         [in]
 *          缓冲区长度
 * @param   int ret                            [in]
 *          异步执行数据写入时的返回值
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_put_async
 */
// TODO
typedef void (*Proc_ahsc_put_callback)(const HbaseKeys *keys, 
    const char *params, int ret, void *reserved);

/** 
 * @brief   ahsc_put_async
 *          将数据写入指定的HBase数据单元, 异步接口, 数据写入缓冲队列后返回,
 *          后台线程执行数据写入后调用指定的回掉函数
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const void *buffer                 [in]
 *          用于存储数据的缓冲区
 * @param   int length                         [in]
 *          缓冲区长度
 * @param   SyncState sync                     [in]
 *          是否进行数据同步
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   Proc_ahsc_put_callback callback    [in]
 *          后台线程执行数据写入后调用的回掉函数指针
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_put
 */
MTR_API int ahsc_put_async(const HbaseKeys *keys, const void *buffer, 
    int length, SyncState sync, Proc_ahsc_put_callback callback, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_put_async)(const HbaseKeys *keys, const void *buffer, 
    int length, SyncState sync, Proc_ahsc_put_callback callback, 
    const char *params, void *reserved);

/** 
 * @brief   ahsc_get
 *          从指定的HBase数据单元获取数据
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const void *buffer                 [out]
 *          用于存储数据的缓冲区
 * @param   int *length                        [in/out]
 *          传入缓冲区长度, 传出实际数据大小
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_get_async
 */
MTR_API int ahsc_get(const HbaseKeys *keys, void **buffer, int *length, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_get)(const HbaseKeys *keys, void **buffer, 
    int *length, const char *params, void *reserved);

/** 
 * @brief   ahsc_get_free_buffer
 *          释放由ahsc_get生成的buffer
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   void *buffer                       [out]
 *          用于存储数据的缓冲区
 * @param   int *length                        [in/out]
 *          传入缓冲区长度, 传出实际数据大小
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_get_async
 */
MTR_API int ahsc_get_free_buffer(void **buffer, void *reserved);
typedef int (*Proc_ahsc_get_free_buffer)(void **buffer, void *reserved);

/** 
 * @brief   Proc_ahsc_get_callback
 *          ahsc_get_async接口的回掉函数指针类型
 * @return  void
 *          无
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   int ret                            [in]
 *          异步执行数据写入时的返回值
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_put_async
 */
typedef void (*Proc_ahsc_get_callback)(const HbaseKeys *keys, 
    const char *params, const void *buffer, int length, int ret, void *reserved);

/** 
 * @brief   ahsc_get_async
 *          从指定的HBase数据单元获取数据
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   Proc_ahsc_get_callback callback    [in]
 *          后台线程执行数据写入后调用的回掉函数指针
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_get_async
 */
MTR_API int ahsc_get_async(const HbaseKeys *keys, Proc_ahsc_get_callback callback, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_get_async)(const HbaseKeys *keys, Proc_ahsc_get_callback callback, 
    const char *params, void *reserved);

/** 
 * @brief   ahsc_delete
 *          删除指定HBase Cell数据
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   SyncState sync                     [in]
 *          是否同步删除数据
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see		ahsc_initialize
 */
MTR_API int ahsc_delete(const HbaseKeys *keys, SyncState sync, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_delete)(const HbaseKeys *keys, SyncState sync, 
    const char *params, void *reserved);

/** 
 * @brief   Proc_ahsc_delete_callback
 *          ahsc_delete_async接口的回掉函数指针类型
 * @return  void
 *          无
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   int ret                            [in]
 *          异步执行数据删除时的返回值
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see     ahsc_initialize, ahsc_put_async
 */
typedef void (*Proc_ahsc_delete_callback)(const HbaseKeys *keys, 
    const char *params, int ret, void *reserved);

/** 
 * @brief   ahsc_delete_async
 *          删除指定的HBase数据单元, 异步接口, 数据删除指令写入缓冲队列后返回,
 *          后台线程执行删除操作后调用指定的回掉函数
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   SyncState sync                     [in]
 *          是否同步删除数据
 * @param   Proc_ahsc_delete_callback callback[in]
 *          后台线程执行数据写入后调用的回掉函数指针
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see		ahsc_initialize
 */
MTR_API int ahsc_delete_async(const HbaseKeys *keys, SyncState sync, 
    Proc_ahsc_delete_callback callback, const char *params, void *reserved);
typedef int (*Proc_ahsc_delete_async)(const HbaseKeys *keys, SyncState sync, 
    Proc_ahsc_delete_callback callback, const char *params, void *reserved);

/** 
 * @brief   ahsc_modify_time
 *          HBase数据单元的修改时间
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   long long *timestamp               [out]
 *          用于存储时间戳的64位整型值
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see		ahsc_initialize
 */
MTR_API int ahsc_modify_time(const HbaseKeys *keys, long long *timestamp, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_modify_time)(const HbaseKeys *keys, long long *timestamp, 
    const char *params, void *reserved);

/** 
 * @brief   ahsc_exists
 *          判断指定的HBase单元是否存在
 * @return  int
 *          数据存在返回0, 数据不存在或调用失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see		ahsc_initialize
 */
MTR_API int ahsc_exists(const HbaseKeys *keys, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_exists)(const HbaseKeys *keys, 
    const char *params, void *reserved);

/** 
 * @brief   ahsc_length
 *          获取指定的HBase单元的长度, 实际会调用一次get, 不推荐先使用ahsc_length, 
 *          再调用ahsc_get, 会造成两倍的耗时
 * @return  int
 *          调用成功返回0, 失败返回相应错误码
 * @param   const HbaseKeys *keys              [in]
 *          用于定位HBase存储单元的key
 * @param   int *length                        [out]
 *          用于存储数据长度的int型指针
 * @param   const char *params                 [in]
 *          用于以key1=value1,key2=value2...的形式传入与存储服务无关的参数, 如会话ID等
 * @param   void *reserved                     [in]
 *          保留以后使用
 * @see		ahsc_initialize
 */
MTR_API int ahsc_length(const HbaseKeys *keys, int *length, 
    const char *params, void *reserved);
typedef int (*Proc_ahsc_length)(const HbaseKeys *keys, int *length, 
    const char *params, void *reserved);

#ifdef __cplusplus
}
#endif

#endif /* __METEOR_AHSC_H__ */
