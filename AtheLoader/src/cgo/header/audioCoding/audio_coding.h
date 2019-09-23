#ifndef __AUDIO_ENCODING_H__
#define __AUDIO_ENCODING_H__

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

typedef int64_t AUDIOCODING_INST;

//#include<stdlib.h>
int AudioCodingInit(const char* libNames, void* pReserved);
typedef int (* Proc_AudioCodingInit)(const char* libNames, void* pReserved);

int AudioCodingStart(AUDIOCODING_INST* codingHandle, const char* algorithmName, const char* params);
typedef int (* Proc_AudioCodingStart)(AUDIOCODING_INST* codingHandle, const char* algorithmName, const char* params);

int AudioCodingEncode(AUDIOCODING_INST codingHandle, const char* speech, unsigned int speechLen, char* compressedAudio, unsigned int* compressedAudioLen, const char* encodeParams);
typedef int (* Proc_AudioCodingEncode)(AUDIOCODING_INST codingHandle, const char* speech, unsigned int speechLen, char* compressedAudio, unsigned int* compressedAudioLen, const char* encodeParams);

int AudioCodingDecode(AUDIOCODING_INST codingHandle, const char* compressedAudio, unsigned int compressedAudioLen, char* speech, unsigned int* speechLen);
typedef int (* Proc_AudioCodingDecode)(AUDIOCODING_INST codingHandle, const char* compressedAudio, unsigned int compressedAudioLen, char* speech, unsigned int* speechLen);

int AudioCodingEnd(AUDIOCODING_INST codingHandle);
typedef int (* Proc_AudioCodingEnd)(AUDIOCODING_INST codingHandle);

int AudioCodingFini();
typedef int (* Proc_AudioCodingFini)();

#ifdef __cplusplus
};

#endif

#endif //__AUDIO_ENCODING_H__