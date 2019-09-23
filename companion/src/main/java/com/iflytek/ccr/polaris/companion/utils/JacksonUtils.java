package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import org.codehaus.jackson.map.ObjectMapper;

import java.io.IOException;

/**
 * Created by eric on 2017/11/21.
 */
public class JacksonUtils {
    private static EasyLogger logger = EasyLoggerFactory.getInstance(JacksonUtils.class);


    public static String toJson(Object o) {
        ObjectMapper mapper = new ObjectMapper();
        try {
            return mapper.writeValueAsString(o);
        } catch (IOException e) {
            logger.error(e);
        }

        return null;
    }

    public static <T> T toObject(String json, Class<T> valueType) {
        ObjectMapper mapper = new ObjectMapper();
        try {
            return mapper.readValue(json, valueType);
        } catch (IOException e) {
            logger.error(e);
        }
        return null;
    }

    public <T> T toObject(String json) {
        return null;
    }
}
