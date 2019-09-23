package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.CopyServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.EditServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.QueryServiceListRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 服务控制器
 *
 * @author sctang2
 * @create 2017-11-16 19:53
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/service")
public class ServiceController {
    @Autowired
    private IService serviceImpl;

    /**
     * 查询最近的服务列表
     *
     * @param body
     * @return
     */
    @RequestMapping(value = "/lastestList", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> lastestList(QueryServiceListRequestBody body) {
        return this.serviceImpl.findLastestList(body);
    }

    /**
     * 查询服务列表
     *
     * @param body
     * @return
     */
    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> list(QueryServiceListRequestBody body) {
        return this.serviceImpl.findList(body);
    }

    /**
     * 查询服务详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<ServiceDetailResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.serviceImpl.find(id);
    }

    /**
     * 新增服务
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public Response<ServiceDetailResponseBody> add(@Validated @RequestBody AddServiceRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceImpl.add(body);
    }

    /**
     * 编辑服务
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<ServiceDetailResponseBody> edit(@Validated @RequestBody EditServiceRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceImpl.edit(body);
    }

    /**
     * 删除服务
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/delete", method = RequestMethod.POST)
    public Response<String> delete(@Validated @RequestBody IdRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceImpl.delete(body);
    }

    /**
     * 复制服务
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/copy", method = RequestMethod.POST)
    public Response<ServiceDetailResponseBody> copy(@Validated @RequestBody CopyServiceRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceImpl.copy(body);
    }
}
