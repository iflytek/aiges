#ifndef UTILS_H
#define UTILS_H
#include <vector>
#include <iconv.h>
#include <string>
#include<string.h>
#include <ctime>
#include <iostream>
#include <dlfcn.h>


void SplitString(const std::string &s, std::vector<std::string> &v, const std::string &c)
{
  std::string::size_type pos1, pos2;
  pos2 = s.find(c);
  pos1 = 0;
  while (std::string::npos != pos2)
  {
    v.push_back(s.substr(pos1, pos2 - pos1));

    pos1 = pos2 + c.size();
    pos2 = s.find(c, pos1);
  }
  if (pos1 != s.length())
    v.push_back(s.substr(pos1));
}

//gbk to utf8
std::string code_convert(const char *from_charset, const char *to_charset, const char *inbuf) //
{
  iconv_t cd;
  int rc;
  char *in = const_cast<char *>(inbuf);
  size_t in_len = strlen(inbuf);

  size_t dst_size = in_len * sizeof(wchar_t);
  char *dst = new char[dst_size];
  if (dst == NULL)
    return "";

  memset(dst, 0, dst_size);
  char *out = dst;
  cd = iconv_open(to_charset, from_charset);
  if (cd == 0)
  {
    delete[] dst;
    return "";
  }

  if (iconv(cd, &in, &in_len, &out, &dst_size) == -1)
  {
    iconv_close(cd);
    delete[] dst;
    return "";
  }
  else
  {
    std::string rslt(dst, strlen(dst));
    delete[] dst;
    iconv_close(cd);
    return rslt;
  }
}

std::string gbk_to_utf8(const char *gb_str)
{
  return code_convert("gbk", "utf-8", gb_str);
}

// @return library handle
void *cLibOpen(const char *libName, char **err)
{
  void *hdl = dlopen(libName, RTLD_NOW);
  if (hdl == NULL)
  {
    *err = (char *)dlerror();
  }
  return hdl;
}
// @return symbol address
void *cLibLoad(void *hdl, const char *sym, char **err)
{
  void *addr = dlsym(hdl, sym);
  if (addr == NULL)
  {
    *err = (char *)dlerror();
  }
  return addr;
}
// @close library
int cLibClose(void *hdl)
{
  int ret = dlclose(hdl);
  if (ret != 0)
    return -1;
  return 0;
}

#endif