package servicefinder;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.service.impl.RouteServiceImpl;
import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.junit.AfterClass;
import org.junit.Assert;
import org.junit.BeforeClass;
import org.junit.Test;
import utils.Md5Util;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.List;

public class RouteServiceTest {


    static RouteService routeService;
    static FinderManager finderManager;
    static String project = "ds";
    static String group = "ds";
    static String service = "se";
    static String version = "3.0";
    static String baseConfigPath;
    static ZkHelper  zkHelper = new ZkHelper("10.1.87.69:2183");

    static List<ServiceInstance> a = new ArrayList<>();
    static List<ServiceInstance> b = new ArrayList<>();
    static   String instanceDataY = "1234567890{\"user\": {\"loadbalance\": \"loadbalance\",\"key1\": \"val\",\"key2\": \"val\"},\"sdk\": {\"isValid\": \"Y\"}}";
    static   String instanceDataN = "1234567890{\"user\": {\"loadbalance\": \"loadbalance\",\"key1\": \"val\",\"key2\": \"val\"},\"sdk\": {\"isValid\": \"N\"}}";

    @BeforeClass
    public static void setUp() throws Exception {
        baseConfigPath = "/polaris/" +
                "service/" +
                Md5Util.getMD5(project+group)+"/" +
                service +"/"+
                version;
        pre();

        finderManager = new FinderManager();
        ServiceMeteData serviceMeteData = new ServiceMeteData(project,group,service,version,"1.1.1.1:1234");
        BootConfig bootConfig = new BootConfig("http://10.1.87.70:6868",serviceMeteData);
        bootConfig.setZkSessionTimeout(1);
        bootConfig.setZkMaxSleepTime(1);
        bootConfig.setZkMaxRetryNum(1);
        bootConfig.setConfigCache(false);
        finderManager.init(bootConfig);
        routeService = new RouteServiceImpl();
        List<String > list = new ArrayList<>();
        list.add("test.txt");

    }

    public static void pre(){
        String routeData = "[{\"id\":\"1\",\"consumer\":[\"11.2.3.4:8080\",\"2.2.3.4:8080\",\"199.99.99.99:99\"],\"provider\":[\"1.1.1.1:1111\",\"1.1.1.2:1111\",\"1.1.1.3:1111\",\"1.1.1.4:1111\"],\"only\":\"Y\"}]";
                     //      ",{\"id\":\"2\",\"consumer\":[\"11.2.3.4:8080\",\"2.2.3.4:8080\",\"199.99.99.99:99\"],\"provider\":[\"1.1.1.1:1111\",\"1.1.1.2:1111\",\"1.1.1.3:1111\",\"1.1.1.4:1111\"],\"only\":\"Y\"}}]";
        String configData = "{\"loadbalance\":\"loadbalance\",\"key1\":\"val\",\"key2\":\"val\"}";
        try {
            byte[] zkBytes = ByteUtil.getZkBytes(routeData.getBytes("utf-8"), "1234567890");
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/route", zkBytes);
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider", "".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/conf", ByteUtil.getZkBytes(configData.getBytes("utf-8"),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider","".getBytes());
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p1",ByteUtil.getZkBytes(instanceDataN.getBytes("utf-8"),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p2",ByteUtil.getZkBytes(instanceDataY.getBytes(),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p3",ByteUtil.getZkBytes(instanceDataY.getBytes(),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p5",ByteUtil.getZkBytes(instanceDataY.getBytes(),"1234567890"));
            zkHelper.addOrUpdatePersistentNode(baseConfigPath + "/provider/p4",ByteUtil.getZkBytes(instanceDataN.getBytes(),"1234567890"));


        } catch (UnsupportedEncodingException e) {
            e.printStackTrace();
        }

        ServiceInstance s1 = new ServiceInstance();s1.setAddr("1.1.1.1:1111");
        ServiceInstance s2 = new ServiceInstance();s2.setAddr("1.1.1.2:1111");
        ServiceInstance s3 = new ServiceInstance();s3.setAddr("1.1.1.3:1111");
        ServiceInstance s4 = new ServiceInstance();s4.setAddr("1.1.1.4:1111");
        ServiceInstance s5 = new ServiceInstance();s5.setAddr("1.1.1.5:1111");

        a.add(s1);
        a.add(s2);
        a.add(s3);
        b.add(s1);
        b.add(s2);
        b.add(s3);
        b.add(s4);
        b.add(s5);


    }

    @AfterClass
    public static void tearDown() throws Exception {

        // System.out.println("tearDown");



    }



    @Test
    public void registerServiceTest() {
        testParserInstanceData();
        testParserRouteData();
        testParserConfData();
        testParseServiceInstanceByPath();
        testCompareServiceInstanceList();
        testFlterServiceInstanceByRoute();
    }


    private void testParserRouteData(){
        List<ServiceRouteValue> serviceRouteValues = routeService.parseRouteData(baseConfigPath + "/route");
        // byte[] byteData = zkHelper.getByteData(baseConfigPath + "/route");
        // System.out.println(new String(byteData));
        // System.out.println(serviceRouteValues);
        String expectedConsumer = "[11.2.3.4:8080, 2.2.3.4:8080, 199.99.99.99:99]";
        String expectedProvider = "[1.1.1.1:1111, 1.1.1.2:1111, 1.1.1.3:1111, 1.1.1.4:1111]";
        String expectedOnly = "Y";
        Assert.assertEquals(expectedConsumer,serviceRouteValues.get(0).getConsumer().toString());
        Assert.assertEquals(expectedProvider,serviceRouteValues.get(0).getProvider().toString());
        Assert.assertEquals(1,serviceRouteValues.size());
        Assert.assertEquals(expectedOnly,serviceRouteValues.get(0).getOnly());
    }

    private void testParserConfData(){
        String s = routeService.parseConfData(baseConfigPath + "/conf");
        String expectedJson = "{\"loadbalance\":\"loadbalance\",\"key1\":\"val\",\"key2\":\"val\"}";
        Assert.assertEquals(expectedJson,s);
    }

    private void testParseServiceInstanceByPath(){
        List<ServiceInstance> serviceInstances = routeService.parseServiceInstanceByPath(finderManager, baseConfigPath + "/provider");
        Assert.assertEquals(3,serviceInstances.size());

    }

    private void testCompareServiceInstanceList(){

        List<InstanceChangedEvent> instanceChangedEvents = routeService.compareServiceInstanceList(b, a);
        //System.out.println(instanceChangedEvents.size());
        for (InstanceChangedEvent e:instanceChangedEvents){
            Assert.assertEquals(InstanceChangedEvent.Type.REMVOE,e.getType());
            Assert.assertEquals(2,e.getServiceInstanceList().size());
        }

    }

    private void testFlterServiceInstanceByRoute(){
        List<ServiceInstance> serviceInstances = routeService.parseServiceInstanceByPath(finderManager, baseConfigPath + "/provider");
        List<ServiceRouteValue> serviceRouteValues = routeService.parseRouteData(baseConfigPath + "/route");
        List<ServiceInstance> rservicer = routeService.filterServiceInstanceByRoute(serviceInstances, serviceRouteValues, "1.1.1.1:1111");

        Assert.assertEquals(0,rservicer.size());
        List<ServiceInstance> rservicer2 = routeService.filterServiceInstanceByRoute(serviceInstances, serviceRouteValues, "11.2.3.4:8080");
        Assert.assertEquals(rservicer2.size(),3);

    }

    private void testParserInstanceData(){
        ServiceInstance instance = routeService.parseServiceInstanceData(finderManager, baseConfigPath + "/provider/p1");
        System.out.println(instance);
    }
}
