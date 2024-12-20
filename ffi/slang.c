#include <stdio.h>
#include <string.h>
#include "slang.h"

int main(int argc, char *argv[]) {

	if (argc < 2) {
		printf("Usage: %s <lang-code>\n", argv[0]);
		return 0;
	}

	size_t langCodeLen = strlen(argv[1]);

    SlangString langCode = {argv[1], langCodeLen};
	struct SlangParseLangResult result = SlangParseLang(langCode);

	printf("SlangParseLang('%s'): code = %d\n%s\n", argv[1], result.errcode, result.tabstr);
	return 0;
}