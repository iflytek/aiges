package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.ConfigCenter;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigHistory;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceConfigMapper;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceConfigRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.AddServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.EditServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.io.UnsupportedEncodingException;
import java.util.*;

/**
 * 服务配置条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 17:32
 **/
@Service
public class ServiceConfigConditionImpl extends BaseService implements IServiceConfigCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceConfigConditionImpl.class);

    @Autowired
    private ServiceConfigMapper serviceConfigMapper;

    @Autowired
    private ConfigCenter configCenter;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceConfigMapper.findTotalCount(map);
    }

    @Override
    public List<ServiceConfig> findList(HashMap<String, Object> map) {
        return this.serviceConfigMapper.findList(map);
    }

    @Override
    public List<ServiceConfig> batchVersionAdd(AddServiceConfigRequestBody body, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<ServiceConfig> serviceConfigList) {

        String userId = this.getUserId();
        String versionId = body.getVersionId();
        String project = addServiceConfig.getProject();
        String cluster = addServiceConfig.getCluster();
        String service = addServiceConfig.getService();
        String version = addServiceConfig.getVersion();
        String desc = body.getDesc();
        Date now = new Date();

        List<ServiceConfig> configs = new ArrayList<>();
        //批量新增,没有拖拽新增配置文件的
        if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                serviceConfig.setId(SnowflakeIdWorker.getId());
                serviceConfig.setUserId(userId);
                serviceConfig.setVersionId(versionId);
                //获取配置路径
                String path = this.configCenter.getConfigPath(project, cluster, service, version, serviceConfig.getName());
                serviceConfig.setPath(path);

                serviceConfig.setCreateTime(now);
                serviceConfig.setUpdateTime(now);
                serviceConfig.setGrayId("0");
                configs.add(serviceConfig);
            }
            this.serviceConfigMapper.batchInsert(configs);
        }

        return configs;
    }

    @Override
    public List<ServiceConfig> batchVersionAddAndFile(AddServiceConfigRequestBody body, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<ServiceConfig> serviceConfigList, List<FileContent> fileContentList) {
        String userId = this.getUserId();
        String versionId = body.getVersionId();
        String project = addServiceConfig.getProject();
        String cluster = addServiceConfig.getCluster();
        String service = addServiceConfig.getService();
        String version = addServiceConfig.getVersion();
        String desc = body.getDesc();
        Date now = new Date();

        List<ServiceConfig> serviceConfigs = new ArrayList<>();

        //批量新增,没有拖拽新增配置文件的
        if (null != serviceConfigList) {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                serviceConfig.setId(SnowflakeIdWorker.getId());
                serviceConfig.setUserId(userId);
                serviceConfig.setVersionId(versionId);

                //获取配置路径
                String path = this.configCenter.getConfigPath(project, cluster, service, version, serviceConfig.getName());
                serviceConfig.setPath(path);

                serviceConfig.setGrayId("0");
                serviceConfig.setCreateTime(now);
                serviceConfig.setUpdateTime(now);
                serviceConfigs.add(serviceConfig);
            }
        }
        //批量新增
        ServiceConfig serviceConfig;
        for (FileContent fileContent : fileContentList) {
            serviceConfig = new ServiceConfig();
            serviceConfig.setId(SnowflakeIdWorker.getId());
            serviceConfig.setUserId(userId);
            serviceConfig.setVersionId(versionId);
            String name = fileContent.getFileName();
            serviceConfig.setName(name);
            serviceConfig.setDescription(desc);

            //获取配置路径
            String path = this.configCenter.getConfigPath(project, cluster, service, version, name);
            serviceConfig.setPath(path);

            byte[] content = fileContent.getContent();
            String md5 = MD5Util.getMD5(content);
            serviceConfig.setContent(content);
            serviceConfig.setMd5(md5);
            serviceConfig.setCreateTime(now);
            serviceConfig.setUpdateTime(now);
            serviceConfig.setGrayId("0");
            serviceConfigs.add(serviceConfig);
        }
//        if (!serviceConfigs.isEmpty()) {
//            for (ServiceConfig serviceConfigCompare : serviceConfigs) {
//                String name = serviceConfigCompare.getName();
//                String grayId = serviceConfigCompare.getGrayGroupId();
//                if (null!=this.findOnlyConfig(name, versionId, grayId)){
//                    return new
//                }
//            }
//        }
        this.serviceConfigMapper.batchInsert(serviceConfigs);
        return serviceConfigs;
    }

    @Override
    public List<ServiceConfig> batchAdd(AddServiceConfigRequestBody body, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<FileContent> fileContentList) {
        if (null == fileContentList || fileContentList.isEmpty()) {
            return null;
        }

        String userId = this.getUserId();
        String versionId = body.getVersionId();
        String project = addServiceConfig.getProject();
        String cluster = addServiceConfig.getCluster();
        String service = addServiceConfig.getService();
        String version = addServiceConfig.getVersion();
        String desc = addServiceConfig.getDesc();
        Date now = new Date();

        //批量新增
        List<ServiceConfig> serviceConfigList = new ArrayList<>();
        ServiceConfig serviceConfig;
        for (FileContent fileContent : fileContentList) {
            serviceConfig = new ServiceConfig();
            serviceConfig.setId(SnowflakeIdWorker.getId());
            serviceConfig.setUserId(userId);
            serviceConfig.setVersionId(versionId);
            String name = fileContent.getFileName();
            serviceConfig.setName(name);
            serviceConfig.setDescription(desc);

            //获取配置路径
            System.out.println("project= "+project);
            System.out.println("cluster= "+cluster);
            System.out.println("service= "+service);
            System.out.println("version= "+version);
            String path = this.configCenter.getConfigPath(project, cluster, service, version, name);
            serviceConfig.setPath(path);

            byte[] content = fileContent.getContent();
            String md5 = MD5Util.getMD5(content);
            serviceConfig.setContent(content);
            serviceConfig.setMd5(md5);
            serviceConfig.setCreateTime(now);
            serviceConfig.setUpdateTime(now);
            serviceConfig.setGrayId("0");
            serviceConfigList.add(serviceConfig);
        }
        this.serviceConfigMapper.batchInsert(serviceConfigList);
        return serviceConfigList;
    }

    @Override
    public List<ServiceConfig> batchUpdate(AddServiceConfigRequestBodyByQuickStart body, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList) {
        String desc = body.getDesc();
        Date now = new Date();

        //批量更新
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            return null;
        }
        for (ServiceConfig serviceConfig : serviceConfigList) {
            for (FileContent fileContent : fileContentList) {
                if (serviceConfig.getName().equals(fileContent.getFileName())) {
                    serviceConfig.setDescription(desc);
                    byte[] content = fileContent.getContent();
                    String md5 = MD5Util.getMD5(content);
                    serviceConfig.setContent(content);
                    serviceConfig.setMd5(md5);
                    serviceConfig.setUpdateTime(now);
                    break;
                }
            }
        }
        this.serviceConfigMapper.batchUpdate(serviceConfigList);
        return serviceConfigList;
    }

    @Override
    public List<ServiceConfig> batchUpdate(List<ServiceConfigHistory> serviceConfigHistoryList, List<ServiceConfig> serviceConfigList) {
        Date now = new Date();

        //批量更新
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            return null;
        }
        for (ServiceConfig serviceConfig : serviceConfigList) {
            for (ServiceConfigHistory serviceConfigHistory : serviceConfigHistoryList) {
                if (serviceConfig.getId().equals(serviceConfigHistory.getConfigId())) {
                    serviceConfig.setDescription(serviceConfigHistory.getDescription());
                    byte[] content = serviceConfigHistory.getContent();
                    String md5 = MD5Util.getMD5(content);
                    serviceConfig.setContent(content);
                    serviceConfig.setMd5(md5);
                    serviceConfig.setUpdateTime(now);
                    break;
                }
            }
        }
        this.serviceConfigMapper.batchUpdate(serviceConfigList);
        return serviceConfigList;
    }

    @Override
    public ServiceConfig add(AddServiceConfigRequestBody body, FileContent fileContent) {
        String userId = this.getUserId();
        String versionId = body.getVersionId();
        String desc = body.getDesc();
        String fileName = fileContent.getFileName();
        byte[] content = fileContent.getContent();
        String path = fileContent.getPath();
        String md5 = MD5Util.getMD5(content);

        //新增
        Date now = new Date();
        ServiceConfig serviceConfig = new ServiceConfig();
        serviceConfig.setId(SnowflakeIdWorker.getId());
        serviceConfig.setUserId(userId);
        serviceConfig.setVersionId(versionId);
        serviceConfig.setName(fileName);
        serviceConfig.setPath(path);
        serviceConfig.setContent(content);
        serviceConfig.setMd5(md5);
        serviceConfig.setDescription(desc);
        serviceConfig.setCreateTime(now);
        serviceConfig.setGrayId("0");
        try {
            this.serviceConfigMapper.insert(serviceConfig);
            return serviceConfig;
        } catch (DuplicateKeyException ex) {
            logger.warn("service config duplicate key " + ex.getMessage());
            return this.find(fileName, versionId);
        }
    }

    @Override
    public ServiceConfig findById(String id) {
        return this.serviceConfigMapper.findById(id);
    }

    @Override
    public List<ServiceConfig> findByIds(List<String> ids) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("ids", ids);
        return this.serviceConfigMapper.findListByMap(map);
    }

    @Override
    public int deleteById(String id) {
        return this.serviceConfigMapper.deleteById(id);
    }

    @Override
    public int deleteByIds(List<String> ids) {
        return this.serviceConfigMapper.deleteByIds(ids);
    }

    @Override
    public ServiceConfig updateById(String id, EditServiceConfigRequestBody body) {
        String desc = body.getDesc();
        byte[] content = new byte[0];
        try {
            content = body.getContent().getBytes("UTF-8");
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        String md5 = MD5Util.getMD5(content);

        //更新
        Date now = new Date();
        ServiceConfig serviceConfig = new ServiceConfig();
        serviceConfig.setId(id);
        serviceConfig.setDescription(desc);
        serviceConfig.setContent(content);
        serviceConfig.setMd5(md5);
        serviceConfig.setUpdateTime(now);
        this.serviceConfigMapper.updateById(serviceConfig);
        return serviceConfig;
    }

    @Override
    public List<ServiceConfig> find(String versionId, List<String> names, String grayId) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("versionId", versionId);
        map.put("names", names);
        map.put("grayId", grayId);
        return this.serviceConfigMapper.findList(map);
    }

    @Override
    public List<ServiceConfig> findConfigsByVersionId(String versionId, String grayId) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("versionId", versionId);
        map.put("grayId", grayId);
        return this.serviceConfigMapper.findConfigsByVersionId(map);
    }

    @Override
    public ServiceConfig find(String name, String versionId) {
        ServiceConfig serviceConfig = new ServiceConfig();
        serviceConfig.setName(name);
        serviceConfig.setVersionId(versionId);
        return this.serviceConfigMapper.find(serviceConfig);
    }

    @Override
    public List<ServiceConfig> findNewList(String serviceId) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("serviceId", serviceId);
        return this.serviceConfigMapper.findNewList(map);
    }

    @Override
    public ServiceConfig findConfigJoinVersionJoinServiceJoinClusterJoinProjectById(String id) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("id", id);
        List<ServiceConfig> serviceConfigList = this.serviceConfigMapper.findConfigJoinVersionJoinServiceJoinClusterJoinProjectByMap(map);
        if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
            return serviceConfigList.get(0);
        }
        return null;
    }

    @Override
    public List<ServiceConfig> findConfigJoinVersionJoinServiceJoinClusterJoinProjectByIds(List<String> ids) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("ids", ids);
        return this.serviceConfigMapper.findConfigJoinVersionJoinServiceJoinClusterJoinProjectByMap(map);
    }

    @Override
    public List<ServiceConfig> findListByIds(List<String> ids) {
        return this.serviceConfigMapper.findListByIds(ids);
    }

    @Override
    public List<ServiceConfig> batchGrayGroupAdd(AddGrayGroupRequestBody body, String grayId, List<ServiceConfig> serviceConfigList) {
        String userId = this.getUserId();
        String versionId = body.getVersionId();
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String version = body.getVersion();
        Date now = new Date();

        List<ServiceConfig> configs = new ArrayList<>();
        //批量新增,没有拖拽新增配置文件的
        if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                serviceConfig.setId(SnowflakeIdWorker.getId());
                serviceConfig.setUserId(userId);
                serviceConfig.setVersionId(versionId);
                serviceConfig.setGrayId(grayId);
                //获取配置路径
                String path = this.configCenter.getConfigPath(project, cluster, service, version, serviceConfig.getName());
                serviceConfig.setPath(path);

                serviceConfig.setCreateTime(now);
                serviceConfig.setUpdateTime(now);
                configs.add(serviceConfig);
            }
            this.serviceConfigMapper.batchInsert(serviceConfigList);
        }
        return serviceConfigList;
    }

    @Override
    public List<ServiceConfig> batchGrayGroupAndFileAdd(AddGrayGroupRequestBody body, String grayId, List<ServiceConfig> serviceConfigList, List<FileContent> fileContentList) {
        String userId = this.getUserId();
        String versionId = body.getVersionId();
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String version = body.getVersion();
        String desc = body.getDesc();
        Date now = new Date();

        List<ServiceConfig> serviceConfigs = new ArrayList<>();

        //批量新增,没有拖拽新增配置文件的
        if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                serviceConfig.setId(SnowflakeIdWorker.getId());
                serviceConfig.setUserId(userId);
                serviceConfig.setVersionId(versionId);
                serviceConfig.setGrayId(grayId);

                //获取配置路径
                String path = this.configCenter.getConfigPath(project, cluster, service, version, serviceConfig.getName());
                serviceConfig.setPath(path);

                serviceConfig.setCreateTime(now);
                serviceConfig.setUpdateTime(now);
                serviceConfigs.add(serviceConfig);
            }
        }

        //批量新增拖拽配置文件
        ServiceConfig serviceConfig;
        if (null != fileContentList && !fileContentList.isEmpty()) {
            for (FileContent fileContent : fileContentList) {
                serviceConfig = new ServiceConfig();
                serviceConfig.setId(SnowflakeIdWorker.getId());
                serviceConfig.setUserId(userId);
                serviceConfig.setVersionId(versionId);
                serviceConfig.setGrayId(grayId);
                String name = fileContent.getFileName();
                serviceConfig.setName(name);
                serviceConfig.setDescription(desc);

                //获取配置路径
                String path = this.configCenter.getConfigPath(project, cluster, service, version, name);
                serviceConfig.setPath(path);

                byte[] content = fileContent.getContent();
                String md5 = MD5Util.getMD5(content);
                serviceConfig.setContent(content);
                serviceConfig.setMd5(md5);
                serviceConfig.setCreateTime(now);
                serviceConfig.setUpdateTime(now);
                serviceConfigs.add(serviceConfig);
            }
            this.serviceConfigMapper.batchInsert(serviceConfigs);
        }
        return serviceConfigs;
    }

    @Override
    public List<ServiceConfig> findListByGrayIds(List<String> ids) {
        return this.serviceConfigMapper.findListByGrayIds(ids);
    }

    @Override
    public List<ServiceConfig> findListByGrayId(String id) {
        return this.serviceConfigMapper.findListByGrayId(id);
    }

    @Override
    public int deleteByGrayId(String grayId) {
        return this.serviceConfigMapper.deleteByGrayId(grayId);
    }

    @Override
    public int findGrayTotalCount(HashMap<String, Object> map) {
        return this.serviceConfigMapper.findGrayTotalCount(map);
    }

    @Override
    public List<ServiceConfig> findGrayList(HashMap<String, Object> map) {
        return this.serviceConfigMapper.findGrayList(map);
    }

    @Override
    public ServiceConfig findOnlyConfig(String name, String versionId, String grayId) {
        ServiceConfig serviceConfig = new ServiceConfig();
        serviceConfig.setName(name);
        serviceConfig.setVersionId(versionId);
        serviceConfig.setGrayId(grayId);
        return this.serviceConfigMapper.findOnlyConfig(serviceConfig);
    }

    @Override
    public List<ServiceConfig> batchAddGrayFile(String versionId, AddGrayConfigRequestBody body, List<FileContent> fileContentList) throws UnsupportedEncodingException {
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String version = body.getVersion();
        String desc = body.getDesc();
        Date now = new Date();
        String grayId = body.getGrayId();

        List<ServiceConfig> serviceConfigList = new ArrayList<>();
        ServiceConfig serviceConfig;
        if (null != fileContentList && !fileContentList.isEmpty()) {
            for (FileContent fileContent : fileContentList) {
                serviceConfig = new ServiceConfig();
                serviceConfig.setId(SnowflakeIdWorker.getId());
                serviceConfig.setUserId(this.getUserId());
                serviceConfig.setVersionId(versionId);
                String name = fileContent.getFileName();
                serviceConfig.setName(name);
                serviceConfig.setDescription(desc);

                //获取配置路径
                String path = this.configCenter.getConfigPath(project, cluster, service, version, name);
                serviceConfig.setPath(path);
                byte[] content = fileContent.getContent();
                String md5 = MD5Util.getMD5(content);
                serviceConfig.setContent(content);
                serviceConfig.setMd5(md5);
                serviceConfig.setCreateTime(now);
                serviceConfig.setUpdateTime(now);
                serviceConfig.setGrayId(grayId);
                serviceConfigList.add(serviceConfig);
            }
        }
        if (!serviceConfigList.isEmpty()) {
            this.serviceConfigMapper.batchInsert(serviceConfigList);
        }
        return serviceConfigList;
    }

    @Override
    public List<ServiceConfig> batchUpdateGrayConfig(AddGrayConfigRequestBody body, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList) {
        String desc = body.getDesc();
        Date now = new Date();

        //批量更新
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            return null;
        }
        for (ServiceConfig serviceConfig : serviceConfigList) {
            for (FileContent fileContent : fileContentList) {
                if (serviceConfig.getName().equals(fileContent.getFileName())) {
                    serviceConfig.setDescription(desc);
                    byte[] content = fileContent.getContent();
                    String md5 = MD5Util.getMD5(content);
                    serviceConfig.setContent(content);
                    serviceConfig.setMd5(md5);
                    serviceConfig.setUpdateTime(now);
                    break;
                }
            }
        }
        this.serviceConfigMapper.batchUpdate(serviceConfigList);
        return serviceConfigList;
    }

    @Override
    public int copyConfigs1(List<ServiceConfig> serviceConfigs, String versionId, Map<String, String> oldGrayId2NewGrayId, String path) {
        String userId = this.getUserId();
        for (ServiceConfig serviceConfig : serviceConfigs) {
            serviceConfig.setId(SnowflakeIdWorker.getId());
            serviceConfig.setUserId(userId);
            serviceConfig.setVersionId(versionId);
            String originalPath = serviceConfig.getPath();
            String[] split = originalPath.split("/");
            String fileName = split[split.length - 1];
            StringBuffer sb = new StringBuffer(path);
            StringBuffer append = sb.append(fileName);
            serviceConfig.setPath(append.toString());
            String grayId = serviceConfig.getGrayId();
            if (oldGrayId2NewGrayId.containsKey(grayId)){
                serviceConfig.setGrayId(oldGrayId2NewGrayId.get(grayId));
            }
            Date date = new Date();
            serviceConfig.setCreateTime(date);
            serviceConfig.setUpdateTime(date);
            try {
                this.serviceConfigMapper.insert(serviceConfig);
            } catch (DuplicateKeyException ex) {
                logger.warn("service config duplicate key " + ex.getMessage());

            }
        }
        return 0;
    }
}
