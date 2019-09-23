package configcenter;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.utils.RemoteUtil;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.CuratorFrameworkFactory;
import org.apache.curator.framework.imps.CuratorFrameworkImpl;
import org.apache.curator.test.TestingServer;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.mockito.runners.MockitoJUnitRunner;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;
import utils.ListUtil;
import utils.Md5Util;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 1、不在灰度组中
 * 1)只订阅一个配置文件，配置文件变更能收到通知
 * 2)订阅二个配置文件，配置文件变更能收到通知
 */
@RunWith(PowerMockRunner.class)
@PrepareForTest(RemoteUtil.class)
public class ConfigTest1 {

    static RouteService routeService;
    static FinderManager finderManager;
    static String project = "att";
    static String group = "att";
    static String service = "att";
    static String version = "1.0.0.2";
    static String baseConfigPath;
    static ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");

    static String a_cfg = "hello world";
    static String b_cfg = "hello thank you";
    static String a_cfg_changed = "I am fine";
    static String b_cfg_changed = "I am ok";

    public static final String pushId = "1234567890";

    @BeforeClass
    public static void setUp() throws Exception {
        baseConfigPath = "/polaris/" +
                "config/" +
                Md5Util.getMD5(project+group)+"/" +
                Md5Util.getMD5(service+version);

        //初始化测试服务器
       // TestingServer server = new TestingServer(2344);
     //   server.start();
//        zkHelper = new ZkHelper(server.getConnectString());
//        zkHelper.addOrUpdatePersistentNode(baseConfigPath,"".getBytes());
//        zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/a.cfg",ByteUtil.getZkBytes(a_cfg.getBytes(),"1234567890"));
//        zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/b.cfg",ByteUtil.getZkBytes(b_cfg.getBytes(),"1234567890"));

        finderManager = new FinderManager();
        //准备mock数据
        CommonResult result = new CommonResult();
        result.setMsg(null);
        result.setRet(Constants.SUCCESS);
        Map<String ,Object > map = new HashMap<>();
        map.put(Constants.CONFIG_PATH,"/polaris/config");
        map.put(Constants.SERVICE_PATH,"/polaris/service");
        map.put(Constants.KEY_ZK_NODE_PATH,"/polaris/zkNode");
        List<String > zkAddrList = new ArrayList<>();
        zkAddrList.add("127.0.0.1:2183");
        map.put(Constants.ZK_ADDR,zkAddrList);
        result.setData(map);
        //初始化sdk
        ServiceMeteData serviceMeteData = new ServiceMeteData(project,group,service,version,"c1");
        BootConfig bootConfig = new BootConfig("http://10.1.87.70:6868",serviceMeteData);
        bootConfig.setZkSessionTimeout(1);
        bootConfig.setZkMaxSleepTime(2);
        bootConfig.setZkMaxRetryNum(1);
        bootConfig.setConfigCache(false);

        PowerMockito.mockStatic(RemoteUtil.class);
        PowerMockito.when(RemoteUtil.queryZkInfo(bootConfig)).thenReturn(result);
        PowerMockito.when(RemoteUtil.registerServiceInfo(bootConfig,version)).thenReturn(null);
        finderManager.init(bootConfig);

        CommonResult commonResult = finderManager.useAndSubscribeConfig(ListUtil.collectAsArrayList("a.cfg", "b.cfg"), new ConfigChangedHandler() {
            @Override
            public boolean onConfigFileChanged(Config config) {
                System.out.println("1");
                changedConfig = config;
                return false;
            }
        });
        System.out.println(commonResult);

    }


    static Config changedConfig;
    @Test
    public void test() throws Exception {



        changeFile();
        Assert.assertEquals(a_cfg_changed,new String (changedConfig.getFile()));
     //   Assert.assertEquals(a_cfg_changed,new String (changedConfig.getFile()));
        Thread.sleep(2000);
    }

    private void changeFile() throws InterruptedException {
        Thread.sleep(2000);
        zkHelper.update(baseConfigPath+"/a.cfg",ByteUtil.getZkBytes(a_cfg_changed.getBytes(),"1234567890"));
      //  zkHelper.update(baseConfigPath+"/b.cfg",ByteUtil.getZkBytes(b_cfg_changed.getBytes(),"1234567890"));
        Thread.sleep(1000);
    }


}
