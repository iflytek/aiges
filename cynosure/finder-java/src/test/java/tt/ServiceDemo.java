package tt;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.value.*;
import tt.ConfigFile;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;

public class ServiceDemo {
    ServiceMeteData serviceMeteData;
    BootConfig bootConfig;
    FinderManager finderManager;
    ConfigChangedHandler configChangedHandler;
    ConfigFile configFile;
    public ServiceDemo(ConfigFile configFile){
        this.configFile = configFile;
        String project = configFile.getProject();
        String group = configFile.getGroup();
        String service = configFile.getService();
        String version = configFile.getVersion();
        String address = configFile.getAddress();
        try {
            serviceMeteData =new ServiceMeteData(project,group,service,version,address);
            bootConfig = new BootConfig(configFile.getCompanionUrl(),serviceMeteData);
        }catch (Exception e){
            e.printStackTrace();
        }
        bootConfig.setZkSessionTimeout(1);
        bootConfig.setZkMaxSleepTime(1);
        bootConfig.setZkMaxRetryNum(1);
        bootConfig.setConfigCache(configFile.isCache());
        finderManager = new FinderManager();
        finderManager.init(bootConfig);
        List<String> configNameList = new ArrayList<String>();
        configChangedHandler = new ConfigChangedHandler() {
           // @Override
            public boolean onConfigFileChanged(Config config) {
                try {

                    System.out.println("配置已改变--"+new Date() + ":" + config.getName() + ":" + new String(config.getFile(), Constants.DEFAULT_CHARSET));

                   if (config.getConfigMap()!=null){
                       System.out.println("configMap:");
                       System.out.println(config.getConfigMap());
                   }
                } catch (UnsupportedEncodingException e) {
                    e.printStackTrace();
                }

                return true;
            }
        };

    }

    public void setConfigChangedHandler(ConfigChangedHandler configChangedHandler) {
        this.configChangedHandler = configChangedHandler;
    }

    public CommonResult<List<Config>> useAndSubscribeConfig(List<String> configNames){
        CommonResult<List<Config>> result = finderManager.useAndSubscribeConfig(configNames, configChangedHandler);
        StringBuilder sb = new StringBuilder();
        sb.append("subscribe result:  "+"ret="+result.getRet()).append(" config{");
        if (result.getData()!=null)
        for (Config c:result.getData()){
            sb.append("\n[name:"+c.getName()).append(",file:"+new String(c.getFile())).append(",configMap="+c.getConfigMap()+"]");
        }
        sb.append("\n}Msg="+result.getMsg());
        System.out.println(sb);

        return result;
    }

    public CommonResult unSubscribeConfig(String configName){
        CommonResult commonResult = finderManager.unSubscribeConfig(configName);
        System.out.println("unsubscribe result:  "+configName+":"+commonResult);
        return commonResult;
    }


    public void registerService(){

        CommonResult commonResult = finderManager.registerService(configFile.getAddress(), configFile.getVersion());
        System.out.println("registerServiceResult: "+commonResult);
    }

    public void unRegisterService(){
        finderManager.unRegisterService(configFile.getAddress(),configFile.getVersion());
    }

    public void subscribleService(){

        ServiceHandle serviceHandle = new ServiceHandle() {
            public boolean onServiceInstanceConfigChanged(String serviceName, String instance, String jsonConfig) {
                System.out.println("instanceConfigChanged:"+serviceName+" "+instance+ " "+jsonConfig);

                return false;
            }

            public boolean onServiceConfigChanged(String serviceName, String jsonConfig) {
                System.out.println("serviceConfigChanged :"+serviceName+" " +jsonConfig);
                return false;
            }

            public boolean onServiceInstanceChanged(String s, List<InstanceChangedEvent> list) {

                for(InstanceChangedEvent event:list){
                    System.out.println(event.getType()+" "+event.getServiceInstanceList());
                }
                return false;
            }
        };

        SubscribeRequestValue requestValue =new SubscribeRequestValue();
        requestValue.setApiVersion(configFile.getVersion());
        requestValue.setServiceName(configFile.getService());
        List<SubscribeRequestValue> list = new ArrayList<>();
        list.add(requestValue);
        CommonResult<List<Service>> listCommonResult = finderManager.useAndSubscribeService(list, serviceHandle);
        System.out.println(listCommonResult);
    }

    public void unSubscribeService(){
        SubscribeRequestValue requestValue =new SubscribeRequestValue();
        requestValue.setApiVersion(configFile.getVersion());
        requestValue.setServiceName(configFile.getService());
        finderManager.unSubscribeService(requestValue);

    }


}
