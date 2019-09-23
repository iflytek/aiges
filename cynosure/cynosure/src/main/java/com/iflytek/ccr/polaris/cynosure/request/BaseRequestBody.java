package com.iflytek.ccr.polaris.cynosure.request;

/**
 * 基类-请求
 *
 * @author sctang2
 * @create 2017-11-14 9:56
 **/
public class BaseRequestBody {
    //当前页
    private int currentPage = 1;

    //每页大小
    private int pagesize = 10;

    //是否需要分页，默认为分页。0.不分页 1.分页
    private int isPage = 1;

    public int getCurrentPage() {
        return currentPage;
    }

    public void setCurrentPage(int currentPage) {
        this.currentPage = currentPage;
    }

    public int getPagesize() {
        return pagesize;
    }

    public void setPagesize(int pagesize) {
        this.pagesize = pagesize;
    }

    public int getIsPage() {
        return isPage;
    }

    public void setIsPage(int isPage) {
        this.isPage = isPage;
    }

    @Override
    public String toString() {
        return "BaseRequestBody{" +
                ", currentPage=" + currentPage +
                ", pagesize=" + pagesize +
                ", isPage=" + isPage +
                '}';
    }
}
