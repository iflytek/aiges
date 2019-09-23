import utils.StreamUtils;

import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.InputStream;
import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

public class SubscribeDataFile {
    private Scanner sc;

    /**
     * 要订阅的文件
     */
    private List<String> configFileNames = new ArrayList<String>();
    private List<String> unsubscribeFiles = new ArrayList<String>();


    private String project;

    private String group;

    private String  service;

    private String version;

   // private tt.ConfigFile configFile = new tt.ConfigFile();
    public SubscribeDataFile(String fileUrl){

        InputStream is = null;
        try {
            is = new FileInputStream(fileUrl);
            sc = new Scanner(is,"gbk");
            System.out.println("success load datafile");
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
            String line = sc.nextLine();
            parserString(line);
        }
    }

    private void parserString(String s){
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
        if ("configFile".equals(ss[0])){
            String files[]=ss[1].split(",");
            if (files == null)return;
            for (String f:files) {
                this.configFileNames.add(f);
            }
            return;
        }
        if ("unscrib".equals(ss[0])){
            String files[]=ss[1].split(",");
            if (files == null) return;
            for (String f:files) {
                this.unsubscribeFiles.add(f);
            }
            return;
        }

    }

    public List<String> getConfigFileNames() {
        return configFileNames;
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

    public List<String> getUnsubscribeFiles() {
        return unsubscribeFiles;
    }

    @Override
    public String toString() {
        return "SubscribeDataFile{" +
                "configFileNames=" + configFileNames +
                ", unsubscribeFiles=" + unsubscribeFiles +
                ", project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                '}';
    }
}
