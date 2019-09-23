package com.iflytek.ccr.finder.service.impl;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.service.GrayConfigService;
import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.utils.JacksonUtils;
import com.iflytek.ccr.finder.utils.StringUtils;
import com.iflytek.ccr.finder.utils.ZkInstanceUtil;
import com.iflytek.ccr.finder.value.ErrorCode;
import com.iflytek.ccr.finder.value.GrayConfigValue;
import com.iflytek.ccr.finder.value.ZkDataValue;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import java.util.Map;

/**
 * 灰度服务实现类
 */
public class GrayConfigServiceImpl implements GrayConfigService {

    private static final Logger logger = LoggerFactory.getLogger(GrayConfigServiceImpl.class);

    @Override
    public List<GrayConfigValue> parseGrayData(String grayConfigPath) {
        List<GrayConfigValue> grayValueList = new ArrayList<>();
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            //获取历史数据
            byte[] data = null;
            if (zkHelper.checkExists(grayConfigPath)) {
                data = zkHelper.getByteData(grayConfigPath);
            }
            if (data == null || data.length == 0) {
                return grayValueList;
            }

            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (zkDataValue.getRet() == ErrorCode.SUCCESS) {
                String grayData = new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET);
                List grayList = JacksonUtils.toObject(grayData, List.class);
                for (Object v : grayList) {
                    Map<String, Object> vMap = (Map<String, Object>) v;
                    GrayConfigValue value = new GrayConfigValue();
                    String tempGroupName = vMap.get(Constants.GROUP_ID).toString();
                    value.setGroupId(tempGroupName);
                    value.setServerList( Arrays.asList(((ArrayList) vMap.get(Constants.SERVER_LIST)).get(0).toString().split(",")));
                    grayValueList.add(value);
                }
            }
        } catch (Exception e) {
            logger.error(String.format("parseRouteData error:%s", e.getMessage()), e);
        }

        return grayValueList;
    }

    @Override
    public GrayConfigValue getGrayServer(FinderManager finderManager, List<GrayConfigValue> grayValueList) {
        GrayConfigValue grayConfigValue = null;
        //获取组件实例标识
        String ipAddr = finderManager.getBootConfig().getMeteData().getAddress();
        if (StringUtils.isNotEmpty(grayValueList) && StringUtils.isNOtNullOrEmpty(ipAddr)) {
            for (GrayConfigValue value : grayValueList) {
                if (StringUtils.isNotEmpty(value.getServerList())) {
                    if (value.getServerList().contains(ipAddr)) {
                        grayConfigValue = value;
                        break;
                    }
                }
            }
        }
        return grayConfigValue;
    }
}
