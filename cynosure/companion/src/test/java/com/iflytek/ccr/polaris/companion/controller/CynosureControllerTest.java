package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.nakedserver.http.ActionResult;
import com.iflytek.ccr.nakedserver.http.HttpContext;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.JsonResult;
import com.iflytek.ccr.polaris.companion.utils.ByteUtil;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.test.TestingServer;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.Map;
import java.util.UUID;

public class CynosureControllerTest {
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
        if(null != zkHelper)
            zkHelper.closeClient();
    }

    @Test
    public void  refresh_provider_statusTest(){
        ConfigManager.getInstance().put(ConfigManager.KEY_ZKSTR,server.getConnectString());

        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        zkHelper.addPersistent("/root/data/a/data","{\"a\":1}");
        zkHelper.addPersistent("/root/data/b/data","{\"a\":2}");
        zkHelper.addPersistent("/root/data/c/data","{\"a\":3}");


        String path = "/root/data";
        String version = "9999";
        String nodeData = "abc123456";
        byte[] dataByte = null;
        try {
            byte[] versionBytes = version.getBytes(Constants.DEFAULT_CHARSET);
            byte[] pre = ByteUtil.intToByteArray(versionBytes.length);
            dataByte = ByteUtil.byteMerge(pre, versionBytes);
            dataByte = ByteUtil.byteMerge(dataByte, nodeData.getBytes());
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }

        zkHelper.addOrUpdatePersistentNode(path, dataByte);

        CynosureController cynosureController = new CynosureController();
        HttpContext context = new HttpContext(UUID.randomUUID().toString());
        Map<String, String> params = new HashMap<>();
        params.put("path","/root");
        context.Request.setParams(params);
        ActionResult result = null;
        result = cynosureController.refresh_provider(context);
        Assert.assertNotNull(result.getData());
        byte[] data = result.getData();
        String str = null;
        try {
            str = new String(data, Constants.DEFAULT_CHARSET);
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }
        JsonResult jsonResult = JacksonUtils.toObject(str,JsonResult.class);
        ArrayList list = (ArrayList) jsonResult.getData().get("pathList");
        Assert.assertEquals(((Map) list.get(0)).get("data"),nodeData);
        Assert.assertEquals(((Map) list.get(0)).get("pushId"),version);
    }
}
