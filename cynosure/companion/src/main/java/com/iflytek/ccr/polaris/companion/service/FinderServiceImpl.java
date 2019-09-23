package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.cache.FeedbackCacheUtil;
import com.iflytek.ccr.polaris.companion.cache.ServiceCacheUtil;
import com.iflytek.ccr.polaris.companion.common.*;
import com.iflytek.ccr.polaris.companion.main.Program;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.HttpsClientUtil;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.http.client.fluent.Request;
import org.apache.http.entity.StringEntity;

import java.util.*;

public class FinderServiceImpl implements FinderService {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(FinderServiceImpl.class);
    private final ConfigService configService = new ConfigServiceImpl();

    @Override
    public JsonResult queryZkPath() {
        JsonResult jsonResult = new JsonResult();
        Map<String, Object> dataMap = new HashMap<>();
        dataMap.put("path", ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZK_NODE_PATH));
        jsonResult.setData(dataMap);
        return jsonResult;
    }

    @Override
    public boolean pushConfigFeedback(ConfigFeedBackValue feedBackValue) {
        logger.info("feedBackValue:" + feedBackValue);
        boolean flag = false;
        try {
            FeedbackCacheUtil.getInstance().add(feedBackValue);
            flag = true;
        } catch (Exception e) {
            logger.error("add ", e);
        }
        return flag;
    }

    @Override
    public boolean pushServiceFeedback(ServiceFeedBackValue serviceFeedBackValue) {
        logger.info("serviceFeedBackValue:" + serviceFeedBackValue);
        boolean flag = false;
        try {
            ServiceCacheUtil.getInstance().add(serviceFeedBackValue);
            flag = true;
        } catch (Exception e) {
            logger.error("add exception ", e);
        }
        return flag;
    }

    @Override
    public JsonResult serviceDiscovery(ServiceValue serviceValue) {
        JsonResult jsonResult = new JsonResult();
        String url = Program.CONFIG_VALUE.getWebsiteUrl() + Constants.DISCOVERY_SERVICE_SITE_URI;
        String reqJson = JacksonUtils.toJson(serviceValue);
        try {
            StringEntity entity = new StringEntity(reqJson);
            String result;
            if (url.startsWith("https")) {
                result = HttpsClientUtil.doPostByStringEntity(url, entity, Constants.DEFAULT_CHARSET);
            } else {
                result = Request.Post(url).setHeader("Content-type", "application/json")
                        .body(entity).socketTimeout(5000).connectTimeout(5000)
                        .execute().returnContent().asString();
            }
            WebsitResult response = JacksonUtils.toObject(result, WebsitResult.class);
            if (!Constants.SUCCESS.equals(response.getCode())) {
                logger.error("request:" + reqJson + ",result:" + result);
            }
            jsonResult.setRet(Integer.parseInt(response.getCode()));
            jsonResult.setMsg(response.getMessage());
        } catch (Exception e) {
            logger.error(e);
        }
        return jsonResult;
    }

    @Override
    public JsonResult queryServicePath(String path) {
        JsonResult jsonResult = new JsonResult();
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        List<String> childrenList = zkHelper.getChildren(path);
        Map<String, Object> dataMap = new HashMap<>();
        dataMap.put("path", childrenList);
        jsonResult.setData(dataMap);
        return jsonResult;
    }

    @Override
    public List<String> getZkrAddrs(String str) {
        List<String> addrs = new ArrayList<>();
        if (str != null && str.length() > 0) {
            addrs.addAll(Arrays.asList(str.split(",")));
        }
        return addrs;
    }

}
