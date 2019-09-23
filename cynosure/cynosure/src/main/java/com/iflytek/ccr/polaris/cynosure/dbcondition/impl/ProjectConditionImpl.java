package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectCondition;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.mapper.ProjectMapper;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.EditProjectRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 项目条件接口实现
 *
 * @author sctang2
 * @create 2017-12-09 16:59
 **/
@Service
public class ProjectConditionImpl extends BaseService implements IProjectCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ProjectConditionImpl.class);

    @Autowired
    private ProjectMapper projectMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.projectMapper.findTotalCount(map);
    }

    @Override
    public List<Project> findList(HashMap<String, Object> map) {
        return this.projectMapper.findList(map);
    }

    @Override
    public Project findById(String id) {
        return this.projectMapper.findById(id);
    }

    @Override
    public Project findByName(String name) {
        return this.projectMapper.findByName(name);
    }

    @Override
    public Project add(AddProjectRequestBody body) {
        String name = body.getName();
        String desc = body.getDesc();
        String userId = this.getUserId();

        //新增
        Date now = new Date();
        Project project = new Project();
        project.setId(SnowflakeIdWorker.getId());
        project.setCreateTime(now);
        project.setDescription(desc);
        project.setName(name);
        project.setUserId(userId);
        try {
            this.projectMapper.insert(project);
            return project;
        } catch (DuplicateKeyException ex) {
            logger.warn("project duplicate key " + ex.getMessage());
            project = this.findByName(name);
            return project;
        }
    }

    @Override
    public Project updateById(String id, EditProjectRequestBody body) {
        String desc = body.getDesc();

        //更新
        Date now = new Date();
        Project project = new Project();
        project.setId(id);
        project.setDescription(desc);
        project.setUpdateTime(now);
        this.projectMapper.updateById(project);
        return project;
    }

    @Override
    public Project findProjectAndClusterListById(String id) {
        return this.projectMapper.findProjectAndClusterListById(id);
    }

    @Override
    public int deleteById(String id) {
        return this.projectMapper.deleteById(id);
    }
}
