package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.ServiceGrayConfigHistoryListRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.GrayConfigListDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.PushServiceConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IQuickStartService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceConfig;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.List;

/**
 * 灰度配置控制器
 * Created by DELL-5490 on 2018/7/2.
 */
@RestController
@RequestMapping(Constant.API + "/{version}/grayConfig")
public class GrayConfigController {

    @Autowired
    private IServiceConfig serviceConfigImpl;

    @Autowired
    private IQuickStartService quickStartServiceImpl;

    /**
     * 新增灰度配置
     *
     * @param file
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/addGrayConfig", method = RequestMethod.POST)
    public Response<GrayConfigListDetailResponseBody> addGrayConfig(MultipartFile[] file, @Validated AddGrayConfigRequestBody body, BindingResult result) throws UnsupportedEncodingException {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }

        //检查是否存在 文件
        int size = file.length;
        if (size <= 0) {
            return new Response<>(SystemErrCode.ERRCODE_NOT_FILE, SystemErrCode.ERRMSG_NOT_FILE);
        }

        List<FileContent> fileContentList = new ArrayList<>();
        List<String> fileNames = new ArrayList<>();
        FileContent fileContent;
        for (MultipartFile item : file) {
            //检查文件名
            String fileName = item.getOriginalFilename();
            if (fileNames.contains(fileName)) {
                return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
            }

            //检查文件内容
            if (item.isEmpty()) {
                return new Response<>(SystemErrCode.ERRCODE_FILE_CONTENT_NOT_NULL, SystemErrCode.ERRMSG_FILE_CONTENT_NOT_NULL);
            }

            //创建文件数据
            fileContent = new FileContent();
            fileContent.setFileName(fileName);
            try {

                byte[] content = item.getBytes();
                fileContent.setContent(content);
            } catch (IOException e) {
                GlobalExceptionUtil.log(e);
            }
            fileNames.add(fileName);
            fileContentList.add(fileContent);
        }
        return this.serviceConfigImpl.addGrayConfig(body, fileContentList);
    }

    /**
     * 查询灰度配置列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/lastestList", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> lastestList(@Validated QueryServiceConfigRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceConfigImpl.findLastestList(body);
    }

    /**
     * 编辑灰度配置
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<ServiceConfigDetailResponseBody> edit(@Validated @RequestBody EditServiceConfigRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceConfigImpl.edit(body);
    }

    /**
     * 查询灰度配置详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<ServiceConfigDetailResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.serviceConfigImpl.find(id);
    }

    /**
     * 推送灰度配置
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/push", method = RequestMethod.POST)
    public Response<PushServiceConfigResponseBody> push(@Validated @RequestBody PushServiceConfigRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        List<String> ids = body.getRegionIds();
        if (null == ids || ids.isEmpty()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_REGION_IDS_NOT_NULL);
        }
        return this.serviceConfigImpl.push(body);
    }

    /**
     * 批量推送配置
     *
     * @param body
     * @return
     */
    @RequestMapping(value = "/batchPush", method = RequestMethod.POST)
    public Response<PushServiceConfigResponseBody> batchPush(@RequestBody BatchPushServiceConfigRequestBody body) {
        List<String> configIds = body.getIds();
        if (null == configIds || configIds.isEmpty()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_SERVICE_CONFIG_IDS_NOT_NULL);
        }
        List<String> ids = body.getRegionIds();
        if (null == ids || ids.isEmpty()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_REGION_IDS_NOT_NULL);
        }
        return this.serviceConfigImpl.batchPush(body);
    }

    /**
     * 查询灰度配置历史列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/historyList", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> historyList(@Validated ServiceGrayConfigHistoryListRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceConfigImpl.findServiceGrayConfigHistoryList(body);
    }

    /**
     * 回滚配置
     *
     * @param body
     * @return
     */
    @RequestMapping(value = "/rollback", method = RequestMethod.POST)
    public Response<List<ServiceConfigDetailResponseBody>> rollback(@RequestBody IdsRequestBody body) {
        List<String> ids = body.getIds();
        if (null == ids || ids.isEmpty()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_SERVICE_CONFIG_IDS_NOT_NULL);
        }
        return this.serviceConfigImpl.rollback(body);
    }

    /**
     * 更新反馈
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/feedback", method = RequestMethod.POST)
    public Response<String> feedback(@Validated @RequestBody ServiceConfigFeedBackRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.serviceConfigImpl.feedback(body);
    }
}