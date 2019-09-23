package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigPushFeedbackCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigPushHistoryCondition;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

/**
 * 服务配置推送事务
 *
 * @author sctang2
 * @create 2017-12-20 9:00
 **/
@Service
public class ServiceConfigPushTransactional {
    @Autowired
    private IServiceConfigPushFeedbackCondition serviceConfigPushFeedbackConditionImpl;

    @Autowired
    private IServiceConfigPushHistoryCondition serviceConfigPushHistoryConditionImpl;

    /**
     * 删除推送
     *
     * @param pushId
     * @return
     */
    @Transactional
    public int deletePush(String pushId) {
        //删除服务配置推送反馈
        this.serviceConfigPushFeedbackConditionImpl.delete(pushId);

        //通过id删除服务推送历史
        return this.serviceConfigPushHistoryConditionImpl.deleteById(pushId);
    }

    /**
     * 批量删除推送
     *
     * @param pushIds
     * @return
     */
    @Transactional
    public int batchDeletePush(List<String> pushIds) {
        //通过pushIds删除服务配置推送反馈
        this.serviceConfigPushFeedbackConditionImpl.deleteByPushIds(pushIds);

        //通过ids删除服务推送历史
        return this.serviceConfigPushHistoryConditionImpl.deleteByIds(pushIds);
    }
}
