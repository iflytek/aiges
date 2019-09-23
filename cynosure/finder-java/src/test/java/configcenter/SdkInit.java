package configcenter;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.value.BootConfig;
import com.iflytek.ccr.finder.value.ServiceMeteData;
import com.iflytek.ccr.zkutil.ZkHelper;

/**
 * 1、不在灰度组中
 * 1)只订阅一个配置文件，配置文件变更能收到通知
 * 2)订阅二个配置文件，配置文件变更能收到通知
 */
public class SdkInit {

    static String a_cfg = "hello world";
    static String a_cfg_gray_1 = "gray hello world";
    static String a_cfg_gray_2 = "gray hello world dsfds";
    static String b_cfg = "hello thank you";
    static String b_cfg_gray = "hello thank you heheh";
    static String a_cfg_changed = "I am fine";
    static String b_cfg_changed = "I am ok";


    public static void init(FinderManager finderManager,String project,String group,String service,String version) throws Exception {
       // System.out.println("basePath:"+baseConfigPath);
      //  finderManager = new FinderManager();
        ServiceMeteData serviceMeteData = new ServiceMeteData(project,group,service,version,"c1");
        BootConfig bootConfig = new BootConfig("http://10.1.87.70:6868",serviceMeteData);
        bootConfig.setZkSessionTimeout(10);
        bootConfig.setZkMaxSleepTime(20);
        bootConfig.setZkMaxRetryNum(1);
        bootConfig.setConfigCache(false);
        finderManager.init(bootConfig);


    }

    public static void initZkData(ZkHelper zkHelper){

    }



}
