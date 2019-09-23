package com.iflytek.ccr.polaris.companion.utils;


import org.apache.commons.io.FileUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;

/**
 * 文件操作工具类
 */
public class CompanionFileUtils {

    private final static Logger logger = LoggerFactory.getLogger(CompanionFileUtils.class);

    /**
     * 创建文件夹
     * @param path
     */
    public static void createFolder(String path) {
        logger.info("createFolder:"+path);
        File file = new File(path);
        if (!file.isDirectory()) {
            file.mkdirs();
        }
    }

    public static void writeByteArrayToFile(String name, byte[] fileContent) {
        try {
            FileUtils.writeByteArrayToFile(new File(name), fileContent);
        } catch (IOException e) {
            logger.error("", e);
        }

    }

    public static byte[] readFileToByteArray(File name) {
        logger.info(String.format("fileName:%s", name));
        try {
            return FileUtils.readFileToByteArray(name);
        } catch (IOException e) {
            logger.error("", e);
        }
        return null;
    }

    public static byte[] readFileToByteArray(String name) {
        logger.info(String.format("fileName:%s", name));
        try {
            return FileUtils.readFileToByteArray(new File(name));
        } catch (IOException e) {
            logger.error("", e);
        }
        return null;
    }
}
