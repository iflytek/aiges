package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.nakedserver.annotation.API;
import com.iflytek.ccr.nakedserver.annotation.Action;
import com.iflytek.ccr.nakedserver.annotation.HttpMethod;
import com.iflytek.ccr.nakedserver.baseface.Controller;
import com.iflytek.ccr.nakedserver.http.ActionResult;
import com.iflytek.ccr.nakedserver.http.HttpBody;
import com.iflytek.ccr.nakedserver.http.HttpContext;
import com.iflytek.ccr.nakedserver.http.HttpDataType;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.JsonResult;
import com.iflytek.ccr.polaris.companion.service.*;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;

import java.io.UnsupportedEncodingException;
import java.util.List;
import java.util.Map;

/**
 *  网站交互入口
 */
@API(name = "cynosure")
public class CynosureController extends Controller {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(CynosureController.class);
    private final ClusterService clusterService = new ClusterServiceImpl();
    private final CynosureService cynosureService = new CynosureServiceImpl();

    /**
     * 推送配置信息
     *
     * @param context
     * @return
     */
    @Action(method = HttpMethod.POST)
    public ActionResult push_config(HttpContext context) {
        Map<String, List<HttpBody>> map = context.Request.getBodies();
        logger.info("bodies:" + JacksonUtils.toJson(map));
        JsonResult jsonResult = clusterService.pushConfig(map);
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

    @Action(method = HttpMethod.GET)
    public ActionResult refresh_conf(HttpContext context) {
        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);
        JsonResult jsonResult = cynosureService.refreshConfStatus(path);
        String result = JacksonUtils.toJson(jsonResult);
        logger.info("日志jsonResult:" + result);
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
    public ActionResult refresh_service(HttpContext context) {
        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);
        JsonResult jsonResult = cynosureService.queryServiceList(path);
        String result = JacksonUtils.toJson(jsonResult);
        logger.info("日志jsonResult:" + result);
        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, result.getBytes());
            logger.error(e);
        }
        return actionResult;
    }

    /**
     * 刷新服务消费者列表信息
     * @param context
     * @return
     */
    @Action(method = HttpMethod.GET)
    public ActionResult refresh_consumer(HttpContext context) {
        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);
        JsonResult jsonResult = cynosureService.queryProviderOrConsumerList(path);
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

    /**
     * 刷新服务提供者列表信息
     *
     * @param context
     * @return
     */
    @Action(method = HttpMethod.GET)
    public ActionResult refresh_provider(HttpContext context) {
        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);
        JsonResult jsonResult = cynosureService.queryProviderOrConsumerList(path);
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

    /**
     * 删除zookeeper数据
     * @param context
     * @return
     */
    @Action(method = HttpMethod.POST)
    public ActionResult del_data(HttpContext context) {
        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);
        JsonResult jsonResult = cynosureService.delZkData(path);
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
