package com.iflytek.ccr.polaris.cynosure.companionservice;

import com.alibaba.fastjson.JSONArray;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushDetailResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.Result;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.network.CustomHttpParams;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.PropUtil;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 配置中心
 *
 * @author sctang2
 * @create 2017-12-10 19:05
 **/
@Service
public class ConfigCenter extends BaseCenter {
    @Autowired
    private PropUtil propUtil;

    /**
     * 获取配置路径
     *
     * @param project
     * @param group
     * @param service
     * @param version
     * @param fileName
     * @return
     */
    public String getConfigPath(String project, String group, String service, String version, String fileName) {
        String projectJointGroup = project + group;
        String projectJointGroupMD5 = MD5Util.getMD5(projectJointGroup.getBytes());//转化为byte数组
        String serviceJointVersion = service + version;
        String serviceJointVersionMD5 = MD5Util.getMD5(serviceJointVersion.getBytes());
        System.out.println("-------------------------");
        System.out.println(projectJointGroupMD5+"   "+serviceJointVersionMD5);
        return propUtil.CONFIG_PATH + projectJointGroupMD5 + "/" + serviceJointVersionMD5 + "/" + fileName;
    }

    /**
     * 通过一对多推送
     *
     * @param path
     * @param bt
     * @param regionList
     * @return
     */
    public PushResult pushByOneToMany(String path, byte[] bt, List<Region> regionList) {
        String pushId = SnowflakeIdWorker.getId();
        //创建
        CustomHttpParams params = this.create(path, pushId, bt);

        //批量post请求
        List<Result> cacheCenterResults = this.batchPost(PUSH_CONFIG_URL, params, regionList);

        //解析推送参数
        List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);
        PushResult result = new PushResult();
        result.setPushId(pushId);
        result.setData(cacheCenterPushDetailResults);
        result.setResult(JSONArray.toJSONString(cacheCenterPushDetailResults));
        return result;
    }

    /**
     * 通过多对多推送
     *
     * @param pathList
     * @param btList
     * @param regionList
     * @return
     */
    public PushResult pushByManyToMany(List<String> pathList, List<byte[]> btList, List<Region> regionList) {
        String pushId = SnowflakeIdWorker.getId();
        //创建
        List<CustomHttpParams> customHttpParamsList = this.create(pathList, pushId, btList);

        //批量post请求
        List<Result> cacheCenterResults = this.batchPost(PUSH_CONFIG_URL, customHttpParamsList, regionList);

        //解析推送参数
        List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);
        PushResult result = new PushResult();
        result.setPushId(pushId);
        result.setData(cacheCenterPushDetailResults);
        result.setResult(JSONArray.toJSONString(cacheCenterPushDetailResults));
        return result;
    }

    /**
     * 删除配置
     *
     * @param path
     * @param regionList
     */
    @Async
    public void deleteConf(String path, List<Region> regionList) {
        Map<String, String> map = new HashMap<>();
        map.put("path", path);
        this.batchPost(DELETE_DATA_URL, map, regionList);
    }
}
