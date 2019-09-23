package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.AddGrayGroupAndConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.List;

/**
 * 灰度组事务操作
 * Created by DELL-5490 on 2018/7/5.
 */
@Service
public class GrayGroupTransactional extends BaseService {

    @Autowired
    private IServiceConfigCondition serviceConfigConditionImpl;

    @Autowired
    private IGrayGroupCondition grayGroupConditionImpl;

    @Autowired
    private IServiceConfigHistoryCondition serviceConfigHistoryConditionImpl;

    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IProjectMemberCondition projectMemberConditionImpl;

    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private IServiceCondition serviceConditionImpl;

    @Autowired
    private IServiceVersionCondition serviceVersionConditionImpl;

    /**
     * 新增灰度组(不含配置文件拖拽上传)
     *
     * @param body
     * @return
     */
    @Transactional
    public Response<AddGrayGroupAndConfigResponseBody> addGrayGroup(AddGrayGroupRequestBody body) {

        List<ServiceConfig> serviceConfigList;
        if (null != body.getIds() && !body.getIds().isEmpty()) {
            //查询服务配置数组信息
            serviceConfigList = this.serviceConfigConditionImpl.findByIds(body.getIds());
        } else {
            serviceConfigList = null;
        }

        //新增灰度组
        GrayGroup grayGroup = this.grayGroupConditionImpl.add(body);

        String grayGroupId = grayGroup.getId();
        //批量新增配置（新增灰度组）
        List<ServiceConfig> newServiceConfigList = this.batchGrayGroupAddServiceConfig(body, grayGroupId, serviceConfigList);

        //创建返回结果
        AddGrayGroupAndConfigResponseBody result = new AddGrayGroupAndConfigResponseBody();

        //创建灰度组结果
        GrayGroupDetailResponseBody grayGroupResult = this.createGrayGroupResult(grayGroup);
        result.setGrayGroup(grayGroupResult);

        //创建灰度配置结果
        List<ServiceConfigDetailResponseBody> configs = new ArrayList<>();
        if (null != newServiceConfigList && !newServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : newServiceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
            result.setConfigs(configs);
        }

        return new Response<>(result);
    }

    /**
     * 新增灰度组(含配置文件拖拽上传)
     *
     * @param body
     * @return
     */
    @Transactional
    public Response<AddGrayGroupAndConfigResponseBody> addGrayGroupAndFile(AddGrayGroupRequestBody body, List<FileContent> fileContentList) {
        //查询完整的配置文件信息
        List<ServiceConfig> serviceConfigList;
        if (null != body.getIds() && !body.getIds().isEmpty()) {
            serviceConfigList = this.serviceConfigConditionImpl.findByIds(body.getIds());
        } else {
            serviceConfigList = null;
        }

        //新增灰度组
        GrayGroup grayGroup = this.grayGroupConditionImpl.add(body);
        String grayGroupId = grayGroup.getId();

        //批量新增配置(含文件)
        List<ServiceConfig> newServiceConfigList = this.batchGrayGroupAddServiceConfigAndFile(body, grayGroupId, fileContentList, serviceConfigList);

        //创建快速结果
        AddGrayGroupAndConfigResponseBody result = new AddGrayGroupAndConfigResponseBody();

        //创建灰度组结果
        GrayGroupDetailResponseBody grayGroupResult = this.createGrayGroupResult(grayGroup);
        result.setGrayGroup(grayGroupResult);

        //创建服务版本配置结果
        List<ServiceConfigDetailResponseBody> configs = new ArrayList<>();
        if (null != newServiceConfigList && !newServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : newServiceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
            result.setConfigs(configs);//放外面就会返回[],放里面就会返回null
        }

        return new Response<>(result);
    }

    /**
     * 删除灰度组和灰度配置文件
     *
     * @param grayGroupId
     * @return
     */
    @Transactional
    public int deleteGroupAndConfig(String grayGroupId, List<String> configIds) {
        //通过grayId删除灰度组
        this.grayGroupConditionImpl.deleteById(grayGroupId);

        //通过grayId删除灰度配置
        this.serviceConfigConditionImpl.deleteByGrayId(grayGroupId);

        //通过configIds删除配置历史
        return this.serviceConfigHistoryConditionImpl.deleteByConfigIds(configIds);
    }

    /**
     * 批量新增配置（无文件上传）
     *
     * @param addGrayGroupRequestBody
     * @param grayGroupId
     * @param serviceConfigList
     * @return
     */
    private List<ServiceConfig> batchGrayGroupAddServiceConfig(AddGrayGroupRequestBody addGrayGroupRequestBody, String grayGroupId, List<ServiceConfig> serviceConfigList) {
        return this.serviceConfigConditionImpl.batchGrayGroupAdd(addGrayGroupRequestBody, grayGroupId, serviceConfigList);
    }

    /**
     * 批量新增配置（含文件上传）
     *
     * @param addGrayGroupRequestBody
     * @param grayGroupId
     * @param fileContentList
     * @param serviceConfigList
     * @return
     */
    private List<ServiceConfig> batchGrayGroupAddServiceConfigAndFile(AddGrayGroupRequestBody addGrayGroupRequestBody, String grayGroupId, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList) {
        return this.serviceConfigConditionImpl.batchGrayGroupAndFileAdd(addGrayGroupRequestBody, grayGroupId, serviceConfigList, fileContentList);
    }
}
