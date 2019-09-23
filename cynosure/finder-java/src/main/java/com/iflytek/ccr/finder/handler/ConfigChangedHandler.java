package com.iflytek.ccr.finder.handler;

import com.iflytek.ccr.finder.value.Config;

/**
 * 配置管理回调函数
 */
public interface ConfigChangedHandler {

    boolean onConfigFileChanged(Config config);
}
