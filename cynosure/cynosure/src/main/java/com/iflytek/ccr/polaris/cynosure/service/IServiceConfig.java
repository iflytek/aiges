package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.domain.DownloadFile;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.ServiceGrayConfigHistoryListRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.GrayConfigListDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.PushServiceConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import org.springframework.http.ResponseEntity;

import javax.servlet.http.HttpServletResponse;
import java.io.UnsupportedEncodingException;
import java.util.List;

/**
 * 配置服务业务逻辑接口
 *
 * @author sctang2
 * @create 2017-11-21 11:46
 **/
public interface IServiceConfig {
    /**
     * 查询最近的配置列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findLastestList(QueryServiceConfigRequestBody body);

    /**
     * 编辑配置
     *
     * @param body
     * @return
     */
    Response<ServiceConfigDetailResponseBody> edit(EditServiceConfigRequestBody body);

    /**
     * 查询配置明细
     *
     * @param id
     * @return
     */
    Response<ServiceConfigDetailResponseBody> find(String id);

    /**
     * 推送
     *
     * @param body
     * @return
     */
    Response<PushServiceConfigResponseBody> push(PushServiceConfigRequestBody body);

    /**
     * 批量推送
     *
     * @param body
     * @return
     */
    Response<PushServiceConfigResponseBody> batchPush(BatchPushServiceConfigRequestBody body);

    /**
     * 删除配置
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);

    /**
     * 批量删除配置
     *
     * @param body
     * @return
     */
    Response<String> batchDelete(IdsRequestBody body);

    /**
     * 查询配置历史列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findServiceConfigHistoryList(ServiceConfigHistoryListRequestBody body);

    /**
     * 回滚配置
     *
     * @param body
     * @return
     */
    Response<List<ServiceConfigDetailResponseBody>> rollback(IdsRequestBody body);

    /**
     * 更新反馈
     *
     * @param body
     * @return
     */
    Response<String> feedback(ServiceConfigFeedBackRequestBody body);

    /**
     * 通过ids查询配置服务列表
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findListByIds(List<String> ids);


    /**
     * 新增灰度配置
     *
     * @param body
     * @param fileContentList
     * @return
     */
    Response<GrayConfigListDetailResponseBody> addGrayConfig(AddGrayConfigRequestBody body, List<FileContent> fileContentList) throws UnsupportedEncodingException;

    /**
     * 查询灰度配置历史列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findServiceGrayConfigHistoryList(ServiceGrayConfigHistoryListRequestBody body);


    /**
     * 下载配置文件
     * @param body
     * @return
     */
    Response<DownloadFile> download(DownloadServiceConfigRequestBody body);
}
