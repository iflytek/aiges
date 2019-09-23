package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.nakedserver.util.StringUtils;

import java.util.Map;

/**
 * Created by eric on 2017/11/21.
 */
public class ArgumentValidator {
    public static void validateKVParams(Map<String, String> kvParams, String... keys) throws IllegalArgumentException {
        if (keys != null) {
            for (int i = 0; i < keys.length; i++) {
                if (!kvParams.containsKey(keys[i]) || kvParams.get(keys[i]) == null || kvParams.get(keys[i]).length() == 0) {
                    throw new IllegalArgumentException(StringUtils.concat("param:", keys[i], " is invalid."));
                }
            }
        }
    }
}
