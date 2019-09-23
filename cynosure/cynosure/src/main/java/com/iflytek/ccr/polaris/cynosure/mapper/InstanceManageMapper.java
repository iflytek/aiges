package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.GrayGroup;
import org.apache.ibatis.annotations.Param;

import java.util.List;

/**
 * Created by DELL-5490 on 2018/7/4.
 */
public interface InstanceManageMapper {
    /**
     * 更新灰度组
     *
     * @param grayGroup
     * @return
     */
    int updateById(GrayGroup grayGroup);

    /**
     * 查明所有推送实例
     *
     * @return
     */
    List<String> findTotal(@Param(value = "versionId") String versionId, @Param(value = "grayGroupId") String grayGroupId);
}