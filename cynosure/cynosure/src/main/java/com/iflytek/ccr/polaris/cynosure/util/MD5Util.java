package com.iflytek.ccr.polaris.cynosure.util;

import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;

import java.io.UnsupportedEncodingException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * @author sctang
 * @version 1.0
 * @create 2017年4月14日 上午10:08:04
 * @description md5
 */
public class MD5Util {
    private static char hexDigits[] = {'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'};

    private final static String[] hexDigitsStr = {"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"};

    /**
     * @param buffer
     * @return
     * @description 对字符数组进行MD5加密
     * @author sctang
     * @create 2015年3月15日下午1:54:02
     * @version 1.0
     */
    public static String getMD5(byte[] buffer) {
        try {
            MessageDigest messagedigest = MessageDigest.getInstance("MD5");
            messagedigest.update(buffer);
            return bufferToHex(messagedigest.digest());
        } catch (NoSuchAlgorithmException e) {
            GlobalExceptionUtil.log(e);
        }
        return "";
    }

    /**
     * @param salt
     * @param val
     * @return
     * @description 对字符串加盐MD5加密
     * @author sctang
     * @create 2015年3月19日下午7:11:27
     * @version 1.0
     */
    public static String getSaltMD5(String salt, String val) {
        try {
            String mergeVal = val + "{" + salt + "}";
            MessageDigest messagedigest = MessageDigest.getInstance("MD5");
            return byteArrayToHexString(messagedigest.digest(mergeVal.getBytes("utf-8")));
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        } catch (NoSuchAlgorithmException e) {
            GlobalExceptionUtil.log(e);
        }
        return "";
    }

    /**
     * @param salt
     * @param val
     * @return
     * @description 对字符串加盐MD5加密，返回16位
     * @author sctang
     * @create 2017年7月17日 上午10:52:16
     * @version 1.0
     */
    public static String getMD5By16Bit(String salt, String val) {
        return getSaltMD5(salt, val).substring(8, 24);
    }

    private static String bufferToHex(byte bytes[]) {
        return bufferToHex(bytes, 0, bytes.length);
    }

    private static String bufferToHex(byte bytes[], int m, int n) {
        StringBuffer stringbuffer = new StringBuffer(2 * n);
        int k = m + n;
        for (int l = m; l < k; l++) {
            char c0 = hexDigits[(bytes[l] & 0xf0) >> 4];
            char c1 = hexDigits[bytes[l] & 0xf];
            stringbuffer.append(c0);
            stringbuffer.append(c1);
        }
        return stringbuffer.toString();
    }

    private static String byteArrayToHexString(byte[] b) {
        StringBuffer resultSb = new StringBuffer();
        for (int i = 0; i < b.length; i++) {
            resultSb.append(byteToHexString(b[i]));
        }
        return resultSb.toString();
    }

    private static String byteToHexString(byte b) {
        int n = b;
        if (n < 0)
            n = 256 + n;
        int d1 = n / 16;
        int d2 = n % 16;
        return hexDigitsStr[d1] + hexDigitsStr[d2];
    }

    public static void main(String[] args) throws Exception {
        System.out.println(getMD5(("project3" + "cluster1").getBytes()));
    }
}
