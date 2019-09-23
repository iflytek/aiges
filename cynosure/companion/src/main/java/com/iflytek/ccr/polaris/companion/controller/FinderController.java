package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.nakedserver.annotation.API;
import com.iflytek.ccr.nakedserver.annotation.Action;
import com.iflytek.ccr.nakedserver.annotation.HttpMethod;
import com.iflytek.ccr.nakedserver.baseface.Controller;
import com.iflytek.ccr.nakedserver.http.ActionResult;
import com.iflytek.ccr.nakedserver.http.HttpContext;
import com.iflytek.ccr.nakedserver.http.HttpDataType;
import com.iflytek.ccr.polaris.companion.common.*;
import com.iflytek.ccr.polaris.companion.service.ConfigService;
import com.iflytek.ccr.polaris.companion.service.ConfigServiceImpl;
import com.iflytek.ccr.polaris.companion.service.FinderService;
import com.iflytek.ccr.polaris.companion.service.FinderServiceImpl;
import com.iflytek.ccr.polaris.companion.utils.ArgumentValidator;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;

import java.io.UnsupportedEncodingException;
import java.util.Map;

/**
 * Created by eric on 2017/11/21.
 */
@API(name = "finder")
public class FinderController extends Controller {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(FinderController.class);
    private final ConfigService configService = new ConfigServiceImpl();
    private final FinderService finderService = new FinderServiceImpl();

    /**
     * 查询zookeeper临时节点路径
     *
     * @return
     */
    @Action(method = HttpMethod.GET)
    public ActionResult query_zk_path(HttpContext context) {
        JsonResult jsonResult = finderService.queryZkPath();
        String result = JacksonUtils.toJson(jsonResult);
        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes());
            logger.error(e);
        }
        return actionResult;
    }

    @Action(method = HttpMethod.GET)
    public ActionResult query_zk_info(HttpContext context) {
        String[] kvParams = {"project", "group", "service", "version"};
        JsonResult jsonResult = new JsonResult();
        try {
            ArgumentValidator.validateKVParams(context.Request.getParams(), kvParams);
            jsonResult.getData().put("config_path", configService.generateConfigPath("/polaris/config/", context.Request.getParams()));
            jsonResult.getData().put("service_path", configService.generateServicePath("/polaris/service/", context.Request.getParams()));
            jsonResult.getData().put("zk_addr", finderService.getZkrAddrs(ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR)));
            jsonResult.getData().put("zk_node_path", ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZK_NODE_PATH));
        } catch (IllegalArgumentException e) {
            jsonResult.setRet(ErrorCode.PARAM_INVALID);
            jsonResult.setMsg(e.getMessage());
        } catch (Exception e) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg(e.getMessage());
        }

        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, JacksonUtils.toJson(jsonResult).getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, JacksonUtils.toJson(jsonResult).getBytes());
            logger.error(e);
        }
        return actionResult;
    }

    @Action(method = HttpMethod.GET)
    public ActionResult query_service_path(HttpContext context) throws UnsupportedEncodingException {
        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);
        JsonResult jsonResult = finderService.queryServicePath(path);
        return new ActionResult(HttpDataType.JSON, JacksonUtils.toJson(jsonResult).getBytes(Constants.DEFAULT_CHARSET));
    }

    @Action(method = HttpMethod.POST)
    public ActionResult push_config_feedback(HttpContext context) {
        JsonResult jsonResult = new JsonResult();
        Map<String, String> paramsMap = context.Request.getParams();
        logger.info("paramsMap:" + JacksonUtils.toJson(paramsMap));

        ConfigFeedBackValue feedBackValue = getFeedBackValue(paramsMap);

        boolean sendSuccess = finderService.pushConfigFeedback(feedBackValue);
        if (!sendSuccess) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("config feedback fail");
            logger.error("paramsMap:" + JacksonUtils.toJson(paramsMap) + ";feedback data fail");
        }

        String result = JacksonUtils.toJson(jsonResult);
        logger.info("jsonResult:" + result);
        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes());
            logger.error(e);
        }
        return actionResult;
    }


    @Action(method = HttpMethod.POST)
    public ActionResult push_service_feedback(HttpContext context) {
        JsonResult jsonResult = new JsonResult();
        Map<String, String> paramsMap = context.Request.getParams();
        logger.info("paramsMap:" + JacksonUtils.toJson(paramsMap));

        ServiceFeedBackValue serviceValue = getServiceValue(paramsMap);

        boolean sendSuccess = finderService.pushServiceFeedback(serviceValue);
        if (!sendSuccess) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("service feedback fail");
            logger.error("paramsMap:" + JacksonUtils.toJson(paramsMap) + ";feedback data fail");
        }

        String result = JacksonUtils.toJson(jsonResult);
        logger.info("jsonResult:" + result);
        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes());
            logger.error(e);
        }
        return actionResult;
    }

    private ServiceFeedBackValue getServiceValue(Map<String, String> paramsMap) {
        ServiceFeedBackValue serviceValue = new ServiceFeedBackValue();
        String pushId = paramsMap.get("push_id");
        String project = paramsMap.get("project");
        String group = paramsMap.get("group");
        String consumer = paramsMap.get("consumer");
        String provider = paramsMap.get("provider");
        String service = paramsMap.get("service");
        String addr = paramsMap.get("addr");
        String consumerVersion = paramsMap.get("consumer_version");
        String providerVersion = paramsMap.get("provider_version");
        String updateStatus = paramsMap.get("update_status");
        String updateTime = paramsMap.get("update_time");
        String loadStatus = paramsMap.get("load_status");
        String loadTime = paramsMap.get("load_time");
        String type = paramsMap.get("type");
        String apiVersion = paramsMap.get("api_version");
        serviceValue.setPushId(pushId);
        serviceValue.setGroup(group);
        serviceValue.setProject(project);
        serviceValue.setLoadStatus(loadStatus);
        serviceValue.setLoadTime(loadTime);
        serviceValue.setUpdateStatus(updateStatus);
        serviceValue.setUpdateTime(updateTime);
        serviceValue.setAddr(addr);
        serviceValue.setConsumerVersion(consumerVersion);
        serviceValue.setConsumer(consumer);
        serviceValue.setProviderVersion(providerVersion);
        serviceValue.setProvider(provider);
        serviceValue.setApiVersion(apiVersion);
        serviceValue.setType(type);
        return serviceValue;
    }

    private ConfigFeedBackValue getFeedBackValue(Map<String, String> paramsMap) {
        ConfigFeedBackValue feedBackValue = new ConfigFeedBackValue();
        String pushId = paramsMap.get("push_id");
        String project = paramsMap.get("project");
        String group = paramsMap.get("group");
        String grayGroupId = paramsMap.get("gray_group_id");
        String service = paramsMap.get("service");
        String version = paramsMap.get("version");
        String config = paramsMap.get("config");
        String addr = paramsMap.get("addr");
        String updateStatus = paramsMap.get("update_status");
        String updateTime = paramsMap.get("update_time");
        String loadStatus = paramsMap.get("load_status");
        String loadTime = paramsMap.get("load_time");
        feedBackValue.setAddr(addr);
        feedBackValue.setGroup(group);
        feedBackValue.setLoadStatus(loadStatus);
        feedBackValue.setLoadTime(loadTime);
        feedBackValue.setPushId(pushId);
        feedBackValue.setService(service);
        feedBackValue.setProject(project);
        feedBackValue.setVersion(version);
        feedBackValue.setUpdateStatus(updateStatus);
        feedBackValue.setUpdateTime(updateTime);
        feedBackValue.setConfig(config);
        feedBackValue.setGrayGroupId(grayGroupId);
        return feedBackValue;
    }

    private ServiceValue getServiceDiscoveryValue(Map<String, String> paramsMap) {
        ServiceValue serviceValue = new ServiceValue();
        String project = paramsMap.get("project");
        String group = paramsMap.get("group");
        String service = paramsMap.get("service");
        String apiVersion = paramsMap.get("api_version");
        serviceValue.setGroup(group);
        serviceValue.setProject(project);
        serviceValue.setService(service);
        serviceValue.setApiVersion(apiVersion);
        return serviceValue;
    }

    @Action(method = HttpMethod.POST)
    public ActionResult register_service_info(HttpContext context) {
        JsonResult jsonResult = new JsonResult();
        Map<String, String> paramsMap = context.Request.getParams();
        logger.info("paramsMap:" + JacksonUtils.toJson(paramsMap));

        ServiceValue serviceValue = getServiceDiscoveryValue(paramsMap);

        jsonResult = finderService.serviceDiscovery(serviceValue);
//        if (!sendSuccess) {
//            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
//            jsonResult.setMsg("register service info fail");
//            logger.error("paramsMap:" + jsonUtils.toJson(paramsMap) + ";register service info fail");
//        }

        String result = JacksonUtils.toJson(jsonResult);
        logger.info("jsonResult:" + result);
        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes());
            logger.error(e);
        }
        return actionResult;
    }

}
