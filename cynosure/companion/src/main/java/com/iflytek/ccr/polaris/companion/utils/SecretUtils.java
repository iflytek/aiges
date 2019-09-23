package com.iflytek.ccr.polaris.companion.utils;

import java.io.UnsupportedEncodingException;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * Created by eric on 2017/11/21.
 */
public class SecretUtils {
    private static final String MD5 = "MD5";
    private static final String DEFAULT_ENCODING = "utf-8";

    public static String getMD5(String str) throws UnsupportedEncodingException, NoSuchAlgorithmException {
        MessageDigest md5 = MessageDigest.getInstance(MD5);

        return bytesToHex(md5.digest(str.getBytes(DEFAULT_ENCODING)));
    }

    public static String bytesToHex(byte[] bytes) {
        StringBuffer md5str = new StringBuffer();
        // 把数组每一字节换成16进制连成md5字符串
        int digital;
        for (int i = 0; i < bytes.length; i++) {
            digital = bytes[i];

            if (digital < 0) {
                digital += 256;
            }
            if (digital < 16) {
                md5str.append("0");
            }
            md5str.append(Integer.toHexString(digital));
        }

        return md5str.toString();
    }
}
