package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.nakedserver.http.HttpBody;
import com.iflytek.ccr.polaris.companion.common.*;
import com.iflytek.ccr.polaris.companion.utils.*;


import java.io.UnsupportedEncodingException;
import java.util.*;

public class ClusterServiceImpl implements ClusterService {

    private final EasyLogger logger = EasyLoggerFactory.getInstance(ClusterServiceImpl.class);

    @Override
    public JsonResult pushConfig(Map<String, List<HttpBody>> map) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        for (Iterator ite = map.entrySet().iterator(); ite.hasNext(); ) {
            Map.Entry entry = (Map.Entry) ite.next();
            ArrayList<HttpBody> octetList = (ArrayList) entry.getValue();
            if (null != octetList && !octetList.isEmpty()) {
                for (HttpBody body : octetList) {
                    String path = body.getParams().get(Constants.ZK_NODE_PATH);
                    String version = body.getParams().get(Constants.ZK_NODE_PUSHID);
                    byte[] dataByte = ByteUtil.getZkBytes(body.getContentData(), version);
                    logger.info("realPath:" + path);
                    boolean flag = ZkInstanceUtil.getInstance().addOrUpdatePersistentNode(path, dataByte);
                    if (!flag) {
                        jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                        jsonResult.setMsg("zookeeper update data fail");
                    }
                }
            }
        }
        logger.info(JacksonUtils.toJson(jsonResult));

        return jsonResult;
    }


    @Override
    public JsonResult pushInstanceConfig(Map<String, List<HttpBody>> map) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        for (Iterator ite = map.entrySet().iterator(); ite.hasNext(); ) {
            Map.Entry entry = (Map.Entry) ite.next();
            ArrayList<HttpBody> octetList = (ArrayList) entry.getValue();
            if (null != octetList && !octetList.isEmpty()) {
                for (HttpBody body : octetList) {
                    String path = body.getParams().get(Constants.ZK_NODE_PATH);
                    String version = body.getParams().get(Constants.ZK_NODE_PUSHID);
                    byte[] dataByte = ByteUtil.getZkBytes(body.getContentData(), version);
                    boolean flag = ZkInstanceUtil.getInstance().addOrUpdateEphemeralNode(path, dataByte);
                    if (!flag) {
                        jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                        jsonResult.setMsg("zookeeper update data fail");
                    }
                }
            }
        }
        logger.info(JacksonUtils.toJson(jsonResult));
        return jsonResult;
    }

    @Override
    public JsonResult grayPushConfig(Map<String, List<HttpBody>> map) {
        JsonResult jsonResult = new JsonResult();
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        String groupId = null;
        String grayServers = null;
        String grayPath = "";
        String version = null;
        for (Iterator ite = map.entrySet().iterator(); ite.hasNext(); ) {
            Map.Entry entry = (Map.Entry) ite.next();
            ArrayList<HttpBody> octetList = (ArrayList) entry.getValue();
            if (null != octetList && !octetList.isEmpty()) {
                for (HttpBody body : octetList) {
                    String path = body.getParams().get(Constants.ZK_NODE_PATH);
                    String fileName = body.getParams().get(Constants.ZK_NODE_FILE_NAME);
                    version = body.getParams().get(Constants.ZK_NODE_PUSHID);
                    String grayGroup = body.getParams().get(Constants.GRAY_GROUP);
                    grayServers = body.getParams().get(Constants.GRAY_SERVERS);
                    grayPath = path + "/" + Constants.ZK_NODE_GRAY;
                    String realPath = grayPath + "/" + grayGroup + "/" + fileName;
                    logger.info("realPath:" + realPath);

                    //保存文件
                    byte[] dataByte = ByteUtil.getZkBytes(body.getContentData(), version);
                    boolean flag = zkHelper.addOrUpdatePersistentNode(realPath, dataByte);
                    if (!flag) {
                        jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                        jsonResult.setMsg("zookeeper update data fail");
                    }
                    groupId = grayGroup;
                }
            }
        }

        List<GrayConfigValue> grayValueList = new ArrayList<>();
        //更新gray节点
        GrayConfigValue value = new GrayConfigValue();
        value.setGroupId(groupId);
        if (null == grayServers) {
            grayServers = "";
        }
        value.setServerList(Arrays.asList(grayServers));
        grayValueList.add(value);
        if (!grayValueList.isEmpty()) {
            updateGrayNode(grayPath, grayValueList, version);
        }
        logger.info(JacksonUtils.toJson(jsonResult));
        return jsonResult;
    }

    @Override
    public JsonResult delGrayGroup(Map<String, List<HttpBody>> map) {
        JsonResult jsonResult = new JsonResult();
        try {
            if (!ZkInstanceUtil.checkZkHelper()) {
                jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                jsonResult.setMsg("Unable to connect to zookeeper");
                return jsonResult;
            }
            for (Iterator ite = map.entrySet().iterator(); ite.hasNext(); ) {
                Map.Entry entry = (Map.Entry) ite.next();
                ArrayList<HttpBody> octetList = (ArrayList) entry.getValue();
                if (null != octetList && !octetList.isEmpty()) {
                    for (HttpBody body : octetList) {
                        String path = body.getParams().get(Constants.ZK_NODE_PATH);
                        String version = body.getParams().get(Constants.ZK_NODE_PUSHID);
                        String grayGroup = body.getParams().get(Constants.GRAY_GROUP);
                        String grayPath = path + "/" + Constants.ZK_NODE_GRAY;
                        delGrayNodeData(grayPath, grayGroup, version);
                    }
                }
            }
            logger.info(JacksonUtils.toJson(jsonResult));
        } catch (Exception e) {
            logger.error(e);
        }
        return jsonResult;
    }

    @Override
    public JsonResult pushServiceConfig(Map<String, String> map) {

        JsonResult jsonResult = new JsonResult();
        try {
            if (!ZkInstanceUtil.checkZkHelper()) {
                jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                jsonResult.setMsg("Unable to connect to zookeeper");
                return jsonResult;
            }
            String basePath = map.get(Constants.ZK_NODE_PATH);
            String version = map.get(Constants.ZK_NODE_PUSHID);
            String userData = map.get(Constants.USER_DATA);
            String sdkData = map.get(Constants.SDK_DATA);

            byte[] userDataByte = ByteUtil.getZkBytes(userData.getBytes(Constants.DEFAULT_CHARSET), version);
            if (null == userDataByte || userDataByte.length == 0) {
                userDataByte = "".getBytes();
            }
            boolean flag = ZkInstanceUtil.getInstance().addOrUpdatePersistentNode(basePath + Constants.ZK_PATH_CONF, userDataByte);
            if (!flag) {
                jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                jsonResult.setMsg("zookeeper update data fail");
            }

            byte[] sdkDataByte = ByteUtil.getZkBytes(sdkData.getBytes(Constants.DEFAULT_CHARSET), version);
            if (null == sdkDataByte || sdkDataByte.length == 0) {
                sdkDataByte = "".getBytes();
            }
            flag = ZkInstanceUtil.getInstance().addOrUpdatePersistentNode(basePath + Constants.ZK_PATH_ROUTE, sdkDataByte);
            if (!flag) {
                jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                jsonResult.setMsg("zookeeper update data fail");
            }
            logger.info(JacksonUtils.toJson(jsonResult));
        } catch (Exception e) {
            logger.error(e);
        }
        return jsonResult;
    }

    @Override
    public JsonResult pushServiceInstanceConfig(Map<String, List<HttpBody>> map) {

        JsonResult jsonResult = new JsonResult();
        try {
            if (!ZkInstanceUtil.checkZkHelper()) {
                jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                jsonResult.setMsg("Unable to connect to zookeeper");
                return jsonResult;
            }
            for (Iterator ite = map.entrySet().iterator(); ite.hasNext(); ) {
                Map.Entry entry = (Map.Entry) ite.next();
                ArrayList<HttpBody> octetList = (ArrayList) entry.getValue();
                if (null != octetList && !octetList.isEmpty()) {
                    for (HttpBody body : octetList) {
                        String basePath = body.getParams().get(Constants.ZK_NODE_PATH);
                        String version = body.getParams().get(Constants.ZK_NODE_PUSHID);
                        String userData = body.getParams().get(Constants.USER_DATA);
                        String sdkData = body.getParams().get(Constants.SDK_DATA);

                        byte[] sdkDataByte = ByteUtil.getZkBytes(sdkData.getBytes(Constants.DEFAULT_CHARSET), version);
                        boolean flag = ZkInstanceUtil.getInstance().addOrUpdatePersistentNode(basePath + Constants.ZK_PATH_CONF, sdkDataByte);
                        if (!flag) {
                            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                            jsonResult.setMsg("zookeeper update data fail");
                        }

                        byte[] userDataByte = ByteUtil.getZkBytes(userData.getBytes(Constants.DEFAULT_CHARSET), version);
                        flag = ZkInstanceUtil.getInstance().addOrUpdatePersistentNode(basePath + Constants.ZK_PATH_ROUTE, userDataByte);
                        if (!flag) {
                            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
                            jsonResult.setMsg("zookeeper update data fail");
                        }
                    }
                }
            }
            logger.info(JacksonUtils.toJson(jsonResult));
        } catch (Exception e) {
            logger.error(e);
        }


        return jsonResult;

    }

    public void delGrayNodeData(String grayPath, String grayGroupId, String pushId) {
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            //获取历史数据
            byte[] data = null;
            if (zkHelper.checkExists(grayPath)) {
                data = zkHelper.getByteData(grayPath);
            }
            if (data == null || data.length == 0) {

                return;
            }

            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            String grayData = new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET);
            //获取当前zookeeper中已经存在的灰度组数据
            List grayList = JacksonUtils.toObject(grayData, List.class);
            Map<String, Object> tempValue = null;
            for (Object v : grayList) {
                Map<String, Object> current = (Map<String, Object>) v;
                String tempGroupId = current.get(Constants.GROUP_ID).toString();
                if (tempGroupId.equals(grayGroupId)) {
                    tempValue = current;
                    break;
                }

            }
            if (tempValue != null) {
                try {
                    grayList.remove(tempValue);
                    byte[] dataByte = ByteUtil.getZkBytes(JacksonUtils.toJson(grayList).getBytes(Constants.DEFAULT_CHARSET), pushId);
                    zkHelper.addOrUpdatePersistentNode(grayPath, dataByte);
                    zkHelper.remove(grayPath + "/" + grayGroupId);
                } catch (UnsupportedEncodingException e) {
                    logger.error(e);
                }
            }
        } catch (Exception e) {
            logger.error(e);
        }
    }

    /**
     * 更新灰度节点数据
     *
     * @param pushId
     * @param grayPath
     * @param currentGrayValueList
     */
    public void updateGrayNode(String grayPath, List<GrayConfigValue> currentGrayValueList, String pushId) {
        try {
            boolean needUpdate = false;
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();

            //获取历史数据
            byte[] data = null;
            if (zkHelper.checkExists(grayPath)) {
                data = zkHelper.getByteData(grayPath);
            }
            if (data == null || data.length == 0) {
                byte[] dataByte = ByteUtil.getZkBytes(JacksonUtils.toJson(currentGrayValueList).getBytes(Constants.DEFAULT_CHARSET), pushId);
                zkHelper.addOrUpdatePersistentNode(grayPath, dataByte);
                return;
            }
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            String grayData = new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET);
            //获取当前zookeeper中已经存在的灰度组数据
            List grayList = JacksonUtils.toObject(grayData, List.class);

            for (GrayConfigValue current : currentGrayValueList) {
                boolean exists = false;
                for (Object v : grayList) {
                    String tempGroupId;
                    List<String> tempList;
                    if (v instanceof GrayConfigValue) {
                        GrayConfigValue temp = (GrayConfigValue) v;
                        tempGroupId = temp.getGroupId();
                        tempList = temp.getServerList();
                    } else {
                        Map<String, Object> vMap = (Map<String, Object>) v;
                        tempGroupId = vMap.get(Constants.GROUP_ID).toString();
                        tempList = (List<String>) vMap.get(Constants.SERVER_LIST);
                    }

                    if (tempGroupId.equals(current.getGroupId())) {
                        exists = true;
                        if (!ArrayUtils.equals(tempList, current.getServerList())) {
                            if (v instanceof GrayConfigValue) {
                                GrayConfigValue temp = (GrayConfigValue) v;
                                temp.setServerList(current.getServerList());
                            } else {
                                Map<String, Object> vMap = (Map<String, Object>) v;
                                vMap.put(Constants.SERVER_LIST, current.getServerList());
                            }

                            needUpdate = true;
                        }
                    }
                }
                if (!exists) {
                    needUpdate = true;
                    grayList.add(current);
                }
            }
            if (needUpdate) {
                try {
                    byte[] dataByte = ByteUtil.getZkBytes(JacksonUtils.toJson(grayList).getBytes(Constants.DEFAULT_CHARSET), pushId);
                    zkHelper.addOrUpdatePersistentNode(grayPath, dataByte);
                } catch (UnsupportedEncodingException e) {
                    logger.error(e);
                }
            }
        } catch (Exception e) {
            logger.error(e);
        }
    }


    public static void main(String[] args) {
        List<GrayConfigValue> currentGrayValueList = new ArrayList<>();
        ConfigManager.getInstance().put(ConfigManager.KEY_ZKSTR, "10.1.87.69:2183");
//        GrayConfigValue grayValue = new GrayConfigValue();
//        grayValue.setGroupId("b");
//        List<String> list = new ArrayList<>();
//        list.add("55.1.2.2:1010");
////        list.add("55.11.2.2:7010");
//        grayValue.setServerList(list);
//        currentGrayValueList.add(grayValue);
//        grayValue = new GrayConfigValue();
//        grayValue.setGroupId("a");
//        list = new ArrayList<>();
//        list.add("88.1.2.2:1010");
//        list.add("88.11.2.2:7030");
//        list.add("55.11.2.2:7020");
//        grayValue.setServerList(list);
//        currentGrayValueList.add(grayValue);
//        ClusterServiceImpl service = new ClusterServiceImpl();
//        String grayPath = "/polaris/config/c94370c924bac3f56e77434935613b23/5870ac10085ea78db59e6a231fefc6cf/gray";


        String pushId = "123456789";
//        service.updateGrayNode(grayPath, currentGrayValueList, pushId);
//
//        service.delGrayNodeData(grayPath, "c", "333331111111");


        String str = "{\"loadbalance\":\"loadbalance\",\"key1\":\"val\",\"key2\":\"val\"}";
        try {
            String str2 = "[{\"routeRuleId\":\"loadbalance\",\"consumer\":\"val\",\"provider\":\"val\",\"only\":\"Y\"}]";
            String path = "/polaris/service/c94370c924bac3f56e77434935613b23/iatExecutor/1.0/conf";
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            byte[] dataByte = ByteUtil.getZkBytes(str.getBytes(Constants.DEFAULT_CHARSET), pushId);
            System.out.println(zkHelper.addOrUpdatePersistentNode(path, dataByte));

            path = "/polaris/service/c94370c924bac3f56e77434935613b23/iatExecutor/1.0/route";
            dataByte = ByteUtil.getZkBytes(str2.getBytes(Constants.DEFAULT_CHARSET), pushId);
            System.out.println(zkHelper.addOrUpdatePersistentNode(path, dataByte));

            zkHelper.closeClient();
        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }

    }
}
