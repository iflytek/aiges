namespace java com.iflytek.comet.rmq.thrift.generated

enum MTRProtocol
{
  RAW = 0,
  PERSONALIZED,
}

struct MTRMessage
{
  // 消息的主题
  1: string topic;
  // 消息内容
  2: binary body;
  // 消息内容序列化版本, 用于构建和解析消息内容
  3: MTRProtocol protocol;
  // 消息的key, 用以查找消息
  4: string key;
  // 用于用户自定义
  5: i32 flag;
  // 消息的tag, 用于过滤消息
  6: string tag;
  // 该条消息在broker中的offset, 仅在消费接口时有效
  7: i64 offset;
}

struct MTRMonitorInfo
{
  // 消息的消费group
  1: string group;
  // 消息积压数
  2: string diff;
  // 消费速度
  3: string tps;
  // 消费的主题
  4: string topic;
}

enum MTRRPCErrorCode
{
	MTR_SUCCESS = 0,
	MTR_RPC_ERROR_BASE = 33000,
	MTR_RPC_ERROR_NO_MORE_DATA = 33001,
	MTR_RPC_MESSAGE_IS_NULL = 33002,
	MTR_PRODUCE_FAILURE = 33005,
	MTR_TOPIC_NOT_EXITST = 33006,
	MTR_GROUP_NOT_EXITST = 33007,
	MTR_TOPIC_ALREADY_EXITST = 33008,
	MTR_GROUP_ALREADY_EXITST = 33009,
	MTR_RECALL_CONSUME_FINISHED = 33010,
	MTR_RECALL_CONSUME_ERROR = 33011,
	MTR_GROUPS_NOT_EXITST = 33012
	
	MTR_ERROR_BASE = 33100,
	MTR_MESSAGE_IS_NULL = 33102,
	MTR_UNKNOW_HOST_EXCEPTION = 33110,
	MTR_CONSUME_NO_MORE_DATA = 33101,
	MTR_PRODUCE_UNKNOW_FAILURE = 33105,
	MTR_PRODUCE_CLIENT_EXCEPTION = 33106,
	MTR_PRODUCE_BROKER_EXCEPTION = 33107,
	MTR_PRODUCE_REMOTING_EXCEPTION = 33108,
	MTR_PRODUCE_INTERRUPTED_EXCEPTION = 33109,
	MTR_CONSUME_SUBSCRIBE_EXCEPTION = 33113,
	MTR_PRODUCE_START_EXCEPTION = 33111,
	MTR_CONSUME_START_EXCEPTION = 33112,	
	MTR_RPT_PRODUCE_CLIENT_EXCEPTION = 33114,
	MTR_RPT_CONSUME_START_EXCEPTION = 33115,
	MTR_RPT_CONSUME_SUBSCRIBE_EXCEPTION = 33116,
	MTR_RPT_PRODUCE_REMOTING_EXCEPTION = 33117,
	MTR_RPT_PRODUCE_BROKER_EXCEPTION = 33118,
	MTR_RPT_PRODUCE_INTERRUPTED_EXCEPTION = 33119
}

exception MTRRPCException
{
  1: required MTRRPCErrorCode id;
  2: string msg;
  3: i32 errorCode;
}

service MTRMessageService {
  i64 produce(1: MTRMessage msg, 2: bool flush) throws (1: MTRRPCException rpcException);
  MTRMessage consume(1: string topic, 2: string group) throws (1: MTRRPCException rpcException);
  MTRMessage consumeWithOffset(1: string topic, 2: i64 offset) throws (1: MTRRPCException rpcException);
  MTRMessage reConsumeByDate(1: string topic, 2: string beginDate, 3: string endDate) throws (1: MTRRPCException rpcException);
  string showTopicInfo(1: string topic) throws (1: MTRRPCException rpcException);
  i64 createTopic(1: string topic)throws (1: MTRRPCException rpcException);
  i64 deleteTopic(1: string topic)throws (1: MTRRPCException rpcException);
  i64 createGroup(1: string topic, 2: string group)throws (1: MTRRPCException rpcException);
  i64 deleteGroup(1: string topic, 2: string group)throws (1: MTRRPCException rpcException);
  list<MTRMonitorInfo> fetchMonitorInfo() throws (1: MTRRPCException rpcException);
}

