package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushFeedbackCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushHistoryCondition;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

/**
 * 服务发现推送事务
 *
 * @author sctang2
 * @create 2017-12-20 9:09
 **/
@Service
public class ServiceDiscoveryPushTransactional {
    @Autowired
    private IServiceDiscoveryPushFeedbackCondition serviceDiscoveryPushFeedbackConditionImpl;

    @Autowired
    private IServiceDiscoveryPushHistoryCondition serviceDiscoveryPushHistoryConditionImpl;

    /**
     * 删除推送
     *
     * @param pushId
     * @return
     */
    @Transactional
    public int deletePush(String pushId) {
        //删除服务发现推送反馈
        this.serviceDiscoveryPushFeedbackConditionImpl.delete(pushId);

        //通过id删除服务推送历史
        return this.serviceDiscoveryPushHistoryConditionImpl.deleteById(pushId);
    }

    /**
     * 批量删除推送
     *
     * @param pushIds
     * @return
     */
    @Transactional
    public int batchDeletePush(List<String> pushIds) {
        // 通过pushIds删除服务配置推送反馈
        this.serviceDiscoveryPushFeedbackConditionImpl.deleteByPushIds(pushIds);

        //通过ids删除服务推送历史
        return this.serviceDiscoveryPushHistoryConditionImpl.deleteByIds(pushIds);
    }
}
