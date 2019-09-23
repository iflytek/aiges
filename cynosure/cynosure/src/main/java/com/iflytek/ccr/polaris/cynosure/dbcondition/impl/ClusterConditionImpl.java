package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IClusterCondition;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.mapper.ClusterMapper;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.EditClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 集群条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 13:51
 **/
@Service
public class ClusterConditionImpl extends BaseService implements IClusterCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ClusterConditionImpl.class);

    @Autowired
    private ClusterMapper clusterMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.clusterMapper.findTotalCount(map);
    }

    @Override
    public List<Cluster> findList(HashMap<String, Object> map) {
        return this.clusterMapper.findList(map);
    }

    @Override
    public List<Cluster> findList(List<String> projectIds) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("projectIds", projectIds);
        return this.clusterMapper.findClusterList(map);
    }

    @Override
    public Cluster findById(String id) {
        return this.clusterMapper.findById(id);
    }

    @Override
    public Cluster find(String projectId, String name) {
        Cluster cluster = new Cluster();
        cluster.setProjectId(projectId);
        cluster.setName(name);
        return this.clusterMapper.find(cluster);
    }

    @Override
    public Cluster add(AddClusterRequestBody body) {
        String name = body.getName();
        String projectId = body.getProjectId();
        String desc = body.getDesc();
        String userId = this.getUserId();

        //新增
        Date now = new Date();
        Cluster cluster = new Cluster();
        cluster.setId(SnowflakeIdWorker.getId());
        cluster.setName(name);
        cluster.setDescription(desc);
        cluster.setCreateTime(now);
        cluster.setUpdateTime(now);
        cluster.setProjectId(projectId);
        cluster.setUserId(userId);
        try {
            this.clusterMapper.insert(cluster);
            return cluster;
        } catch (DuplicateKeyException ex) {
            logger.warn("cluster duplicate key " + ex.getMessage());
            return this.find(projectId, name);
        }
    }

    @Override
    public Cluster updateById(String id, EditClusterRequestBody body) {
        String desc = body.getDesc();

        //更新
        Date now = new Date();
        Cluster cluster = new Cluster();
        cluster.setId(id);
        cluster.setDescription(desc);
        cluster.setUpdateTime(now);
        this.clusterMapper.updateById(cluster);
        return cluster;
    }

    @Override
    public Cluster findServiceListById(String id) {
        return this.clusterMapper.findServiceListById(id);
    }

    @Override
    public int deleteById(String id) {
        return this.clusterMapper.deleteById(id);
    }
}
