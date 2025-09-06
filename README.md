### 创建菜单
/api/wx/push?signature=f58a9a20be241f8c85df2445948a32096f449329&timestamp=1757051260&nonce=1740609658&openid=oIin3168TLKg1X8OU2xBBWLlMEdI

```json
{
  "button": [
    {
      "type": "click",
      "name": "考勤打卡",
      "key": "16:00:6f:83:35:e1"
    },
    {
      "name": "功能",
      "sub_button": [
        {
          "type": "view",
          "name": "504.7000",
          "key": "WORK_003",
          "url": "http://uuxia.cn:7000?auth_code=oIin3168TLKg1X8OU2xBBWLlMEdI"
        },
        {
          "type": "view",
          "name": "gz.7000",
          "key": "WORK_004",
          "url": "http://uuxia.cn:6633?auth_code=oIin3168TLKg1X8OU2xBBWLlMEdI"
        },
        {
          "type": "view",
          "name": "baidu",
          "key": "WORK_005",
          "url": "https://www.baidu.com"
        },
        {
          "type": "click",
          "name": "获取链接",
          "key": "WORK_007"
        },
        {
          "type": "view",
          "name": "clife.7000",
          "key": "WORK_002",
          "url": "http://uuxia.cn:6615?auth_code=oIin3168TLKg1X8OU2xBBWLlMEdI"
        }
      ]
    }
  ]
}
```