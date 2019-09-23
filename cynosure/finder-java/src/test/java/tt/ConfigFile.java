package tt;

import utils.NumberUtils;
import utils.StreamUtils;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

public class ConfigFile {
    private File file;

    private Scanner sc;

    private String companionUrl;
    private String address;
   // private String dataFileUrl;
    private boolean cache =true;
    private long sleepTimeSecond=10000;
    private String project;
    private String group;
    private String service;
    private String version;
    private int registerService;
    private List<String> configFileNames = new ArrayList<String>();
    private List<String> unsubscribeFiles = new ArrayList<String>();


    public ConfigFile(String fileName) {
        file = new File(fileName);
        InputStream is=null;
        try {
            is = new FileInputStream(file);
            sc = new Scanner(is,"gbk");
            System.out.println("success load config file!");
            init();
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }finally {
            StreamUtils.closeInputStream(is);
            sc.close();
        }
    }

    private void init(){
        while (sc.hasNext()){
            String configLine = sc.nextLine();
            configLine = configLine.replace(" ","");
            if ("".equals(configLine))continue;
            parserString(configLine);
            parserString1(configLine);
        }

    }

    private void parserString(String  s){
        if (s == null || s.startsWith("#"))return;
        String []strs = s.split("=");
        if (strs == null && strs.length!=2)return;

        if ("companionUrl".equals(strs[0])){
            this.companionUrl=strs[1];
        } else if("address".equals(strs[0])){
            this.address=strs[1];
        } else if ("dataFileUrl".equals(strs[0])){
           // this.dataFileUrl=strs[1];
        } else if ("cache".equals(strs[0])){
            if ("false".equals(strs[1])){
                this.cache = false;
            }
        }
        else if ("sleep".equals(strs[0])){
            if (NumberUtils.isNumber(strs[1])){
                sleepTimeSecond = 1000*Integer.parseInt(strs[1]);
            }
        }
    }
    private void parserString1(String s){
        if (s == null || s.startsWith("#"))return;
        s = s.replace(" ","");
        String ss [] = s.split("=");
        if (ss == null || ss.length != 2){
            return;
        }
        if ("project".equals(ss[0])){
            this.project=ss[1];
            return;
        }
        if ("group".equals(ss[0])){
            this.group=ss[1];
            return;
        }
        if ("service".equals(ss[0])){
            this.service=ss[1];
            return;
        }
        if ("version".equals(ss[0])){
            this.version=ss[1];
            return;
        }
        if ("subscribeFile".equals(ss[0])){
            String files[]=ss[1].split(",");
            if (files == null)return;
            for (String f:files) {
                this.configFileNames.add(f);
            }
            return;
        }
        if ("unsubscribeFile".equals(ss[0])){
            String files[]=ss[1].split(",");
            if (files == null) return;
            for (String f:files) {
                this.unsubscribeFiles.add(f);
            }
            return;
        }

        if ("opt".equals(ss[0])){
            this.registerService =Integer.parseInt(ss[1]);
            return;
        }



    }

    public boolean isCache() {
        return cache;
    }

    public String getAddress() {
        return address;
    }

    public String getCompanionUrl() {
        return companionUrl;
    }


    public String getProject() {
        return project;
    }

    public String getGroup() {
        return group;
    }

    public String getService() {
        return service;
    }

    public String getVersion() {
        return version;
    }

    public int getRegister() {
        return registerService;
    }

    public List<String> getConfigFileNames() {
        return configFileNames;
    }

    public List<String> getUnsubscribeFiles() {
        return unsubscribeFiles;
    }

    @Override
    public String toString() {
        return "tt.ConfigFile{" +
                "companionUrl='" + companionUrl + '\'' +
                ", address='" + address + '\'' +
                //", dataFileUrl='" +// dataFileUrl + '\'' +
                ", cache=" + cache +
                ", sleepTimeSecond=" + sleepTimeSecond/1000 +"S\n"+
                ", project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", subscribleFiles=" + configFileNames +
                ", unsubscribeFiles=" + unsubscribeFiles +
                ",cache="+isCache()+
                '}';
    }

    public long getSleepTimeSecond() {
        return sleepTimeSecond;
    }
}
