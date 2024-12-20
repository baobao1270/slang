#ifndef SLANG_GO_H
#define SLANG_GO_H
	#include <stddef.h>

	#ifdef __cplusplus
	extern "C" {
	#endif
		#ifndef SLANG_GO_CFFI_DEF_START
		#define SLANG_GO_CFFI_DEF_START
		#endif
			typedef signed char SlangInt8;
			typedef struct {
				const char *p;
				ptrdiff_t n;
			} SlangString;

			typedef struct SlangParseLangResult {
				SlangInt8 errcode;
				char*     tabstr;
			} SlangParseLangResult;

			extern SlangParseLangResult SlangParseLang(SlangString langCode);
		#ifndef SLANG_GO_CFFI_DEF_END
		#define SLANG_GO_CFFI_DEF_END
		#endif
	#ifdef __cplusplus
	}
	#endif
#endif