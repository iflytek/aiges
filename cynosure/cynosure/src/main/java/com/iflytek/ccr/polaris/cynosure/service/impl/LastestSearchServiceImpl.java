package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.ILastestSearchCondition;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

/**
 * 最近搜索业务逻辑接口实现
 *
 * @author sctang2
 * @create 2018-02-05 15:51
 **/
@Service
public class LastestSearchServiceImpl extends BaseService implements ILastestSearchService {
    @Autowired
    private ILastestSearchCondition lastestSearchConditionImpl;

    @Override
    public SearchCondition find(String projectName) {
        SearchCondition searchCondition = new SearchCondition();
        if (StringUtils.isNotBlank(projectName)) {
            searchCondition.setProject(projectName);
            return searchCondition;
        }

        //查询搜索条件
        String url = this.getRequestURI();
        SearchCondition result = this.lastestSearchConditionImpl.find(url);
        if (null == result) {
            return searchCondition;
        }
        projectName = result.getProject();
        searchCondition.setProject(projectName);
        return searchCondition;
    }

    @Override
    public String saveLastestSearch(String projectName) {
        SearchCondition searchCondition = new SearchCondition();
        searchCondition.setProject(projectName);
        String condition = JSON.toJSONString(searchCondition);
        String url = this.getRequestURI();
        //同步搜索条件
        this.lastestSearchConditionImpl.syncSearchCondition(url, condition);
        return condition;
    }

    @Override
    public SearchCondition find(String projectName, String clusterName) {
        SearchCondition searchCondition = new SearchCondition();
        if (StringUtils.isNotBlank(projectName) && StringUtils.isNotBlank(clusterName)) {
            searchCondition.setProject(projectName);
            searchCondition.setCluster(clusterName);
            return searchCondition;
        }
        if (StringUtils.isBlank(projectName) && StringUtils.isBlank(clusterName)) {
            //查询搜索条件
            String url = this.getRequestURI();
            SearchCondition result = this.lastestSearchConditionImpl.find(url);
            if (null == result) {
                return searchCondition;
            } else {
                projectName = result.getProject();
                clusterName = result.getCluster();
            }
        }
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        return searchCondition;
    }

    @Override
    public String saveLastestSearch(String projectName, String clusterName) {
        SearchCondition searchCondition = new SearchCondition();
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        String condition = JSON.toJSONString(searchCondition);
        String url = this.getRequestURI();
        //同步搜索条件
        this.lastestSearchConditionImpl.syncSearchCondition(url, condition);
        return condition;
    }

    @Override
    public SearchCondition find(String projectName, String clusterName, String serviceName) {
        SearchCondition searchCondition = new SearchCondition();
        if (StringUtils.isNotBlank(projectName) && StringUtils.isNotBlank(clusterName) && StringUtils.isNotBlank(serviceName)) {
            searchCondition.setProject(projectName);
            searchCondition.setCluster(clusterName);
            searchCondition.setService(serviceName);
            return searchCondition;
        }
        if (StringUtils.isBlank(projectName) && StringUtils.isBlank(clusterName) && StringUtils.isBlank(serviceName)) {
            //查询搜索条件
            String url = this.getRequestURI();
            SearchCondition result = this.lastestSearchConditionImpl.find(url);
            if (null == result) {
                return searchCondition;
            } else {
                projectName = result.getProject();
                clusterName = result.getCluster();
                serviceName = result.getService();
            }
        }
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        searchCondition.setService(serviceName);
        return searchCondition;
    }

    @Override
    public String saveLastestSearch(String projectName, String clusterName, String serviceName) {
        SearchCondition searchCondition = new SearchCondition();
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        searchCondition.setService(serviceName);
        String condition = JSON.toJSONString(searchCondition);
        String url = this.getRequestURI();
        //同步搜索条件
        this.lastestSearchConditionImpl.syncSearchCondition(url, condition);
        return condition;
    }

    @Override
    public SearchCondition find(String projectName, String clusterName, String serviceName, String versionName) {
        SearchCondition searchCondition = new SearchCondition();
        if (StringUtils.isNotBlank(projectName) && StringUtils.isNotBlank(clusterName) && StringUtils.isNotBlank(serviceName) && StringUtils.isNotBlank(versionName)) {
            searchCondition.setProject(projectName);
            searchCondition.setCluster(clusterName);
            searchCondition.setService(serviceName);
            searchCondition.setVersion(versionName);
            return searchCondition;
        }
        /**
         * 1.此处代码的意思是，当输入的参数（project，cluster，service，version）有一个是空格时，
         * 从数据库中查询本用户在最近一次查询记录，将最近一次查询记录作为查询条件返回。
         */
        if (StringUtils.isBlank(projectName) && StringUtils.isBlank(clusterName) && StringUtils.isBlank(serviceName) && StringUtils.isBlank(versionName)) {
            //查询搜索条件
            String url = this.getRequestURI();
            SearchCondition result = this.lastestSearchConditionImpl.find(url);
            if (null == result) {
                return searchCondition;
            } else {
                projectName = result.getProject();
                clusterName = result.getCluster();
                serviceName = result.getService();
                versionName = result.getVersion();
            }
        }
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        searchCondition.setService(serviceName);
        searchCondition.setVersion(versionName);
        return searchCondition;
    }

    @Override
    public String saveLastestSearch(String projectName, String clusterName, String serviceName, String versionName) {
        SearchCondition searchCondition = new SearchCondition();
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        searchCondition.setService(serviceName);
        searchCondition.setVersion(versionName);
        String condition = JSON.toJSONString(searchCondition);
        String url = this.getRequestURI();
        //同步搜索条件
        this.lastestSearchConditionImpl.syncSearchCondition(url, condition);
        return condition;
    }

    @Override
    public SearchCondition find(String projectName, String clusterName, String serviceName, String versionName, String grayName) {
        SearchCondition searchCondition = new SearchCondition();
        if (StringUtils.isNotBlank(projectName) && StringUtils.isNotBlank(clusterName) && StringUtils.isNotBlank(serviceName) && StringUtils.isNotBlank(versionName) && StringUtils.isNotBlank(versionName)) {
            searchCondition.setProject(projectName);
            searchCondition.setCluster(clusterName);
            searchCondition.setService(serviceName);
            searchCondition.setVersion(versionName);
            searchCondition.setGray(grayName);
            return searchCondition;
        }
        if (StringUtils.isBlank(projectName) && StringUtils.isBlank(clusterName) && StringUtils.isBlank(serviceName) && StringUtils.isBlank(versionName) && StringUtils.isNotBlank(versionName)) {
            //查询搜索条件
            String url = this.getRequestURI();
            SearchCondition result = this.lastestSearchConditionImpl.find(url);
            if (null == result) {
                return searchCondition;
            } else {
                projectName = result.getProject();
                clusterName = result.getCluster();
                serviceName = result.getService();
                versionName = result.getVersion();
                grayName = result.getGray();
            }
        }
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        searchCondition.setService(serviceName);
        searchCondition.setVersion(versionName);
        searchCondition.setGray(grayName);
        return searchCondition;
    }

    @Override
    public String saveLastestSearch(String projectName, String clusterName, String serviceName, String versionName, String grayName) {
        SearchCondition searchCondition = new SearchCondition();
        searchCondition.setProject(projectName);
        searchCondition.setCluster(clusterName);
        searchCondition.setService(serviceName);
        searchCondition.setVersion(versionName);
        searchCondition.setGray(grayName);
        String condition = JSON.toJSONString(searchCondition);
        String url = this.getRequestURI();
        //同步搜索条件
        this.lastestSearchConditionImpl.syncSearchCondition(url, condition);
        return condition;
    }
}
