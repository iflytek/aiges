package configcenter;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.value.BootConfig;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.Config;
import com.iflytek.ccr.finder.value.ServiceMeteData;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.junit.Assert;
import org.testng.annotations.BeforeClass;
import org.junit.Test;
import utils.ListUtil;
import utils.Md5Util;

/**
 *2、在灰度组中
 * 1)只订阅一个配置文件，配置文件变更能收到通知
 * 2)订阅二个配置文件，配置文件变更能收到通知
 */
public class ConfigTest2 {

    static RouteService routeService;
    static FinderManager finderManager;
    static String project = "att";
    static String group = "att";
    static String service = "att";
    static String version = "1.0.0.1";
    static String baseConfigPath;
    static ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");

    static String a_cfg = "hello world";
    static String a_cfg_gray_1 = "gray hello world";
    static String a_cfg_gray_2 = "gray hello world dsfds";
    static String b_cfg = "hello thank you";
    static String b_cfg_gray = "hello thank you heheh";
    static String a_cfg_changed = "I am fine";
    static String b_cfg_changed = "I am ok";

    static String gray_group_1="gray1";
    static String gray_group_2="gray2";
    static String gray_data = "[{\"group_id\":\"gray1\",\"server_list\":[\"c1,c2\"]}]";
    static String changed_gray_data = "[{\"group_id\":\"gray2\",\"server_list\":[\"c1,c2\"]}]";
    public static final String pushId = "1234567890";

    @BeforeClass
    public static void setUp() throws Exception {
        baseConfigPath = "/polaris/" +
                "config/" +
                Md5Util.getMD5(project+group)+"/" +
                Md5Util.getMD5(service+version);

        System.out.println("basePath:"+baseConfigPath);
        finderManager = new FinderManager();
        SdkInit.init(finderManager,project,group,service,version);
        //准备数据

       zkHelper.addOrUpdatePersistentNode(baseConfigPath,"".getBytes());
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/a.cfg",ByteUtil.getZkBytes(a_cfg.getBytes(),pushId));
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/b.cfg",ByteUtil.getZkBytes(b_cfg.getBytes(),pushId));
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/gray/gray1/b.cfg",ByteUtil.getZkBytes(b_cfg_gray.getBytes(),pushId));
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/gray/gray1/a.cfg",ByteUtil.getZkBytes(a_cfg_gray_1.getBytes(),pushId));
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/gray/gray2/a.cfg",ByteUtil.getZkBytes(a_cfg_gray_2.getBytes(),pushId));
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/gray/gray3/a.cfg",ByteUtil.getZkBytes(a_cfg_gray_1.getBytes(),pushId));
       zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/gray",ByteUtil.getZkBytes(gray_data.getBytes(),pushId));


    }


    Config changedConfig;
    @Test
    public void test() throws Exception {

        CommonResult commonResult = finderManager.useAndSubscribeConfig(ListUtil.collectAsArrayList("a.cfg"), new ConfigChangedHandler() {
            @Override
            public boolean onConfigFileChanged(Config config) {
                System.out.println("1");
                changedConfig = config;
                return false;
            }
        });

        System.out.println(commonResult);
        changeFile();
        Assert.assertEquals(a_cfg_gray_2,new String (changedConfig.getFile()));
     //   Assert.assertEquals(a_cfg_changed,new String (changedConfig.getFile()));
        Thread.sleep(2000);

    }

    private void changeFile() throws InterruptedException {
        Thread.sleep(2000);
      //  zkHelper.update(baseConfigPath+"/gray/a.cfg",ByteUtil.getZkBytes(a_cfg_changed.getBytes(),"1234567890"));
      //  zkHelper.update(baseConfigPath+"/b.cfg",ByteUtil.getZkBytes(b_cfg_changed.getBytes(),"1234567890"));
        zkHelper.update(baseConfigPath+"/gray",ByteUtil.getZkBytes(changed_gray_data.getBytes(),pushId));
        Thread.sleep(1000);
    }
}
