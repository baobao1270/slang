from __future__ import print_function
import os
import sys
from cffi import FFI

CFFI_DEF_START       = "#define SLANG_GO_CFFI_DEF_START"
CFFI_DEF_END         = "#ifndef SLANG_GO_CFFI_DEF_END"
CFFI_DEF_START_CLOSE = "#endif"

ffi = FFI()
with open(os.path.join(os.path.dirname(__file__), "slang.h")) as f:
	f = f.read()
	start = f.find(CFFI_DEF_START) + len(CFFI_DEF_START)
	end = f.find(CFFI_DEF_END)
	f = f[start:end].strip()
	if f.startswith(CFFI_DEF_START_CLOSE):
		f = f[len(CFFI_DEF_START_CLOSE):].strip()
	ffi.cdef(f)
	lib_slang = ffi.dlopen(os.path.join(os.path.dirname(__file__), "../bin/slang.so"))

if __name__ == "__main__":
	if len(sys.argv) != 2:
		print("Usage: %s <lang>" % sys.argv[0])
		sys.exit(1)
	lang_code = sys.argv[1].encode("utf-8")
	result = lib_slang.SlangParseLang((ffi.new("char[]", lang_code), len(lang_code)))
	error_code = int(result.errcode)
	if error_code != 0:
		print("error code: 0x%02x" % error_code)
		sys.exit(1)
	else:
		tabstr = ffi.string(result.tabstr).decode("utf-8")
		name, location, clidHex, bcp47, winid, iso639_1, iso639_2, iso639_3 = tabstr.split("\t")
		print({
			"name": name,
			"location": location,
			"clidHex": clidHex,
			"clid": int(clidHex, 16),
			"bcp47": bcp47,
			"winid": winid,
			"iso639_1": iso639_1,
			"iso639_2": iso639_2,
			"iso639_3": iso639_3,
		})
