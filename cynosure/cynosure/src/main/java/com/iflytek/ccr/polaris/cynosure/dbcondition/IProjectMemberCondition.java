package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.ProjectMember;

import java.util.HashMap;
import java.util.List;

/**
 * 项目成员条件接口
 *
 * @author sctang2
 * @create 2018-01-17 10:27
 **/
public interface IProjectMemberCondition {
    /**
     * 查询项目成员总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询项目成员列表
     *
     * @param map
     * @return
     */
    List<ProjectMember> findList(HashMap<String, Object> map);

    /**
     * 根据userId，projectId查询项目成员
     *
     * @param userId
     * @param projectId
     * @return
     */
    ProjectMember find(String userId, String projectId);

    /**
     * 根据userId，projectId删除项目成员
     *
     * @param userId
     * @param projectId
     * @return
     */
    int delete(String userId, String projectId);

    /**
     * 通过项目id删除项目成员
     *
     * @param projectId
     * @return
     */
    int deleteByProjectId(String projectId);

    /**
     * 新增项目创建者
     *
     * @param projectId
     * @return
     */
    ProjectMember addCreator(String projectId);

    /**
     * 新增项目成员
     *
     * @param userId
     * @param projectId
     * @return
     */
    ProjectMember add(String userId, String projectId);
}
