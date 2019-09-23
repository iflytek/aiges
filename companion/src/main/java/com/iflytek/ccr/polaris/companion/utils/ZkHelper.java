package com.iflytek.ccr.polaris.companion.utils;


import org.apache.curator.RetryPolicy;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.CuratorFrameworkFactory;
import org.apache.curator.framework.recipes.cache.*;
import org.apache.curator.framework.recipes.queue.DistributedQueue;
import org.apache.curator.framework.recipes.queue.QueueBuilder;
import org.apache.curator.framework.recipes.queue.QueueConsumer;
import org.apache.curator.framework.recipes.queue.QueueSerializer;
import org.apache.curator.retry.ExponentialBackoffRetry;
import org.apache.curator.utils.CloseableUtils;
import org.apache.zookeeper.CreateMode;
import org.apache.zookeeper.data.ACL;
import org.apache.zookeeper.data.Stat;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.concurrent.TimeUnit;

/**
 * Created by Administrator on 2017/11/14.
 */
public class ZkHelper {


    private static final Logger logger = LoggerFactory.getLogger(ZkHelper.class);
    private static String CHARSET = "utf-8";

    private CuratorFramework client = null;

    public ZkHelper(String connectionString) {
        client = createSimpleClient(connectionString);
    }

    /**
     * @param connectionString    连接字符串，eg：127.0.0.1:2181
     * @param connectionTimeoutMs connection timeout
     * @param sessionTimeoutMs    session timeout
     */
    public ZkHelper(String connectionString, int connectionTimeoutMs, int sessionTimeoutMs) {
        ExponentialBackoffRetry retryPolicy = new ExponentialBackoffRetry(1000, 3);
        client = createClientWithOptions(connectionString, retryPolicy, connectionTimeoutMs, sessionTimeoutMs);
    }

    /**
     * 构造函数
     *
     * @param connectionString    连接字符串，eg：127.0.0.1:2181
     * @param baseSleepTimeMs     initial amount of time to wait between retries
     * @param maxRetries          max number of times to retry
     * @param connectionTimeoutMs connection timeout
     * @param sessionTimeoutMs    session timeout
     */
    public ZkHelper(String connectionString, int baseSleepTimeMs, int maxRetries, int connectionTimeoutMs, int sessionTimeoutMs) {
        ExponentialBackoffRetry retryPolicy = new ExponentialBackoffRetry(baseSleepTimeMs, maxRetries);
        client = createClientWithOptions(connectionString, retryPolicy, connectionTimeoutMs, sessionTimeoutMs);
    }

    /**
     * 检测连接是否可用
     *
     * @param maxWaitTime 连接等待的最大时间
     * @param units       时间单位
     */
    public boolean canConnect(int maxWaitTime, TimeUnit units) throws InterruptedException {
        return client.blockUntilConnected(maxWaitTime, units);
    }

    /**
     * 检测连接是否可用：一直等待，直到连接可用
     */
    public void blockUntilConnected() throws InterruptedException {
        client.blockUntilConnected();
    }

    /**
     * 关闭连接
     */
    public void closeClient() {
        CloseableUtils.closeQuietly(client);
    }

    /**
     * 创建一个客户端连接
     *
     * @param connectionString
     * @return
     */
    private CuratorFramework createSimpleClient(String connectionString) {
        // these are reasonable arguments for the ExponentialBackoffRetry.
        // The first retry will wait 1 second - the second will wait up to 2 seconds - the
        // third will wait up to 4 seconds.
        ExponentialBackoffRetry retryPolicy = new ExponentialBackoffRetry(1000, 3);
        // The simplest way to get a CuratorFramework instance. This will use default values.
        // The only required arguments are the connection string and the retry policy
        CuratorFramework client = CuratorFrameworkFactory.newClient(connectionString, retryPolicy);
        client.start();
        return client;
    }


    /**
     * 创建一个客户端连接
     *
     * @param connectionString
     * @param baseSleepTimeMs
     * @param maxRetries
     * @return
     */
    private CuratorFramework createSimpleClient(String connectionString, int baseSleepTimeMs, int maxRetries) {
        // these are reasonable arguments for the ExponentialBackoffRetry.
        // The first retry will wait 1 second - the second will wait up to 2 seconds - the
        // third will wait up to 4 seconds.
        ExponentialBackoffRetry retryPolicy = new ExponentialBackoffRetry(baseSleepTimeMs, maxRetries);
        // The simplest way to get a CuratorFramework instance. This will use default values.
        // The only required arguments are the connection string and the retry policy
        CuratorFramework client = CuratorFrameworkFactory.newClient(connectionString, retryPolicy);

        client.start();
        return client;
    }

    private CuratorFramework createClientWithOptions(String connectionString, RetryPolicy retryPolicy, int connectionTimeoutMs, int sessionTimeoutMs) {
        // using the CuratorFrameworkFactory.builder() gives fine grained control
        // over creation options. See the CuratorFrameworkFactory.Builder javadoc details


        CuratorFramework client = CuratorFrameworkFactory.builder().connectString(connectionString)
                .retryPolicy(retryPolicy)
                .connectionTimeoutMs(connectionTimeoutMs)
                .sessionTimeoutMs(sessionTimeoutMs)
                // etc. etc.
                .build();
        client.start();
        return client;
    }


    /**
     * 获取节点的数据
     *
     * @param path
     * @return
     */
    public String getData(String path) {
        try {
            return new String(client.getData().forPath(path), CHARSET);
        } catch (Exception e) {
            logger.error("getData error", e);
        }
        return null;
    }

    /**
     * 获取节点的数据
     *
     * @param path
     * @return
     */
    public byte[] getByteData(String path) {
        long start = System.currentTimeMillis();
        try {
            byte[] result = client.getData().forPath(path);
            System.out.println("getByteData:"+(System.currentTimeMillis() - start));
            return result;
        } catch (Exception e) {
            logger.error("getByteData error", e);
        }
        return null;
    }

    /**
     * 获取节点的数据
     *
     * @param path
     * @return
     */
    public String getData(String path, String charset) {
        try {
            return new String(client.getData().forPath(path), charset);
        } catch (Exception e) {
            logger.error("getData error", e);
        }
        return null;
    }

    /**
     * 获取子节点
     *
     * @param path
     * @return
     */
    public List<String> getChildren(String path) {
        long start = System.currentTimeMillis();
        try {
            List<String> result = client.getChildren().forPath(path);
            System.out.println("getChildren:"+(System.currentTimeMillis() - start));
            return result;
        } catch (Exception e) {
            logger.error("getChildren error", e);
        }
        return null;
    }

    /**
     * 获取客户端连接状态
     * LATENT    not init
     * STARTED    has been init
     * STOPPED  close has been called
     *
     * @return
     */
    public String getClientState() {
        try {
            return client.getState().name();
        } catch (Exception e) {
            logger.error("getClientState error", e);
        }
        return null;
    }

    /**
     * 检查节点是否存在
     *
     * @param path
     * @return
     */
    public boolean checkExists(String path) {
        long start = System.currentTimeMillis();
        try {
            boolean flag = null != client.checkExists().forPath(path);
            System.out.println("checkExists:"+(System.currentTimeMillis() - start));
            return flag;
        } catch (Exception e) {
            logger.error("checkExists error", e);
        }
        return false;
    }

    /**
     * 获取stat对象
     *
     * @param path
     * @return
     */
    public Stat getStat(String path) {
        try {
            return client.checkExists().forPath(path);
        } catch (Exception e) {
            logger.error("getStat error", e);
        }
        return null;
    }

    /**
     * 设置acl
     *
     * @return
     */
    public Stat setACL(List<ACL> aclList, String path) {
        try {
            return client.setACL().withACL(aclList).forPath(path);
        } catch (Exception e) {
            logger.error("setACL error", e);
        }
        return null;
    }

    /**
     * 设置acl
     *
     * @return
     */
    public List<ACL> getACL(String path) {
        try {
            return client.getACL().forPath(path);
        } catch (Exception e) {
            logger.error("getACL error", e);
        }
        return null;
    }

    /**
     * 增加临时节点
     *
     * @param path
     * @param data
     * @return
     */
    public String addEphemeral(String path, String data) {
        long start = System.currentTimeMillis();
        try {
            String result = client.create().creatingParentsIfNeeded().withMode(CreateMode.EPHEMERAL).forPath(path, data.getBytes(CHARSET));
            System.out.println("addEphemeral:"+(System.currentTimeMillis() - start));
            return result;
        } catch (Exception e) {
            logger.error("addEphemeral error", e);
        }
        return null;
    }

    public String addEphemeral(String path, byte[] data) {
        long start = System.currentTimeMillis();
        try {
            String result = client.create().creatingParentsIfNeeded().withMode(CreateMode.EPHEMERAL).forPath(path, data);
            System.out.println("addEphemeral:"+(System.currentTimeMillis() - start));
            return result;
        } catch (Exception e) {
            logger.error("addEphemeral error", e);
        }
        return null;
    }

    /**
     * 增加临时有序节点
     *
     * @param path
     * @param data
     * @return
     */
    public String addEphemeralSequential(String path, String data) {
        try {
            return client.create().creatingParentsIfNeeded().withMode(CreateMode.EPHEMERAL_SEQUENTIAL).forPath(path, data.getBytes(CHARSET));
        } catch (Exception e) {
            logger.error("addEphemeralSequential error", e);
        }
        return null;
    }

    public String addEphemeralSequential(String path, byte[] data) {
        try {
            return client.create().creatingParentsIfNeeded().withMode(CreateMode.EPHEMERAL_SEQUENTIAL).forPath(path, data);
        } catch (Exception e) {
            logger.error("addEphemeralSequential error", e);
        }
        return null;
    }

    /**
     * 增加持久节点
     *
     * @param path
     * @param data
     * @return
     */
    public String addPersistent(String path, String data) {
        try {
            return client.create().creatingParentsIfNeeded().withMode(CreateMode.PERSISTENT).forPath(path, data.getBytes(CHARSET));
        } catch (Exception e) {
            logger.error("addPersistent error", e);
        }
        return null;
    }

    public String addPersistent(String path, byte[] data) {
        try {
            return client.create().creatingParentsIfNeeded().withMode(CreateMode.PERSISTENT).forPath(path, data);
        } catch (Exception e) {
            logger.error("addPersistent error", e);
        }
        return null;
    }

    /**
     * 增加持久有序的节点
     *
     * @param path
     * @param data
     * @return
     */
    public String addPersistentSequential(String path, String data) {
        try {
            return client.create().creatingParentsIfNeeded().withMode(CreateMode.PERSISTENT_SEQUENTIAL).forPath(path, data.getBytes(CHARSET));
        } catch (Exception e) {
            logger.error("addPersistentSequential error", e);
        }
        return null;
    }

    public String addPersistentSequential(String path, byte[] data) {
        try {
            return client.create().creatingParentsIfNeeded().withMode(CreateMode.PERSISTENT_SEQUENTIAL).forPath(path, data);
        } catch (Exception e) {
            logger.error("addPersistentSequential error", e);
        }
        return null;
    }

    /**
     * 更新节点内容
     *
     * @param path
     * @param data
     * @return
     */
    public Stat update(String path, String data) {
        try {
            return client.setData().forPath(path, data.getBytes(CHARSET));
        } catch (Exception e) {
            logger.error("update error", e);
        }
        return null;
    }

    /**
     * 更新节点数据，data统一使用utf-8 编码
     *
     * @param path
     * @param data
     * @return
     */
    public Stat update(String path, byte[] data) {
        try {
            return client.setData().forPath(path, data);
        } catch (Exception e) {
            logger.error("update error", e);
        }
        return null;
    }

    /**
     * @param path
     */
    public void remove(String path) {
        try {
            client.delete().deletingChildrenIfNeeded().forPath(path);
        } catch (Exception e) {
            logger.error("remove error", e);
        }
    }

    /**
     * @param path
     */
    public void removeChildrenIfNeeded(String path) {

        try {
            client.delete().deletingChildrenIfNeeded().forPath(path);
        } catch (Exception e) {
            logger.error("removeChildrenIfNeeded error", e);
        }
    }

    /**
     * 增加或者更新节点数据
     *
     * @param path
     * @param data
     * @return
     */
    public boolean addOrUpdatePersistentNode(String path, byte[] data) {
        boolean flag = false;
        if (checkExists(path)) {
            Stat stat = update(path, data);
            if (stat != null) {
                flag = true;
            }
        } else {
            String result = addPersistent(path, data);
            if (null != result && result.length() > 0) {
                flag = true;
            }
        }
        return flag;
    }

    /**
     * 增加或者更新节点数据
     *
     * @param path
     * @param data
     * @return
     */
    public boolean addOrUpdatePersistentNode(String path, String data) {
        boolean flag = false;
        if (checkExists(path)) {
            Stat stat = update(path, data);
            if (stat != null) {
                flag = true;
            }
        } else {
            String result = addPersistent(path, data);
            if (null != result && result.length() > 0) {
                flag = true;
            }
        }
        return flag;
    }

    /**
     * 增加或者更新节点数据
     *
     * @param path
     * @param data
     * @return
     */
    public boolean addOrUpdateEphemeralNode(String path, String data) {
        boolean flag = false;
        if (checkExists(path)) {
            Stat stat = update(path, data);
            if (stat != null) {
                flag = true;
            }
        } else {
            String result = addEphemeral(path, data);
            if (null != result && result.length() > 0) {
                flag = true;
            }
        }
        return flag;
    }

    /**
     * 增加或者更新节点数据
     *
     * @param path
     * @param data
     * @return
     */
    public boolean addOrUpdateEphemeralNode(String path, byte[] data) {
        boolean flag = false;
        if (checkExists(path)) {
            Stat stat = update(path, data);
            if (stat != null) {
                flag = true;
            }
        } else {
            String result = addEphemeral(path, data);
            if (null != result && result.length() > 0) {
                flag = true;
            }
        }
        return flag;
    }

    public DistributedQueue getDistributedQueue(QueueConsumer<String> consumer, QueueSerializer<String> serializer, String queuePath) throws Exception {
        DistributedQueue distributedQueue = QueueBuilder.builder(client, consumer, serializer, queuePath).buildQueue();
        distributedQueue.start();
        return distributedQueue;
    }

    public DistributedQueue getDistributedQueue4Byte(QueueConsumer<byte[]> consumer, QueueSerializer<byte[]> serializer, String queuePath) throws Exception {
        DistributedQueue distributedQueue = QueueBuilder.builder(client, consumer, serializer, queuePath).buildQueue();
        distributedQueue.start();
        return distributedQueue;
    }

    /**
     * 添加监听
     *
     * @param listener
     * @param path
     * @return
     */
    public TreeCache addListener(TreeCacheListener listener, String path) {
        TreeCache treeCache = new TreeCache(client, path);
        treeCache.getListenable().addListener(listener);
        try {
            treeCache.start();
        } catch (Exception e) {
            logger.error("addListener error", e);
        }
        return treeCache;
    }

    /**
     * 添加监听
     *
     * @param listener
     * @param path
     * @return
     */
    public PathChildrenCache addListener(PathChildrenCacheListener listener, String path) {
        PathChildrenCache pathChildrenCache = new PathChildrenCache(client, path, true);
        pathChildrenCache.getListenable().addListener(listener);
        try {
            pathChildrenCache.start();
        } catch (Exception e) {
            logger.error("addListener error", e);
        }
        return pathChildrenCache;
    }

    /**
     * 添加监听
     *
     * @param listener
     * @param path             the full path to the node to cache
     * @param dataIsCompressed if true, data in the path is compressed
     * @return
     */
    public NodeCache addListener(NodeCacheListener listener, String path, boolean dataIsCompressed) {
        NodeCache nodeCache = new NodeCache(client, path, dataIsCompressed);
        nodeCache.getListenable().addListener(listener);
        try {
            nodeCache.start();
        } catch (Exception e) {
            logger.error("addListener error", e);
        }
        return nodeCache;
    }


    /**
     * 移除监听
     *
     * @param treeCache
     */
    public void removeListener(TreeCache treeCache) {
        if (null != treeCache) {
            treeCache.close();
        }
    }
}
