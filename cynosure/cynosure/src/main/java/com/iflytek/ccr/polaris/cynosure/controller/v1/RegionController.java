package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.annotation.Access;
import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.AddRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.EditRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.region.RegionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IRegionService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

/**
 * 区域控制器
 *
 * @author sctang2
 * @create 2017-11-14 20:51
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/region")
public class RegionController {
    @Autowired
    private IRegionService regionServiceImpl;

    /**
     * 查询区域列表
     *
     * @param body
     * @return
     */
    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> list(BaseRequestBody body) {
        return this.regionServiceImpl.findList(body);
    }

    /**
     * 查询区域详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<RegionDetailResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.regionServiceImpl.find(id);
    }

    /**
     * 新增区域
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<RegionDetailResponseBody> add(@Validated @RequestBody AddRegionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.regionServiceImpl.add(body);
    }

    /**
     * 编辑区域
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    @Access(authorities = "admin")
    public Response<RegionDetailResponseBody> edit(@Validated @RequestBody EditRegionRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.regionServiceImpl.edit(body);
    }

    /**
     * 删除区域
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
        return this.regionServiceImpl.delete(body);
    }
}
