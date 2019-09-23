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
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.JsonResult;
import com.iflytek.ccr.polaris.companion.task.RestoreZkData;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;

import java.io.UnsupportedEncodingException;

@API(name = "restore")
public class RestoreController extends Controller {

    private final EasyLogger logger = EasyLoggerFactory.getInstance(RestoreController.class);

    @Action(method = HttpMethod.GET)
    public ActionResult restore_zk_data(HttpContext context) {

        String path = context.Request.getParams().get("path");
        logger.info("path:" + path);

        new RestoreZkData().restore(path);
        JsonResult jsonResult = new JsonResult();
        ActionResult actionResult = null;
        try {
            actionResult = new ActionResult(HttpDataType.JSON, JacksonUtils.toJson(jsonResult).getBytes(Constants.DEFAULT_CHARSET));
        } catch (UnsupportedEncodingException e) {
            actionResult = new ActionResult(HttpDataType.JSON, JacksonUtils.toJson(jsonResult).getBytes());
            logger.error(e);
        }
        return actionResult;

    }
}
