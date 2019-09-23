package com.iflytek.ccr.finder.utils;


import com.iflytek.ccr.finder.constants.Constants;
import org.apache.commons.io.FileUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.*;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 文件处理工具类
 */
public class FinderFileUtils {

    private final static Logger logger = LoggerFactory.getLogger(FinderFileUtils.class);

    public static Map<String, Object> parseTomlFile(byte[] file) {
        Map<String, Object> map = new HashMap<>();
        try {
            String content = new String(file, Constants.DEFAULT_CHARSET);
            String[] result = content.split("\r\n");
            List<String> temp = new ArrayList<>();
            for(String st : result){
                for(String t :st.split("\n")){
                    temp.add(t);
                }
            }

            String currentGroup = null;
            for (String st : temp) {
                st = st.trim();
                if (st.startsWith("#")) {
                    continue;
                } else if (StringUtils.isNullOrEmpty(st)) {
                    continue;
                } else if (st.startsWith("[") && st.endsWith("]")) {
                    currentGroup = st.substring(1, st.length() - 1).trim();
                } else if (!st.contains("=")) {
                    continue;
                } else {
                    if (st.contains("#")) {
                        String[] ss = st.trim().split("#");
                        String[] sss = ss[0].trim().split("=");
                        if (StringUtils.isNullOrEmpty(currentGroup)) {
                            map.put(sss[0], sss[1]);
                        } else {
                            map.put(currentGroup + "." + sss[0].trim(), sss[1].trim());
                        }
                    } else {
                        String[] ss = st.trim().split("=");
                        if (StringUtils.isNullOrEmpty(currentGroup)) {
                            map.put(ss[0].trim(), ss[1].trim());
                        } else {
                            map.put(currentGroup + "." + ss[0].trim(), ss[1].trim());
                        }
                    }
                }
            }
        } catch (Exception e) {
            logger.error("parseTomlFile error,", e);
        }
        return map;
    }

    public static void createFolder(String path) {
        logger.info(path);
        File file = new File(path);
        if (!file.isDirectory()) {
            file.mkdirs();
            logger.info("absolutePath:" + file.getAbsolutePath());
        }
    }

    public static void writeObjectToFile(String fileName, Object object) {
        try {
            File file = new File(fileName);
            File parent = new File(file.getParent());
            if (!parent.isDirectory()) {
                parent.mkdirs();
            }
            ObjectOutputStream os = new ObjectOutputStream(
                    new FileOutputStream(fileName));
            os.writeObject(object);// 将User对象写进文件
            os.close();
        } catch (FileNotFoundException e) {
            logger.error("", e);
        } catch (IOException e) {
            logger.error("", e);
        }
    }

    public static Object readObjectFromFile(String fileName) {
        Object result = null;
        try {
            ObjectInputStream is = new ObjectInputStream(new FileInputStream(
                    fileName));
            result = is.readObject();// 从流中读取Object的数据
            is.close();
        } catch (FileNotFoundException e) {
            logger.error("", e);
        } catch (IOException e) {
            logger.error("", e);
        } catch (ClassNotFoundException e) {
            logger.error("", e);
        } catch (Exception e) {
            logger.error("", e);
        }
        return result;
    }

    public static void writeByteArrayToFile(String name, byte[] fileContent) {
        logger.info(String.format("fileName:%s,filecontent len:%s", name, fileContent.length));
        try {
            FileUtils.writeByteArrayToFile(new File(name), fileContent);
        } catch (IOException e) {
            logger.error("", e);
        }
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
