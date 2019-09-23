package com.iflytek.ccr.finder.utils;

/**
 * Created by eric on 2017/11/21.
 */
public interface JsonUtils {
      String toJson(Object o);

      <T> T toObject(String json);

      <T> T toObject(String json, Class<T> valueType);
}
