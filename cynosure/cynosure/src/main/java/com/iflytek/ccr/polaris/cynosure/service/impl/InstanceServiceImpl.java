package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.GrayConfigCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceProviderConsumerResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceResult;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IGrayGroupCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IRegionCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.InstanceManageCondition;
import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.AddGrayGroupInstanceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody.EditInstanceRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.InstanceService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashSet;
import java.util.List;

/**
 * Created by DELL-5490 on 2018/7/7.
 */
@Service
public class InstanceServiceImpl extends BaseService implements InstanceService {
    @Autowired
    private IGrayGroupCondition grayGroupConditionImpl;

    @Autowired
    private InstanceManageCondition instanceManageConditionImpl;

    @Autowired
    private IRegionCondition regionConditionImpl;

    @Autowired
    private GrayConfigCenter grayConfigCenter;

    @Override
    public Response<String> findById(String id) {
        //根据id查询灰度组
        GrayGroup grayGroup = this.grayGroupConditionImpl.findById(id);
        if (null == grayGroup) {
            //不存在该灰度组
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
        }
        //创建推送实例结果
        String content = grayGroup.getContent();

        return new Response<>(content);
    }

    @Override
    public Response<String> findList(String versionId) {
        //获取推送组全部实例
        List<String> contentList = this.instanceManageConditionImpl.findTotal(versionId, null);
        List<String> oldGroupList = new ArrayList<>();
        for (String content : contentList) {
            oldGroupList.addAll(Arrays.asList(StringUtils.split(content, ",")));
        }
        //用逗号分隔转化为字符串，未去重
        String contentJoint = String.join(",", oldGroupList);

        //hashset去重
        HashSet<String> contentSet = new HashSet<>(Arrays.asList(contentJoint.split(",")));

        //转化为字符串返回
        String contentDistinct = String.join(",", contentSet);
        return new Response<>(contentDistinct);
    }

    @Override
    public Response<GrayGroupDetailResponseBody> edit(EditInstanceRequestBody body) {
        //根据id查询灰度组
        String id = body.getGrayId();
        GrayGroup grayGroup = this.grayGroupConditionImpl.findById(id);
        if (null == grayGroup) {
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
        }

        //校验灰度组的推送实例是否有重复内容
        String instanceContent = body.getContent();
        if (StringUtils.isNotBlank(instanceContent)) {
            List<String> contentList = Arrays.asList(instanceContent.split(","));
            boolean isRepeat = contentList.size() == new HashSet<>(contentList).size();
            if (!isRepeat) {
                return new Response<>(SystemErrCode.ERRCODE_GRAY_INSTANCE_REPEAT, SystemErrCode.ERRMSG_GRAY_INSTANCE_REPEAT);
            }
            //查询的是不含该灰度组自身的推送实例
            List<String> versionContent = this.instanceManageConditionImpl.findTotal(body.getVersionId(), body.getGrayId());
            if (null != versionContent && !versionContent.isEmpty()) {
                //校验老的版本是下的推送实例是否重复
                List<String> oldGroupList = new ArrayList<>();
                for (String content : versionContent) {
                    oldGroupList.addAll(Arrays.asList(StringUtils.split(content, ",")));
                }
                oldGroupList.addAll(contentList);
                HashSet<String> compareSet = new HashSet<>(oldGroupList);
                if (oldGroupList.size() != compareSet.size()) {
                    return new Response<>(SystemErrCode.ERRCODE_GRAY_INSTANCE_ARE_USED, SystemErrCode.ERRMSG_GRAY_INSTANCE_ARE_USED);
                }
            }
        }
        //根据id更新灰度组
        GrayGroup updateGrayGroup = this.instanceManageConditionImpl.updateById(id, body);

        //创建灰度组结果
        updateGrayGroup.setId(grayGroup.getId());
        updateGrayGroup.getContent();
        updateGrayGroup.getUpdateTime();
        updateGrayGroup.setUserId(grayGroup.getUserId());
        updateGrayGroup.getVersionId();
        updateGrayGroup.setName(grayGroup.getName());
        updateGrayGroup.setDescription(grayGroup.getDescription());
        updateGrayGroup.setCreateTime(grayGroup.getCreateTime());
        GrayGroupDetailResponseBody result = this.createGrayGroupResult(updateGrayGroup);
        return new Response<>(result);
    }

    @Override
    public Response<String> appointList(AddGrayGroupInstanceRequestBody body) {
        //通过区域名称查询区域信息
        List<Region> regionList = this.regionConditionImpl.findList(null);
        if (null == regionList) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }
        //获取服务消费者path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String version = body.getVersion();
        String path = this.grayConfigCenter.getServiceConsumerPath(project, cluster, service, version, "0");

        //查询服务消费端
        int startIndex = PagingUtil.getStartIndex(body);
        int endIndex = PagingUtil.getEndIndex(body);
        List<String> totalResult = new ArrayList<>();
        for (Region region : regionList) {
            ServiceResult serviceResult = this.grayConfigCenter.findConsumersByPaging(path, region, false, startIndex, endIndex);
            //创建消费端
            List<ServiceProviderConsumerResult> serviceProviderConsumerResults = serviceResult.getResults();
            if (null != serviceProviderConsumerResults && !serviceProviderConsumerResults.isEmpty()) {
                List<String> serviceDiscoveryConsumerList = this.createAddGroupConsumer(serviceProviderConsumerResults);
                totalResult.addAll(serviceDiscoveryConsumerList);
            }
        }

        //当前版本下已经使用的推送实例，此时没有创建灰度组，所以不用删除
        String grayId = body.getGrayId();
        String versionId = body.getVersionId();

        List<String> oldContentList = this.instanceManageConditionImpl.findTotal(versionId, null);
        //将所有实例切分装在list里
        if (null != oldContentList && !oldContentList.isEmpty()) {
            List<String> oldGroupList = new ArrayList<>();
            for (String content : oldContentList) {
                oldGroupList.addAll(Arrays.asList(StringUtils.split(content, ",")));
            }
            totalResult.removeAll(oldGroupList);
        }

        //组合的实例列表去重
        HashSet<String> NoRepeatSet = new HashSet<>(totalResult);

        //逗号分隔数据拼接字符串返回前端
        String LastContent = String.join(",", NoRepeatSet);
        return new Response<>(LastContent);
    }

    /**
     * 创建消费端
     *
     * @param cacheCenterServiceProviderConsumerResults
     * @return
     */
    private List<String> createAddGroupConsumer(List<ServiceProviderConsumerResult> cacheCenterServiceProviderConsumerResults) {
        List<String> results = new ArrayList<>();
        cacheCenterServiceProviderConsumerResults.forEach(x -> {
            String result = new String();
            result = x.getAddr();
            results.add(result);
        });
        return results;
    }
}