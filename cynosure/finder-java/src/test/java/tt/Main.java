package tt;

public class Main {


    public static void main(String[] args) {
        if (args.length!=1){
            System.out.println("args conflict");return;
        }

        ConfigFile configFile = new ConfigFile(args[0]);
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

       // System.out.println("正在监听配置文件......................");
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
                for (String file:configFile.getUnsubscribeFiles()){
                    serviceDemo.unSubscribeConfig(file);
                }
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        for (String s:configFile.getUnsubscribeFiles()){
            serviceDemo.unSubscribeConfig(s);
        }
    }

}
