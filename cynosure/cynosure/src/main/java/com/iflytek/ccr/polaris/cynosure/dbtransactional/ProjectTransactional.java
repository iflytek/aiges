package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectMemberCondition;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectRequestBody;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

/**
 * 项目事务操作
 *
 * @author sctang2
 * @create 2018-01-24 9:36
 **/
@Service
public class ProjectTransactional {
    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IProjectMemberCondition projectMemberConditionImpl;

    /**
     * 新增项目
     *
     * @param body
     * @return
     */
    @Transactional
    public Project add(AddProjectRequestBody body) {
        //新增项目
        Project project = this.projectConditionImpl.add(body);

        //新增项目创建者
        String projectId = project.getId();
        this.projectMemberConditionImpl.addCreator(projectId);
        return project;
    }

    /**
     * 删除项目
     *
     * @param projectId
     * @return
     */
    @Transactional
    public int delete(String projectId) {
        //删除项目成员
        this.projectMemberConditionImpl.deleteByProjectId(projectId);

        //删除项目
        return this.projectConditionImpl.deleteById(projectId);
    }
}
