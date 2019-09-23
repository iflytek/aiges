package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.annotation.Access;
import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.QueryProjectMemberResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IProjectService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 项目控 制器
 *
 * @author sctang2
 * @create 2017-11-19 21:04
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/project")
public class ProjectController {
    @Autowired
    private IProjectService projectServiceImpl;

    /**
     * 查询项目列表
     *
     * @param body
     * @return
     */
    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> list(QueryProjectRequestBody body) {
        return this.projectServiceImpl.findList(body);
    }

    /**
     * 查询项目详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<ProjectDetailResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.projectServiceImpl.find(id);
    }

    /**
     * 新增项目
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<ProjectDetailResponseBody> add(@Validated @RequestBody AddProjectRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.projectServiceImpl.add(body);
    }

    /**
     * 编辑项目
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<ProjectDetailResponseBody> edit(@Validated @RequestBody EditProjectRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.projectServiceImpl.edit(body);
    }

    /**
     * 删除项目
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/delete", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<String> delete(@Validated @RequestBody IdRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.projectServiceImpl.delete(body);
    }

    /**
     * 新增项目成员
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/member/add", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<QueryProjectMemberResponseBody> addMember(@Validated @RequestBody AddProjectMemberRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.projectServiceImpl.addMember(body);
    }

    /**
     * 删除项目成员
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/member/delete", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<String> deleteMember(@Validated @RequestBody DeleteProjectMemberRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.projectServiceImpl.deleteMember(body);
    }

    /**
     * 查询项目成员列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/member/list", method = RequestMethod.GET)
    @Access(authorities = "admin")
    public Response<QueryPagingListResponseBody> findMemberList(@Validated QueryProjectMemberRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.projectServiceImpl.findMemberList(body);
    }
}
