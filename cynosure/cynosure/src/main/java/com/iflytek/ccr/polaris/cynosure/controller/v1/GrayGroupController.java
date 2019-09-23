package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.AddGrayGroupAndConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IGrayService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceConfig;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.util.*;

/**
 * 灰度组控制器
 * Created by DELL-5490 on 2018/7/4.
 */
@RestController
@RequestMapping(Constant.API + "/{version}/gray")
public class GrayGroupController {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(GrayGroupController.class);
    @Autowired
    private IServiceConfig serviceConfigImpl;

    @Autowired
    private IGrayService grayServiceImpl;

    /**
     * 新增灰度组
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public Response<AddGrayGroupAndConfigResponseBody> addGrayGroup(MultipartFile[] file, @Validated AddGrayGroupRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }

        // 取出配置Name的列表
        if (null != body.getIds() && !body.getIds().isEmpty()) {
            List<ServiceConfig> configList = this.serviceConfigImpl.findListByIds(body.getIds());
            List<String> serviceConfigNames = new ArrayList<>();
            for (ServiceConfig serviceConfig : configList) {
                serviceConfigNames.add(serviceConfig.getName());
            }
            Set<String> setServiceConfigName = new HashSet<>();
            if (serviceConfigNames.size() > 0) {
                //检查配置文件名是否重复
                for (String name : serviceConfigNames) {
                    if (!setServiceConfigName.add(name)) {
                        return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
                    }
                }
            }
        }

        //校验灰度组的推送实例是否有重复内容
        String instanceContent = body.getContent();
        if (StringUtils.isNotBlank(instanceContent)) {
            List<String> contentList = Arrays.asList(instanceContent.split(","));
            boolean isRepeat = contentList.size() != new HashSet<>(contentList).size();
            if (isRepeat) {
                return new Response<>(SystemErrCode.ERRCODE_GRAY_INSTANCE_REPEAT, SystemErrCode.ERRMSG_GRAY_INSTANCE_REPEAT);
            }
        }

        if (null != file && file.length > 0) {
            List<FileContent> fileContentList = new ArrayList<>();
            List<String> fileNames = new ArrayList<>();
            FileContent fileContent;

            //检查上传的文件中是否有重名的，同时还要排除是不是有和最近配置一样的文件名
            for (MultipartFile item : file) {
                //检查文件名
                String fileName = item.getOriginalFilename();
                if (fileNames.contains(fileName)) {
                    return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
                }

                //检查上传文件名和最近的配置文件是否有重复
                if (null != body.getIds() && !body.getIds().isEmpty()) {
                    List<ServiceConfig> configList = this.serviceConfigImpl.findListByIds(body.getIds());
                    List<String> serviceConfigNames = new ArrayList<>();
                    for (ServiceConfig serviceConfig : configList) {
                        serviceConfigNames.add(serviceConfig.getName());
                    }
                    if (serviceConfigNames.contains(fileName)) {
                        return new Response<>(SystemErrCode.ERRCODE_FILE_NAME_BATCH_SAME, SystemErrCode.ERRMSG_FILE_NAME_BATCH_SAME);
                    }
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
            if (fileContentList.size() > 10) {
                return new Response<>(SystemErrCode.ERRCODE_GRAY_CONFIGS_MAX_SIZE, SystemErrCode.ERRMSG_GRAY_CONFIGS_MAX_SIZE);
            }
            //调用含文件上传的新建灰度组业务
            return this.grayServiceImpl.addAndFile(body, fileContentList);
        } else {
            //调用不含文件上传的新建灰度组业务
            return this.grayServiceImpl.add(body);
        }
    }

    /**
     * 查询灰度组列表
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> list(@Validated QueryGrayGroupListRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.grayServiceImpl.findList(body);
    }

    /**
     * 查询灰度组详情
     *
     * @param id
     * @return
     */
    @RequestMapping(value = "/detail", method = RequestMethod.GET)
    public Response<AddGrayGroupAndConfigResponseBody> find(@RequestParam("id") String id) {
        if (StringUtils.isEmpty(id)) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
        }
        return this.grayServiceImpl.findById(id);
    }

    /**
     * 删除灰度组
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/delete", method = RequestMethod.POST)
    public Response<String> delete(@Validated @RequestBody DeleteGrayGroupRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.grayServiceImpl.delete(body);
    }

    /**
     * 编辑灰度组
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/edit", method = RequestMethod.POST)
    public Response<GrayGroupDetailResponseBody> edit(@Validated @RequestBody EditGrayGroupRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.grayServiceImpl.edit(body);
    }

    /**
     * 服务订阅查看
     *
     * @param body
     * @param result
     * @return
     */
    @RequestMapping(value = "/consumer", method = RequestMethod.GET)
    public Response<QueryPagingListResponseBody> consumer(@Validated QueryCustomDetailRequestBody body, BindingResult result) {
        if (result.hasErrors()) {
            return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
        }
        return this.grayServiceImpl.consumer(body);
    }
}