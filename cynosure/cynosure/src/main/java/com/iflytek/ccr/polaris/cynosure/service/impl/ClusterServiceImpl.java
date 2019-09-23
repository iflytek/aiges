package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IClusterCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IProjectCondition;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.CopyAndAddTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.CopyClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.EditClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.QueryClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IClusterService;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;

/**
 * 集群业务接口实现
 *
 * @author sctang2
 * @create 2017-11-15 17:39
 **/
@Service
public class ClusterServiceImpl extends BaseService implements IClusterService {
    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private CopyAndAddTransactional copyAndAddTransactional;

    @Override
    public Response<QueryPagingListResponseBody> findLastestList(QueryClusterRequestBody body) {
        QueryPagingListResponseBody result;
        String projectName = body.getProject();

        //查询最近的搜索
        SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName);
        projectName = searchCondition.getProject();
        if (StringUtils.isBlank(projectName)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);

        //查询总数
        int totalCount = this.clusterConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ClusterDetailResponseBody> list = new ArrayList<>();
        Optional<List<Cluster>> serviceGroupList = Optional.ofNullable(this.clusterConditionImpl.findList(map));
        serviceGroupList.ifPresent(x -> {
            x.forEach(y -> {
                //创建集群结果
                ClusterDetailResponseBody serviceGroupDetail = this.createClusterResult(y);
                list.add(serviceGroupDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList(QueryClusterRequestBody body) {
        String projectName = body.getProject();
        if (StringUtils.isBlank(projectName)) {
            return new Response<>(PagingUtil.createResult(body, 0));
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);

        //查询总数
        int totalCount = this.clusterConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ClusterDetailResponseBody> list = new ArrayList<>();
        Optional<List<Cluster>> serviceGroupList = Optional.ofNullable(this.clusterConditionImpl.findList(map));
        serviceGroupList.ifPresent(x -> {
            x.forEach(y -> {
                //创建集群结果
                ClusterDetailResponseBody serviceGroupDetail = this.createClusterResult(y);
                list.add(serviceGroupDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<ClusterDetailResponseBody> find(String id) {
        //根据id查询集群信息
        Cluster cluster = this.clusterConditionImpl.findById(id);
        if (null == cluster) {
            //不存在该集群
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }

        //创建集群结果
        ClusterDetailResponseBody result = this.createClusterResult(cluster);
        return new Response<>(result);
    }

    @Override
    public Response<ClusterDetailResponseBody> add(AddClusterRequestBody body) {
        //通过id查询项目信息
        String projectId = body.getProjectId();
        Project project = this.projectConditionImpl.findById(projectId);
        if (null == project) {
            //不存在项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        String name = body.getName();
        //根据id和集群名称查询集群信息
        Cluster cluster = this.clusterConditionImpl.find(projectId, name);
        if (null != cluster) {
            //已存在该集群
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_EXISTS, SystemErrCode.ERRMSG_CLUSTER_EXISTS);
        }

        //创建集群
        Cluster newCluster = this.clusterConditionImpl.add(body);

        //创建集群结果
        ClusterDetailResponseBody result = this.createClusterResult(newCluster);
        return new Response<>(result);
    }

    @Override
    public Response<ClusterDetailResponseBody> edit(EditClusterRequestBody body) {
        //根据id查询集群信息
        String id = body.getId();
        Cluster cluster = this.clusterConditionImpl.findById(id);
        if (null == cluster) {
            //不存在该集群
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }
        String projectId = cluster.getProjectId();

        //根据id更新集群
        Cluster updateCluster = this.clusterConditionImpl.updateById(id, body);

        //创建集群结果
        updateCluster.setName(cluster.getName());
        updateCluster.setCreateTime(cluster.getCreateTime());
        updateCluster.setProjectId(projectId);
        ClusterDetailResponseBody result = this.createClusterResult(updateCluster);
        return new Response<>(result);
    }

    @Override
    public Response<String> delete(IdRequestBody body) {
        //通过id查询服务列表
        String id = body.getId();
        Cluster cluster = this.clusterConditionImpl.findServiceListById(id);
        if (null == cluster) {
            //不存在该集群
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }

        List<com.iflytek.ccr.polaris.cynosure.domain.Service> serviceList = cluster.getServiceList();
        if (null != serviceList && !serviceList.isEmpty()) {
            //该用户已创建服务
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CREATE, SystemErrCode.ERRMSG_SERVICE_CREATE);
        }

        //通过id删除集群
        this.clusterConditionImpl.deleteById(id);
        return new Response<>(null);
    }

    @Override
    public Response<ClusterDetailResponseBody> copy(CopyClusterRequestBody body) {
        //根据projectId查询tb_project表，若该project不存在，直接返回
        String projectId = body.getProjectId();
        Project project = this.projectConditionImpl.findById(projectId);
        if (null == project) {
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //根据集群名字clusterName和项目projectId查询tb_service_group表，若已存在该集群，直接返回
        String clusterName = body.getClusterName();
        Cluster cluster = this.clusterConditionImpl.find(projectId, clusterName);
        if (null != cluster) {
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_EXISTS, SystemErrCode.ERRMSG_CLUSTER_EXISTS);
        }

        //根据oldClusterId查询被复制的集群是否存在，若不存在，直接返回
        String oldClusterId = body.getOldClusterId();
        Cluster clusterCopy = this.clusterConditionImpl.findById(oldClusterId);
        if (null == clusterCopy) {
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_COPY_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_COPY_NOT_EXISTS);
        }

        //复制集群
        Cluster newCluster = this.copyAndAddTransactional.copyAndAddCluster(body);

        //创建复制结果
        ClusterDetailResponseBody result = this.createClusterResult(newCluster);
        return new Response<>(result);
    }
}
