package com.iflytek.ccr.polaris.cynosure.companionservice;

import com.alibaba.fastjson.JSONArray;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.*;
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
 * 灰度配置中心
 *
 * @author sctang2
 * @create 2017-12-10 19:05
 **/
@Service
public class GrayConfigCenter extends BaseCenter {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(GrayConfigCenter.class);
    @Autowired
    private PropUtil propUtil;

    @Autowired
    private CompanionCacheCenter companionCacheCenter;

    /**
     * 获取配置服务path
     *
     * @param project
     * @param group
     * @return
     */
    public String getConfigPath(String project, String group, String service, String version) {
        String projectJointGroup = project + group;
        String projectJointGroupMD5 = MD5Util.getMD5(projectJointGroup.getBytes());
        String serviceJointVersion = service + version;
        String serviceJointVersionMD5 = MD5Util.getMD5(serviceJointVersion.getBytes());
        return propUtil.CONFIG_PATH + projectJointGroupMD5 + "/" + serviceJointVersionMD5 + "/";
    }

    /**
     * 获取配置消费者path
     *
     * @param project
     * @param group
     * @param service
     * @param grayGroupId
     * @return
     */
    public String getServiceConsumerPath(String project, String group, String service, String version, String grayGroupId) {
        String configPath = this.getConfigPath(project, group, service, version);
        if ("0".equals(grayGroupId)) {
            return configPath + "consumer" + "/normal";
        } else {
            return configPath + "consumer" + "/gray" + "/" + grayGroupId;
        }
    }

    /**
     * 通过一对多推送
     *
     * @param path
     * @param bt
     * @param regionList
     * @return
     */
    public PushResult pushByOneToMany(String path, String fileName, String grayGroupId, String grayServers, byte[] bt, List<Region> regionList) {
        String pushId = SnowflakeIdWorker.getId();
        //创建
        CustomHttpParams params = this.create(path, fileName, grayGroupId, grayServers, pushId, bt);

        //批量post请求
        List<Result> cacheCenterResults = this.batchPost(PUSH_GRAY_CONFIG_URL, params, regionList);

        //解析推送返回参数
        List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);

        //构造返回结果
        PushResult result = new PushResult( pushId, JSONArray.toJSONString(cacheCenterPushDetailResults), cacheCenterPushDetailResults);
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
    public PushResult pushByManyToMany(List<String> pathList, List<String> fileNameList, String grayGroupId, String grayServers, List<byte[]> btList, List<Region> regionList) {
        String pushId = SnowflakeIdWorker.getId();
        //创建
        List<CustomHttpParams> customHttpParamsList = this.create(pathList, fileNameList, grayGroupId, grayServers, pushId, btList);

        //批量post请求
        List<Result> cacheCenterResults = this.batchPost(PUSH_GRAY_CONFIG_URL, customHttpParamsList, regionList);

        //解析推送参数
        List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);
        PushResult result = new PushResult();
        result.setPushId(pushId);
        result.setData(cacheCenterPushDetailResults);
        result.setResult(JSONArray.toJSONString(cacheCenterPushDetailResults));
        return result;
    }

    /**
     * 查询消费端列表
     *
     * @param path
     * @param region
     * @param isProvider
     * @param startIndex
     * @param endIndex
     * @return
     */
    public ServiceResult findConsumersByPaging(String path, Region region, boolean isProvider, int startIndex, int endIndex) {
        //查询提供端、消费端列表
        String type = "consumer";
        ServiceResult serviceResult = this.companionCacheCenter.findConfigsConsumers(path, region, type);
        int totalCount = serviceResult.getTotalCount();
        if (0 == totalCount) {
            return serviceResult;
        }
        List<ServiceProviderConsumerResult> serviceProviderConsumerResults = serviceResult.getResults();
        if (startIndex >= totalCount) {
            startIndex = 0;
        }
        if (endIndex >= totalCount) {
            endIndex = totalCount;
        }
        List<ServiceProviderConsumerResult> newServiceProviderConsumerResults = serviceProviderConsumerResults.subList(startIndex, endIndex);
        ServiceResult newServiceResult = new ServiceResult();
        newServiceResult.setTotalCount(totalCount);
        newServiceResult.setResults(newServiceProviderConsumerResults);
        return newServiceResult;
    }

    /**
     * 删除灰度配置
     *
     * @param path
     * @param regionList
     */
    @Async
    public void deleteGrayConfig(String path, String grayGroupId, List<Region> regionList) {
        String pushId = SnowflakeIdWorker.getId();
        Map<String, String> map = new HashMap<>();
        map.put("path", path);
        map.put("pushId", pushId);
        map.put("grayGroupId", grayGroupId);
        this.batchPost(DELETE_DATA_URL, map, regionList);
    }

    /**
     * 删除灰度组
     */
    public void deleteGrayGroup(String path, String grayGroupId, List<Region> regionList) {
        String pushId = SnowflakeIdWorker.getId();
        //创建
        byte[] bt = {};//发送空字符bt
        CustomHttpParams params = this.create(path, grayGroupId, pushId, bt);
        //批量post请求
        this.batchPost(DELETE_GRAY_GROUP_URL, params, regionList);
    }
}