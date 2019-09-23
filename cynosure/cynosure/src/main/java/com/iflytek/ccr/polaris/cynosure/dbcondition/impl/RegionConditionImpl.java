package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IRegionCondition;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.mapper.RegionMapper;
import com.iflytek.ccr.polaris.cynosure.request.region.AddRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.region.EditRegionRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 区域条件接口实现
 *
 * @author sctang2
 * @create 2017-12-09 15:09
 **/
@Service
public class RegionConditionImpl extends BaseService implements IRegionCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(RegionConditionImpl.class);

    @Autowired
    private RegionMapper regionMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.regionMapper.findTotalCount(map);
    }

    @Override
    public List<Region> findList(HashMap<String, Object> map) {
        return this.regionMapper.findList(map);
    }

    @Override
    public Region findById(String id) {
        return this.regionMapper.findById(id);
    }

    @Override
    public Region findByName(String name) {
        return this.regionMapper.findByName(name);
    }

    @Override
    public Region findByPushUrl(String pushUrl) {
        return this.regionMapper.findByPushUrl(pushUrl);
    }

    @Override
    public Region add(AddRegionRequestBody body) {
        Date now = new Date();
        String name = body.getName();
        String pushUrl = body.getPushUrl();
        String userId = this.getUserId();

        //新增
        Region region = new Region();
        region.setCreateTime(now);
        region.setUpdateTime(now);
        region.setId(SnowflakeIdWorker.getId());
        region.setName(name);
        region.setPushUrl(pushUrl);
        region.setUserId(userId);
        try {
            this.regionMapper.insert(region);
            return region;
        } catch (DuplicateKeyException ex) {
            logger.warn("region duplicate key " + ex.getMessage());
            return this.findByName(name);
        }
    }

    @Override
    public Region updateById(String id, EditRegionRequestBody body) {
        String pushUrl = body.getPushUrl();

        //更新
        Date now = new Date();
        Region region = new Region();
        region.setId(id);
        region.setPushUrl(pushUrl);
        region.setUpdateTime(now);
        this.regionMapper.updateById(region);
        return region;
    }

    @Override
    public int deleteById(String id) {
        return this.regionMapper.deleteById(id);
    }

    @Override
    public List<Region> findListByIds(List<String> ids) {
        if (null == ids || ids.isEmpty()) {
            return null;
        }
        HashMap<String, Object> map = new HashMap<>();
        map.put("ids", ids);
        return this.regionMapper.findList(map);
    }
}
