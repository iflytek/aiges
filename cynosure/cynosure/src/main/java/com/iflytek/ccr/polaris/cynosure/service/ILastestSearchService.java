package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;

/**
 * 最近搜索业务逻辑接口
 *
 * @author sctang2
 * @create 2018-02-05 15:49
 **/
public interface ILastestSearchService {
    /**
     * 查询最近的搜索
     *
     * @param projectName
     * @return
     */
    SearchCondition find(String projectName);

    /**
     * 保存最近的搜索
     *
     * @param projectName
     * @return
     */
    String saveLastestSearch(String projectName);

    /**
     * 查询最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @return
     */
    SearchCondition find(String projectName, String clusterName);

    /**
     * 保存最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @return
     */
    String saveLastestSearch(String projectName, String clusterName);

    /**
     * 查询最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @param serviceName
     * @return
     */
    SearchCondition find(String projectName, String clusterName, String serviceName);

    /**
     * 保存最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @param serviceName
     * @return
     */
    String saveLastestSearch(String projectName, String clusterName, String serviceName);

    /**
     * 查询最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @param serviceName
     * @param versionName
     * @return
     */
    SearchCondition find(String projectName, String clusterName, String serviceName, String versionName);

    /**
     * 保存最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @param serviceName
     * @param versionName
     * @return
     */
    String saveLastestSearch(String projectName, String clusterName, String serviceName, String versionName);

    /**
     * 查询最近的搜索
     *
     * @param projectName
     * @param clusterName
     * @param serviceName
     * @param versionName
     * @param grayName
     * @return
     */
    SearchCondition find(String projectName, String clusterName, String serviceName, String versionName, String grayName);

    /**
     * @param projectName
     * @param clusterName
     * @param serviceName
     * @param versionName
     * @param grayName
     * @return
     */
    String saveLastestSearch(String projectName, String clusterName, String serviceName, String versionName, String grayName);

}