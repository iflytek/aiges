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
import com.iflytek.ccr.polaris.companion.service.ClusterService;
import com.iflytek.ccr.polaris.companion.service.ClusterServiceImpl;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;

import java.io.UnsupportedEncodingException;
import java.util.List;
import java.util.Map;

/**
 * 灰度组配置推送入口
 */
@API(name = "gray")
public class GrayController extends Controller {

    private final EasyLogger logger = EasyLoggerFactory.getInstance(GrayController.class);
    private final ClusterService clusterService = new ClusterServiceImpl();

    /**
     * 集群配置信息推送(灰度组下面的配置文件推送)
     *
     * @param context
     * @return
     */
    @Action(method = HttpMethod.POST)
    public ActionResult push_cluster_config(HttpContext context) {
        Map<String, List<HttpBody>> map = context.Request.getBodies();
        logger.info("bodies:" + JacksonUtils.toJson(map));
        JsonResult jsonResult = clusterService.grayPushConfig(map);
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
     * 灰度组数据删除
     *
     * @param context
     * @return
     */
    @Action(method = HttpMethod.POST)
    public ActionResult del_gray_group(HttpContext context) {
        Map<String, List<HttpBody>> map = context.Request.getBodies();
        logger.info("bodies:" + JacksonUtils.toJson(map));
        JsonResult jsonResult = clusterService.delGrayGroup(map);
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
