package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.CopyServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.QueryServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IServiceVersion;
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
@RequestMapping(Constant.API + "/{version}/service/version")
public class ServiceVersionController {
    @Autowired
    private IServiceVersion serviceVersionImpl;

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
        return this.serviceVersionImpl.findLastestList(body);
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
        return this.serviceVersionImpl.findList(body);
    }

    /**
     * 查询版本详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<ServiceVersionDetailResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.serviceVersionImpl.find(id);
    }

    /**
     * 新增版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public Response<ServiceVersionDetailResponseBody> add(@Validated @RequestBody AddServiceVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceVersionImpl.add(body);
    }

    /**
     * 编辑版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<ServiceVersionDetailResponseBody> edit(@Validated @RequestBody EditServiceVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceVersionImpl.edit(body);
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
        return this.serviceVersionImpl.delete(body);
    }

    /**
     * 复制版本
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/copy", method = RequestMethod.POST)
    public Response<ServiceVersionDetailResponseBody> copy(@Validated @RequestBody CopyServiceVersionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceVersionImpl.copy(body);
    }
}
