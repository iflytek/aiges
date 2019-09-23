package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.InstanceManageCondition;
import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import com.iflytek.ccr.polaris.cynosure.mapper.InstanceManageMapper;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.EditInstanceRequestBody;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.List;

/**
 * 服务版本条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 16:25
 **/
@Service
public class InstanceManageConditionImpl extends BaseService implements InstanceManageCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(InstanceManageConditionImpl.class);

    @Autowired
    private InstanceManageMapper instanceManageMapper;

    @Override
    public GrayGroup updateById(String id, EditInstanceRequestBody body) {
        String versionId = body.getVersionId();
        String content = body.getContent();
        id = body.getGrayId();
        //更新
        Date now = new Date();
        GrayGroup grayGroup = new GrayGroup();
        grayGroup.setId(id);
        grayGroup.setContent(content);
        grayGroup.setVersionId(versionId);
        grayGroup.setUpdateTime(now);
        this.instanceManageMapper.updateById(grayGroup);
        return grayGroup;
    }

    @Override
    public List<String> findTotal(String versionId, String grayGroupId) {
        return instanceManageMapper.findTotal(versionId, grayGroupId);
    }
}
