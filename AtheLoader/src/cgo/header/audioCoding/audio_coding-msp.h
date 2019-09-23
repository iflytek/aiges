/** 
 * @file	coding_define.h
 * @brief	
 * 
 *  This file defined several interfaces for audio encoding and decoding.
 * 
 * @author	hlli
 * @version	1.0
 * @date	2010-03-22
 * 
 * @see		
 * 
 * <b>History:</b><br>
 * <table>
 *  <tr> <th>Version	<th>Date		<th>Author	<th>Notes</tr>
 *  <tr> <td>1.0		<td>2010-03-22	<td>hlli	<td>Create this file</tr>
 * </table>
 * 
 */

#ifndef __AUDIO_CODING_H__
#define __AUDIO_CODING_H__

static float amr_compres_ratio[] = { 24.615385f			//mode=0, MR475
								   , 22.857143f			//mode=1, MR515
								   , 20.0f				//mode=2, MR59
								   , 17.777778f			//mode=3, MR67
								   , 16.0f				//mode=4, MR74
								   , 15.238095f			//mode=5, MR795
								   , 11.851852f			//mode=6, MR102
								   , 10.0f				//mode=7, MR122
								   , 1.0f };			//MRDTX

static float amrwb_compres_ratio[] = { 35.555556f		//mode=0
								     , 26.666667f		//mode=1
									 , 19.393940f		//mode=2
									 , 17.297297f		//mode=3
									 , 15.609756f		//mode=4
									 , 13.617021f		//mode=5
									 , 12.549020f		//mode=6
									 , 10.847458f		//mode=7
									 , 1.0f };			//MRDTX

static float speex_compres_ratio[] = { 45.71f		//mode=0
									 , 29.09f		//mode=1
									 , 20.0f		//mode=2
									 , 15.24f		//mode=3
									 , 15.24f		//mode=4
									 , 11.03f		//mode=5
									 , 11.03f		//mode=6
									 , 8.21f		//mode=7
									 , 8.21f		//mode=8
									 , 6.81f		//mode = 9
									 , 5.08f };		//mode = 10

static float speexwb_compres_ratio[] = { 58.18f		//mode=0
									   , 40.0f		//mode=1
									   , 30.48f		//mode=2
									   , 24.62f		//mode=3
									   , 19.39f		//mode=4
									   , 14.88f		//mode=5
									   , 12.08f		//mode=6
									   , 10.49f		//mode=7
									   , 9.01f		//mode=8
									   , 7.36f		//mode = 9
									   , 5.98f };	//mode = 10

#ifdef __cplusplus
extern "C" 
{
#endif /* C++ */

/** 
 * @fn		AudioCodingInit
 * @brief   
 * 
 *  initialize the audio encoder/decoder.
 * 
 * @return	int							- Return 0 in success, otherwise return error code.
 * @param	const char* libNames		- [in] name of the librarys.
 * @param	void * pReserved			- [in,out] Reserved, must be NULL.
 * @see		
 */
int AudioCodingInit(const char* libNames, void* pReserved);
typedef int (* Proc_AudioCodingInit)(const char* libNames, void* pReserved);

/** 
 * @fn		AudioCodingStart
 * @brief   
 * 
 *  start to encode or decode speech.
 * 
 * @return	int							- Return 0 in success, otherwise return error code.
 * @param	void** codingHanle			- [out]handle of instance for encoding/decoding; if failed, return NULL.
 * @param	const char* algorithmName	- [in] name of the algorithm, such as amr, amr-wb-fx, etc.
 * @param	const char* params			- [in] parameters for encoding/decoding, using ",;\n" as a spliter.
 * @see		
 */
int AudioCodingStart(void** codingHandle, const char* algorithmName, const char* params);
typedef int (* Proc_AudioCodingStart)(void** codingHandle, const char* algorithmName, const char* params);

/** 
 * @fn		AudioCodingEncode
 * @brief	encode speech
 * 
 *  Encode the speech data using algorithm specified in AudioCodingStart function.
 * 
 * @return	int									- Return 0 in success, otherwise return error code.
 * @param	void* codingHandle					- [in] handle of the encoder, returned by AudioCodingInit function.
 * @param	const char* speech					- [in] raw-format speech data buffer.
 * @param	unsigned int speechLen				- [in] length of raw-format speech data.
 * @param	char* compressedAudio				- [in] buffer for compressed audio data.
 * @param	unsigned int* compressedAudioLen	- [in/out] in: length of compressed-audio buffer, out: length of compressed-audio data.
 * @param	const char* encodeParams			- [in] params for encoding, using ",;\n" as a spliter.
 * @see		
 */
int AudioCodingEncode(void* codingHandle, const char* speech, unsigned int speechLen, char* compressedAudio, unsigned int* compressedAudioLen, const char* encodeParams);
typedef int (* Proc_AudioCodingEncode)(void* codingHandle, const char* speech, unsigned int speechLen, char* compressedAudio, unsigned int* compressedAudioLen, const char* encodeParams);

/** 
 * @fn		AudioCodingDecode
 * @brief	decode speech
 * 
 *  Decode the speech data using algorithm specified in AudioCodingStart function.
 * 
 * @return	int							- Return 0 in success, otherwise return error code.
 * @param	void* codingHandle			- [in] handle of the decoder, returned by AudioCodingInit().
 * @param	const char* compressedAudio	- [in] compressed audio data.
 * @param	unsigned int amrLen			- [in] length of compressed-audio data.
 * @param	char* speech				- [in] buffer for raw-format speech.
 * @param	unsigned int* speechLen		- [in/out] in: length of raw-format speech buffer, out: length of raw-format speech.

 * @see		
 */
int AudioCodingDecode(void* codingHandle, const char* compressedAudio, unsigned int compressedAudioLen, char* speech, unsigned int* speechLen);
typedef int (* Proc_AudioCodingDecode)(void* codingHandle, const char* compressedAudio, unsigned int compressedAudioLen, char* speech, unsigned int* speechLen);

/** 
 * @fn		AudioCodingEnd
 * @brief
 * 
 *  end to encode and decode speech.
 * 
 * @return	int					- Return 0 in success, otherwise return error code.
 * @param	void* codingHandle	- [in] handle of encoder/decoder, returned by AudioCodingInit().
 * @see		
 */
int AudioCodingEnd(void* codingHandle);
typedef int (* Proc_AudioCodingEnd)(void* codingHandle);

/** 
 * @fn		AudioCodingFini
 * @brief
 * 
 *  uninitialize the amr encoder or decoder.
 * 
 * @return	int			- Return 0 in success, otherwise return error code.
 * @see		
 */
int AudioCodingFini();
typedef int (* Proc_AudioCodingFini)();

#ifdef __cplusplus
} /* extern "C" */
#endif /* C++ */

#endif	/* __AUDIO_CODING_H__ */
