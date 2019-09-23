![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/polaris-w.png)
# Written for Polaris
## Core Functions
 - Dynamic Config
 - Register Service & Discover Service
 
## Features
 - Full category for managing service 
 - Support multi service versions
 - Support roll back config
 - Support feedback for pushing config
 - Support management for provider and consumer online
 - High available for some not expected cases
 - Easy Integration
 - Support Delivery by docker 
 
## Architecture
 ![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/15138432239107.jpg)

## Core Modules
- **Cynosure**

    `This is a module which used for managing some basic data and config data. For example, you can use this for creating region info、project info、cluster info、service info、config info and push the config info to client.`
    
   ![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/15138461728383.jpg)
 
 
- **Companion**

    `This module is used to operating zookeeper for the Cynosure and receiving feedback from client which has integrated with the Finder.`
    
    **case #1:**
    
    ![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/15138469634925.jpg)

    **case #2:**
    
    ![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/15138470006112.jpg)

    **case #3**
    
    ![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/15138470283686.jpg)

- **Finder(SDK with go/java/c++)**

    `As you expected, it's used for integrating with your program. You can do this by calling some functions easily.`

  
## How to use in your project?

- **Install**

    `You can get from `[install.md](http://git.xfyun.cn/AIaaS/polaris/src/master/install.md)


- **SDK**[supported]

    `Beginning it with the Finder SDK as you know, here are some examples in golang. More detail code has uploaded and you can view from `[finder-go/example/demo.go](http://git.xfyun.cn/AIaaS/finder-go/src/master/example/demo.go)
    
    ```go
    package main
    import (
	       "encoding/json"
	       "finder-go"
	       "finder-go/common"
	       "finder-go/utils/httputil"
	       "fmt"
	       "net"
	       "net/http"
	       "os"
	       "time"
    )    
    func main() {
       cachePath, err := os.Getwd()
	   if err != nil {
		  return
       }        
        
	   cachePath += "/findercache"
	   config := common.BootConfig{
		  CompanionUrl:     "http://    10.1.86.223:9080",
		  CachePath:        cachePath,
		  TickerDuration:   5000,
		  ZkSessionTimeout: 1000 * time.Second,
		  ZkConnectTimeout: 300 * time.Second,
		  ZkMaxSleepTime:   15 * time.Second,
		  ZkMaxRetryNum:    3,
		  MeteData: &common.ServiceMeteData{
			 Project: "test",
			 Group:   "default",
			 Service: "xsf",
			 Version: "1.0.0",
			 Address: "192.168.1.2:9091",
		  },
	   }	
       
	   f, err := finder.NewFinder(config)
	   if err != nil {
		  fmt.Println(err)
	   }	   
       
	   // use config with watcher
	   handler := new(ConfigChangedHandle)
	   configFiles, err := f.ConfigFinder.UseAndSubscribeConfig([]string{"default.cfg", "xsfc.tmol"}, handler)	   
       
	   // register service 
	   err = f.ServiceFinder.RegisterService()
	   if err != nil {
		  fmt.Println(err)
	   } 	   
       
	   // describe service
	   handler := new(ServiceChangedHandle)
	   serviceList, err = f.ServiceFinder.UseAndSubscribeService([]string{"xsf"}, handler)
	   if err != nil {
		  fmt.Println(err)
	   }	   
       
	   //todo business...
    }
    ```
       
- **Agent**[planning]

    `Also,you can have an integration with our agent without coding. This is in the planning stage already.`

![](http://git.xfyun.cn/ypwu/markdown-resource/raw/master/15138545464059.jpg)


