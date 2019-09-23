package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigPushHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigHistory;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushHistory;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.PushServiceConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigHistoryResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.track.TrackConfig;
import com.iflytek.ccr.polaris.cynosure.response.track.TrackRegion;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.List;

/**
 * 服务配置事务操作
 *
 * @author sctang2
 * @create 2017-12-10 22:25
 **/
@Service
public class ServiceConfigTransactional {
    @Autowired
    private IServiceConfigCondition serviceConfigConditionImpl;

    @Autowired
    private IServiceConfigHistoryCondition serviceConfigHistoryConditionImpl;

    @Autowired
    private IServiceConfigPushHistoryCondition servicePushHistoryConditionImpl;

    /**
     * 新增服务配置和推送历史
     *
     * @param serviceConfig
     * @param cacheCenterPushResult
     * @return
     */
    @Transactional
    public PushServiceConfigResponseBody addServiceConfigAndPushHistory(ServiceConfig serviceConfig, PushResult cacheCenterPushResult) {

        //新增服务配置历史
        ServiceConfigHistory serviceConfigHistory = this.serviceConfigHistoryConditionImpl.add(serviceConfig);

        //新增服务配置推送历史
        ServiceConfigPushHistory serviceConfigPushHistory = this.servicePushHistoryConditionImpl.add(serviceConfig, cacheCenterPushResult);

        //创建推送服务配置
        List<ServiceConfigHistory> serviceConfigHistories = new ArrayList<>();
        serviceConfigHistories.add(serviceConfigHistory);
        List<ServiceConfig> serviceConfigList = new ArrayList<>();
        serviceConfigList.add(serviceConfig);
        PushServiceConfigResponseBody result = this.createPushServiceConfig(serviceConfigPushHistory, serviceConfigHistories, serviceConfigList);
        return result;
    }

    /**
     * 新增服务配置列表和推送历史
     *
     * @param serviceConfigList
     * @param cacheCenterPushResult
     */
    @Transactional
    public PushServiceConfigResponseBody addServiceConfigsAndPushHistory(List<ServiceConfig> serviceConfigList, PushResult cacheCenterPushResult) {
        //新增服务配置历史
        List<ServiceConfigHistory> serviceConfigHistories = this.serviceConfigHistoryConditionImpl.batchAdd(serviceConfigList);

        //新增服务配置推送历史
        ServiceConfigPushHistory serviceConfigPushHistory = this.servicePushHistoryConditionImpl.add(serviceConfigList, cacheCenterPushResult);

        //创建推送服务配置
        PushServiceConfigResponseBody result = this.createPushServiceConfig(serviceConfigPushHistory, serviceConfigHistories, serviceConfigList);
        return result;
    }

    /**
     * 删除配置和历史
     *
     * @param configId
     * @return
     */
    @Transactional
    public int deleteConfigAndHistory(String configId) {
        //通过configId删除配置历史
        this.serviceConfigHistoryConditionImpl.deleteByConfigId(configId);

        //通过id删除服务配置
        return this.serviceConfigConditionImpl.deleteById(configId);
    }

    /**
     * 批量删除配置和历史
     *
     * @param configIds
     * @return
     */
    @Transactional
    public int batchDeleteConfigAndHistory(List<String> configIds) {
        //通过configIds删除配置历史
        this.serviceConfigHistoryConditionImpl.deleteByConfigIds(configIds);

        //通过ids删除服务配置
        return this.serviceConfigConditionImpl.deleteByIds(configIds);
    }

    /**
     * 创建推送服务配置
     *
     * @param serviceConfigPushHistory
     * @param serviceConfigHistories
     * @param serviceConfigList
     * @return
     */
    private PushServiceConfigResponseBody createPushServiceConfig(ServiceConfigPushHistory serviceConfigPushHistory, List<ServiceConfigHistory> serviceConfigHistories, List<ServiceConfig> serviceConfigList) {
        PushServiceConfigResponseBody result = new PushServiceConfigResponseBody();
        result.setId(serviceConfigPushHistory.getId());
        String configText = serviceConfigPushHistory.getServiceConfigText();
        List<TrackConfig> trackConfigs = JSON.parseArray(configText, TrackConfig.class);
        result.setConfigs(trackConfigs);
        String regionText = serviceConfigPushHistory.getClusterText();
        List<TrackRegion> trackRegions = JSON.parseArray(regionText, TrackRegion.class);
        result.setRegions(trackRegions);
        result.setPushTime(serviceConfigPushHistory.getPushTime());
        if (null != serviceConfigHistories && !serviceConfigHistories.isEmpty()) {
            List<ServiceConfigHistoryResponseBody> histories = new ArrayList<>();
            ServiceConfigHistoryResponseBody history;
            for (ServiceConfigHistory serviceConfigHistory : serviceConfigHistories) {
                history = new ServiceConfigHistoryResponseBody();
                String configId = serviceConfigHistory.getConfigId();
                for (ServiceConfig serviceConfig : serviceConfigList) {
                    if (configId.equals(serviceConfig.getId())) {
                        history.setName(serviceConfig.getName());
                        break;
                    }
                }
                history.setId(serviceConfigHistory.getId());
                history.setConfigId(configId);
                String content = null;
                try {
                    content = new String(serviceConfigHistory.getContent(), "utf-8");
                } catch (UnsupportedEncodingException e) {
                    GlobalExceptionUtil.log(e);
                }
                history.setContent(content);
                history.setCreateTime(serviceConfigHistory.getCreateTime());
                history.setDesc(serviceConfigHistory.getDescription());
                history.setPushVersion(serviceConfigHistory.getPushVersion());
                histories.add(history);
            }
            result.setHistories(histories);
        }
        return result;
    }
}
