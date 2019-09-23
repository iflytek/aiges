package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.EditServiceDiscoveryRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;

public interface IMyServiceDiscovery {
    Response<String> edit1(EditServiceDiscoveryRequestBody body);
}
