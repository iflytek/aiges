package com.iflytek.ccr.polaris.cynosure.base;

import com.iflytek.ccr.polaris.cynosure.domain.*;
import com.iflytek.ccr.polaris.cynosure.enums.DBEnumInt;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceApiVersionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;

/**
 * 业务逻辑基类
 *
 * @author sctang2
 * @create 2017-11-15 10:27
 **/
@Service
public class BaseService {
    @Autowired
    protected HttpServletRequest httpServletRequest;

    /**
     * 获取用户id
     *
     * @return
     */
    protected String getUserId() {
        User user = (User) httpServletRequest.getSession().getAttribute("user");
        if (null == user) {
            return "0";
        }
        return user.getId();
    }

    /**
     * 判断是否为管理员
     *
     * @return
     */
    protected boolean isAdmin() {
        User user = (User) httpServletRequest.getSession().getAttribute("user");
        int roleType = user.getRoleType();
        if (DBEnumInt.ROLE_TYPE_ADMIN.getIndex() == roleType) {
            return true;
        }
        return false;
    }

    /**
     * 获取请求的uri
     *
     * @return
     */
    protected String getRequestURI() {
        return this.httpServletRequest.getRequestURI();
    }

    /**
     * 创建项目结果
     *
     * @param project
     * @return
     */
    protected ProjectDetailResponseBody createProjectResult(Project project) {
        ProjectDetailResponseBody result = new ProjectDetailResponseBody();
        result.setCreateTime(project.getCreateTime());
        result.setDesc(project.getDescription());
        result.setId(project.getId());
        result.setName(project.getName());
        result.setUpdateTime(project.getUpdateTime());
        return result;
    }

    /**
     * 创建集群结果
     *
     * @param cluster
     * @return
     */
    protected ClusterDetailResponseBody createClusterResult(Cluster cluster) {
        ClusterDetailResponseBody result = new ClusterDetailResponseBody();
        result.setCreateTime(cluster.getCreateTime());
        result.setId(cluster.getId());
        result.setName(cluster.getName());
        result.setDesc(cluster.getDescription());
        result.setUpdateTime(cluster.getUpdateTime());
        result.setProjectId(cluster.getProjectId());
        return result;
    }

    /**
     * 创建服务结果
     *
     * @param service
     * @return
     */
    protected ServiceDetailResponseBody createServiceResult(com.iflytek.ccr.polaris.cynosure.domain.Service service) {
        ServiceDetailResponseBody result = new ServiceDetailResponseBody();
        result.setId(service.getId());
        result.setName(service.getName());
        result.setDesc(service.getDescription());
        result.setClusterId(service.getGroupId());
        result.setCreateTime(service.getCreateTime());
        result.setUpdateTime(service.getUpdateTime());
        return result;
    }

    /**
     * 创建服务版本结果
     *
     * @param serviceVersion
     * @return
     */
    protected ServiceVersionDetailResponseBody createServiceVersionResult(ServiceVersion serviceVersion) {
        ServiceVersionDetailResponseBody result = new ServiceVersionDetailResponseBody();
        result.setId(serviceVersion.getId());
        result.setVersion(serviceVersion.getVersion());
        result.setServiceId(serviceVersion.getServiceId());
        result.setDesc(serviceVersion.getDescription());
        result.setCreateTime(serviceVersion.getCreateTime());
        result.setUpdateTime(serviceVersion.getUpdateTime());
        return result;
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

    /**
     * 创建服务配置结果
     *
     * @param serviceConfig
     * @return
     */
    protected ServiceConfigDetailResponseBody createServiceConfigResult(ServiceConfig serviceConfig) {
        ServiceConfigDetailResponseBody result = new ServiceConfigDetailResponseBody();
        result.setId(serviceConfig.getId());
        result.setVersionId(serviceConfig.getVersionId());
        result.setGrayId(serviceConfig.getGrayId());
        result.setName(serviceConfig.getName());
        String groupId = serviceConfig.getGrayId();
        String path = serviceConfig.getPath();
        //前端展示数据，若是灰度配置，则在path中拼接...../gray/grayId/....
        if (!("0".equals(groupId)) && groupId!=null){
            String[] split = path.split("/");
            String tmpStr = StringUtils.substringBeforeLast(path, "/");
            StringBuffer sb = new StringBuffer(tmpStr);
            String newPath = sb.append("/").append("gray/").append(groupId).append("/").append(split[split.length - 1]).toString();
            result.setPath(newPath);
        }else {
            result.setPath(serviceConfig.getPath());
        }
        try {
            result.setContent(new String(serviceConfig.getContent(), "UTF-8"));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        result.setDesc(serviceConfig.getDescription());
        result.setCreateTime(serviceConfig.getCreateTime());
        result.setUpdateTime(serviceConfig.getUpdateTime());
        return result;
    }

    /**
     * 创建灰度组结果
     *
     * @param grayGroup
     * @return
     */
    protected GrayGroupDetailResponseBody createGrayGroupResult(GrayGroup grayGroup) {
        GrayGroupDetailResponseBody result = new GrayGroupDetailResponseBody();
        result.setId(grayGroup.getId());
        result.setVersionId(grayGroup.getVersionId());
        result.setUserId(grayGroup.getUserId());
        result.setName(grayGroup.getName());
        result.setContent(grayGroup.getContent());
        result.setDesc(grayGroup.getDescription());
        result.setCreateTime(grayGroup.getCreateTime());
        result.setUpdateTime(grayGroup.getUpdateTime());
        return result;
    }
}