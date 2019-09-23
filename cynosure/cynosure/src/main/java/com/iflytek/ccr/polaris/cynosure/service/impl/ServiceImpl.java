package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IClusterCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceCondition;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.CopyAndAddTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.CopyServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.EditServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.QueryServiceListRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.service.IService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;

/**
 * 服务业务逻辑接口实现
 *
 * @author sctang2
 * @create 2017-11-16 19:55
 **/
@Service
public class ServiceImpl extends BaseService implements IService {
    @Autowired
    private IServiceCondition serviceConditionImpl;

    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private CopyAndAddTransactional copyAndAddTransactional;

    @Override
    public Response<QueryPagingListResponseBody> findLastestList(QueryServiceListRequestBody body) {
        QueryPagingListResponseBody result;
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        //查询最近的搜索
        SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName);
        projectName = searchCondition.getProject();
        clusterName = searchCondition.getCluster();
        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);
        map.put("clusterName", clusterName);

        //查询总数
        int totalCount = this.serviceConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ServiceDetailResponseBody> list = new ArrayList<>();
        Optional<List<com.iflytek.ccr.polaris.cynosure.domain.Service>> serviceList = Optional.ofNullable(this.serviceConditionImpl.findList(map));
        serviceList.ifPresent(x -> {
            x.forEach(y -> {
                //创建服务结果
                ServiceDetailResponseBody serviceDetail = this.createServiceResult(y);
                list.add(serviceDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList(QueryServiceListRequestBody body) {
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName)) {
            return new Response<>(PagingUtil.createResult(body, 0));
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);
        map.put("clusterName", clusterName);

        //查询总数
        int totalCount = this.serviceConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ServiceDetailResponseBody> list = new ArrayList<>();
        Optional<List<com.iflytek.ccr.polaris.cynosure.domain.Service>> serviceList = Optional.ofNullable(this.serviceConditionImpl.findList(map));
        serviceList.ifPresent(x -> {
            x.forEach(y -> {
                //创建服务结果
                ServiceDetailResponseBody serviceDetail = this.createServiceResult(y);
                list.add(serviceDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<ServiceDetailResponseBody> find(String id) {
        //通过id查询服务
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.findById(id);
        if (null == service) {
            //不存在该服务
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
        }

        //创建服务结果
        ServiceDetailResponseBody result = this.createServiceResult(service);
        return new Response<>(result);
    }

    @Override
    public Response<ServiceDetailResponseBody> add(AddServiceRequestBody body) {
        //根据id查询集群信息
        String clusterId = body.getClusterId();
        Cluster cluster = this.clusterConditionImpl.findById(clusterId);
        if (null == cluster) {
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }

        //根据服务名和集群id查询服务
        String name = body.getName();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.find(name, clusterId);
        if (null != service) {
            //已存在该服务
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_EXISTS, SystemErrCode.ERRMSG_SERVICE_EXISTS);
        }

        //创建服务
        com.iflytek.ccr.polaris.cynosure.domain.Service newService = this.serviceConditionImpl.add(body);

        //创建服务结果
        ServiceDetailResponseBody result = this.createServiceResult(newService);
        return new Response<>(result);
    }

    @Override
    public Response<ServiceDetailResponseBody> edit(EditServiceRequestBody body) {
        //通过id查询服务
        String id = body.getId();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.findById(id);
        if (null == service) {
            //不存在该服务
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
        }

        //修改服务信息
        com.iflytek.ccr.polaris.cynosure.domain.Service updateService = this.serviceConditionImpl.updateById(id, body);

        //创建服务结果
        updateService.setName(service.getName());
        updateService.setGroupId(service.getGroupId());
        updateService.setCreateTime(service.getCreateTime());
        ServiceDetailResponseBody result = this.createServiceResult(updateService);
        return new Response<>(result);
    }

    @Override
    public Response<String> delete(IdRequestBody body) {
        //通过id查询版本列表
        String id = body.getId();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.findServiceVersionListById(id);
        if (null == service) {
            //不存在该服务
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
        }

        List<ServiceVersion> serviceVersionList = service.getServiceVersionList();
        if (null != serviceVersionList && !serviceVersionList.isEmpty()) {
            //该用户已经创建版本
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_CREATE, SystemErrCode.ERRMSG_SERVICE_VERSION_CREATE);
        }

        //根据id删除服务
        this.serviceConditionImpl.deleteById(id);
        return new Response<>(null);
    }

    @Override
    public Response<ServiceDetailResponseBody> copy(CopyServiceRequestBody body) {
        //根据id查询集群，若不存在直接返回
        String clusterId = body.getClusterId();
        Cluster cluster = this.clusterConditionImpl.findById(clusterId);
        if (null == cluster) {
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }

        //根据服务名字serviceName和集群clusterId查询表tb_service,
        // 判断要新增的服务是否已经存在，若是，直接返回
        String serviceName = body.getServiceName();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.find(serviceName, clusterId);
        if (null != service) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_EXISTS, SystemErrCode.ERRMSG_SERVICE_EXISTS);
        }

        //查询被复制的服务是否存在，若不存在，直接返回
        String oldServiceId = body.getOldServiceId();
        com.iflytek.ccr.polaris.cynosure.domain.Service serviceCopy = this.serviceConditionImpl.findById(oldServiceId);
        if (null == serviceCopy) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_COPY_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_COPY_NOT_EXISTS);
        }

        //复制服务
        com.iflytek.ccr.polaris.cynosure.domain.Service newService = this.copyAndAddTransactional.copyAndAddService(body);

        //返回复制的结果
        ServiceDetailResponseBody result = this.createServiceResult(newService);
        return new Response<>(result);
    }
}
