package com.iflytek.ccr.finder.service.impl;

import com.iflytek.ccr.finder.service.CommonService;
import com.iflytek.ccr.finder.utils.ZkInstanceUtil;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.ErrorCode;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class CommonServiceImpl implements CommonService {

    private static final Logger logger = LoggerFactory.getLogger(CommonServiceImpl.class);

    @Override
    public CommonResult unRegisterConsumer(String path) {
        CommonResult commonResult = new CommonResult();
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            zkHelper.remove(path);
        } catch (Exception e) {
            logger.error(String.format("unSubscribeConfig error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }
        return commonResult;
    }

    @Override
    public CommonResult registerConsumer(String path) {
        CommonResult commonResult = new CommonResult();
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            zkHelper.remove(path);
            zkHelper.addEphemeral(path, "");
        } catch (Exception e) {
            logger.error(String.format("unSubscribeConfig error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }
        return commonResult;
    }
}
