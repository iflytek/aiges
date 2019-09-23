package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.nakedserver.http.ActionResult;
import com.iflytek.ccr.nakedserver.http.HttpContext;
import com.iflytek.ccr.polaris.companion.common.JsonResult;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.test.TestingServer;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

public class FinderControllerTest {

    private static ZkHelper zkHelper = null;
    private static TestingServer server = null;

    @BeforeClass
    public static void setUp() throws Exception {
        if (null == server) {
            server = new TestingServer();
        }
        zkHelper = new ZkHelper(server.getConnectString());
    }

    @AfterClass
    public static void tearDown() throws Exception {
        if (null != zkHelper)
            zkHelper.closeClient();
    }

    @Test
    public void testPush_config_feedback(){
        ConfigManager.getInstance().put(ConfigManager.KEY_ZKSTR,server.getConnectString());

        HttpContext context = new HttpContext(UUID.randomUUID().toString());
        Map<String, String> params = new HashMap<>();
        params.put("pushId","/root");
        params.put("project","/root");
        params.put("serviceGroup","/root");
        params.put("service","/root");
        params.put("version","/root");
        params.put("config","/root");
        params.put("addr","/root");
        params.put("updateStatus","/root");
        params.put("updateTime","/root");
        context.Request.setParams(params);

        FinderController finderController = new FinderController();
        ActionResult result = null;
        result = finderController.push_config_feedback(context);

        Assert.assertNotNull(result.getData());
        byte[] data = result.getData();
        String str = new String(data);
        JsonResult jsonResult = JacksonUtils.toObject(str,JsonResult.class);
        System.out.println( JacksonUtils.toJson(jsonResult));
        Assert.assertNotNull(result.getData());
        Assert.assertEquals(0,jsonResult.getRet());
    }
}
