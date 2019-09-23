package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.nakedserver.util.StringUtils;
import com.iflytek.ccr.polaris.companion.utils.SecretUtils;

import java.io.UnsupportedEncodingException;
import java.security.NoSuchAlgorithmException;
import java.util.Map;

/**
 * Created by eric on 2017/11/21.
 */
public class ConfigServiceImpl implements ConfigService {
    public String generateConfigPath(String prefix, Map<String,String> params) throws UnsupportedEncodingException, NoSuchAlgorithmException {
        if (params != null && params.size() > 0) {
            return StringUtils.concat(prefix, SecretUtils.getMD5(StringUtils.concat(params.get("project"),params.get("group"))),"/",SecretUtils.getMD5(StringUtils.concat(params.get("service"),params.get("version"))));
        } else {
            throw new IllegalArgumentException("items");
        }
    }

    @Override
    public String generateServicePath(String prefix, Map<String, String> params) throws Exception {
        if (params != null && params.size() > 0) {
            return StringUtils.concat(prefix, SecretUtils.getMD5(StringUtils.concat(params.get("project"),params.get("group"))));
        } else {
            throw new IllegalArgumentException("items");
        }
    }
}
