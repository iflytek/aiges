package test;

import tt.ConfigFile;
import tt.ServiceDemo;

public class Main2 {
    public static void newS() {
        //加载并打相关印配置文件

        ConfigFile configFile = new ConfigFile("F:\\config.cfg");
        System.out.println(configFile.toString());
//        SubscribeDataFile dataFile = new SubscribeDataFile(configFile.getDataFileUrl());
//        System.out.println(dataFile.toString());

        ServiceDemo serviceDemo = new ServiceDemo(configFile);

        //设置回调函数，配置文件改变时促发
//        serviceDemo.setConfigChangedHandler(new ConfigChangedHandler() {
//            @Override
//            public boolean onConfigFileChanged(Config config) {
//                return false;
//            }
//        });

        //订阅
        serviceDemo.useAndSubscribeConfig(configFile.getConfigFileNames());
        System.out.println("正在监听配置文件......................");
        try {
            Thread.sleep(configFile.getSleepTimeSecond());
            System.out.println("开始取消订阅......................");
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        for (String s:configFile.getUnsubscribeFiles()){
            serviceDemo.unSubscribeConfig(s);
        }

        //取消订阅

    }


    public static void main(String[] args) {

        ConfigFile configFile = new ConfigFile("F:/config2.cfg");
        System.out.println(configFile.toString());
//        SubscribeDataFile dataFile = new SubscribeDataFile(configFile.getDataFileUrl());
//        System.out.println(dataFile.toString());

        ServiceDemo serviceDemo = new ServiceDemo(configFile);

        //设置回调函数，配置文件改变时促发
//        serviceDemo.setConfigChangedHandler(new ConfigChangedHandler() {
//            @Override
//            public boolean onConfigFileChanged(Config config) {
//                return false;
//            }
//        });

        //订阅

        serviceDemo.useAndSubscribeConfig(configFile.getConfigFileNames());
        System.out.println("正在监听配置文件......................");
        try {
            Thread.sleep(configFile.getSleepTimeSecond());
            System.out.println("开始取消订阅......................");
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

        for (String s:configFile.getUnsubscribeFiles()){
            serviceDemo.unSubscribeConfig(s);
        }
    }

}
