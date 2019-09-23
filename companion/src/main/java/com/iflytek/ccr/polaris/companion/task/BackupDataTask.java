package com.iflytek.ccr.polaris.companion.task;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.utils.CompanionFileUtils;
import com.iflytek.ccr.polaris.companion.utils.FileNameComparator;
import com.iflytek.ccr.polaris.companion.utils.ZkHelper;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import org.apache.http.client.utils.DateUtils;

import java.io.File;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Date;
import java.util.List;

public class BackupDataTask implements Runnable {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(BackupDataTask.class);

    String backupPath = Constants.BACKUP_PRE;

    public static void deleteFolder(String path) {
        File file = new File(path);
        // 判断是否是一个目录, 不是的话跳过, 直接删除; 如果是一个目录, 先将其内容清空.
        if (file.isDirectory()) {
            // 获取子文件/目录
            File[] subFiles = file.listFiles();
            // 遍历该目录
            for (File subFile : subFiles) {
                // 递归调用删除该文件: 如果这是一个空目录或文件, 一次递归就可删除. 如果这是一个非空目录, 多次
                // 递归清空其内容后再删除
                deleteFolder(subFile.getPath());
            }
        }
        // 删除空目录或文件
        file.delete();
    }

    public static void main(String[] args) {
        BackupDataTask backupDataTask = new BackupDataTask();
        String path = "E:\\opt\\server\\backup\\20180111092103";
        backupDataTask.deleteFolder(path);


        List<String> list = new ArrayList<>();
        list.add("2018010501");
        list.add("2018010602");
        list.add("2018011203");
        list.add("2018011104");
        list.add("2018010909");
        list.add("2018020502");
        list.add("2018110608");
        list.add("2018040107");
        String[] nameArray = new String[list.size()];
        list.toArray(nameArray);
        Arrays.sort(nameArray, new FileNameComparator());
//        Comparator<String> comparator = new Comparator<Date>();

        for (String str : nameArray) {
            System.out.println(str);
        }
    }

    public void backup() {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        if (!ZkInstanceUtil.checkZkHelper()) {
            logger.error("Unable to connect to zookeeper");
            return;
        }
        CompanionFileUtils.createFolder(backupPath);
        saveByPath(zkHelper, Constants.CONFIG_PATH_PREFIX);
        saveByPath(zkHelper, Constants.SERVICE_PATH_PREFIX);
    }

    public void saveByPath(ZkHelper zkHelper, String zkPath) {
        saveData(zkHelper, zkPath);
        if (hasChildren(zkHelper, zkPath)) {
            CompanionFileUtils.createFolder(zkPath);
            List<String> list = zkHelper.getChildren(zkPath);
            for (String temp : list) {
                if (Constants.PROVIDER.equals(temp) || Constants.CONSUMER.equals(temp)) {
                    continue;
                }
                String tempZkPath = zkPath + "/" + temp;
                saveByPath(zkHelper, tempZkPath);
            }
        }
    }

    public boolean hasChildren(ZkHelper zkHelper, String zkPath) {
        return zkHelper.checkExists(zkPath);
    }

    /**
     * 保存数据
     *
     * @param zkHelper
     * @param zkPath
     */
    public void saveData(ZkHelper zkHelper, String zkPath) {
        if(zkHelper.checkExists(zkPath)){
            byte[] data = zkHelper.getByteData(zkPath);
            String filePath = covertZkpath2FilePath(zkPath) + ".data";
            CompanionFileUtils.writeByteArrayToFile(filePath, data);
        }else {
            logger.error(String.format("%s not exists",zkPath));
        }

    }

    /**
     * 将zk路径转换为文件路径
     *
     * @param zkPath
     * @return
     */
    public String covertZkpath2FilePath(String zkPath) {
        return backupPath + File.separator + zkPath;
    }

    @Override
    public void run() {
        deleteOld();
        String date = DateUtils.formatDate(new Date(), Constants.DATE_PATTERN);
        backupPath = Constants.BACKUP_PRE + date + File.separator;
        File file = new File(backupPath);
        if (file.exists()) {
            return;
        }
        backup();
    }

    public void deleteOld() {
        List<String> list = new ArrayList<>();
        File file = new File(Constants.BACKUP_PRE);
        if (file.isDirectory()) {
            File[] fileArray = file.listFiles();
            for (File temp : fileArray) {
                if (temp.isDirectory()) {
                    list.add(temp.getName());
                    logger.info("fileName:" + temp);
                }
            }
        }
        logger.info("listFiles:" + list.size());
        if (list.size() > 10) {
            String[] nameArray = new String[list.size()];
            list.toArray(nameArray);
            Arrays.sort(nameArray, new FileNameComparator());

            for (int i = 10; i < list.size(); i++) {
                String path = Constants.BACKUP_PRE + nameArray[i];
                deleteFolder(path);
                logger.info("deleteFolder:" + path);
            }

        }
    }
}
