package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectMemberCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IUserCondition;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.ProjectTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.domain.ProjectMember;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.QueryProjectMemberResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IProjectService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.*;

/**
 * 项目业务逻辑接口实现
 *
 * @author sctang2
 * @create 2017-11-19 21:03
 **/
@Service
public class ProjectServiceImpl extends BaseService implements IProjectService {
    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IProjectMemberCondition projectMemberConditionImpl;

    @Autowired
    private IUserCondition userConditionImpl;

    @Autowired
    private ProjectTransactional projectTransactionalImpl;

    @Override
    public Response<QueryPagingListResponseBody> findList(QueryProjectRequestBody body) {
        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        String userId = this.getUserId();
        map.put("userId", userId);
        String name = body.getName();
        if (StringUtils.isNotBlank(name)) {
            map.put("name", name);
        }

        //查询总数
        int totalCount = this.projectConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ProjectDetailResponseBody> list = new ArrayList<>();
        Optional<List<Project>> projectList = Optional.ofNullable(this.projectConditionImpl.findList(map));
        projectList.ifPresent(x -> {
            x.forEach(y -> {
                //创建项目结果
                ProjectDetailResponseBody projectDetail = this.createProjectResult(y);
                list.add(projectDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<ProjectDetailResponseBody> find(String id) {
        //通过id查询项目信息
        Project project = this.projectConditionImpl.findById(id);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //创建项目结果
        ProjectDetailResponseBody result = this.createProjectResult(project);
        return new Response<>(result);
    }

    @Override
    public Response<ProjectDetailResponseBody> add(AddProjectRequestBody body) {
        //通过项目名称查询项目信息
        String name = body.getName();
        Project project = this.projectConditionImpl.findByName(name);
        if (null != project) {
            //已存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_EXISTS, SystemErrCode.ERRMSG_PROJECT_EXISTS);
        }

        //新增项目
        Project newProject = this.projectTransactionalImpl.add(body);

        //创建项目结果
        ProjectDetailResponseBody result = this.createProjectResult(newProject);
        return new Response<>(result);
    }

    @Override
    public Response<ProjectDetailResponseBody> edit(EditProjectRequestBody body) {
        //通过id查询项目信息
        String id = body.getId();
        Project project = this.projectConditionImpl.findById(id);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //根据id更新项目
        Project updateProject = this.projectConditionImpl.updateById(id, body);
        updateProject.setName(project.getName());
        updateProject.setCreateTime(project.getCreateTime());

        //创建项目结果
        ProjectDetailResponseBody result = this.createProjectResult(updateProject);
        return new Response<>(result);
    }

    @Override
    public Response<String> delete(IdRequestBody body) {
        //通过id查询集群和成员列表
        String id = body.getId();
        Project project = this.projectConditionImpl.findProjectAndClusterListById(id);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //检查是否创建集群
        List<Cluster> clusterList = project.getClusterList();
        if (null != clusterList && !clusterList.isEmpty()) {
            //该用户已创建集群
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_CREATE, SystemErrCode.ERRMSG_CLUSTER_CREATE);
        }

        //通过id删除项目
        this.projectTransactionalImpl.delete(id);
        return new Response<>(null);
    }

    @Override
    public Response<QueryProjectMemberResponseBody> addMember(AddProjectMemberRequestBody body) {
        //通过id查询项目信息
        String id = body.getId();
        Project project = this.projectConditionImpl.findById(id);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //通过账号查询用户信息
        String account = body.getAccount();
        User user = this.userConditionImpl.findByAccount(account);
        if (null == user) {
            //不存在该用户
            return new Response<>(SystemErrCode.ERRCODE_USER_NOT_EXISTS, SystemErrCode.ERRMSG_USER_NOT_EXISTS);
        }

        //根据userId，projectId查询用户所在的项目
        String memberId = user.getId();
        ProjectMember projectMember = this.projectMemberConditionImpl.find(memberId, id);
        if (null != projectMember) {
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_MEMBER_EXIEST, SystemErrCode.ERRMSG_PROJECT_MEMBER_EXIEST);
        }

        //新增项目成员
        ProjectMember newProjectMember = this.projectMemberConditionImpl.add(memberId, id);

        //创建项目成员
        QueryProjectMemberResponseBody result = this.createProjectMember(newProjectMember, memberId, account);
        return new Response<>(result);
    }

    @Override
    public Response<String> deleteMember(DeleteProjectMemberRequestBody body) {
        //通过id查询项目信息
        String id = body.getId();
        Project project = this.projectConditionImpl.findById(id);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //根据userId，projectId删除项目成员
        String memberId = body.getUserId();
        int success = this.projectMemberConditionImpl.delete(memberId, id);
        if (success <= 0) {
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_MEMBER_NOT_EXIEST, SystemErrCode.ERRMSG_PROJECT_MEMBER_NOT_EXIEST);
        }
        return new Response<>(null);
    }

    @Override
    public Response<QueryPagingListResponseBody> findMemberList(QueryProjectMemberRequestBody body) {
        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        String id = body.getId();
        map.put("projectId", id);

        //查询总数
        int totalCount = this.projectMemberConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<QueryProjectMemberResponseBody> list = new ArrayList<>();
        Optional<List<ProjectMember>> projectMemberList = Optional.ofNullable(this.projectMemberConditionImpl.findList(map));
        projectMemberList.ifPresent(x -> {
            x.forEach(y -> {
                //创建项目成员
                String userId = y.getUserId();
                String account = y.getUser().getAccount();
                QueryProjectMemberResponseBody projectMember = this.createProjectMember(y, userId, account);
                list.add(projectMember);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    /**
     * 创建项目成员
     *
     * @param projectMember
     * @param id
     * @param account
     * @return
     */
    private QueryProjectMemberResponseBody createProjectMember(ProjectMember projectMember, String id, String account) {
        QueryProjectMemberResponseBody result = new QueryProjectMemberResponseBody();
        result.setId(id);
        result.setAccount(account);
        result.setCreateTime(projectMember.getCreateTime());
        return result;
    }
}