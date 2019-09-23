package servicefinder;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.service.impl.RouteServiceImpl;
import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import configcenter.SdkInit;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;
import utils.ListUtil;
import utils.Md5Util;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.List;

public class ModuleTest {


    static RouteService routeService;
    static FinderManager finderManager;
    static String project = "ds";
    static String group = "ds";
    static String service = "se";
    static String version = "2.0";
    static String baseConfigPath;
    static ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");

    static List<ServiceInstance> a = new ArrayList<>();
    static List<ServiceInstance> b = new ArrayList<>();
    static String routeData;
    //改变后的配置节点
    String changedData = "{\"loace\":\"loance\",\"key1\":\"val\",\"key2\":\"val\"}";
    String instanceData = "{\"loacdffdgfde\":\"loagfdgfdgnce\",\"kgdfgey1\":\"vagfl\",\"kegdfgy2\":\"gdfgval\"}";
    @BeforeClass
    public static void setUp() throws Exception {
        baseConfigPath = "/polaris/" +
                "service/" +
                Md5Util.getMD5(project+group)+"/" +
                service +"/"+
                version;
        pre();

        finderManager = new FinderManager();
        SdkInit.init(finderManager,project,group,service,version);
        routeService = new RouteServiceImpl();


    }

    public static void pre(){
        routeData = "[{\"id\":\"1\",\"consumer\":[\"11.2.3.4:8080\",\"2.2.3.4:8080\",\"199.99.99.99:99\"],\"provider\":[\"1.1.1.1:1111\",\"1.1.1.2:1111\",\"1.1.1.3:1111\",\"1.1.1.4:1111\"],\"only\":\"Y\"}]";
        //    ",{\"id\":\"2\",\"consumer\":[\"11.2.3.4:8080\",\"2.2.3.4:8080\",\"199.99.99.99:99\"],\"provider\":[\"1.1.1.1:1111\",\"1.1.1.2:1111\",\"1.1.1.3:1111\",\"1.1.1.4:1111\"],\"only\":\"Y\"}}]";
        String configData = "{\"loadbalance\":\"loadbalance\",\"key1\":\"val\",\"key2\":\"val\"}";
        try {
            byte[] zkBytes = ByteUtil.getZkBytes(routeData.getBytes("UTF-8"), "1234567890");
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/route", zkBytes);
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider", "".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/conf", ByteUtil.getZkBytes(configData.getBytes("utf-8"),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/1.1.1.1:1111",ByteUtil.getZkBytes(configData.getBytes("utf-8"),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/1.1.1.2:1111","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/1.1.1.3:1111","".getBytes());

        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }



//        ServiceInstance s1 = new ServiceInstance();s1.setAddr("1.1.1.1:1111");
//        ServiceInstance s2 = new ServiceInstance();s2.setAddr("1.1.1.2:1111");
//        ServiceInstance s3 = new ServiceInstance();s3.setAddr("1.1.1.3:1111");
//        ServiceInstance s4 = new ServiceInstance();s4.setAddr("1.1.1.4:1111");
//        ServiceInstance s5 = new ServiceInstance();s5.setAddr("1.1.1.5:1111");
//
//        a.add(s1);
//        a.add(s2);
//        a.add(s3);
//        b.add(s1);
//        b.add(s2);
//        b.add(s3);
//        b.add(s4);
//        b.add(s5);


    }

    @AfterClass
    public static void tearDown() throws Exception {

        // System.out.println("tearDown");



    }
//不在路由规则中，提供者被加入到一个 only=N 的路由规则
    String changedRouteData = "[{\"id\":\"1\",\"consumer\":[\"11.2.3.4:8080\",\"2.2.3.4:8080\",\"199.99.99.99:99\"],\"provider\":" +
            "[\"1.1.1.1:1111\",\"1.1.1.2:1111\",\"1.1.1.3:1111\",\"1.1.1.4:1111\"],\"only\":\"Y\"},{\"id\":\"2\",\"consumer\":[\"11.2.3.4:8080\"],\"provider\":[\"1.1.1.4:1111\"],\"only\":\"N\"}]";


    @Test
    public void moduleTest() {
        try {
            testSubscribeService();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


    private void testSubscribeService() throws Exception {
        SubscribeRequestValue value = new SubscribeRequestValue();
        value.setServiceName(service);
        value.setApiVersion(version);

        CommonResult<List<Service>> listCommonResult = finderManager.useAndSubscribeService(ListUtil.collectAsList(value), new ServiceHandle() {
            @Override
            public boolean onServiceInstanceConfigChanged(String serviceName, String instance, String jsonConfig) {
                System.out.println("1");
                Assert.assertEquals(instanceData,jsonConfig);
                return false;
            }

            @Override
            public boolean onServiceConfigChanged(String serviceName, String jsonConfig) {
                System.out.println("2");
                System.out.println("jsonConfig:"+jsonConfig);

                Assert.assertEquals(changedData,jsonConfig);
                return false;
            }

            @Override
            public boolean onServiceInstanceChanged(String serviceName, List<InstanceChangedEvent> eventList) {
                System.out.println("3");

                return false;
            }
        });
      //  changeConfigFile();
       // removeInstance();
        changeInstanceData();
        //changeRouteDataFromYtoN();
        System.out.println(listCommonResult);
        Thread.sleep(3000);
    }

    /**
     * 1 改变配置文件
     * @throws Exception
     */

    private void changeConfigFile() throws Exception{
        Thread.sleep(2000);
        zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/conf",ByteUtil.getZkBytes(changedData.getBytes(),"1234567890"));
    }

    /**
     * 2 下线一个服务提供者
     */
    private void removeInstance() throws InterruptedException {
        Thread.sleep(2000);
        zkHelper.remove(baseConfigPath+"/provider/1.1.1.1:1111");
    }

    /**
     * 3 改变实例配置文件
     * @throws InterruptedException
     */
    private void changeInstanceData() throws InterruptedException {
        Thread.sleep(2000);
        zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/provider/1.1.1.1:1111",ByteUtil.getZkBytes(instanceData.getBytes(),"01234567890"));
    }

    /**
     * 不在路由规则中，提供者被加入到一个 only=N 的路由规则
     */
    private void changeRouteDataFromYtoN() throws InterruptedException {
        Thread.sleep(1000);
        zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/route",ByteUtil.getZkBytes(changedRouteData.getBytes(),"1234567890"));
        System.out.println(changedRouteData);
    }
}
