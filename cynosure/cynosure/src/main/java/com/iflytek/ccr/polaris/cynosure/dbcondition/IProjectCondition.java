package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.EditProjectRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 项目条件接口
 *
 * @author sctang2
 * @create 2017-12-09 16:58
 **/
public interface IProjectCondition {
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
    List<Project> findList(HashMap<String, Object> map);

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
     * 新增项目
     *
     * @param body
     * @return
     */
    Project add(AddProjectRequestBody body);

    /**
     * 根据id更新项目
     *
     * @param id
     * @param body
     * @return
     */
    Project updateById(String id, EditProjectRequestBody body);

    /**
     * 通过id查询集群和成员列表
     *
     * @param id
     * @return
     */
    Project findProjectAndClusterListById(String id);

    /**
     * 通过id删除项目
     *
     * @param id
     * @return
     */
    int deleteById(String id);
}
