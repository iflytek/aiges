import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.utils.RemoteUtil;
import com.iflytek.ccr.finder.value.BootConfig;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.ServiceMeteData;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.test.TestingServer;
import org.junit.AfterClass;
import org.junit.BeforeClass;
import org.junit.runner.RunWith;
import org.powermock.api.mockito.PowerMockito;
import org.powermock.core.classloader.annotations.PowerMockIgnore;
import org.powermock.core.classloader.annotations.PrepareForTest;
import org.powermock.modules.junit4.PowerMockRunner;

@RunWith(PowerMockRunner.class)
@PrepareForTest(RemoteUtil.class)
@PowerMockIgnore({"com.iflytek.ccr.zkutil.*", "javax.xml.*"})
public class Test {
    static TestingServer server = null;
    static String project = "att";
    static String group = "att";
    static String service = "att";
    static String version = "1.0.0.2";
    static String a_cfg = "hello world";
    static String b_cfg = "hello thank you";
    static String a_cfg_changed = "I am fine";
    static String b_cfg_changed = "I am ok";
    @BeforeClass
    public static void setUp() throws Exception {
        System.out.println("BeforeClass");
//          server = new TestingServer();
    }

    @AfterClass
    public static void tearDown() throws Exception {
        System.out.println("tearDown");
    }

    @org.junit.Test
    public void registerServiceTest() {

        CommonResult result = new CommonResult();
        result.setMsg(null);
        result.setRet(Constants.SUCCESS);

        ServiceMeteData serviceMeteData = new ServiceMeteData(project,group,service,version,"c1");
        BootConfig bootConfig = new BootConfig("http://10.1.87.70:6868",serviceMeteData);
        bootConfig.setZkSessionTimeout(1);
        bootConfig.setZkConnectTimeout(1);
        bootConfig.setZkMaxSleepTime(2);
        bootConfig.setZkMaxRetryNum(1);
        bootConfig.setConfigCache(false);
        System.out.println("registerServiceTest");
        PowerMockito.mockStatic(RemoteUtil.class);
        PowerMockito.when(RemoteUtil.queryZkInfo(bootConfig)).thenReturn(result);
        ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");
        System.out.println(zkHelper.checkExists("/"));
        System.out.println(zkHelper.addOrUpdatePersistentNode("/a/b",""));
        System.out.println(RemoteUtil.queryZkInfo(bootConfig));
        PowerMockito.when(RemoteUtil.registerServiceInfo(bootConfig,version)).thenReturn(null);
        PowerMockito.when(RemoteUtil.pushConfigFeedback(null,null,null,"","","","","")).thenReturn(null);
        System.out.println(zkHelper.checkExists("/"));
    }
}
