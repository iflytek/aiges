package com.iflytek.ccr.polaris.cynosure.response;

import java.io.Serializable;
import java.util.List;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

/**
 * 查询分页列表-响应
 *
 * @author sctang2
 * @create 2017-11-15 8:51
 **/
@ApiModel("分页页表响应")
public class QueryPagingListResponseBody<T> implements Serializable {
	private static final long serialVersionUID = -103524281611219672L;

	// 当前页
	@ApiModelProperty("当前页")
	private int currentPage;

	// 总数
	@ApiModelProperty("总页数")
	private int totalCount;

	// 列表
	@ApiModelProperty("列表")
	private List<T> list;

	// 搜索条件
	@ApiModelProperty("搜索条件")
	private String condition;

	public int getCurrentPage() {
		return currentPage;
	}

	public void setCurrentPage(int currentPage) {
		this.currentPage = currentPage;
	}

	public int getTotalCount() {
		return totalCount;
	}

	public void setTotalCount(int totalCount) {
		this.totalCount = totalCount;
	}

	public List<T> getList() {
		return list;
	}

	public void setList(List<T> list) {
		this.list = list;
	}

	public String getCondition() {
		return condition;
	}

	public void setCondition(String condition) {
		this.condition = condition;
	}

	@Override
	public String toString() {
		return "QueryPagingListResponseBody{" + "currentPage=" + currentPage + ", totalCount=" + totalCount + ", list=" + list + ", condition='" + condition + '\'' + '}';
	}
}
