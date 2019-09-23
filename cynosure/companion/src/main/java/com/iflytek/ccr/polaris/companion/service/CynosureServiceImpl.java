package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.common.*;
import com.iflytek.ccr.polaris.companion.utils.ByteUtil;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class CynosureServiceImpl implements CynosureService {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(CynosureServiceImpl.class);

    @Override
    public JsonResult queryServiceList(String path) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();

        Map<String, Object> dataMap = new HashMap<>();
        jsonResult.setData(dataMap);
        List<ServieValue> list = new ArrayList<>();
        dataMap.put("pathList", list);
        if (zkHelper.checkExists(path)) {
            List<String> childrenList = zkHelper.getChildren(path);
            if (null != childrenList) {
                for (String childStr : childrenList) {
                    ServieValue zkData = new ServieValue();
                    list.add(zkData);
                    zkData.setPath(childStr);
                    String versionPath = path + "/" + childStr;
                    if (zkHelper.checkExists(versionPath)) {
                        List<String> versionList = zkHelper.getChildren(versionPath);
                        zkData.setVersionList(versionList);
                    }
                }
            }
        } else {
            jsonResult.setRet(ErrorCode.PATH_NOT_EXISTS);
            jsonResult.setMsg("Path does not exist");
        }
        logger.info(JacksonUtils.toJson(jsonResult));
        return jsonResult;
    }

    @Override
    public JsonResult queryProviderOrConsumerList(String path) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();

        Map<String, Object> dataMap = new HashMap<>();
        jsonResult.setData(dataMap);
        List<ZkData> list = new ArrayList<>();
        dataMap.put("pathList", list);
        if (zkHelper.checkExists(path)) {
            List<String> childrenList = zkHelper.getChildren(path);
            if (null != childrenList) {
                for (String childStr : childrenList) {
                    ZkData zkData = new ZkData();
                    list.add(zkData);
                    zkData.setChildPath(childStr);
                    String tempPath = path + "/" + childStr;
                    if (zkHelper.checkExists(tempPath)) {
                        byte[] data = zkHelper.getByteData(tempPath);
                        byte[] preByte = new byte[4];
                        if (data.length <= 4) {
                            logger.error("path  is invalid:" + childStr);
                            continue;
                        }
                        System.arraycopy(data, 0, preByte, 0, 4);
                        int versionLength = ByteUtil.byteArrayToInt(preByte);
                        if (data.length < versionLength) {
                            logger.error("versionLength is invalid:" + versionLength);
                            continue;
                        }
                        byte[] verByte = new byte[versionLength];
                        System.arraycopy(data, 4, verByte, 0, versionLength);
                        String pushId = new String(verByte);

                        byte[] realData = new byte[data.length - 4 - versionLength];
                        System.arraycopy(data, 4 + versionLength, realData, 0, data.length - 4 - versionLength);
                        try {
                            zkData.setData(new String(realData, Constants.DEFAULT_CHARSET));
                        } catch (UnsupportedEncodingException e) {
                            logger.error(e);
                        }
                        zkData.setPushId(pushId);
                    }
                }
            }
        } else {
            jsonResult.setRet(ErrorCode.PATH_NOT_EXISTS);
            jsonResult.setMsg("Path does not exist");
        }
        logger.info(JacksonUtils.toJson(jsonResult));
        return jsonResult;
    }

    @Override
    public JsonResult refreshConfStatus(String path) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        Map<String, Object> dataMap = new HashMap<>();
        jsonResult.setData(dataMap);
        dataMap.put("user", getZkData(path + "/conf"));
        dataMap.put("sdk", getZkData(path + "/route"));
        logger.info(JacksonUtils.toJson(jsonResult));
        return jsonResult;
    }

    private ZkData getZkData(String path) {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        ZkData zkData = new ZkData();
        if (zkHelper.checkExists(path)) {
            byte[] data = zkHelper.getByteData(path);
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                zkData.setPushId(zkDataValue.getPushId());
                zkData.setChildPath(path.substring(path.lastIndexOf("/") + 1));
                try {
                    zkData.setData(new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET));
                } catch (UnsupportedEncodingException e) {
                    logger.error(e);
                }
            } else {
                zkData.setData("parse data error");
            }
        } else {
            zkData.setData("path does not exists");
        }
        return zkData;
    }

    @Override
    public JsonResult refreshServiceConfStatus(String path) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        Map<String, Object> dataMap = new HashMap<>();
        jsonResult.setData(dataMap);
        List<ZkData> list = new ArrayList<>();
        dataMap.put("pathList", list);
        list.add(getZkData(path));
        logger.info(JacksonUtils.toJson(jsonResult));
        return jsonResult;
    }

    @Override
    public JsonResult delZkData(String path) {
        JsonResult jsonResult = new JsonResult();
        if (!ZkInstanceUtil.checkZkHelper()) {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("Unable to connect to zookeeper");
            return jsonResult;
        }
        if (path.startsWith(Constants.CONFIG_PATH_PREFIX) || path.startsWith(Constants.SERVICE_PATH_PREFIX)) {
            boolean flag = path.length() > Constants.CONFIG_PATH_PREFIX.length() + 5 || path.length() > Constants.SERVICE_PATH_PREFIX.length() + 5;
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            if (zkHelper.checkExists(path) && flag) {
                zkHelper.remove(path);
            }
        } else {
            jsonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            jsonResult.setMsg("path is invalid");
        }
        return jsonResult;
    }
}
