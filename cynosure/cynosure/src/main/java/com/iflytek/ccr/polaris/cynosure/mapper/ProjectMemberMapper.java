package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ProjectMember;

import java.util.HashMap;
import java.util.List;

/**
 * 项目成员持久化接口
 *
 * @author sctang2
 * @create 2018-01-15 20:05
 **/
public interface ProjectMemberMapper {
    /**
     * 新增
     *
     * @param projectMember
     * @return
     */
    int insert(ProjectMember projectMember);

    /**
     * 删除
     *
     * @param projectMember
     * @return
     */
    int delete(ProjectMember projectMember);

    /**
     * 通过项目id删除项目成员
     *
     * @param projectId
     * @return
     */
    int deleteByProjectId(String projectId);

    /**
     * 查询
     *
     * @param projectMember
     * @return
     */
    ProjectMember find(ProjectMember projectMember);

    /**
     * 查询总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询列表
     *
     * @param map
     * @return
     */
    List<ProjectMember> findList(HashMap<String, Object> map);
}
