/*
	数据处理模块,用于进行数据预处理&后处理操作;
	包含：音频重采样,音频头处理等操作;
*/
package dp

//	音频重采样接口
func NewResampler() *AudioResampler {
	return &AudioResampler{}
}
