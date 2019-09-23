package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.AddServiceApiVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.QueryServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceApiVersionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IServiceApiVersion;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 服务版本控制器
 *
 * @author sctang2
 * @create 2017-11-18 10:23
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/service/apiVersion")
public class ServiceApiVersionController {
    @Autowired
    private IServiceApiVersion serviceApiVersionImpl;

    /**
     * 查询最近的版本列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/lastestList", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> lastestList(@Validated QueryServiceVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceApiVersionImpl.findLastestList(body);
    }

    /**
     * 查询服务版本列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> list(@Validated QueryServiceVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceApiVersionImpl.findList(body);
    }

    /**
     * 查询版本详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<ServiceApiVersionDetailResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.serviceApiVersionImpl.find(id);
    }

    /**
     * 新增版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public Response<ServiceApiVersionDetailResponseBody> add(@Validated @RequestBody AddServiceApiVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceApiVersionImpl.add(body);
    }

    /**
     * 编辑版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<ServiceApiVersionDetailResponseBody> edit(@Validated @RequestBody EditServiceVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceApiVersionImpl.edit(body);
    }

    /**
     * 删除版本
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
        return this.serviceApiVersionImpl.delete(body);
    }
}
