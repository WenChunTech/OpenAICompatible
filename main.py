import hashlib
import time
import requests
import json

import secrets

def r(e=21):
    try:
        e = int(e)  # 模拟 JS 中的 e|=0 (强制转为整数)
    except (TypeError, ValueError):
        e = 0

    # 生成加密安全的随机字节数组 (等效 crypto.getRandomValues)
    n = secrets.token_bytes(e)

    # 字符集 (与 JS 完全一致)
    char_set = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict"

    # 倒序遍历字节数组 (模拟 JS 中 e-- 的递减遍历)
    result = []
    for i in range(e-1, -1, -1):
        # 63 & n[i] 等效于 n[i] % 64 (取低6位)
        index = n[i] & 63
        result.append(char_set[index])

    return ''.join(result)

url = "https://ai-api.dangbei.net/ai-search/configApi/v1/updateRecentModel"

payload = '{"model":"qwen3-235b-a22b"}'

timstamp = str(int(time.time()))
print(timstamp)
nonce = r()
s = timstamp + payload.replace(" ","") + 'vGU2OOUTVhBYHU1t4TazA'
# s = timstamp + payload.replace(" ","") + nonce
# print(s)
sign = hashlib.md5(s.encode('utf-8')).hexdigest().upper()
headers = {
  'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0',
  'appType': '6',
  'appVersion': '1.1.16',
  'client-ver': '1.0.2',
  'content-type': 'application/json',
  'deviceId': 'f94745995ac935e16ffa37b2dd449f2b_kM4BW_OO',
#   'nonce': nonce,
  'nonce': 'vGU2OOUTVhBYHU1t4TazA',
  'sign': sign,
  'timestamp': timstamp,
#   'timestamp': '1753807133',
  'token': ''
}

response = requests.request("POST", url, headers=headers, data=payload)

print(response.text)
