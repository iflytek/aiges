package com.iflytek.ccr.polaris.companion.task;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.utils.CompanionFileUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkHelper;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;

import java.io.File;

public class RestoreZkData {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(RestoreZkData.class);
    String prePath = null;

    public static void main(String[] args) {
        String backupPath = "E:/aa/aa3";
        new RestoreZkData().restore(backupPath);
    }

    public void restore(String filePath) {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        if (!ZkInstanceUtil.checkZkHelper()) {
            logger.error("Unable to connect to zookeeper");
            return;
        }
        prePath = filePath;

        File file = new File(prePath);
        if (!file.exists()) {
            logger.error(String.format("Path %s is not exists!", filePath));
            return;
        }
        restoreData(zkHelper, filePath);
    }

    /**
     * 恢复数据
     *
     * @param zkHelper
     * @param path
     */
    public void restoreData(ZkHelper zkHelper, String path) {
        File file = new File(path);
        if (file.isFile()) {
            String zkPath = convert2ZkPath(file);
            byte[] data = CompanionFileUtils.readFileToByteArray(file);
            zkPath = zkPath.replaceAll("\\\\", "/");
            zkHelper.addOrUpdatePersistentNode(zkPath, data);
            System.out.println(zkPath);
        } else if (file.isDirectory()) {
            for (File temp : file.listFiles()) {
                restoreData(zkHelper, temp.getPath());
            }
        }
    }

    /**
     * 根据文件路径获取zk路径
     *
     * @param file
     * @return
     */
    public String convert2ZkPath(File file) {
        String path = file.getPath().substring(prePath.length(), file.getPath().length() - 5);
        return path.replace("polaris", "polaris");
    }

}
