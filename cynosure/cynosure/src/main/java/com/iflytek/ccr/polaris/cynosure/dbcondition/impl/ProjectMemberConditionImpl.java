package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectMemberCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ProjectMember;
import com.iflytek.ccr.polaris.cynosure.enums.DBEnumInt;
import com.iflytek.ccr.polaris.cynosure.mapper.ProjectMemberMapper;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 用户所在项目条件接口实现
 *
 * @author sctang2
 * @create 2018-01-17 10:28
 **/
@Service
public class ProjectMemberConditionImpl extends BaseService implements IProjectMemberCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ProjectMemberConditionImpl.class);

    @Autowired
    private ProjectMemberMapper projectMemberMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.projectMemberMapper.findTotalCount(map);
    }

    @Override
    public List<ProjectMember> findList(HashMap<String, Object> map) {
        return this.projectMemberMapper.findList(map);
    }

    @Override
    public ProjectMember find(String userId, String projectId) {
        ProjectMember projectMember = new ProjectMember();
        projectMember.setUserId(userId);
        projectMember.setProjectId(projectId);
        return this.projectMemberMapper.find(projectMember);
    }

    @Override
    public int delete(String userId, String projectId) {
        ProjectMember projectMember = new ProjectMember();
        projectMember.setUserId(userId);
        projectMember.setProjectId(projectId);
        return this.projectMemberMapper.delete(projectMember);
    }

    @Override
    public int deleteByProjectId(String projectId) {
        return this.projectMemberMapper.deleteByProjectId(projectId);
    }

    @Override
    public ProjectMember addCreator(String projectId) {
        String userId = this.getUserId();
        ProjectMember projectMember = new ProjectMember();
        projectMember.setId(SnowflakeIdWorker.getId());
        projectMember.setUserId(userId);
        projectMember.setCreator((byte) DBEnumInt.CREATOR_Y.getIndex());
        projectMember.setProjectId(projectId);
        projectMember.setCreateTime(new Date());
        try {
            this.projectMemberMapper.insert(projectMember);
            return projectMember;
        } catch (DuplicateKeyException ex) {
            logger.warn("project member duplicate key " + ex.getMessage());
            return this.find(userId, projectId);
        }
    }

    @Override
    public ProjectMember add(String userId, String projectId) {
        ProjectMember projectMember = new ProjectMember();
        projectMember.setId(SnowflakeIdWorker.getId());
        projectMember.setUserId(userId);
        projectMember.setCreator((byte) DBEnumInt.CREATOR_N.getIndex());
        projectMember.setProjectId(projectId);
        projectMember.setCreateTime(new Date());
        try {
            this.projectMemberMapper.insert(projectMember);
            return projectMember;
        } catch (DuplicateKeyException ex) {
            logger.warn("project member duplicate key " + ex.getMessage());
            return this.find(userId, projectId);
        }
    }
}
