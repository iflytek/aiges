package testt;

import tt.ConfigFile;
import tt.ServiceDemo;

public class consumer1 {


    public static void main(String[] args) {


        ConfigFile configFile = new ConfigFile("F:/cfg/cons1.cfg");
        System.out.println(configFile.toString());
        ServiceDemo serviceDemo = new ServiceDemo(configFile);
        if (configFile.getRegister() == 0){
            serviceDemo.registerService();
        }
        else if (configFile.getRegister() == 1){
            serviceDemo.subscribleService();
        }
        else{
            serviceDemo.useAndSubscribeConfig(configFile.getConfigFileNames());
        }

        //System.out.println("正在监听配置文件......................");
        try {
            Thread.sleep(configFile.getSleepTimeSecond());
            System.out.println("开始取消订阅......................");
            if (configFile.getRegister() == 0){
                serviceDemo.unRegisterService();
            }
            else if (configFile.getRegister() == 1){
                serviceDemo.unSubscribeService();
            }
            else{
                serviceDemo.unSubscribeConfig(configFile.getUnsubscribeFiles().get(0));
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        for (String s:configFile.getUnsubscribeFiles()){
            serviceDemo.unSubscribeConfig(s);
        }
    }

}
