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

/**
 * 3/4）在路由规则中，路由规则新增一个存在的 服务实例
 * 预期： 收到实例上线通知
 */
public class RouteModuleTest4_3 {

    static RouteService routeService;
    static FinderManager finderManager;
    static String project = "ds";
    static String group = "ds";
    static String service = "se";
    static String version = "1.0";
    static String baseConfigPath;
    static ZkHelper zkHelper = new ZkHelper("10.1.87.69:2183");

    static List<ServiceInstance> a = new ArrayList<>();
    static List<ServiceInstance> b = new ArrayList<>();

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
        //    ",{\"id\":\"2\",\"consumer\":[\"11.2.3.4:8080\",\"2.2.3.4:8080\",\"199.99.99.99:99\"],\"provider\":[\"1.1.1.1:1111\",\"1.1.1.2:1111\",\"1.1.1.3:1111\",\"1.1.1.4:1111\"],\"only\":\"Y\"}}]";
        String configData = "{\"loadbalance\":\"loadbalance\",\"key1\":\"val\",\"key2\":\"val\"}";
        try {
            byte[] zkBytes = ByteUtil.getZkBytes(routeData.getBytes("UTF-8"), "1234567890");
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/route", zkBytes);
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider", "".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/conf", ByteUtil.getZkBytes(configData.getBytes("utf-8"),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p1",ByteUtil.getZkBytes(configData.getBytes("utf-8"),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p2","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p3","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p5","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p4","".getBytes());

        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }

    }

    @AfterClass
    public static void tearDown() throws Exception {

        // System.out.println("tearDown");



    }
    //不在路由规则中，提供者被加入到一个 only=N 的路由规则
    static final int ROUTE_CHANGED_FROM_YTON = 1;
    static String  routeData = "[{\"id\":\"1\",\"consumer\":[\"c1\",\"c3\",\"c2\"],\"provider\":[\"p1\",\"p2\",\"p3\",\"p4\"],\"only\":\"N\"}]";

    String changedRouteDataAdd = "[{\"id\":\"1\",\"consumer\":[\"c1\",\"c3\",\"c2\"],\"provider\":[\"p1\",\"p2\",\"p3\",\"p4\",\"p5\"],\"only\":\"N\"}]";

    String changedRouteDataRemove = "[{\"id\":\"1\",\"consumer\":[\"c1\",\"c3\",\"c2\"],\"provider\":[\"p1\",\"p2\",\"p3\"],\"only\":\"N\"}]";


    @Test
    public void moduleTest() {
        try {
            testSubscribeService();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
    int size= 0;
    List<InstanceChangedEvent> eventList;
    private void testSubscribeService() throws Exception {
        SubscribeRequestValue value = new SubscribeRequestValue();
        value.setServiceName(service);
        value.setApiVersion(version);

        CommonResult<List<Service>> listCommonResult = finderManager.useAndSubscribeService(ListUtil.collectAsList(value), new ServiceHandle() {
            @Override
            public boolean onServiceInstanceConfigChanged(String serviceName, String instance, String jsonConfig) {
                System.out.println("1 called");
                Assert.assertEquals(instanceData, jsonConfig);

                return false;
            }

            @Override
            public boolean onServiceConfigChanged(String serviceName, String jsonConfig) {
                System.out.println("2 called");
                System.out.println("jsonConfig:" + jsonConfig);

                Assert.assertEquals(changedData, jsonConfig);
                return false;
            }

            @Override
            public boolean onServiceInstanceChanged(String serviceName, List<InstanceChangedEvent> eventList) {
                System.out.println("3 called");
                System.out.println(eventList.size());
                size = eventList.size();
                System.out.println(eventList);
                //Assert.assertEquals(eventList.size(),23);
                RouteModuleTest4_3.this.eventList=eventList;
                System.out.println(eventList.get(0).getServiceInstanceList());
                return false;
            }
        });

        changeRouteFromNtoY();

        Thread.sleep(3000);

        List expected = ListUtil.collectAsArrayList( "p4");
        List<ServiceInstance> serviceInstanceList = eventList.get(0).getServiceInstanceList();
        List acture = ListUtil.collectAsArrayList(serviceInstanceList.get(0).getAddr());
        Assert.assertEquals(1,size);
        Assert.assertEquals(InstanceChangedEvent.Type.REMVOE,eventList.get(0).getType());
        Assert.assertEquals(true,ListUtil.equals(expected,acture));


    }

    private void changeRouteFromNtoY(){
        try {
            Thread.sleep(1000);
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        zkHelper.addOrUpdatePersistentNode(baseConfigPath+"/route",ByteUtil.getZkBytes(changedRouteDataRemove.getBytes(),"1234567890"));
    }
}
