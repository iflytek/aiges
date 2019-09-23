package com.iflytek.ccr.polaris.companion.service;

import java.util.Map;

/**
 * Created by eric on 2017/11/17.
 */
public interface ConfigService {
    /**
     * 生成配置路径
     * @param prefix
     * @param params
     * @return
     * @throws Exception
     */
    String generateConfigPath(String prefix, Map<String, String> params) throws Exception;

    /**
     *  生成服务路径
     * @param prefix
     * @param params
     * @return
     * @throws Exception
     */
    String generateServicePath(String prefix, Map<String, String> params) throws Exception;
}

