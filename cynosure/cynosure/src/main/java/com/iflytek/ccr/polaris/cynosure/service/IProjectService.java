package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.QueryProjectMemberResponseBody;

/**
 * 项目业务逻辑接口
 *
 * @author sctang2
 * @create 2017-11-19 21:03
 **/
public interface IProjectService {
    /**
     * 查询项目列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(QueryProjectRequestBody body);

    /**
     * 查询项目详情
     *
     * @param id
     * @return
     */
    Response<ProjectDetailResponseBody> find(String id);

    /**
     * 新增项目
     *
     * @param body
     * @return
     */
    Response<ProjectDetailResponseBody> add(AddProjectRequestBody body);

    /**
     * 编辑项目
     *
     * @param body
     * @return
     */
    Response<ProjectDetailResponseBody> edit(EditProjectRequestBody body);

    /**
     * 删除项目
     *
     * @param body
     * @return
     */
    Response<String> delete(IdRequestBody body);

    /**
     * 新增项目成员
     *
     * @param body
     * @return
     */
    Response<QueryProjectMemberResponseBody> addMember(AddProjectMemberRequestBody body);

    /**
     * 删除项目成员
     *
     * @param body
     * @return
     */
    Response<String> deleteMember(DeleteProjectMemberRequestBody body);

    /**
     * 查询项目成员列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findMemberList(QueryProjectMemberRequestBody body);
}
