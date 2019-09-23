package servicefinder;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.service.ConfigFinder;
import com.iflytek.ccr.finder.service.impl.ConfigFinderImpl;
import com.iflytek.ccr.finder.value.BootConfig;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.Config;
import com.iflytek.ccr.finder.value.ServiceMeteData;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;

public class ConfigFinderTest {
    static ConfigFinder configFinder;
    static FinderManager finderManager;
    static String project = "jj";
    static String group = "jj";
    static String service = "jj";
    static String version = "1.0";
    static String baseConfigPath;

    static String data="dsfffffffffffffffffffff";
    @BeforeClass
    public static void setUp() throws Exception {

        //     System.out.println("BeforeClass");


        //http://10.1.86.223:9080/finder/query_zk_info?project=AIaaS&group=aitest&service=iatExecutor&version=2.0.0
        baseConfigPath = "/polaris/config/3b6281fa2ce2b6c20669490ef4b026a4/2f45195dae8ecdb5cb6afc1763579105";
//
       // baseConfigPath = "/polaris/config/c94370c924bac3f56e77434935613b23/5870ac10085ea78db59e6a231fefc6cf/gray/a";
      //  baseConfigPath = "/polaris/config/c94370c924bac3f56e77434935613b23/5870ac10085ea78db59e6a231fefc6cf/gray/b";

        ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");

        //zkHelper.addPersistent(baseConfigPath+"/gray/test.txt","".getBytes());

        //zkHelper.addPersistent(baseConfigPath+"/test.txt",data);


        finderManager = new FinderManager();
        ServiceMeteData serviceMeteData = new ServiceMeteData(project,group,service,version,"1.1.1.1:1234");
        BootConfig bootConfig = new BootConfig("http://10.1.87.70:6868",serviceMeteData);
        bootConfig.setZkSessionTimeout(1);
        bootConfig.setZkMaxSleepTime(1);
        bootConfig.setZkMaxRetryNum(1);
        bootConfig.setConfigCache(false);
        finderManager.init(bootConfig);

        List<String > list = new ArrayList<>();
        list.add("test.txt");
        //finderManager.useAndSubscribeConfig()

    }

    @AfterClass
    public static void tearDown() throws Exception {

        // System.out.println("tearDown");



    }



    @Test
    public void registerServiceTest() {

        testGetCurrentConfig();

    }

    public void testGetCurrentConfig(){
        List<String> list = new ArrayList<>();
        list.add("test.txt");
        configFinder = new ConfigFinderImpl();
        CommonResult<List<Config>> result = configFinder.getCurrentConfig(finderManager, baseConfigPath, list);
        System.out.println(result);
        Assert.assertEquals(data,new String (result.getData().get(0).getFile()));

    }

    public void test(){

    }
}
