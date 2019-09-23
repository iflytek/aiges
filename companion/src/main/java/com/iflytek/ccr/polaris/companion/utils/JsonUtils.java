package com.iflytek.ccr.polaris.companion.utils;

/**
 * Created by eric on 2017/11/21.
 */
public interface JsonUtils {
      String toJson(Object o);

      <T> T toObject(String json);

      <T> T toObject(String json,Class<T> valueType);
}
