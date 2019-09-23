package fluent.httpclient;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import org.apache.http.HttpHost;
import org.apache.http.NameValuePair;
import org.apache.http.message.BasicNameValuePair;
import org.junit.Test;

import java.io.File;
import java.io.IOException;
import java.net.URISyntaxException;
import java.util.HashMap;
import java.util.List;
import java.util.logging.Logger;

public class UnitTest {
    private String url = "http://10.1.87.70:6868/service/push_cluster_config";
    private HttpFluentService fluentService;
    private Logger logger;
    public void init(){
        fluentService = new HttpFluentService();
        logger = Logger.getLogger("");
    }
    @Test
    public void testGet() {
        String result = fluentService.doGet(url);
        logger.info("返回结果："+result);
    }

    @Test
    public void testGetParam() {
        HashMap<String, String> map = Maps.newHashMap();
        map.put("name", "wang");
        map.put("value", "hello");
        String result = fluentService.doGet(url, map);
        logger.info("返回结果："+result);
    }

    @Test
    public void testGetProxy() {
        HashMap<String, String> map = Maps.newHashMap();
        map.put("name", "wang");
        map.put("value", "hello");
        String result = fluentService.doGet(url, "localhost", 8888, null, map);
        logger.info("返回结果："+result);
    }

    @Test
    public void testPost() {
        String result = fluentService.doPost(url, null);
        logger.info("返回结果："+result);
    }

    @Test
    public void testPostProxy() {
        String result = fluentService.doPost(url, "119.28.99.246", 8998, null, null, null);
        logger.info("返回结果："+result);
    }

    @Test
    public void testPostParam() throws IOException, URISyntaxException {
        List<NameValuePair> nameValuePairs = Lists.newArrayList();
        nameValuePairs.add(new BasicNameValuePair("name", "王志雄"));
//        List<File> fileList = Lists.newArrayList();
//        fileList.add(new File("D:\\response.html"));
        String result = fluentService.doPost(url, "localhost", 8888, HttpHost.DEFAULT_SCHEME_NAME, nameValuePairs, null);
        logger.info("返回结果："+result);
    }
}
