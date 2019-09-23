package com.iflytek.ccr.polaris.cynosure.request.serviceversion;

import java.util.Date;

public class CopyAddServiceVersionRequestBody extends AddServiceVersionRequestBody {

	private static final long serialVersionUID = -5170693675561014373L;
	
	
	public CopyAddServiceVersionRequestBody() {
		super();
	}

	public CopyAddServiceVersionRequestBody(String version, String desc, String serviceId) {
		super(version, desc, serviceId);
	}

	private Date updateTime;

	public Date getUpdateTime() {
		return updateTime;
	}

	public void setUpdateTime(Date updateTime) {
		this.updateTime = updateTime;
	}

}
