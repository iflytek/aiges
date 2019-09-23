package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.log.util.StringUtil;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.ErrorCode;
import com.iflytek.ccr.polaris.companion.common.ZkDataValue;

import java.io.UnsupportedEncodingException;

public class ByteUtil {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ByteUtil.class);

    //byte 数组与 int 的相互转换
    public static int byteArrayToInt(byte[] b) {
        return b[3] & 0xFF |
                (b[2] & 0xFF) << 8 |
                (b[1] & 0xFF) << 16 |
                (b[0] & 0xFF) << 24;
    }

    public static byte[] intToByteArray(int a) {
        return new byte[]{
                (byte) ((a >> 24) & 0xFF),
                (byte) ((a >> 16) & 0xFF),
                (byte) ((a >> 8) & 0xFF),
                (byte) (a & 0xFF)
        };
    }

    public static byte[] byteMerge(byte[] byte_1, byte[] byte_2) {
        byte[] byte_3 = new byte[byte_1.length + byte_2.length];
        System.arraycopy(byte_1, 0, byte_3, 0, byte_1.length);
        System.arraycopy(byte_2, 0, byte_3, byte_1.length, byte_2.length);
        return byte_3;
    }

    public static void main(String[] args) {
        System.out.println(intToByteArray(99999).length);
        System.out.println(byteArrayToInt(intToByteArray(9999)));
    }


    /**
     * 生成zk中保存的数据格式
     *
     * @param contentData
     * @param version
     * @return
     */
    public static byte[] getZkBytes(byte[] contentData, String version) {
        if (StringUtil.isNullOrEmpty(version)) {
            version = "0";
        }
        byte[] dataByte = null;
        try {
            byte[] versionBytes = version.getBytes(Constants.DEFAULT_CHARSET);
            byte[] pre = ByteUtil.intToByteArray(versionBytes.length);
            dataByte = ByteUtil.byteMerge(pre, versionBytes);
            dataByte = ByteUtil.byteMerge(dataByte, contentData);
        } catch (UnsupportedEncodingException e) {
            logger.error("", e);
        }
        return dataByte;
    }

    /**
     * 解析zk中保存的数据
     * @param data
     * @return
     */
    public static ZkDataValue parseZkData(byte[] data) {
        ZkDataValue zkDataValue = new ZkDataValue();
        zkDataValue.setRet(ErrorCode.SUCCESS);
        byte[] preByte = new byte[4];
        if (data.length <= 4) {
            logger.error("data length invalid:" + data);
            zkDataValue.setDesc("zk data is invalid");
            zkDataValue.setRet(ErrorCode.INTERNAL_EXCEPTION);
            return zkDataValue;
        }
        System.arraycopy(data, 0, preByte, 0, 4);
        int versionLength = ByteUtil.byteArrayToInt(preByte);
        if (data.length < versionLength) {
            String desc = "version Length is invalid:" + versionLength;
            logger.error(desc);
            zkDataValue.setRet(ErrorCode.INTERNAL_EXCEPTION);
            zkDataValue.setDesc(desc);
        }
        byte[] verByte = new byte[versionLength];
        System.arraycopy(data, 4, verByte, 0, versionLength);
        String pushId = new String(verByte);

        byte[] realData = new byte[data.length - 4 - versionLength];
        System.arraycopy(data, 4 + versionLength, realData, 0, data.length - 4 - versionLength);

        zkDataValue.setPushId(pushId);
        zkDataValue.setRealData(realData);
        return zkDataValue;
    }

}
