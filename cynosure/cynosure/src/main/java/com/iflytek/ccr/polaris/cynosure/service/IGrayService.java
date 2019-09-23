package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.AddGrayGroupAndConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;

import java.util.List;

/**
 * Created by DELL-5490 on 2018/7/4.
 */
public interface IGrayService {
    /**
     * 新增灰度组（不含文件拖拽上传）
     *
     * @param body
     * @return
     */
    Response<AddGrayGroupAndConfigResponseBody> add(AddGrayGroupRequestBody body);

    /**
     * 新增灰度组（含文件拖拽上传）
     *
     * @param body
     * @return
     */
    Response<AddGrayGroupAndConfigResponseBody> addAndFile(AddGrayGroupRequestBody body, List<FileContent> fileContentList);

    /**
     * 查询灰度组列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(QueryGrayGroupListRequestBody body);

    /**
     * 根据灰度组id查询灰度组详情
     *
     * @param id
     * @return
     */
    Response<AddGrayGroupAndConfigResponseBody> findById(String id);

    /**
     * 删除灰度组
     *
     * @param body
     * @return
     */
    Response<String> delete(DeleteGrayGroupRequestBody body);

    /**
     * 编辑灰度组
     *
     * @param body
     * @return
     */
    Response<GrayGroupDetailResponseBody> edit(EditGrayGroupRequestBody body);

    /**
     * 查询服务定阅列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> consumer(QueryCustomDetailRequestBody body);

}