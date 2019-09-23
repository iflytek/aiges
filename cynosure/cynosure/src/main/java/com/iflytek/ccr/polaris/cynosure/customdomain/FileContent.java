package com.iflytek.ccr.polaris.cynosure.customdomain;

/**
 * 文件内容
 *
 * @author sctang2
 * @create 2017-12-10 18:23
 **/
public class FileContent {
    //文件路径
    private String path;

    //文件名
    private String fileName;

    //内容
    private byte[] content;

    public String getFileName() {
        return fileName;
    }

    public void setFileName(String fileName) {
        this.fileName = fileName;
    }

    public byte[] getContent() {
        return content;
    }

    public void setContent(byte[] content) {
        this.content = content;
    }

    public String getPath() {
        return path;
    }

    public void setPath(String path) {
        this.path = path;
    }
}
