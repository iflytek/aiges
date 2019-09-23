/** 
 * @file    meteor_error.h
 * @brief   个性化存储平台 2.0 (meteor) 错误定义
 * 
 * 本文件定义了科大讯飞云平台研发部云计算组个性化存储平台2.0的错误代码及其文字说明
 * 
 * @author  mingzhang2
 * @version 1.0
 * @date    2013-11-15
 * 
 * @see
 * 
 * <b>History:</b><br>
 * <table>
 *  <tr> <th>Version    <th>Date        <th>Author       <th>Notes</tr>
 *  <tr> <td>2.0.11.011 <td>2010-11-15  <td>mingzhang2   <td>Create this file</tr>
 * </table>
 * 
 */
#ifndef __METEOR_ERROR_H__
#define __METEOR_ERROR_H__

enum
{
    MTR_SUCCESS             =  0,
    MTR_ERROR_FATAL         = -1,
    MTR_ERROR_EXCEPTION     = -2,

    MTR_ERROR_BASE          = 30000,

    MTR_ERROR_UNKNOWN       = 30001,
    MTR_ALREADY_INITILIZED  = 30002,
    MTR_NOT_INITILIZED      = 30003,
    MTR_CFG_ERROR           = 30004,

    MTR_DATA_NOT_FOUND      = 30102,
    MTR_INVALID_PARA        = 30106,
    MTR_INVALID_PARA_VALUE  = 30107,
    MTR_INVALID_HANDLE      = 30108,
    MTR_INVALID_DATA        = 30109,
    
    MTR_TIME_OUT            = 30114,
    MTR_NO_RESPONSE         = 30115,
    MTR_NO_ENOUGH_BUFFER	= 30117,

    MTR_CACHE_FULL          = 30134,
    MTR_LOCAL_FILE_ERROR    = 30136,

    /* Network Error 30200 */
    MTR_NET_GENERAL         = 30200,
    MTR_NET_OPENSOCK        = 30201,
    MTR_NET_CONNRUNOUT      = 30215,

    /* ERROR FROM SERVER */
    MTR_SERVER_ERROR_BASE   = 33000,
    MTR_SERVER_BUILDMESSAGE = 33001,
    MTR_SERVER_NULL_MSG     = 33002,
    MTR_SERVER_PUT          = 33003,
    MTR_TABLE_NOT_EXIST     = 33004,
};

#endif /* __METEOR_ERROR_H__ */
