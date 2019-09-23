package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IGrayGroupCondition;
import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import com.iflytek.ccr.polaris.cynosure.mapper.GrayGroupMapper;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.EditGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 服务版本条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 16:25
 **/
@Service
public class GrayGroupConditionImpl extends BaseService implements IGrayGroupCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(GrayGroupConditionImpl.class);

    @Autowired
    private GrayGroupMapper grayGroupMapper;

    @Override
    public GrayGroup add(AddGrayGroupRequestBody body) {
        String versionId = body.getVersionId();
        String userId = this.getUserId();
        String desc = body.getDesc();
        String name = body.getName();
        String content = body.getContent();
        //新增
        Date now = new Date();
        GrayGroup grayGroup = new GrayGroup();
        grayGroup.setId(SnowflakeIdWorker.getId());
        grayGroup.setVersionId(versionId);
        grayGroup.setUserId(userId);
        grayGroup.setName(name);
        grayGroup.setDescription(desc);
        grayGroup.setContent(content);
        grayGroup.setCreateTime(now);

        this.grayGroupMapper.insert(grayGroup);
        return grayGroup;
    }

    @Override
    public GrayGroup find(String name, String versionId) {
        GrayGroup grayGroup = new GrayGroup();
        grayGroup.setName(name);
        grayGroup.setVersionId(versionId);
        return this.grayGroupMapper.find(grayGroup);
    }

    @Override
    public List<GrayGroup> findList(HashMap<String, Object> map) {
        return this.grayGroupMapper.findList(map);
    }

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.grayGroupMapper.findTotalCount(map);
    }

    @Override
    public GrayGroup findById(String id) {
        return this.grayGroupMapper.findById(id);
    }

    @Override
    public int deleteById(String id) {
        return this.grayGroupMapper.deleteById(id);
    }

    @Override
    public GrayGroup updateById(String id, EditGrayGroupRequestBody body) {
        String desc = body.getDesc();

        //更新
        Date now = new Date();
        GrayGroup grayGroup = new GrayGroup();
        grayGroup.setId(id);
        grayGroup.setDescription(desc);
        grayGroup.setUpdateTime(now);
        this.grayGroupMapper.updateById(grayGroup);
        return grayGroup;
    }

    @Override
    public Map<String, String> copy1(List<GrayGroup> grayGroupList, String versionId) {
        Map<String, String> map = new HashMap<>();
        String userId = this.getUserId();
        for (GrayGroup grayGroup : grayGroupList) {
            String oldId = grayGroup.getId();
            String newId = SnowflakeIdWorker.getId();
            map.put(oldId, newId);
            grayGroup.setId(newId);
            grayGroup.setVersionId(versionId);
            grayGroup.setUserId(userId);
            Date date = new Date();
            grayGroup.setCreateTime(date);
            grayGroup.setUpdateTime(date);
            this.grayGroupMapper.insert(grayGroup);
        }
        return map;
    }
}
