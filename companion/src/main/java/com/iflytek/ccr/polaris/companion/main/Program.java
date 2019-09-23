package com.iflytek.ccr.polaris.companion.main;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.log.Log4J2LoggerFactory;
import com.iflytek.ccr.nakedserver.route.UrlTemplate;
import com.iflytek.ccr.nakedserver.server.NakedHttpServer;
import com.iflytek.ccr.polaris.companion.common.ConfigValue;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.task.BackupDataTask;
import com.iflytek.ccr.polaris.companion.task.ConfigFeedbackQueueConsumerMonitorTask;
import com.iflytek.ccr.polaris.companion.task.MonitorZkTask;
import com.iflytek.ccr.polaris.companion.task.ServiceFeedbackQueueConsumerMonitorTask;
import com.iflytek.ccr.polaris.companion.utils.ConfigManager;
import com.iflytek.ccr.polaris.companion.utils.ZkHelper;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import org.apache.commons.cli.*;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

/**
 * Created by eric on 2017/11/17.
 */
public class Program {
    public static ConfigValue CONFIG_VALUE = null;
    private static EasyLogger logger = EasyLoggerFactory.getInstance(Program.class);

    static {
        EasyLoggerFactory.setDefaultFactory(Log4J2LoggerFactory.INSTANCE);
    }

    public static void main(String[] args) {
        if (null != args) {
            for (String str : args) {
                System.out.println(str);
                logger.info("args:" + str);
            }
        }
        CONFIG_VALUE = initConfig(args);
        logger.info(CONFIG_VALUE.toString());

        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        if (!zkHelper.checkExists(Constants.QUEUE_PATH_ZK_NODE)) {
            zkHelper.addPersistent(Constants.QUEUE_PATH_ZK_NODE, "");
        }

        String zkNodePath = zkHelper.addEphemeralSequential(Constants.QUEUE_PATH_ZK_NODE + "/zk", ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR));
        ConfigManager.getInstance().put("zkNodePath", zkNodePath);
        try {
//            ServiceCacheUtil.getInstance();
//            FeedbackCacheUtil.getInstance();
            ScheduledExecutorService scheduledExecutorService = Executors.newSingleThreadScheduledExecutor();
            scheduledExecutorService.scheduleAtFixedRate(new BackupDataTask(), 0, 1, TimeUnit.HOURS);
            ExecutorService executorService = Executors.newSingleThreadExecutor();
            executorService.submit(new MonitorZkTask());

            ExecutorService distributedQueueManagerService = Executors.newSingleThreadScheduledExecutor();
            distributedQueueManagerService.submit(new ConfigFeedbackQueueConsumerMonitorTask());
            distributedQueueManagerService.submit(new ServiceFeedbackQueueConsumerMonitorTask());

            NakedHttpServer server = new NakedHttpServer();
            server.setListenWorkers(1).setNetWorkers(16).setTaskWorkers(64)
                    .setTcpTimeout(3000).listen(ConfigManager.getStringConfigByKey(ConfigManager.KEY_HOST), ConfigManager.getIntConfigByKey(ConfigManager.KEY_PORT)).setUrlTemplate(UrlTemplate.DEFAULT)
                    .start();
        } catch (Exception e) {
            logger.error(e);
            System.out.println("start error");
            e.printStackTrace();
            System.exit(-1);
        }

    }

    private static ConfigValue initConfig(String[] args) {

        CommandLineParser parser = new DefaultParser();
        ConfigValue configValue = new ConfigValue();

        Options options = buildCommandlineOptions();

        try {
            CommandLine cmd = parser.parse(options, args);
            if (cmd.hasOption("h")) {
                String host = cmd.getOptionValue("h");
                configValue.setIpAddr(host);
                ConfigManager.getInstance().put("host", host);
            }
            if (cmd.hasOption("p")) {
                String port = cmd.getOptionValue("p");
                configValue.setPort(Integer.parseInt(port));
                ConfigManager.getInstance().put("port", Integer.parseInt(port));
            }

            if (cmd.hasOption("z")) {
                String zkStr = cmd.getOptionValue("z");
                configValue.setZkStr(zkStr);
                ConfigManager.getInstance().put("zkStr", zkStr);
            }

            if (cmd.hasOption("w")) {
                String websiteUrl = cmd.getOptionValue("w");
                configValue.setWebsiteUrl(websiteUrl);
                ConfigManager.getInstance().put("websiteUrl", websiteUrl);
            }
        } catch (Exception e) {
            logger.error(e);
        }
        if (null == ConfigManager.getStringConfigByKey(ConfigManager.KEY_HOST) || ConfigManager.getStringConfigByKey(ConfigManager.KEY_HOST).isEmpty()) {
            ConfigManager.getInstance().put("host", "0.0.0.0");
            configValue.setIpAddr("0.0.0.0");
        }
        if (0 == ConfigManager.getIntConfigByKey(ConfigManager.KEY_PORT)) {
            configValue.setPort(9090);
            ConfigManager.getInstance().put("host", "9090");
        }
        if (null == ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR) || ConfigManager.getStringConfigByKey(ConfigManager.KEY_ZKSTR).isEmpty()) {
            logger.error("plase config zookeeper!");
            System.out.println("plase config zookeeper!");
            System.exit(0);
        }
        return configValue;
    }

    public static Options buildCommandlineOptions() {
        // 创建 Options 对象
        Options options = new Options();
        Option opt = new Option("d", "configDir", true, "config directory");
        opt.setRequired(false);
        options.addOption(opt);

        opt = new Option("h", "host", true, "host config");
        opt.setRequired(false);
        options.addOption(opt);

        opt = new Option("p", "port", true, "port config");
        opt.setRequired(false);
        options.addOption(opt);

        opt = new Option("w", "websiteUrl", true, "website Url");
        opt.setRequired(false);
        options.addOption(opt);

        opt = new Option("z", "zkStr", true, "zookeeper str");
        opt.setRequired(true);
        options.addOption(opt);
        return options;
    }
}
