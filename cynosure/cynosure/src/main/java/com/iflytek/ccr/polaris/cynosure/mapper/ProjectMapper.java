package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.Project;

import java.util.HashMap;
import java.util.List;

/**
 * 项目持久化接口
 *
 * @author sctang2
 * @create 2017-11-20 9:05
 **/
public interface ProjectMapper {
    /**
     * 新增项目
     *
     * @param project
     * @return
     */
    int insert(Project project);

    /**
     * 删除项目
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新项目
     *
     * @param project
     * @return
     */
    int updateById(Project project);

    /**
     * 通过id查询项目信息
     *
     * @param id
     * @return
     */
    Project findById(String id);

    /**
     * 通过项目名称查询项目信息
     *
     * @param name
     * @return
     */
    Project findByName(String name);

    /**
     * 查询项目总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询项目列表
     *
     * @param map
     * @return
     */
    List<Project> findList(HashMap<String, Object> map);

    /**
     * 通过id查询集群和成员列表
     *
     * @param id
     * @return
     */
    Project findProjectAndClusterListById(String id);
}
