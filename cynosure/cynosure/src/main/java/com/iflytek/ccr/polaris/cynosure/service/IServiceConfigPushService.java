package com.iflytek.ccr.polaris.cynosure.service;

/**
 * create by ygli3
 */

import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushHistory;

import java.util.List;

public interface IServiceConfigPushService {
    List<ServiceConfigPushHistory> getConfigPushHistoryList(String projectName, String clusterName, String serviceName, String version, Integer filterGray, Integer currentPage, Integer pageSize);
}
