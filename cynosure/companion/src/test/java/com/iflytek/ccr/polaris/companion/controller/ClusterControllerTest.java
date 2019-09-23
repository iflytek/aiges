package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.nakedserver.http.ActionResult;
import com.iflytek.ccr.nakedserver.http.HttpBody;
import com.iflytek.ccr.nakedserver.http.HttpContext;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.JsonResult;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.test.TestingServer;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.util.*;

public class ClusterControllerTest {

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
    public void push_cluster_configTest() {
        ConfigManager.getInstance().put(ConfigManager.KEY_ZKSTR, server.getConnectString());
        System.out.println(ZkInstanceUtil.getInstance().getClientState());
        ClusterController clusterController = new ClusterController();
        HttpContext context = new HttpContext(UUID.randomUUID().toString());
        Map<String, String> params = new HashMap<>();
        params.put(Constants.ZK_NODE_PATH, "/root");
        params.put(Constants.ZK_NODE_PUSHID, "1111");
        context.Request.setParams(params);

        HttpBody httpBody = new HttpBody();
        httpBody.setParams(params);
        httpBody.setContentData("test".getBytes());

        List<HttpBody> list = new ArrayList<>();

        Map<String, List<HttpBody>> bodies = new HashMap<>();
        bodies.put(Constants.CONTENT_TYPE_OCTET_STREAM, list);
        context.Request.setBodies(bodies);
        ActionResult result = null;
        result = clusterController.push_cluster_config(context);
        Assert.assertNotNull(result.getData());
        byte[] data = result.getData();
        String str = new String(data);
        JsonResult jsonResult = JacksonUtils.toObject(str, JsonResult.class);
        Assert.assertEquals(0, jsonResult.getRet());
        ZkInstanceUtil.checkZkHelper();
    }

    @Test
    public void push_instance_configTest(){
        ConfigManager.getInstance().put(ConfigManager.KEY_ZKSTR, server.getConnectString());
        System.out.println(ZkInstanceUtil.getInstance().getClientState());
        ClusterController clusterController = new ClusterController();
        HttpContext context = new HttpContext(UUID.randomUUID().toString());
        Map<String, String> params = new HashMap<>();
        params.put(Constants.ZK_NODE_PATH, "/root");
        params.put(Constants.ZK_NODE_PUSHID, "1111");
        context.Request.setParams(params);

        HttpBody httpBody = new HttpBody();
        httpBody.setParams(params);
        httpBody.setContentData("test".getBytes());

        List<HttpBody> list = new ArrayList<>();

        Map<String, List<HttpBody>> bodies = new HashMap<>();
        bodies.put(Constants.CONTENT_TYPE_OCTET_STREAM, list);
        context.Request.setBodies(bodies);
        ActionResult result = null;
        result = clusterController.push_instance_config(context);
        Assert.assertNotNull(result.getData());
        byte[] data = result.getData();
        String str = new String(data);
        JsonResult jsonResult = JacksonUtils.toObject(str, JsonResult.class);
        Assert.assertEquals(0, jsonResult.getRet());
        System.out.println( ZkInstanceUtil.checkZkHelper());
    }

    @Test
    public void updateGrayNodeTest(){

    }
}
