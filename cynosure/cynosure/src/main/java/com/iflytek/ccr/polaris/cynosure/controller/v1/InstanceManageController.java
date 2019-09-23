package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.AddGrayGroupInstanceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.EditInstanceRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IGrayService;
import com.iflytek.ccr.polaris.cynosure.service.InstanceService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * Created by DELL-5490 on 2018/7/7.
 */
@RestController
@RequestMapping(Constant.API + "/{version}/instanceManage")
public class InstanceManageController {

    @Autowired
    private IGrayService grayServiceImpl;

    @Autowired
    private InstanceService instanceServiceImpl;

    /**
     * 查询推送实例详情
     *
     * @param grayId
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<String> find(@Validated @RequestParam(value = "grayId") String grayId) {
        if (StringUtils.isEmpty(grayId)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.instanceServiceImpl.findById(grayId);
    }

    /**
     * 查询全部推送实例
     *
     * @return
     */
    @RequestMapping(value = "/List", method = RequestMethod.GET)
    public Response<String> list() {
        return this.instanceServiceImpl.findList(null);
    }

    /**
     * 编辑推送实例
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<GrayGroupDetailResponseBody> edit(@Validated @RequestBody EditInstanceRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.instanceServiceImpl.edit(body);
    }

    /**
     * 新增灰度组时查询实例列表
     *
     * @return
     */
    @RequestMapping(value = "/appointList", method = RequestMethod.GET)
    public Response<String> appointList(@Validated AddGrayGroupInstanceRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.instanceServiceImpl.appointList(body);
    }
}