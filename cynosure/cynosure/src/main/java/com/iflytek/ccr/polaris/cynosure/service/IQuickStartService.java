package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceConfigRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddServiceConfigResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddServiceResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddVersionResponseBodyByQuickStart;

import java.util.List;

/**
 * 快速创建业务逻辑接口
 *
 * @author sctang2
 * @create 2018-01-29 14:03
 **/
public interface IQuickStartService {
    /**
     * 新增服务
     *
     * @param body
     * @return
     */
    Response<AddServiceResponseBodyByQuickStart> addService(AddServiceRequestBodyByQuickStart body);

    /**
     * 新增服务
     *
     * @param body
     * @return
     */
    Response<AddVersionResponseBodyByQuickStart> addVersion(AddVersionRequestBodyByQuickStart body);

    /**
     * 新增服务版本(不含配置文件拖拽上传)
     *
     * @param body
     * @return
     */
    Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersion(AddServiceVersionRequestBodyByQuickStart body);

    /**
     * 新增服务版本（含配置文件拖拽上传）
     *
     * @param body
     * @return
     */
    Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersionAndFile(AddServiceVersionRequestBodyByQuickStart body, List<FileContent> fileContentList);

    /**
     * 新增服务配置
     *
     * @param body
     * @param fileContentList
     * @return
     */
    Response<AddServiceConfigResponseBodyByQuickStart> addServiceConfig(AddServiceConfigRequestBodyByQuickStart body, List<FileContent> fileContentList);

    /**
     * 快速查询
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(BaseRequestBody body);

    /**
     * 快速查询1
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList1(BaseRequestBody body);

}
