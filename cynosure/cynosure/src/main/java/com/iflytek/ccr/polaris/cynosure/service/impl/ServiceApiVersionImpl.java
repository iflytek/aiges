package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceApiVersionCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceApiVersionMapper;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.AddServiceApiVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.QueryServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceApiVersionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;

/**
 * 服务历史业务接口实现
 *
 * @author sctang2
 * @create 2017-11-17 16:03
 **/
@Service
public class ServiceApiVersionImpl extends BaseService implements IServiceApiVersion {

    @Autowired
    private IServiceApiVersionCondition serviceApiVersionConditionImpl;

    @Autowired
    private IServiceCondition serviceConditionimpl;

    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private ServiceApiVersionMapper apiVersionMapper;

    @Override
    public Response<ServiceApiVersionDetailResponseBody> add(AddServiceApiVersionRequestBody body) {
        //通过id查询服务
        String serviceId = body.getServiceId();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionimpl.findById(serviceId);
        if (null == service) {
            //不存在该服务
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
        }

        //根据版本和服务id查询服务版本
        String apiVersion = body.getApiVersion();
        ServiceApiVersion serviceApiVersion = this.serviceApiVersionConditionImpl.find(apiVersion, serviceId);
        if (null != serviceApiVersion) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_EXISTS);
        }

        //创建版本
        ServiceApiVersion newServiceApiVersion = this.serviceApiVersionConditionImpl.add(body);

        //创建版本结果
        ServiceApiVersionDetailResponseBody result = this.createServiceApiVersionResult(newServiceApiVersion);
        return new Response<>(result);
    }
    

    @Override
    public Response<QueryPagingListResponseBody> findLastestList(QueryServiceVersionRequestBody body) {
        QueryPagingListResponseBody result;
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        //查询最近的搜索
        SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName);
        projectName = searchCondition.getProject();
        clusterName = searchCondition.getCluster();
        serviceName = searchCondition.getService();
        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);
        map.put("clusterName", clusterName);
        map.put("serviceName", serviceName);

        //查询总数
        int totalCount = this.serviceApiVersionConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ServiceApiVersionDetailResponseBody> list = new ArrayList<>();
        Optional<List<ServiceApiVersion>> serviceApiVersionList = Optional.ofNullable(this.serviceApiVersionConditionImpl.findList(map));
        serviceApiVersionList.ifPresent(x -> {
            x.forEach(y -> {
                //创建版本结果
                ServiceApiVersionDetailResponseBody serviceApiVersionDetail = this.createServiceApiVersionResult(y);
                list.add(serviceApiVersionDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList(QueryServiceVersionRequestBody body) {
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
            return new Response<>(PagingUtil.createResult(body, 0));
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);
        map.put("clusterName", clusterName);
        map.put("serviceName", serviceName);

        //查询总数
        int totalCount = this.serviceApiVersionConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<ServiceApiVersionDetailResponseBody> list = new ArrayList<>();
        Optional<List<ServiceApiVersion>> serviceApiVersionList = Optional.ofNullable(this.serviceApiVersionConditionImpl.findList(map));
        serviceApiVersionList.ifPresent(x -> {
            x.forEach(y -> {
                //创建版本结果
                ServiceApiVersionDetailResponseBody serviceApiVersionDetail = this.createServiceApiVersionResult(y);
                list.add(serviceApiVersionDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public List<ServiceApiVersion> findList1(List<String> serviceIds) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("serviceIds", serviceIds);
        return this.apiVersionMapper.findApiVersionList(map);
    }


    @Override
    public Response<ServiceApiVersionDetailResponseBody> find(String id) {
        //根据id查询版本
        ServiceApiVersion serviceApiVersion = this.serviceApiVersionConditionImpl.findById(id);
        if (null == serviceApiVersion) {
            //不存在该版本
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //创建版本结果
        ServiceApiVersionDetailResponseBody result = this.createServiceApiVersionResult(serviceApiVersion);
        return new Response<>(result);
    }

    @Override
    public Response<ServiceApiVersionDetailResponseBody> edit(EditServiceVersionRequestBody body) {
        //根据id查询版本
        String id = body.getId();
        ServiceApiVersion serviceApiVersion = this.serviceApiVersionConditionImpl.findById(id);
        if (null == serviceApiVersion) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //根据id更新版本
        ServiceApiVersion updateServiceApiVersion = this.serviceApiVersionConditionImpl.updateById(id, body);

        //创建服务版本结果
        updateServiceApiVersion.setServiceId(id);
        updateServiceApiVersion.setCreateTime(serviceApiVersion.getCreateTime());
        updateServiceApiVersion.setApiVersion(serviceApiVersion.getApiVersion());
        ServiceApiVersionDetailResponseBody result = this.createServiceApiVersionResult(updateServiceApiVersion);
        return new Response<>(result);
    }

    @Override
    public Response<String> delete(IdRequestBody body) {
        //通过id查询服务配置列表
        String id = body.getId();
        ServiceApiVersion serviceApiVersion = this.serviceApiVersionConditionImpl.findServiceConfigListById(id);
        if (null == serviceApiVersion) {
            //不存在该版本
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //根据id删除版本
        this.serviceApiVersionConditionImpl.deleteById(id);
        return new Response<>(null);
    }
    /**
     * 创建服务版本结果
     *
     * @param serviceApiVersion
     * @return
     */
    protected ServiceApiVersionDetailResponseBody createServiceApiVersionResult(ServiceApiVersion serviceApiVersion) {
        ServiceApiVersionDetailResponseBody result = new ServiceApiVersionDetailResponseBody();
        result.setId(serviceApiVersion.getId());
        result.setApiVersion(serviceApiVersion.getApiVersion());
        result.setServiceId(serviceApiVersion.getServiceId());
        result.setDesc(serviceApiVersion.getDescription());
        result.setCreateTime(serviceApiVersion.getCreateTime());
        result.setUpdateTime(serviceApiVersion.getUpdateTime());
        return result;
    }
}